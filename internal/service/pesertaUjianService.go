package service

import(
	"fmt"
	"database/sql"
	"strings"
	// "strconv"
	// "time"

	"hongde_backend/internal/model"
	"hongde_backend/internal/database"
	"hongde_backend/internal/middleware"
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

func GetJawabanPesertaUjian(pesertaId string)([]model.SoalPesertaUjian, error) {
	QueryData := `SELECT jwbpst_peserta, jwbpst_soal,jwbpst_urutan,jwbpst_jawaban FROM hd_jawabanpeserta WHERE jwbpst_peserta = ? ORDER BY jwbpst_urutan ASC`
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
		err = rows.Scan(&rowData.PesertaId, &rowData.SoalId, &rowData.UrutanSoal, &jawabanValue)
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