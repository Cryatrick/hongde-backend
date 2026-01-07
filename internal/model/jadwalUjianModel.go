package model

type JadwalUjian struct {
	JadwalId int `json:"jdw_id"`
	NamaJadwal string `json:"nama_jadwal", validate:"required"`
	BankSoal int `json:"bank_soal", validate:"required"`
	NamaBankSoal string
	JenisSoal string `json:"jenis_soal", validate:"required"`
	JumlahSoal int `json:"jumlah_soal", validate:"required"`
	TanggalMulai string `json:"tanggal_mulai", validate:"required"`
	JamMulai string `json:"jam_mulai", validate:"required"`
	DurasiUjian int `json:"durasi_ujian", validate:"required"`
	ToleransiTerlambat int `json:"toleransi_terlambat", validate:"required"`
	TokenUjian string 
	KetentuanWaktu string `json:"ketentuan_waktu", validate:"required"`
	UserUpdate string 
	LastUpdate string
}

type JadwalUjianSiswa struct {
	JadwalId int 
	NamaJadwal string 
	TanggalMulai string 
	JamMulai string 
	DurasiUjian int 
	ToleransiTerlambat int
	KetentuanWaktu string
	PesertaId string
	StartPeserta string
	EndPeserta string
	ExamStatus string
}