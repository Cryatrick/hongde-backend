package service

import(
	"fmt"
	"database/sql"
	"strings"
	"strconv"
	// "time"

	"hongde_backend/internal/model"
	"hongde_backend/internal/database"
	"hongde_backend/internal/middleware"
	// "hongde_backend/internal/config"
)

// Function Untuk Peserta Ujian
func CheckSudahMengikuti(jawalId int, siswaId string) (*model.PesertaUjian, error) {
	row := database.DbMain.QueryRow(`SELECT ps_id,ps_siswa,sw_nama,ps_jadwal_ujian,jdw_namajadwal,ps_jam_mulai,ps_islate,ps_jam_selesai,ps_final_score FROM hd_pesertaujian JOIN hd_siswa ON hd_siswa.sw_id = hd_pesertaujian.ps_siswa JOIN hd_jadwalujian ON hd_jadwalujian.jdw_id = hd_pesertaujian.ps_jadwal_ujian WHERE hd_pesertaujian.ps_siswa = ? AND ps_jadwal_ujian = ?`,siswaId,jawalId)
	var returnData = &model.PesertaUjian{}
	var jamSelsesai model.NullString
	var finalScore model.NullFloat
	err := row.Scan(&returnData.PesertaId, &returnData.SiswaId, &returnData.SiswaNama, &returnData.UjianId, &returnData.UjianNama, &returnData.StartedAt, &returnData.IsLate, &jamSelsesai, &finalScore)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		middleware.LogError(err, "Data Scan Error")
		return nil, err
	}
	if jamSelsesai.Valid == true {
		returnData.EndAt = jamSelsesai.Str
	}
	if finalScore.Valid == true {
		returnData.FinalScore = finalScore.Float
	}
	returnData.SoalArray, err = GetJawabanPesertaUjian(returnData.PesertaId)
	if err != nil {
		middleware.LogError(err, "Data Scan Error")
		return nil, err
	}
	return returnData, nil
}

func GetDetailPesertaId(pesertaId string)(*model.PesertaUjian, error) {
	row := database.DbMain.QueryRow(`SELECT ps_id,ps_siswa,sw_nama,ps_jadwal_ujian,jdw_namajadwal,ps_jam_mulai,ps_islate,ps_jam_selesai,ps_final_score FROM hd_pesertaujian JOIN hd_siswa ON hd_siswa.sw_id = hd_pesertaujian.ps_siswa JOIN hd_jadwalujian ON hd_jadwalujian.jdw_id = hd_pesertaujian.ps_jadwal_ujian WHERE hd_pesertaujian.ps_id = ? `,pesertaId)
	var returnData = &model.PesertaUjian{}
	var jamSelsesai model.NullString
	var finalScore model.NullFloat
	err := row.Scan(&returnData.PesertaId, &returnData.SiswaId, &returnData.SiswaNama, &returnData.UjianId, &returnData.UjianNama, &returnData.StartedAt, &returnData.IsLate, &jamSelsesai, &finalScore)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		middleware.LogError(err, "Data Scan Error")
		return nil, err
	}
	if jamSelsesai.Valid == true {
		returnData.EndAt = jamSelsesai.Str
	}
	if finalScore.Valid == true {
		returnData.FinalScore = finalScore.Float
	}
	returnData.SoalArray, err = GetJawabanPesertaUjian(pesertaId)
	if err != nil {
		middleware.LogError(err, "Data Scan Error")
		return nil, err
	}
	return returnData, nil
}

func GetJawabanPesertaUjian(pesertaId string)([]model.SoalPesertaUjian, error) {
	QueryData := `SELECT jwbpst_peserta, jwbpst_soal,jwbpst_urutan,jwbpst_jawaban,COALESCE(jwbpst_isbenar,"")AS is_benar, COALESCE(jwbpst_bobot,0) AS bobot_jawaban FROM hd_jawabanpeserta WHERE jwbpst_peserta = ? ORDER BY jwbpst_urutan ASC`
	var returnData []model.SoalPesertaUjian

	rows, err := database.DbMain.Query(QueryData,pesertaId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		middleware.LogError(err,"Query Error")
		return nil,err
	}
	defer rows.Close()
	for rows.Next() {
		var rowData model.SoalPesertaUjian
		var jawabanValue model.NullString
		err = rows.Scan(&rowData.PesertaId, &rowData.SoalId, &rowData.UrutanSoal, &jawabanValue, &rowData.IsRight, &rowData.BobotSoal)
		if err != nil {
			middleware.LogError(err,"Data Scan Error")
			return nil,err
		}
		if jawabanValue.Valid == true {
			rowData.JawabanSiswa = jawabanValue.Str
		}
		returnData = append(returnData, rowData)
	}
	return returnData,nil
}

