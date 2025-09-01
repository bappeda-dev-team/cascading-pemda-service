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
	UrusanPokin       []Urusan            `json:"urusan_pokin,omitempty"`
	BidangUrusanPokin []BidangUrusan      `json:"bidang_urusan_pokin,omitempty"`
	ProgramPokin      []Program           `json:"program_pokin,omitempty"`
	KegiatanPokin     []Kegiatan          `json:"kegiatan_pokin,omitempty"`
	Pagu              Pagu                `json:"pagu"`
	Indikators        []IndikatorPohon    `json:"indikator,omitempty"`
	Childs            []PohonKinerjaPemda `json:"childs,omitempty"`
	RencanaKinerjas   []RencanaKinerjaAsn `json:"rencana_kinerjas,omitempty"`

	Status string `json:"-"`
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
