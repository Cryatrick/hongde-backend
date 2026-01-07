package model

type PesertaUjian struct{
	PesertaId string
	SiswaId string 
	SiswaNama string
	UjianId int
	UjianNama string
	StartedAt string
	EndAt string
	IsLate string
	FinalScore float64
	SoalArray []SoalPesertaUjian
	UserUpdate string
	LastUpdate string
} 

type SoalPesertaUjian struct {
	PesertaId string
	SoalId int
	SoalData SoalList
	UrutanSoal int
	JawabanSiswa string
}