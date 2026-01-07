package service

import(
	"fmt"
	"database/sql"
	"strconv"
	// "time"

	"hongde_backend/internal/model"
	"hongde_backend/internal/database"
	"hongde_backend/internal/middleware"
)

// Function Untuk Manajemen Jadwal Ujian
func GetJadwalUjian() ([]model.JadwalUjian, error) {
	QueryData := `SELECT jdw_id,jdw_namajadwal,jdw_banksoal,bs_namabank,jdw_jenissoal,jdw_jumlahsoal,jdw_tgl_mulai,jdw_jam_mulai,jdw_toleransiketerlambatan,jdw_durasi,jdw_token,jdw_ketentuan_waktu,jdw_userupdate,jdw_lastupdate FROM hd_jadwalujian JOIN hd_banksoal ON hd_banksoal.bs_id = hd_jadwalujian.jdw_banksoal WHERE jdw_is_delete = "n" `
	var returnData []model.JadwalUjian

	rows, err := database.DbMain.Query(QueryData)
	if err != nil {
		middleware.LogError(err,"Query Error")
		return nil,err
	}
	defer rows.Close()
	for rows.Next() {
		var rowData model.JadwalUjian
		err = rows.Scan(&rowData.JadwalId, &rowData.NamaJadwal, &rowData.BankSoal, &rowData.NamaBankSoal, &rowData.JenisSoal, &rowData.JumlahSoal, &rowData.TanggalMulai, &rowData.JamMulai, &rowData.ToleransiTerlambat, &rowData.DurasiUjian, &rowData.TokenUjian, &rowData.KetentuanWaktu, &rowData.UserUpdate, &rowData.LastUpdate)
		if err != nil {
			middleware.LogError(err,"Data Scan Error")
			return nil,err
		}
		returnData = append(returnData, rowData)
	}
	return returnData,nil
}

func GetJadwalById(jdwId string)(*model.JadwalUjian, error) {
	row := database.DbMain.QueryRow(`SELECT jdw_id,jdw_namajadwal,jdw_banksoal,bs_namabank,jdw_jenissoal,jdw_jumlahsoal,jdw_tgl_mulai,jdw_jam_mulai,jdw_toleransiketerlambatan,jdw_durasi,jdw_token,jdw_ketentuan_waktu,jdw_userupdate,jdw_lastupdate FROM hd_jadwalujian JOIN hd_banksoal ON hd_banksoal.bs_id = hd_jadwalujian.jdw_banksoal WHERE (jdw_id = ? OR jdw_token = ?) AND jdw_is_delete = "n"`,jdwId,jdwId)
	var returnData = &model.JadwalUjian{}
	err := row.Scan(&returnData.JadwalId, &returnData.NamaJadwal, &returnData.BankSoal, &returnData.NamaBankSoal, &returnData.JenisSoal, &returnData.JumlahSoal, &returnData.TanggalMulai, &returnData.JamMulai, &returnData.ToleransiTerlambat, &returnData.DurasiUjian, &returnData.TokenUjian, &returnData.KetentuanWaktu, &returnData.UserUpdate, &returnData.LastUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		// Log and return actual error
		middleware.LogError(err, "Data Scan Error")
		return nil, err
	}
	return returnData, nil
}

