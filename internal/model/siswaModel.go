package model

type SiswaList struct{
	SiswaId string `json:"siswa_id"`
	NamaSiswa string `json:"nama_siswa",validate:"required"`
	EmailSiswa string `json:"email_siswa",validate:"required"`
	NamaMandarin string `json:"nama_mandarin"`
	JenisIdentitas string `json:"jenis_identitas", validate:"required"`
	NoIdentitas string `json:"no_identitas", validate:"required"`
	TempatLahir string `json:"tempat_lahir", validate:"required"`
	TanggalLahir string `json:"tanggal_lahir", validate:"required"`
	TempatTinggal string `json:"tempat_tinggal", validate:"required"`
	NoKontak string `json:"no_kontak", validate:"required"`
	JenisSiswa string `json:"jenis_siswa", validate:"required"`
	PasswordSiswa string  
	UserUpdate string 
	LastUpdate string
} 