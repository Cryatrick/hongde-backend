package controller

import(
	_ "fmt"
	"time"
	"os"
	"path/filepath"
	"strings"
	"strconv"

	// IMPORTANT IMPORT FOR ALL CONTROLLER
	"net/http"
	"github.com/gin-gonic/gin"

	"hongde_backend/internal/config"
	"hongde_backend/internal/service"
	"hongde_backend/internal/model"
	"hongde_backend/internal/middleware"
	"hongde_backend/internal/thirdparty"
)

// Function For Bank Soal
func GetBankSoal(c *gin.Context) {
	bankId := c.Param("bankId")
	if bankId != "" {
		bankData, err := service.GetBankId(bankId)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Error detected : "+err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":    http.StatusOK,
			"bank_data": bankData,
		})
	} else {
		allBank, err := service.GetBankSoal()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Error detected : "+err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":    http.StatusOK,
			"bank_data": allBank,
		})
	}
}
func InsertBankSoal(c *gin.Context) {
	validatedInput, _ := c.Get("validatedInput")
	bankInput := validatedInput.(*model.BankSoal)

	userId,_ := c.Get("userID")
	bankInput.UserUpdate = userId.(string)
	bankInput.LastUpdate = time.Now().In(config.TimeZone).Format("2006-01-02 15:04:05")

	statusInsert, err := service.InsertBankSoal(bankInput) 
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error While Insert Data detected : "+err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"statusInsert":  statusInsert,
	})
	return
}
func UpdateBankSoal(c *gin.Context) {
	validatedInput, _ := c.Get("validatedInput")
	bankInput := validatedInput.(*model.BankSoal)
	userId,_ := c.Get("userID")
	bankInput.UserUpdate,_ = userId.(string)
	bankInput.LastUpdate = time.Now().In(config.TimeZone).Format("2006-01-02 15:04:05")
	statusUpdate, err := service.UpdateBankSoal(bankInput) 
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error While Update Data detected : "+err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"statusUpdate":  statusUpdate,
	})
	return
}
func DeleteBankSoal(c *gin.Context) {
	bankId := c.Param("bankId")

	if bankId == "" {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "BankSoal ID Is Not Provided",
		})
		return
	}
	statusDelete, err := service.DeleteBankSoal(bankId)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error detected : "+err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"statusDelete":  statusDelete,
	})
}
// End Function For Bank Soal