func GetJawabanPesertaWithSoalData(pesertaId string)([]model.SoalPesertaUjian, error) {
	QueryData := `
	SELECT
	jwbpst_peserta,
	jwbpst_soal,
	jwbpst_urutan,
	jwbpst_jawaban,
	(CASE WHEN hd_soalujian.soal_tipe = "essay" THEN "w" WHEN hd_jawabanpeserta.jwbpst_jawaban = hd_soalujian.soal_jwbbenar THEN "y" ELSE "n" END ) as is_benar,
	hd_soalujian.soal_bobot
	FROM
	hd_jawabanpeserta
	JOIN hd_soalujian ON hd_soalujian.soal_id = hd_jawabanpeserta.jwbpst_soal
	WHERE
	jwbpst_peserta = ?
	ORDER BY
	jwbpst_urutan ASC
	`
	var returnData []model.SoalPesertaUjian

	rows, err := database.DbMain.Query(QueryData,pesertaId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		middleware.LogError(err,"Query Error")
		return nil,err
	}
	defer rows.Close()
	for rows.Next() {
		var rowData model.SoalPesertaUjian
		var jawabanValue model.NullString
		err = rows.Scan(&rowData.PesertaId, &rowData.SoalId, &rowData.UrutanSoal, &jawabanValue, &rowData.IsRight, &rowData.BobotSoal)
		if err != nil {
			middleware.LogError(err,"Data Scan Error")
			return nil,err
		}
		if jawabanValue.Valid == true {
			rowData.JawabanSiswa = jawabanValue.Str
		}
		returnData = append(returnData, rowData)
	}
	return returnData,nil
}

