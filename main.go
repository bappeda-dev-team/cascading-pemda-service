package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {
	dsn := os.Getenv("PERENCANAAN_DB_URL")
	if dsn == "" {
		log.Fatal("PERENCANAAN_DB_URL env tidak terdefinisi")
	}

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("[FATAL] Error connecting to db: %v", err)
	}

	log.Printf("koneksi ke database berhasil")
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Gagal terhubung ke database dalam 10 detik: %v", err)
		log.Printf("Mencoba koneksi ulang...")

		// Coba lagi dengan timeout yang lebih lama
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err = db.PingContext(ctx)
		if err != nil {
			db.Close()
			log.Fatalf("Koneksi database gagal setelah percobaan ulang: %v", err)
		}
	}

	log.Print("Berhasil terhubung ke database")
	log.Printf("Max Open Connections: %d", db.Stats().MaxOpenConnections)
	log.Printf("Open Connections: %d", db.Stats().OpenConnections)
	log.Printf("In Use Connections: %d", db.Stats().InUse)
	log.Printf("Idle Connections: %d", db.Stats().Idle)
}

func getRencanaKinerjaPokin(idPokin int) ([]RencanaKinerjaAsn, error) {
	query := `
		SELECT rekin.id,
		       rekin.nama_rencana_kinerja,
		       pegawai.nama,
		       pegawai.nip,
		       subkegiatan.kode_subkegiatan,
		       subkegiatan.nama_subkegiatan,
		       rinbel.anggaran
		FROM tb_rencana_kinerja rekin
		JOIN tb_pegawai pegawai ON pegawai.nip = rekin.pegawai_id
		JOIN tb_subkegiatan_terpilih sub_rekin ON sub_rekin.rekin_id = rekin.id
		LEFT JOIN tb_subkegiatan subkegiatan ON subkegiatan.kode_subkegiatan = sub_rekin.kode_subkegiatan
		JOIN tb_rencana_aksi renaksi ON renaksi.rencana_kinerja_id = rekin.id
		JOIN tb_rincian_belanja rinbel ON rinbel.renaksi_id = renaksi.id
		JOIN tb_pohon_kinerja pokin ON rekin.id_pohon = pokin.id
		WHERE pokin.id = ?`

	rows, err := db.Query(query, idPokin)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var rekins []RencanaKinerjaAsn

	for rows.Next() {
		var rekin RencanaKinerjaAsn
		var kodeSub, namaSub sql.NullString
		var pagu sql.NullInt64

		if err := rows.Scan(
			&rekin.IdRekin,
			&rekin.RencanaKinerja,
			&rekin.NamaPelaksana,
			&rekin.NIPPelaksana,
			&kodeSub,
			&namaSub,
			&pagu,
		); err != nil {
			log.Printf("[ERROR] scan rekin error: %v", err)
			return nil, fmt.Errorf("scan error: %w", err)
		}

		// Handle NULL dengan NullString/NullInt64
		if kodeSub.Valid {
			rekin.SubkegiatanRekin.KodeSubkegiatan = kodeSub.String
		}
		if namaSub.Valid {
			rekin.SubkegiatanRekin.NamaSubkegiatan = namaSub.String
		}
		if pagu.Valid {
			rekin.Pagu = Pagu(pagu.Int64)
		}

		rekins = append(rekins, rekin)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return rekins, nil
}

func findPokinById(idPokin int) (PohonKinerjaPemda, error) {
	query := `SELECT id, tahun, nama_pohon, kode_opd, jenis_pohon, keterangan, status
			  FROM tb_pohon_kinerja
			  WHERE tahun = ? AND clone_from = ? LIMIT 1`

	var pokin PohonKinerjaPemda
	err := db.QueryRow(query, 2025, idPokin).Scan(
		&pokin.IdPohon,
		&pokin.Tahun,
		&pokin.NamaPohon,
		&pokin.KodeOpd,
		&pokin.JenisPohon,
		&pokin.Keterangan,
		&pokin.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return PohonKinerjaPemda{}, nil
		}
		return PohonKinerjaPemda{}, fmt.Errorf("query error: %w", err)
	}

	sasarans, err := getRencanaKinerjaPokin(pokin.IdPohon)
	if err != nil {
		log.Printf("[ERROR] Get Rekin Pokin %d error: %v", idPokin, err)
		return pokin, fmt.Errorf("getRencanaKinerjaPokin(%d): %w", pokin.IdPohon, err)
	}

	pokin.RencanaKinerjas = sasarans

	return pokin, nil
}