// Function For Soal
func GetSoalList(c *gin.Context) {
	bankId := c.Param("bankId")
	soalId := c.Param("soalId")

	if soalId != "" {
		soaldata, err := service.GetSingleSoal(soalId)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Error detected : "+err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":    http.StatusOK,
			"soal_data": soaldata,
		})
		return
	}else {
		allSoal, err := service.GetSoalBank(bankId)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Error detected : "+err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":    http.StatusOK,
			"all_soal": allSoal,
		})
		return
	}
}
func SaveSoal(c *gin.Context) {
	soalId := c.PostForm("soal_id")
	bankId := c.PostForm("bank_id")
	urutanSoal := c.PostForm("urutan_soal")
	jawabanBenar := c.PostForm("jawaban_benar")
	bobotSoal := c.PostForm("bobot_soal")
	tipeSoal := c.PostForm("tipe_soal")
	queryType := "update"

	if bankId == "" || urutanSoal == "" || jawabanBenar == "" || bobotSoal == "" || tipeSoal == "" {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Please Fill All of The Required Fields",
		})
		return
	}
	if soalId == "" {
		generatedId,err := service.GenerateSoalId()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Error detected While Generating Soal ID : "+err.Error(),
			})
			return
		}
		queryType = "insert"
		soalId = generatedId
	}

	const maxFileSize = 3 * 1024 * 1024 // Max file size in bytes (3 MB)
	// Allowed file extensions
	var allowedExtensions = []string{".jpg",".png",".jpeg",".webp"}

	pertanyaanGambar := "n"
	fileNameGambarPertanyaan := ""
	pertanyaanText := c.PostForm("pertanyaan_text")
	pertanyaanFile, err :=c.FormFile("pertanyaan_gambar")
	if err == nil  {
		ext := strings.ToLower(filepath.Ext(pertanyaanFile.Filename))
		fileStatus, errMsg := middleware.ValidateFile(maxFileSize, pertanyaanFile.Size, ext, allowedExtensions)
		if fileStatus == false {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  errMsg,
			})
			return
		}
		f, err := pertanyaanFile.Open()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Could not Open Pertanyaan File : "+err.Error(),
			})
			return
		}
		// Save the file locally
		resultName, err := thirdparty.SaveImageAsWebp(f,soalId+"_SOAL", "./web/soal_images")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Error Detected While Converting To Webp : "+err.Error(),
			})
			return
		}
		pertanyaanGambar = "y"
		fileNameGambarPertanyaan = resultName
	}
	if pertanyaanGambar == "n" && pertanyaanText == "" {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Please Fill All of The Required Fields",
		})
		return
	}

	soalAGambar := "n"
	fileNameGambarSoalA := ""
	soalAText := c.PostForm("soal_a_text")
	soalAFile, err :=c.FormFile("soal_a_gambar")
	if err == nil  {
		ext := strings.ToLower(filepath.Ext(soalAFile.Filename))
		fileStatus, errMsg := middleware.ValidateFile(maxFileSize, soalAFile.Size, ext, allowedExtensions)
		if fileStatus == false {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  errMsg,
			})
			return
		}
		f, err := soalAFile.Open()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Could not Open A File : "+err.Error(),
			})
			return
		}
		// Save the file locally
		resultName, err := thirdparty.SaveImageAsWebp(f,soalId+"_JWBA", "./web/soal_images")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Error Detected While Converting To Webp : "+err.Error(),
			})
			return
		}
		soalAGambar = "y"
		fileNameGambarSoalA = resultName
	}
	if soalAGambar == "n" && soalAText == "" && tipeSoal != "essay" {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Please Fill All of The Required Fields",
		})
		return
	}

	soalBGambar := "n"
	fileNameGambarSoalB := ""
	soalBText := c.PostForm("soal_b_text")
	soalBFile, err :=c.FormFile("soal_b_gambar")
	if err == nil  {
		ext := strings.ToLower(filepath.Ext(soalBFile.Filename))
		fileStatus, errMsg := middleware.ValidateFile(maxFileSize, soalBFile.Size, ext, allowedExtensions)
		if fileStatus == false {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  errMsg,
			})
			return
		}
		f, err := soalBFile.Open()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Could not Open B File : "+err.Error(),
			})
			return
		}
		// Save the file locally
		resultName, err := thirdparty.SaveImageAsWebp(f,soalId+"_JWBB", "./web/soal_images")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Error Detected While Converting To Webp : "+err.Error(),
			})
			return
		}
		soalBGambar = "y"
		fileNameGambarSoalB = resultName
	}
	if soalBGambar == "n" && soalBText == "" && tipeSoal != "essay" {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Please Fill All of The Required Fields",
		})
		return
	}

	soalCGambar := "n"
	fileNameGambarSoalC := ""
	soalCText := c.PostForm("soal_c_text")
	soalCFile, err :=c.FormFile("soal_c_gambar")
	if err == nil  {
		ext := strings.ToLower(filepath.Ext(soalCFile.Filename))
		fileStatus, errMsg := middleware.ValidateFile(maxFileSize, soalCFile.Size, ext, allowedExtensions)
		if fileStatus == false {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  errMsg,
			})
			return
		}
		f, err := soalCFile.Open()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Could not Open C File : "+err.Error(),
			})
			return
		}
		// Save the file locally
		resultName, err := thirdparty.SaveImageAsWebp(f,soalId+"_JWBC", "./web/soal_images")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Error Detected While Converting To Webp : "+err.Error(),
			})
			return
		}
		soalCGambar = "y"
		fileNameGambarSoalC = resultName
	}
	if soalCGambar == "n" && soalCText == "" && tipeSoal == "pilihan_ganda"  {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Please Fill All of The Required Fields",
		})
		return
	}

	soalDGambar := "n"
	fileNameGambarSoalD := ""
	soalDText := c.PostForm("soal_d_text")
	soalDFile, err :=c.FormFile("soal_d_gambar")
	if err == nil  {
		ext := strings.ToLower(filepath.Ext(soalDFile.Filename))
		fileStatus, errMsg := middleware.ValidateFile(maxFileSize, soalDFile.Size, ext, allowedExtensions)
		if fileStatus == false {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  errMsg,
			})
			return
		}
		f, err := soalDFile.Open()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Could not Open D File : "+err.Error(),
			})
			return
		}
		// Save the file locally
		resultName, err := thirdparty.SaveImageAsWebp(f,soalId+"_JWBD", "./web/soal_images")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Error Detected While Donverting To Webp : "+err.Error(),
			})
			return
		}
		soalDGambar = "y"
		fileNameGambarSoalD = resultName
	}
	if soalDGambar == "n" && soalDText == "" && tipeSoal == "pilihan_ganda" {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Please Fill All of The Required Fields",
		})
		return
	}

	userId,_ := c.Get("userID")
	intId,_ :=  strconv.Atoi(soalId)
	intBankId,_ :=  strconv.Atoi(bankId)
	intUrutan,_ :=  strconv.Atoi(urutanSoal)
	floatBobot,_ :=  strconv.ParseFloat(bobotSoal,64)
	soalData := model.SoalList{
		SoalId:intId,
		BankId:intBankId,
		UrutanSoal:intUrutan,
		PertanyaanSoal:pertanyaanText,
		GambarSoal:fileNameGambarPertanyaan,
		JawabanA:soalAText,
		GambarA:fileNameGambarSoalA,
		JawabanB:soalBText,
		GambarB:fileNameGambarSoalB,
		JawabanC:soalCText,
		GambarC:fileNameGambarSoalC,
		JawabanD:soalDText,
		GambarD:fileNameGambarSoalD,
		JawabanBenar:jawabanBenar,
		BobotSoal:floatBobot,
		TipeSoal:tipeSoal,
		UserUpdate : userId.(string),
		LastUpdate : time.Now().In(config.TimeZone).Format("2006-01-02 15:04:05"),
	}
	var processResult bool
	if queryType == "insert" {
		res, err := service.InsertSoalUjian(&soalData)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Error Inserting : "+err.Error(),
			})
			return
		}
		processResult = res
	}else {
		res, err := service.UpdateSoalUjian(&soalData)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Error Updating : "+err.Error(),
			})
			return
		}
		processResult = res
	}
	c.JSON(http.StatusOK, gin.H{
		"status" :http.StatusOK,
		"message": "Check Data Please",
		"soal_data":soalData,
		"processResult":processResult,
	})
	return
}
func DeleteSoal(c *gin.Context) {
	soalId := c.Param("soalId")

	if soalId == "" {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Soal ID Is Not Provided",
		})
		return
	}
	statusDelete, err := service.DeleteSoalUjian(soalId)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error detected : "+err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"statusDelete":  statusDelete,
	})
}
func ReadExcelSoal(c *gin.Context) {
	bankId := c.PostForm("bank_id")
	if bankId == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Bank Id is requred",
		})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"error":  "file is required",
		})
		return
	}

	const maxFileSize = 3 * 1024 * 1024 // Max file size in bytes (3 MB)
	// Allowed file extensions
	var allowedExtensions = []string{".xlsx"}
	const uploadDir = "./temp/" // Directory to save uploaded files locally

	// Check file type (by extension)
	ext := strings.ToLower(filepath.Ext(file.Filename))

	fileStatus, errMsg := middleware.ValidateFile(maxFileSize, file.Size, ext, allowedExtensions)
	if fileStatus == false {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"error":  errMsg,
		})
		return
	}

	// Save the file locally
	localFileName := "temp" + ext
	localFilePath := filepath.Join(uploadDir, localFileName)
	if err := c.SaveUploadedFile(file, localFilePath); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "failed to save file locally",
		})
		return
	}
	startId, err := service.GenerateSoalId()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error Generating Soal ID : "+err.Error(),
		})
		return
	}
	intStartId,_ := strconv.Atoi(startId)
	intBankId,_ := strconv.Atoi(bankId)
	var insertSoalList []model.SoalList
	userId,_ := c.Get("userID")

	// To Start Process Import Pilihan Ganda
	_, excelData, err := thirdparty.ReadExcelFile(`pilihan_ganda`, localFilePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed To Read Excel : " + err.Error(),
		})
		return
	}

	for index, excelRow := range excelData {
		if excelRow["Urutan soal"] == nil || excelRow["Pertanyaan"] == nil || excelRow["Bobot Benar"] == nil || excelRow["Jawaban A"] == nil || excelRow["Jawaban B"] == nil || excelRow["Jawaban C"] == nil || excelRow["Jawaban D"] == nil || excelRow["Jawaban Benar"] == nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Detected Null For Required Field On Row - "+strconv.Itoa(index+1) + " At pilihan_ganda Sheet",
			})
			return
		}
		var soalRow model.SoalList
		soalRow.SoalId = intStartId
		soalRow.BankId = intBankId

		soalRow.PertanyaanSoal = excelRow["Pertanyaan"].(string)
		soalRow.UrutanSoal,_ = strconv.Atoi(excelRow["Urutan soal"].(string))
		soalRow.JawabanA = excelRow["Jawaban A"].(string)
		soalRow.JawabanB = excelRow["Jawaban B"].(string)
		soalRow.JawabanC = excelRow["Jawaban C"].(string)
		soalRow.JawabanD = excelRow["Jawaban D"].(string)
		soalRow.JawabanBenar = excelRow["Jawaban Benar"].(string)
		soalRow.BobotSoal,_ = strconv.ParseFloat(excelRow["Bobot Benar"].(string),64)
		soalRow.TipeSoal = "pilihan_ganda"
		soalRow.UserUpdate,_ = userId.(string)
		soalRow.LastUpdate = time.Now().In(config.TimeZone).Format("2006-01-02 15:04:05")

		insertSoalList = append(insertSoalList,soalRow)

		intStartId += 1
	}
	// End Process Import Pilihan Ganda

	// To Start Process Import True Or False
	_, excelTrue, err := thirdparty.ReadExcelFile(`true_or_false`, localFilePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed To Read Excel : " + err.Error(),
		})
		return
	}

	for index, rowTrue := range excelTrue {
		if rowTrue["Urutan soal"] == nil || rowTrue["Pertanyaan"] == nil || rowTrue["Bobot Benar"] == nil ||  rowTrue["Jawaban Benar"] == nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Detected Null For Required Field On Row - "+strconv.Itoa(index+1) + " At true_or_false Sheet",
			})
			return
		}
		var soalRow model.SoalList
		soalRow.SoalId = intStartId
		soalRow.BankId = intBankId

		soalRow.PertanyaanSoal = rowTrue["Pertanyaan"].(string)
		soalRow.UrutanSoal,_ = strconv.Atoi(rowTrue["Urutan soal"].(string))
		soalRow.JawabanA = "True"
		soalRow.JawabanB = "False"
		if(strings.ToLower(rowTrue["Jawaban Benar"].(string)) == "true") {
			soalRow.JawabanBenar = "a"
		}else {
			soalRow.JawabanBenar = "b"
		}
		soalRow.BobotSoal,_ = strconv.ParseFloat(rowTrue["Bobot Benar"].(string),64)
		soalRow.TipeSoal = "true_false"
		soalRow.UserUpdate,_ = userId.(string)
		soalRow.LastUpdate = time.Now().In(config.TimeZone).Format("2006-01-02 15:04:05")

		insertSoalList = append(insertSoalList,soalRow)

		intStartId += 1
	}
	// End Process Import True Or False

	// To Start Process Import Essay
	_, excelEssay, err := thirdparty.ReadExcelFile(`essay`, localFilePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed To Read Excel : " + err.Error(),
		})
		return
	}

	for index, rowEssay := range excelEssay {
		if rowEssay["Urutan soal"] == nil || rowEssay["Pertanyaan"] == nil  {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Detected Null For Required Field On Row - "+strconv.Itoa(index+1) + " At essay Sheet",
			})
			return
		}
		var soalRow model.SoalList
		soalRow.SoalId = intStartId
		soalRow.BankId = intBankId

		soalRow.PertanyaanSoal = rowEssay["Pertanyaan"].(string)
		soalRow.UrutanSoal,_ = strconv.Atoi(rowEssay["Urutan soal"].(string))
		soalRow.TipeSoal = "essay"
		soalRow.UserUpdate,_ = userId.(string)
		soalRow.LastUpdate = time.Now().In(config.TimeZone).Format("2006-01-02 15:04:05")

		insertSoalList = append(insertSoalList,soalRow)

		intStartId += 1
	}
	// End Process Import Essay


	// Delete the local file after successful upload
	if err := os.Remove(localFilePath); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "failed to clean up local file",
		})
		return
	}

 	insertStatus, err := service.InsertBulkSoalujian(insertSoalList)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error While Inserting Soal : "+err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"insertStatus": insertStatus,
		"excelData": gin.H{
			"data":   insertSoalList,
		},
	})
}


