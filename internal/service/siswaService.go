package service

import(
	"fmt"
	"database/sql"
	"strconv"
	"time"

	"hongde_backend/internal/model"
	"hongde_backend/internal/database"
	"hongde_backend/internal/middleware"
)


func LoginSiswa(userName, userPassword string) (*model.UserLogin, error) {
	row := database.DbMain.QueryRow(`SELECT sw_id, sw_nama,"9" AS role_id  FROM hd_siswa WHERE sw_id = ? AND sw_password = ? AND sw_isdelete = "n"`,userName, userPassword)
	var user = &model.UserLogin{}
	err := row.Scan(&user.UserId, &user.UserNama, &user.UserRole)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		// Log and return actual error
		middleware.LogError(err, "Data Scan Error")
		return nil, err
	}
	return user, nil
}

func GetSiswaId(idSiswa string)(*model.SiswaList, error) {
	row := database.DbMain.QueryRow(`SELECT sw_id, sw_nama, sw_email, CASE WHEN sw_namamandarin IS NULL THEN "" ELSE sw_namamandarin END AS sw_namamandarin, sw_jenisidentitas, sw_noidentitas, sw_tempatlahir, sw_tanggallahir, sw_tempattinggal, sw_nokontak, sw_jenissiswa, sw_userupdate, sw_lastupdate FROM hd_siswa WHERE sw_id = ? AND sw_isdelete = "n"`,idSiswa)
	var siswa = &model.SiswaList{}
	err := row.Scan(&siswa.SiswaId,&siswa.NamaSiswa, &siswa.EmailSiswa, &siswa.NamaMandarin, &siswa.JenisIdentitas, &siswa.NoIdentitas, &siswa.TempatLahir, &siswa.TanggalLahir, &siswa.TempatTinggal, &siswa.NoKontak, &siswa.JenisSiswa, &siswa.UserUpdate, &siswa.LastUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		// Log and return actual error
		middleware.LogError(err, "Data Scan Error")
		return nil, err
	}
	return siswa, nil
}

func GetAllSiswa()([]model.SiswaList,error) {
	QueryData := `SELECT sw_id, sw_nama, sw_email, CASE WHEN sw_namamandarin IS NULL THEN "" ELSE sw_namamandarin END AS sw_namamandarin, sw_jenisidentitas, sw_noidentitas, sw_tempatlahir, sw_tanggallahir, sw_tempattinggal, sw_nokontak, sw_jenissiswa, sw_userupdate, sw_lastupdate FROM hd_siswa WHERE sw_isdelete = "n"`
	var returnData []model.SiswaList

	rows, err := database.DbMain.Query(QueryData)
	if err != nil {
		middleware.LogError(err,"Query Error")
		return nil,err
	}
	defer rows.Close()
	for rows.Next() {
		var siswaRow model.SiswaList
		err = rows.Scan(&siswaRow.SiswaId,&siswaRow.NamaSiswa, &siswaRow.EmailSiswa, &siswaRow.NamaMandarin, &siswaRow.JenisIdentitas, &siswaRow.NoIdentitas, &siswaRow.TempatLahir, &siswaRow.TanggalLahir, &siswaRow.TempatTinggal, &siswaRow.NoKontak, &siswaRow.JenisSiswa, &siswaRow.UserUpdate, &siswaRow.LastUpdate)
		if err != nil {
			middleware.LogError(err,"Data Scan Error")
			return nil,err
		}
		returnData = append(returnData, siswaRow)
	}
	return returnData,nil
}

func GenerateSiswaId() (string, error) {
	codeId := time.Now().Format("0601")

	var lastNumber model.NullString
	QueryData := `SELECT MAX((REPLACE(hd_siswa.sw_id,?,''))) as total FROM hd_siswa WHERE sw_id LIKE ?`
	err := database.DbMain.QueryRow(QueryData,codeId,codeId+"%").Scan(&lastNumber)
	if err != nil {
		middleware.LogError(err, "Data Scan Error")
		return "",err
	}

	finalNumber := 1;

	if lastNumber.Valid == true {
		finalNumber,_ = strconv.Atoi(lastNumber.Str)
		finalNumber = finalNumber + 1
	}

	siswaId := fmt.Sprintf("%v%03d",codeId,finalNumber)

	return siswaId, nil
}