func getIndikators(idPokin int) ([]IndikatorPohon, error) {
	indTematikRows, err := db.Query(`SELECT id, indikator FROM tb_indikator WHERE pokin_id = ?`, idPokin)
	if err != nil {
		return nil, fmt.Errorf("query error %v", err)
	}
	var indPt []IndikatorPohon
	for indTematikRows.Next() {
		var ind IndikatorPohon
		if err := indTematikRows.Scan(&ind.IdIndikator, &ind.Indikator); err != nil {
			return nil, fmt.Errorf("query error %v", err)
		}
		// targets
		indTargetRows, err := db.Query(`SELECT id, target, satuan, tahun FROM tb_target WHERE indikator_id = ?`, ind.IdIndikator)
		if err != nil {
			return nil, fmt.Errorf("query error %v", err)
		}
		var tarPt []TargetIndikator
		for indTargetRows.Next() {
			var tar TargetIndikator
			if err := indTargetRows.Scan(&tar.IdTarget, &tar.Target, &tar.Satuan, &tar.Tahun); err != nil {
				return nil, fmt.Errorf("query error %v", err)
			}
			tarPt = append(tarPt, tar)
		}
		// end targets
		ind.Target = tarPt

		indPt = append(indPt, ind)
	}
	return indPt, nil
}

func getChildPokins(parentId int) ([]PohonKinerjaPemda, Pagu, error) {
	rows, err := db.Query(`SELECT id, tahun, nama_pohon, kode_opd, jenis_pohon, keterangan, status
		FROM tb_pohon_kinerja
		WHERE tahun = 2025 AND parent = ?`, parentId)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var childs []PohonKinerjaPemda
	var totalPagu Pagu = 0

	for rows.Next() {
		var pt PohonKinerjaPemda
		if err := rows.Scan(&pt.IdPohon, &pt.Tahun, &pt.NamaPohon, &pt.KodeOpd,
			&pt.JenisPohon, &pt.Keterangan, &pt.Status); err != nil {
			return nil, 0, err
		}

		// ambil indikator
		indCt, err := getIndikators(pt.IdPohon)
		if err != nil {
			return nil, 0, err
		}
		pt.Indikators = indCt

		// operational pemda → ambil rencana kinerja langsung pakai IdPohon
		if pt.JenisPohon == "Operational Pemda" && pt.Status == "disetujui" {
			sourcePokin, err := findPokinById(pt.IdPohon)
			if err != nil {
				return nil, 0, fmt.Errorf("findPokinById(%d): %w", pt.IdPohon, err)
			}
			pt.RencanaKinerjas = sourcePokin.RencanaKinerjas
		}

		// rekursif ambil anaknya
		childTematiks, childPagu, err := getChildPokins(pt.IdPohon)
		if err != nil {
			return nil, 0, err
		}
		pt.Childs = childTematiks

		// hitung pagu node ini sendiri
		var nodePagu Pagu = 0
		for _, rekin := range pt.RencanaKinerjas {
			nodePagu += rekin.Pagu
		}

		// tambahkan pagu anak
		nodePagu += childPagu

		// set pagu node ini sendiri
		pt.Pagu = nodePagu

		// tambahkan ke total pagu parent
		totalPagu += nodePagu

		childs = append(childs, pt)
	}

	return childs, totalPagu, nil
}

func cascadingHandler(w http.ResponseWriter, r *http.Request) {
	// hanya terima GET method
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed, pakai GET", http.StatusMethodNotAllowed)
		return
	}

	// parameter for tematik
	tematikIdStr := r.URL.Query().Get("tematikId")
	if tematikIdStr == "" {
		http.Error(w, "params tematikId is required, misal: ?tematikId=123", http.StatusBadRequest)
		return
	}

	tematikId, err := strconv.Atoi(tematikIdStr)
	if err != nil {
		http.Error(w, "invalid tematikId", http.StatusBadRequest)
		return
	}

	// query pohon tematik
	rows, err := db.Query(`SELECT id, tahun, nama_pohon, kode_opd, jenis_pohon, keterangan,  status
                           FROM tb_pohon_kinerja
                           WHERE level_pohon = 0 AND parent = 0 AND jenis_pohon = 'Tematik' AND id = ? LIMIT 1`, tematikId)
	if err != nil {
		http.Error(w, "query error: "+err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	var totalPagu Pagu
	var list []PohonKinerjaPemda
	for rows.Next() {
		var pt PohonKinerjaPemda
		if err := rows.Scan(&pt.IdPohon, &pt.Tahun, &pt.NamaPohon, &pt.KodeOpd, &pt.JenisPohon, &pt.Keterangan, &pt.Status); err != nil {
			http.Error(w, "scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		indList, err := getIndikators(pt.IdPohon)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pt.Indikators = indList

		childs, pagu, err := getChildPokins(pt.IdPohon)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		totalPagu += pagu

		pt.Childs = childs
		pt.Pagu = totalPagu

		list = append(list, pt)
	}

	response := CascadingPemda{
		Status:  http.StatusOK,
		Message: "Laporan Cascading Pemda Tahun 2025",
		Tematik: list}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	log.Print("CASCADING PEMDA 2025")

	initDB()

	http.HandleFunc("/laporan/cascading_pemda", cascadingHandler)
	http.ListenAndServe(":8080", nil)
}