func ProcessInsertPesertaUjian(pesertaData *model.PesertaUjian)(bool, error) {
	QueryInsertPeserta := `INSERT INTO hd_pesertaujian (ps_id,ps_siswa,ps_jadwal_ujian,ps_jam_mulai,ps_islate,ps_userupdate,ps_lastupdate) VALUES (?,?,?,?,?,?,?)`

	QueryInsertSoalPeserta := `INSERT INTO hd_jawabanpeserta (jwbpst_peserta,jwbpst_soal,jwbpst_urutan) VALUES `
	InsertSoalPesertaPlaceholder := ``
	InsertSoalPesertaValues := []interface{}{}

	for i, record := range pesertaData.SoalArray {
		if i > 0 {
			InsertSoalPesertaPlaceholder += ", "
		}
		InsertSoalPesertaPlaceholder += "(?,?,?)"
		InsertSoalPesertaValues = append(InsertSoalPesertaValues,record.PesertaId,record.SoalId,record.UrutanSoal)
	}

	QueryInsertSoalPeserta += InsertSoalPesertaPlaceholder

	tx, err := database.DbMain.Begin()
	if err != nil {
		middleware.LogError(err, "Failed to start transaction")
		return false, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	_, err = tx.Exec(QueryInsertPeserta, pesertaData.PesertaId,pesertaData.SiswaId,pesertaData.UjianId,pesertaData.StartedAt,pesertaData.IsLate,pesertaData.UserUpdate,pesertaData.LastUpdate)
	if err != nil {
		middleware.LogError(err, "Insert Peserta Ujian Query Error")
		return false, err 
	}

	_, err = tx.Exec(QueryInsertSoalPeserta,InsertSoalPesertaValues...)
	if err != nil {
		middleware.LogError(err, "Insert Soal Peserta Ujian Query Error")
		return false, err 
	}

	err = tx.Commit()
	if err != nil {
		middleware.LogError(err, "Commit Error")
		return false, err
	}

	return true, nil
}

func SaveJawabanPesertaUjian(pesertaId string,jawabanArray map[string]string)(bool, error) {
	var (
		caseParts []string
		args      []interface{}
		ids       []string
	)

	for soalId, jawabanValue := range jawabanArray {
		caseParts = append(caseParts, "WHEN ? THEN ?")
		args = append(args, soalId, jawabanValue)
		ids = append(ids, soalId)
	}

	query := fmt.Sprintf(`UPDATE hd_jawabanpeserta SET jwbpst_jawaban = CASE jwbpst_soal %s END WHERE jwbpst_peserta = ? AND jwbpst_soal IN (%s)`, strings.Join(caseParts, " "), strings.Join(ids, ","))

	args = append(args, pesertaId)

	_, err := database.DbMain.Exec(query, args...)
	if err != nil {
		middleware.LogError(err, "Failed To Save Jawaban Peserta")
		return false, err
	}
	return true, nil
}

func SaveSingleJawabanPeserta(pesertaId string, soalId string, savedJawaban string)(bool, error) {
	query := `UPDATE hd_jawabanpeserta SET jwbpst_jawaban = ? WHERE jwbpst_peserta = ? AND jwbpst_soal = ?`

	_, err := database.DbMain.Exec(query, savedJawaban,pesertaId,soalId)
	if err != nil {
		middleware.LogError(err, "Failed To Save Jawaban Peserta")
		return false, err
	}
	return true, nil
}

func FinalizePesertaUjian(pesertaId string, finishTime string,jawabanArray []model.SoalPesertaUjian)(bool, error) {
	var (
		caseParts []string
		args      []interface{}
		ids       []string
	)

	for _, jawabanValue := range jawabanArray {
		caseParts = append(caseParts, "WHEN ? THEN ?")
		args = append(args, jawabanValue.SoalId, jawabanValue.IsRight)
		strId := strconv.Itoa(jawabanValue.SoalId)
		ids = append(ids, strId)
	}
	for _, jawabanValue := range jawabanArray {
		args = append(args, jawabanValue.SoalId, jawabanValue.BobotSoal)
	}
	queryJawaban := fmt.Sprintf(`UPDATE hd_jawabanpeserta SET jwbpst_isbenar = CASE jwbpst_soal %s END,jwbpst_bobot = CASE jwbpst_soal %s END WHERE jwbpst_peserta = ? AND jwbpst_soal IN (%s)`, strings.Join(caseParts, " "), strings.Join(caseParts, " "), strings.Join(ids, ","))

	args = append(args, pesertaId)

	tx, err := database.DbMain.Begin()
	if err != nil {
		middleware.LogError(err, "Failed to start transaction")
		return false, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	_, err = tx.Exec(queryJawaban,args...)
	if err != nil {
		middleware.LogError(err, "Save Jawaban Peserta Query Error")
		return false, err 
	}

	if finishTime != "" {
		_, err = tx.Exec(`UPDATE hd_pesertaujian SET ps_jam_selesai = ? WHERE ps_id = ?`,finishTime, pesertaId)
		if err != nil {
			middleware.LogError(err, "Update Peserta Ujian Query Error")
			return false, err 
		}
	}

	err = tx.Commit()
	if err != nil {
		middleware.LogError(err, "Commit Error")
		return false, err
	}

	return true, nil
}

// End Function Untuk Peserta Ujian

// Function Untuk Hasil Ujian
func GetAllPesertaUjianJadwal(jadwalId string)([]model.PesertaUjian, error) {
	query := `SELECT
	hd_pesertaujian.ps_id,
	hd_pesertaujian.ps_siswa,
	hd_siswa.sw_nama,
	hd_pesertaujian.ps_jadwal_ujian,
	hd_jadwalujian.jdw_namajadwal,
	hd_pesertaujian.ps_jam_mulai,
	COALESCE(hd_pesertaujian.ps_jam_selesai,"") AS jam_selesai,
	hd_pesertaujian.ps_islate,
	SUM(CASE WHEN hd_jawabanpeserta.jwbpst_isbenar = "y" THEN 1 ELSE 0 END) AS total_right,
	SUM(CASE WHEN hd_jawabanpeserta.jwbpst_isbenar = "w" THEN 1 ELSE 0 END) AS total_waiting,
	SUM(CASE WHEN hd_jawabanpeserta.jwbpst_isbenar = "n" THEN 1 ELSE 0 END) AS total_wrong,
	COALESCE(ps_final_score,SUM(CASE WHEN hd_jawabanpeserta.jwbpst_isbenar = "y" THEN jwbpst_bobot ELSE 0 END),0) AS current_score,
	(CASE WHEN ps_final_score IS NOT NULL THEN "y" ELSE "n" END) AS is_recap
	FROM
	hd_pesertaujian
	JOIN hd_siswa ON hd_siswa.sw_id = hd_pesertaujian.ps_siswa
	JOIN hd_jadwalujian ON hd_jadwalujian.jdw_id = hd_pesertaujian.ps_jadwal_ujian
	JOIN hd_jawabanpeserta ON hd_jawabanpeserta.jwbpst_peserta = hd_pesertaujian.ps_id
	WHERE
	hd_pesertaujian.ps_jadwal_ujian = ?
	GROUP BY
	hd_pesertaujian.ps_id`
	var returnData []model.PesertaUjian
	rows, err := database.DbMain.Query(query,jadwalId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		middleware.LogError(err,"Query Error")
		return nil,err
	}
	defer rows.Close()
	for rows.Next() {
		var rowData model.PesertaUjian
		err = rows.Scan(&rowData.PesertaId, &rowData.SiswaId, &rowData.SiswaNama, &rowData.UjianId, &rowData.UjianNama, &rowData.StartedAt, &rowData.EndAt, &rowData.IsLate, &rowData.TotalRight, &rowData.TotalWaiting, &rowData.TotalWrong, &rowData.FinalScore, &rowData.IsRecapped)
		if err != nil {
			middleware.LogError(err,"Data Scan Error")
			return nil,err
		}
		returnData = append(returnData, rowData)
	}
	return returnData,nil
}

func UpdateSinglePenilaianJawaban(jawabanData *model.SoalPesertaUjian)(bool, error) {
	query := `UPDATE hd_jawabanpeserta SET`
	updateString := ""
	var args []interface{}

	if jawabanData.IsRight != "" {
		if updateString != "" {
			updateString += `, `
		}
		updateString += `jwbpst_isbenar = ? `
		args = append(args,jawabanData.IsRight)
	}
	if jawabanData.BobotSoal > 0 {
		if updateString != "" {
			updateString += `, `
		}
		updateString += `jwbpst_bobot = ? `
		args = append(args,jawabanData.BobotSoal)
	}
	if len(updateString) > 0 {
		query = fmt.Sprintf("%s %s", query, updateString)
	}
	query += ` WHERE jwbpst_peserta = ? AND jwbpst_soal = ?`
	args = append(args, jawabanData.PesertaId, jawabanData.SoalId)

	_, err := database.DbMain.Exec(query, args...)
	if err != nil {
		middleware.LogError(err, "Update Data Failed")
		return false, err
	}
	return true, nil
}

func FinalizeHasilUjianJadwal(pesertaData []model.PesertaUjian)(bool, error) {
	queryUpdate := "UPDATE hd_pesertaujian SET ps_final_score = ?,ps_userupdate = ?, ps_lastupdate = ? WHERE ps_id = ?"

	tx, err := database.DbMain.Begin()
	if err != nil {
		middleware.LogError(err, "Failed to start transaction")
		return false, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	
	stmt, err := tx.Prepare(queryUpdate)
	if err != nil {
		middleware.LogError(err, "Falied To Prepare Query")
		return false,err
	}
	defer stmt.Close()
	for _, d := range pesertaData {
		_, err = stmt.Exec(d.FinalScore, d.UserUpdate, d.LastUpdate, d.PesertaId)
		if err != nil {
			middleware.LogError(err, "Update Data Failed")
			return false,err
		}
	}
	err = tx.Commit()
	if err != nil {
		middleware.LogError(err, "Commit Error")
		return false, err
	}
	return true, nil
}
// End Function Untuk Hasil Ujian