func InsertSiswa(siswaData *model.SiswaList)(bool, error) {
	QueryInsert := `INSERT INTO hd_siswa VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_, err := database.DbMain.Exec(QueryInsert, siswaData.SiswaId,siswaData.NamaSiswa,siswaData.EmailSiswa,siswaData.NamaMandarin,siswaData.JenisIdentitas,siswaData.NoIdentitas,siswaData.TempatLahir,siswaData.TanggalLahir,siswaData.TempatTinggal,siswaData.NoKontak,siswaData.PasswordSiswa,"n",siswaData.JenisSiswa,siswaData.UserUpdate,siswaData.LastUpdate)
	if err != nil {
		middleware.LogError(err, "Insert Data Failed")
		return false, err
	}
	return true, nil
}

func UpdateSiswa(siswaData *model.SiswaList) (bool,error) {
	QueryUpdate := `UPDATE hd_siswa SET `
	UpdateString := ""
	var args []interface{}

	if siswaData.NamaSiswa != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `sw_nama = ? `
		args = append(args,siswaData.NamaSiswa)
	}
	if siswaData.EmailSiswa != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `sw_email = ? `
		args = append(args,siswaData.EmailSiswa)
	}
	if siswaData.NamaMandarin != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `sw_namamandarin = ? `
		args = append(args,siswaData.NamaMandarin)
	}
	if siswaData.JenisIdentitas != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `sw_jenisidentitas = ? `
		args = append(args,siswaData.JenisIdentitas)
	}
	if siswaData.NamaSiswa != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `sw_noidentitas = ? `
		args = append(args,siswaData.NoIdentitas)
	}
	if siswaData.TempatLahir != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `sw_tempatlahir = ? `
		args = append(args,siswaData.TempatLahir)
	}
	if siswaData.TanggalLahir != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `sw_tanggallahir = ? `
		args = append(args,siswaData.TanggalLahir)
	}
	if siswaData.NoKontak != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `sw_nokontak = ? `
		args = append(args,siswaData.NoKontak)
	}
	if siswaData.JenisSiswa != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `sw_jenissiswa = ? `
		args = append(args,siswaData.JenisSiswa)
	}
	if siswaData.PasswordSiswa != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `sw_password = ? `
		args = append(args,siswaData.PasswordSiswa)
	}

	if siswaData.UserUpdate != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `sw_userupdate = ? `
		args = append(args,siswaData.UserUpdate)
	}

	if siswaData.LastUpdate != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `sw_lastupdate = ? `
		args = append(args,siswaData.LastUpdate)
	}

	if len(UpdateString) > 0 {
		QueryUpdate = fmt.Sprintf("%s %s", QueryUpdate, UpdateString)
	}
	QueryUpdate += ` WHERE sw_id = ?`
	args = append(args, siswaData.SiswaId)

	_, err := database.DbMain.Exec(QueryUpdate, args...)
	if err != nil {
		middleware.LogError(err, "Update Data Failed")
		return false, err
	}
	return true, nil
}

func DeleteSiswa(siswaId string) (bool, error) {
	QueryUpdate := `UPDATE hd_siswa SET sw_isdelete = "y" WHERE sw_id = ?`
	_, err := database.DbMain.Exec(QueryUpdate, siswaId)
	if err != nil {
		middleware.LogError(err, "Update Data Failed")
		return false, err
	}
	return true, nil
}

func InsertBulkSiswa(insertData []model.SiswaList)(bool, error) {
	InsertQuery := "INSERT INTO hd_siswa VALUES "
	InsertPlaceholder := ""
	InsertValues := []interface{}{}

	for i,record := range insertData {
		if i > 0 {
			InsertPlaceholder += ", "
		}
		InsertPlaceholder += "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
		InsertValues = append(InsertValues,record.SiswaId,record.NamaSiswa,record.EmailSiswa,record.NamaMandarin,record.JenisIdentitas,record.NoIdentitas,record.TempatLahir,record.TanggalLahir,record.TempatTinggal,record.NoKontak,record.PasswordSiswa,"n",record.JenisSiswa,record.UserUpdate,record.LastUpdate)
	}
	InsertQuery += InsertPlaceholder
	_, err := database.DbMain.Exec(InsertQuery, InsertValues...)
	if err != nil {
		middleware.LogError(err, "Insert Batch Failed")
		return false, err
	}
	return true, nil
}