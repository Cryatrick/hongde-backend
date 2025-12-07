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

// Function Untuk Bank Soal
func GetBankSoal() ([]model.BankSoal, error) {
	QueryData := `SELECT bs_id, bs_namabank, COUNT(soal_id) AS total_soal, bs_userupdate, bs_lastupdate FROM hd_banksoal LEFT JOIN hd_soalujian ON hd_soalujian.soal_bank = hd_banksoal.bs_id WHERE bs_isdelete = "n" `
	var returnData []model.BankSoal

	rows, err := database.DbMain.Query(QueryData)
	if err != nil {
		middleware.LogError(err,"Query Error")
		return nil,err
	}
	defer rows.Close()
	for rows.Next() {
		var rowData model.BankSoal
		var checkId model.NullInt
		var checkNama model.NullString
		var checkUserupdate model.NullString
		var checkLastupdate model.NullString
		err = rows.Scan(&checkId,&checkNama, &rowData.JumlahSoal, &checkUserupdate, &checkLastupdate)
		if err != nil {
			middleware.LogError(err,"Data Scan Error")
			return nil,err
		}
		if checkId.Valid == true {
			rowData.BankId = checkId.Int
			rowData.NamaBank = checkNama.Str
			rowData.UserUpdate = checkUserupdate.Str
			rowData.LastUpdate = checkLastupdate.Str
			returnData = append(returnData, rowData)
		}else {
			break
		}
	}
	return returnData,nil
}

func GetBankId(idBank string)(*model.BankSoal, error) {
	row := database.DbMain.QueryRow(`SELECT bs_id, bs_namabank, COUNT(soal_id) AS total_soal, bs_userupdate, bs_lastupdate FROM hd_banksoal LEFT JOIN hd_soalujian ON hd_soalujian.soal_bank = hd_banksoal.bs_id WHERE bs_id = ? AND bs_isdelete = "n"`,idBank)
	var returnData = &model.BankSoal{}
	err := row.Scan(&returnData.BankId,&returnData.NamaBank, &returnData.JumlahSoal, &returnData.UserUpdate, &returnData.LastUpdate)
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

func InsertBankSoal(bankData *model.BankSoal)(bool, error) {
	QueryInsert := `INSERT INTO hd_banksoal (bs_namabank, bs_userupdate, bs_lastupdate, bs_isdelete) VALUES(?,?,?,?)`
	_, err := database.DbMain.Exec(QueryInsert, bankData.NamaBank,bankData.UserUpdate,bankData.LastUpdate,"n")
	if err != nil {
		middleware.LogError(err, "Insert Data Failed")
		return false, err
	}
	return true, nil
}

func UpdateBankSoal(bankData *model.BankSoal) (bool,error) {
	QueryUpdate := `UPDATE hd_banksoal SET `
	UpdateString := ""
	var args []interface{}

	if bankData.NamaBank != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `bs_namabank = ? `
		args = append(args,bankData.NamaBank)
	}
	if bankData.UserUpdate != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `bs_userupdate = ? `
		args = append(args,bankData.UserUpdate)
	}

	if bankData.LastUpdate != "" {
		if UpdateString != "" {
			UpdateString += `, `
		}
		UpdateString += `bs_lastupdate = ? `
		args = append(args,bankData.LastUpdate)
	}

	if len(UpdateString) > 0 {
		QueryUpdate = fmt.Sprintf("%s %s", QueryUpdate, UpdateString)
	}
	QueryUpdate += ` WHERE bs_id = ?`
	args = append(args, bankData.BankId)

	_, err := database.DbMain.Exec(QueryUpdate, args...)
	if err != nil {
		middleware.LogError(err, "Update Data Failed")
		return false, err
	}
	return true, nil
}

func DeleteBankSoal(bankId string) (bool, error) {
	QueryUpdate := `UPDATE hd_banksoal SET bs_isdelete = "y" WHERE bs_id = ?`
	_, err := database.DbMain.Exec(QueryUpdate, bankId)
	if err != nil {
		middleware.LogError(err, "Update Data Failed")
		return false, err
	}
	return true, nil
}
// End Function Untuk Bank Soal