func InsertJadwalUjian(jadwalData *model.JadwalUjian)(bool, error) {
	QueryInsert := `INSERT INTO hd_jadwalujian (jdw_namajadwal,jdw_banksoal,jdw_jenissoal,jdw_jumlahsoal,jdw_tgl_mulai,jdw_jam_mulai,jdw_toleransiketerlambatan,jdw_durasi,jdw_token,jdw_ketentuan_waktu,jdw_userupdate,jdw_lastupdate,jdw_is_delete) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_, err := database.DbMain.Exec(QueryInsert, jadwalData.NamaJadwal, jadwalData.BankSoal, jadwalData.JenisSoal, jadwalData.JumlahSoal, jadwalData.TanggalMulai, jadwalData.JamMulai, jadwalData.ToleransiTerlambat, jadwalData.DurasiUjian, jadwalData.TokenUjian, jadwalData.KetentuanWaktu, jadwalData.UserUpdate, jadwalData.LastUpdate,"n")
	if err != nil {
		middleware.LogError(err, "Insert Data Failed")
		return false, err
	}
	return true, nil
}

func UpdateJadwalUjian(jadwalData *model.JadwalUjian) (bool,error) {
	QueryUpdate := `UPDATE hd_jadwalujian SET `
	UpdateString := ""
	var args []interface{}

	if jadwalData.NamaJadwal != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `jdw_namajadwal = ? `
		args = append(args,jadwalData.NamaJadwal)
	}

	tempString := strconv.Itoa(jadwalData.BankSoal)
	if tempString != "" && tempString != "0" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `jdw_banksoal = ? `
		args = append(args,jadwalData.BankSoal)
	}

	if jadwalData.JenisSoal != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `jdw_jenissoal = ? `
		args = append(args,jadwalData.JenisSoal)
	}

	tempString = strconv.Itoa(jadwalData.JumlahSoal)
	if tempString != "" && tempString != "0" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `jdw_jumlahsoal = ? `
		args = append(args,jadwalData.JumlahSoal)
	}

	if jadwalData.TanggalMulai != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `jdw_tgl_mulai = ? `
		args = append(args,jadwalData.TanggalMulai)
	}

	if jadwalData.JamMulai != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `jdw_jam_mulai = ? `
		args = append(args,jadwalData.JamMulai)
	}

	tempString = strconv.Itoa(jadwalData.ToleransiTerlambat)
	if tempString != "" && tempString != "0" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `jdw_toleransiketerlambatan = ? `
		args = append(args,jadwalData.ToleransiTerlambat)
	}

	tempString = strconv.Itoa(jadwalData.DurasiUjian)
	if tempString != "" && tempString != "0" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `jdw_durasi = ? `
		args = append(args,jadwalData.DurasiUjian)
	}

	if jadwalData.TokenUjian != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `jdw_token = ? `
		args = append(args,jadwalData.TokenUjian)
	}

	if jadwalData.KetentuanWaktu != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `jdw_ketentuan_waktu = ? `
		args = append(args,jadwalData.KetentuanWaktu)
	}

	if jadwalData.UserUpdate != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `jdw_userupdate = ? `
		args = append(args,jadwalData.UserUpdate)
	}

	if jadwalData.LastUpdate != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `jdw_lastupdate = ? `
		args = append(args,jadwalData.LastUpdate)
	}

	if len(UpdateString) > 0 {
		QueryUpdate = fmt.Sprintf("%s %s", QueryUpdate, UpdateString)
	}
	QueryUpdate += ` WHERE jdw_id = ?`
	args = append(args, jadwalData.JadwalId)

	_, err := database.DbMain.Exec(QueryUpdate, args...)
	if err != nil {
		middleware.LogError(err, "Update Data Failed")
		return false, err
	}
	return true, nil
}

func DeleteJadwalUjian(jdwId string) (bool, error) {
	QueryUpdate := `UPDATE hd_jadwalujian SET jdw_is_delete = "y" WHERE jdw_id = ?`
	_, err := database.DbMain.Exec(QueryUpdate, jdwId)
	if err != nil {
		middleware.LogError(err, "Update Data Failed")
		return false, err
	}
	return true, nil
}
// End Function Untuk Manajemen Jadwal Ujian

// Function Untuk Peserta Ujian
func GetJadwalUjianPesertaByDate(siswaId string, tanggal string) ([]model.JadwalUjianSiswa, error) {
	QueryData := `SELECT jdw_id,jdw_namajadwal,jdw_tgl_mulai,jdw_jam_mulai,jdw_durasi,jdw_toleransiketerlambatan,jdw_ketentuan_waktu,ps_id,ps_jam_mulai,ps_jam_selesai FROM hd_jadwalujian LEFT JOIN hd_pesertaujian ON hd_pesertaujian.ps_jadwal_ujian = hd_jadwalujian.jdw_id AND hd_pesertaujian.ps_siswa = ? WHERE jdw_tgl_mulai = ? AND jdw_is_delete = "n"`
	var returnData []model.JadwalUjianSiswa

	rows, err := database.DbMain.Query(QueryData, siswaId, tanggal)
	if err != nil {
		middleware.LogError(err,"Query Error")
		return nil,err
	}
	defer rows.Close()
	for rows.Next() {
		var rowData model.JadwalUjianSiswa
		var psId model.NullString
		var psJamMulai model.NullString
		var psJamSelesai model.NullString
		err = rows.Scan(&rowData.JadwalId, &rowData.NamaJadwal, &rowData.TanggalMulai, &rowData.JamMulai, &rowData.DurasiUjian, &rowData.ToleransiTerlambat, &rowData.KetentuanWaktu, &psId, &psJamMulai, &psJamSelesai)
		if err != nil {
			middleware.LogError(err,"Data Scan Error")
			return nil,err
		}
		rowData.PesertaId = psId.Str
		rowData.StartPeserta = psJamMulai.Str
		rowData.EndPeserta = psJamSelesai.Str
		returnData = append(returnData, rowData)
	}
	return returnData,nil
}