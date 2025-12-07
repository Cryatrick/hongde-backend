package model

type BankSoal struct {
	BankId int `json:"bank_id"`
	NamaBank string `json:"nama_bank", validate:"required"`
	JumlahSoal int 
	UserUpdate string 
	LastUpdate string
}

type SoalList struct {
	SoalId int `json:"soal_id"`
	BankId int `json:"bank_id"`
	PertanyaanSoal string `json:"pertanyaan_soal"`
	GambarSoal string `json:"gambar_soal"`
	UrutanSoal int `json:"urutan_soal"`
	JawabanA string `json:"jawaban_a"`
	GambarA string `json:"gambar_a"`
	JawabanB string `json:"jawaban_b"`
	GambarB string `json:"gambar_b"`
	JawabanC string `json:"jawaban_c"`
	GambarC string `json:"gambar_c"`
	JawabanD string `json:"jawaban_d"`
	GambarD string `json:"gambar_d"`
	JawabanBenar string `json:"jawaban_benar"`
	BobotSoal float64 `json:"bobot_soal"`
	TipeSoal string `json:"tipe_soal"`
	UserUpdate string 
	LastUpdate string
}	