// Function Untuk Soal
func GenerateSoalId() (string, error) {
	var lastNumber model.NullInt
	QueryData := `SELECT MAX(soal_id) as total FROM hd_soalujian `
	err := database.DbMain.QueryRow(QueryData).Scan(&lastNumber)
	if err != nil {
		middleware.LogError(err, "Data Scan Error")
		return "",err
	}

	finalNumber := 1;

	if lastNumber.Valid == true {
		finalNumber = lastNumber.Int
		finalNumber = finalNumber + 1
	}

	returnNumber := strconv.Itoa(finalNumber)

	return returnNumber, nil
}
func GetSoalBank(idBank string)([]model.SoalList, error) {
	Query := `SELECT soal_id,soal_bank,soal_gambar_pertanyaan,soal_pertanyaan,soal_urutan,soal_gambarjwba,soal_jwba,soal_gambarjwbb,soal_jwbb,soal_gambarjwbc,soal_jwbc,soal_gambarjwbd,soal_jwbd,soal_jwbbenar,soal_bobot,soal_tipe,soal_userupdate,soal_lastupdate FROM hd_soalujian WHERE soal_is_delete = "n"`
	var args []interface{}

	if idBank != `` {
		Query = Query + " AND soal_bank = ?"
		args  = append(args,idBank)
	}
	Query = Query + ` ORDER BY soal_urutan ASC`
	rows, err := database.DbMain.Query(Query,args...)
	if err != nil {
		middleware.LogError(err,"Query Error")
		return nil, err
	}
	defer rows.Close()
	var returnDataList []model.SoalList
	for rows.Next() {
		var rowData model.SoalList
		err = rows.Scan(&rowData.SoalId,&rowData.BankId, &rowData.GambarSoal, &rowData.PertanyaanSoal, &rowData.UrutanSoal, &rowData.GambarA, &rowData.JawabanA, &rowData.GambarB, &rowData.JawabanB, &rowData.GambarC, &rowData.JawabanC, &rowData.GambarD, &rowData.JawabanD, &rowData.JawabanBenar, &rowData.BobotSoal,&rowData.TipeSoal,&rowData.UserUpdate,&rowData.LastUpdate)
		if err != nil {
			middleware.LogError(err,"Data Scan Error")
			return nil,err
		}
		if rowData.GambarSoal != "" {
			rowData.GambarSoal = BuildImageURL(rowData.GambarSoal)
		}
		if rowData.GambarA != "" {
			rowData.GambarA = BuildImageURL(rowData.GambarA)
		}
		if rowData.GambarB != "" {
			rowData.GambarB = BuildImageURL(rowData.GambarB)
		}
		if rowData.GambarC != "" {
			rowData.GambarC = BuildImageURL(rowData.GambarC)
		}
		if rowData.GambarD != "" {
			rowData.GambarD = BuildImageURL(rowData.GambarD)
		}
		returnDataList = append(returnDataList,rowData)
	}
	return returnDataList,nil
}
func GetSingleSoal(idSoal string)(*model.SoalList, error) {
	Query := `SELECT soal_id,soal_bank,soal_gambar_pertanyaan,soal_pertanyaan,soal_urutan,soal_gambarjwba,soal_jwba,soal_gambarjwbb,soal_jwbb,soal_gambarjwbc,soal_jwbc,soal_gambarjwbd,soal_jwbd,soal_jwbbenar,soal_bobot,soal_tipe,soal_userupdate,soal_lastupdate FROM hd_soalujian WHERE soal_is_delete = "n" AND soal_id = ?`

	row := database.DbMain.QueryRow(Query,idSoal)
	var rowData = &model.SoalList{}
	err := row.Scan(&rowData.SoalId,&rowData.BankId, &rowData.GambarSoal, &rowData.PertanyaanSoal, &rowData.UrutanSoal, &rowData.GambarA, &rowData.JawabanA, &rowData.GambarB, &rowData.JawabanB, &rowData.GambarC, &rowData.JawabanC, &rowData.GambarD, &rowData.JawabanD, &rowData.JawabanBenar, &rowData.BobotSoal,&rowData.TipeSoal,&rowData.UserUpdate,&rowData.LastUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		// Log and return actual error
		middleware.LogError(err, "Data Scan Error")
		return nil, err
	}
	if rowData.GambarSoal != "" {
		rowData.GambarSoal = BuildImageURL(rowData.GambarSoal)
	}
	if rowData.GambarA != "" {
		rowData.GambarA = BuildImageURL(rowData.GambarA)
	}
	if rowData.GambarB != "" {
		rowData.GambarB = BuildImageURL(rowData.GambarB)
	}
	if rowData.GambarC != "" {
		rowData.GambarC = BuildImageURL(rowData.GambarC)
	}
	if rowData.GambarD != "" {
		rowData.GambarD = BuildImageURL(rowData.GambarD)
	}
	return rowData, nil
}
func InsertSoalUjian(soalData *model.SoalList)(bool, error) {
	QueryInsert := `INSERT INTO hd_soalujian (soal_id,soal_bank, soal_gambar_pertanyaan, soal_pertanyaan, soal_urutan, soal_gambarjwba, soal_jwba,soal_gambarjwbb,soal_jwbb,soal_gambarjwbc,soal_jwbc,soal_gambarjwbd,soal_jwbd,soal_jwbbenar,soal_bobot,soal_tipe,soal_is_delete,soal_userupdate,soal_lastupdate) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_, err := database.DbMain.Exec(QueryInsert, soalData.SoalId, soalData.BankId, &soalData.GambarSoal, &soalData.PertanyaanSoal, &soalData.UrutanSoal, &soalData.GambarA, &soalData.JawabanA, &soalData.GambarB, &soalData.JawabanB, &soalData.GambarC, &soalData.JawabanC, &soalData.GambarD, &soalData.JawabanD, &soalData.JawabanBenar, &soalData.BobotSoal,&soalData.TipeSoal,"n",soalData.UserUpdate,soalData.LastUpdate)
	if err != nil {
		middleware.LogError(err, "Insert Data Failed")
		return false, err
	}
	return true, nil
}
func UpdateSoalUjian(soalData *model.SoalList) (bool,error) {
	QueryUpdate := `UPDATE hd_soalujian SET `
	UpdateString := ""
	var args []interface{}

	// if soalData.PertanyaanSoal != "" {
	if UpdateString != "" {
		UpdateString += `, `
	}
	UpdateString += `soal_pertanyaan = ? `
	args = append(args,soalData.PertanyaanSoal)
	// }
	// if soalData.GambarSoal != "" {
	if UpdateString != "" {
		UpdateString += `, `
	}
	UpdateString += `soal_gambar_pertanyaan = ? `
	args = append(args,soalData.GambarSoal)
	// }

	if UpdateString != "" {
		UpdateString += `, `
	}
	UpdateString += `soal_urutan = ? `
	args = append(args,soalData.UrutanSoal)

	// if soalData.GambarA != "" {
	if UpdateString != "" {
		UpdateString += `, `
	}
	UpdateString += `soal_gambarjwba = ? `
	args = append(args,soalData.GambarA)
	// }
	// if soalData.JawabanA != "" {
	if UpdateString != "" {
		UpdateString += `, `
	}
	UpdateString += `soal_jwba = ? `
	args = append(args,soalData.JawabanA)
	// }
	// if soalData.GambarB != "" {
	if UpdateString != "" {
		UpdateString += `, `
	}
	UpdateString += `soal_gambarjwbb = ? `
	args = append(args,soalData.GambarB)
	// }
	// if soalData.JawabanB != "" {
	if UpdateString != "" {
		UpdateString += `, `
	}
	UpdateString += `soal_jwbb = ? `
	args = append(args,soalData.JawabanB)
	// }
	// if soalData.GambarC != "" {
	if UpdateString != "" {
		UpdateString += `, `
	}
	UpdateString += `soal_gambarjwbc = ? `
	args = append(args,soalData.GambarC)
	// }
	// if soalData.JawabanC != "" {
	if UpdateString != "" {
		UpdateString += `, `
	}
	UpdateString += `soal_jwbc = ? `
	args = append(args,soalData.JawabanC)
	// }
	// if soalData.GambarD != "" {
	if UpdateString != "" {
		UpdateString += `, `
	}
	UpdateString += `soal_gambarjwbd = ? `
	args = append(args,soalData.GambarD)
	// }
	// if soalData.JawabanD != "" {
	if UpdateString != "" {
		UpdateString += `, `
	}
	UpdateString += `soal_jwbd = ? `
	args = append(args,soalData.JawabanD)
	// }
	// if soalData.JawabanBenar != "" {
	if UpdateString != "" {
		UpdateString += `, `
	}
	UpdateString += `soal_jwbbenar = ? `
	args = append(args,soalData.JawabanBenar)
	// }
	
	if UpdateString != "" {
		UpdateString += `, `
	}
	UpdateString += `soal_bobot = ? `
	args = append(args,soalData.BobotSoal)
	
	// if soalData.TipeSoal != "" {
	if UpdateString != "" {
		UpdateString += `, `
	}
	UpdateString += `soal_tipe = ? `
	args = append(args,soalData.TipeSoal)
	// }
	// if soalData.UserUpdate != "" {
	if UpdateString != "" {
		UpdateString += `, `
	}
	UpdateString += `soal_userupdate = ? `
	args = append(args,soalData.UserUpdate)
	// }

	// if soalData.LastUpdate != "" {
	if UpdateString != "" {
		UpdateString += `, `
	}
	UpdateString += `soal_lastupdate = ? `
	args = append(args,soalData.LastUpdate)
	// }

	if len(UpdateString) > 0 {
		QueryUpdate = fmt.Sprintf("%s %s", QueryUpdate, UpdateString)
	}
	QueryUpdate += ` WHERE soal_id = ?`
	args = append(args, soalData.SoalId)

	_, err := database.DbMain.Exec(QueryUpdate, args...)
	if err != nil {
		middleware.LogError(err, "Update Data Failed")
		return false, err
	}
	return true, nil
}
func DeleteSoalUjian(soalId string) (bool, error) {
	QueryUpdate := `UPDATE hd_soalujian SET soal_is_delete = "y" WHERE soal_id = ?`
	_, err := database.DbMain.Exec(QueryUpdate, soalId)
	if err != nil {
		middleware.LogError(err, "Update Data Failed")
		return false, err
	}
	return true, nil
}
func InsertBulkSoalujian(insertData []model.SoalList)(bool,string, error) {
	InsertQuery := "INSERT INTO hd_soalujian VALUES "
	InsertPlaceholder := ""
	InsertValues := []interface{}{}

	for i,record := range insertData {
		if i > 0 {
			InsertPlaceholder += ", "
		}
		InsertPlaceholder += `(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
		InsertValues = append(InsertValues, record.SoalId, record.BankId, &record.GambarSoal, &record.PertanyaanSoal, &record.UrutanSoal, &record.GambarA, &record.JawabanA, &record.GambarB, &record.JawabanB, &record.GambarC, &record.JawabanC, &record.GambarD, &record.JawabanD, &record.JawabanBenar, &record.BobotSoal,&record.TipeSoal,"n",record.UserUpdate,record.LastUpdate)
	}
	InsertQuery += InsertPlaceholder
	_, err := database.DbMain.Exec(InsertQuery,InsertValues...)
	if err != nil {
		middleware.LogError(err, "Insert Batch Failed")
		return false,InsertQuery, err
	}
	return true,InsertQuery, nil
}
// End Function Untuk Soal


