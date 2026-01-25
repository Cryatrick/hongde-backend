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
	TotalRight int
	TotalWaiting int
	TotalWrong int
	FinalScore float64
	IsRecapped string
	SoalArray []SoalPesertaUjian
	UserUpdate string
	LastUpdate string
} 

type SoalPesertaUjian struct {
	PesertaId string
	SoalId int
	UrutanSoal int
	JawabanSiswa string
	IsRight string
	BobotSoal float64
}