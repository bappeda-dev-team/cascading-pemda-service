package main

type JenisPohon string
type Keterangan string
type Tahun int
type Pagu int

type CascadingPemda struct {
	Status  int                 `json:"status"`
	Message string              `json:"message"`
	Tematik []PohonKinerjaPemda `json:"data"`
}

type PohonKinerjaPemda struct {
	IdPohon           int                 `json:"id_pohon"`
	Parent            int                 `json:"parent"`
	Tahun             Tahun               `json:"tahun"`
	NamaPohon         string              `json:"nama_pohon"`
	KodeOpd           string              `json:"kode_opd,omitempty"`
	LevelPohon        int                 `json:"level_pohon"`
	JenisPohon        JenisPohon          `json:"jenis_pohon"`
	Keterangan        Keterangan          `json:"keterangan"`
	TujuanPemda       []TujuanPemda       `json:"tujuan_pemda,omitempty"`
	SasaranPemda      []SasaranPemda      `json:"sasaran_pemda,omitempty"`
	UrusanPokin       []Urusan            `json:"urusan_pokin,omitempty"`
	BidangUrusanPokin []BidangUrusan      `json:"bidang_urusan_pokin,omitempty"`
	ProgramPokin      []Program           `json:"program_pokin,omitempty"`
	KegiatanPokin     []Kegiatan          `json:"kegiatan_pokin,omitempty"`
	Pagu              Pagu                `json:"pagu"`
	Tagging           []TaggingPokin      `json:"tagging"`
	RencanaKinerjas   []RencanaKinerjaAsn `json:"pelaksana,omitempty"`

	Indikators []IndikatorPohon    `json:"indikator,omitempty"`
	Childs     []PohonKinerjaPemda `json:"childs,omitempty"`
	Status     string              `json:"-"`
}

type Urusan struct {
	KodeUrusan string `json:"kode_urusan"`
	NamaUrusan string `json:"nama_urusan"`
}

type BidangUrusan struct {
	KodeBidangUrusan string `json:"kode_bidang_urusan"`
	NamaBidangUrusan string `json:"nama_bidang_urusan"`
}

type Program struct {
	KodeProgram string `json:"kode_program"`
	NamaProgram string `json:"nama_program"`
}

type Kegiatan struct {
	KodeKegiatan string `json:"kode_kegiatan"`
	NamaKegiatan string `json:"nama_kegiatan"`
}

type Subkegiatan struct {
	KodeSubkegiatan string `json:"kode_subkegiatan"`
	NamaSubkegiatan string `json:"nama_subkegiatan"`
}

type RencanaKinerjaAsn struct {
	IdRekin         string `json:"id_rekin"`
	RencanaKinerja  string `json:"rencana_kinerja"`
	NamaPelaksana   string `json:"nama_pelaksana"`
	NIPPelaksana    string `json:"nip_pelaksana"`
	KodeSubkegiatan string `json:"kode_subkegiatan"`
	NamaSubkegiatan string `json:"nama_subkegiatan"`
	Pagu            Pagu   `json:"pagu"`
}

type IndikatorPohon struct {
	IdIndikator string            `json:"id_indikator"`
	IdPokin     string            `json:"id_pokin"`
	Indikator   string            `json:"nama_indikator"`
	Target      []TargetIndikator `json:"targets"`
}

type TargetIndikator struct {
	IdTarget    string `json:"id_target"`
	IndikatorId string `json:"indikator_id"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
	Tahun       Tahun  `json:"tahun"`
}

type TaggingPokin struct {
	Id                int    `json:"id"`
	IdPokin           int    `json:"id_pokin"`
	NamaTagging       string `json:"nama_tagging"`
	KeteranganTagging string `json:"keterangan_tagging"`
	CloneFrom         int    `json:"clone_from"`
}

type TujuanPemda struct {
	IdTujuanPemda int           `json:"id_tujuan_pemda"`
	TujuanPemda   string        `json:"tujuan_pemda"`
	TematikId     int           `json:"tematik_id,omitempty"`
	PeriodeId     int           `json:"periode_id,omitempty"`
	Periode       PeriodeTujuan `json:"periode"`
}

type PeriodeTujuan struct {
	TahunAwal    string `json:"tahun_awal"`
	TahunAkhir   string `json:"tahun_akhir"`
	JenisPeriode string `json:"jenis_periode"`
}

type SasaranPemda struct {
	IdSasaranPemda int           `json:"id_sasaran_pemda"`
	SubtemaId      int           `json:"subtema_id,omitempty"`
	SasaranPemda   string        `json:"sasaran_pemda"`
	PeriodeId      int           `json:"periode_id,omitempty"`
	Periode        PeriodeTujuan `json:"periode"`
}
