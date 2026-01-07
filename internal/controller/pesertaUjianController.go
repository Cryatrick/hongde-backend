package controller

import(
	"fmt"
	"time"
	"strconv"
	"math/rand"

	// IMPORTANT IMPORT FOR ALL CONTROLLER
	"net/http"
	"github.com/gin-gonic/gin"

	"hongde_backend/internal/service"
	"hongde_backend/internal/model"
	"hongde_backend/internal/config"
)

// Function For Basic Peserta Ujian
func GetJadwalUjianSiswaToday(c *gin.Context) {
	// userId,_ := c.Get("userID")
	siswaid := c.Param("siswaId")
	currDate := time.Now().In(config.TimeZone).Format("2006-01-02")

	allJadwal,err := service.GetJadwalUjianPesertaByDate(siswaid,currDate)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error detected : "+err.Error(),
		})
		return
	}
	now := time.Now().In(config.TimeZone)
	for index, rowData := range allJadwal {
		allJadwal[index].ExamStatus = `allowed`
		ujianTimeFormat,_ := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%v %v",rowData.TanggalMulai,rowData.JamMulai), config.TimeZone)
		_, isClosed := CheckExamTiming(ujianTimeFormat,rowData.DurasiUjian,now,rowData.ToleransiTerlambat)
		if isClosed == true {
			allJadwal[index].ExamStatus = `closed`
		}
		if rowData.PesertaId != "" {
			allJadwal[index].ExamStatus = `taken`
			pesertaStartTime,_ := time.ParseInLocation("2006-01-02 15:04:05", rowData.StartPeserta, config.TimeZone)
			canResumeExam,_ := CanResumeExam(ujianTimeFormat,rowData.DurasiUjian,pesertaStartTime,now, rowData.KetentuanWaktu)
			if canResumeExam == false {
				allJadwal[index].ExamStatus = `finished`
			}

			if rowData.EndPeserta != "" {
				allJadwal[index].ExamStatus = `finished`
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"all_exam": allJadwal,
	})
}

func ProcessTokenExan(c *gin.Context) {
	var postData struct {
		SiswaId string `json:"siswa_id"`
		TokenUjian string `json:"token_ujian"`
	}

	if c.Bind(&postData) != nil {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Failed To read Body",
		})
		return
	}
	jadwalData, err := service.GetJadwalById(postData.TokenUjian)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error detected : "+err.Error(),
		})
		return
	}
	if jadwalData == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Token Ujian Tidak Ditemukan",
		})
		return
	}
	now := time.Now().In(config.TimeZone)
	ujianTimeFormat,_ := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%v %v",jadwalData.TanggalMulai,jadwalData.JamMulai), config.TimeZone)
	isLate, isClosed := CheckExamTiming(ujianTimeFormat,jadwalData.DurasiUjian,now,jadwalData.ToleransiTerlambat)
	if isClosed == true {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Ujian Telah Ditutup",
		})
		return
	}
	checkPeserta,err := service.CheckSudahMengikuti(jadwalData.JadwalId, postData.SiswaId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error detected : "+err.Error(),
		})
		return
	}
	if checkPeserta != nil {
		pesertaStartTime,_ := time.ParseInLocation("2006-01-02 15:04:05", checkPeserta.StartedAt, config.TimeZone)
		canResumeExam,remainingDuration := CanResumeExam(ujianTimeFormat,jadwalData.DurasiUjian,pesertaStartTime,now, jadwalData.KetentuanWaktu)
		if canResumeExam == false || checkPeserta.EndAt != "" {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Tidak Bisa Melanjutkan Ujian Lagi",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":    http.StatusOK,
			"jadwal_data":jadwalData,
			"peserta_data": checkPeserta,
			"sisa_waktu":remainingDuration,
		})
		return
	}
	pesertaId := fmt.Sprintf("%v%v",postData.SiswaId,jadwalData.JadwalId)
	var soalData []int
	var soalPeserta []model.SoalPesertaUjian
	if jadwalData.JenisSoal == `acak`{
		bankData, err := service.GetBankId(strconv.Itoa(jadwalData.BankSoal))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Error detected : "+err.Error(),
			})
			return
		}
		offset := rand.Intn(bankData.JumlahSoal - jadwalData.JumlahSoal)
		soalData,err = service.GetSoalIdArray(strconv.Itoa(jadwalData.BankSoal),strconv.Itoa(jadwalData.JumlahSoal), strconv.Itoa(offset))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Error detected : "+err.Error(),
			})
			return
		}
		r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
		// shuffled := make([]int, len(soalData))
		// copy(shuffled, soalData)

		r.Shuffle(len(soalData), func(i, j int) {
			soalData[i], soalData[j] = soalData[j], soalData[i]
		})
		// soalData = shuffled
	}else {
		soalData,err = service.GetSoalIdArray(strconv.Itoa(jadwalData.BankSoal),strconv.Itoa(jadwalData.JumlahSoal),"")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Error detected : "+err.Error(),
			})
			return
		}
	}
	for index,soalId := range soalData {
		soalPeserta = append(soalPeserta,model.SoalPesertaUjian{
			PesertaId:pesertaId,
			SoalId:soalId,
			UrutanSoal:(index+1),
		})
	}
	userId,_ := c.Get("userID")
	lateValue := "n"
	if isLate == true {
		lateValue = "y"
	}
	pesertaData := model.PesertaUjian{
		PesertaId:pesertaId,
		SiswaId:postData.SiswaId,
		UjianId:jadwalData.JadwalId,
		IsLate:lateValue,
		StartedAt:time.Now().In(config.TimeZone).Format("2006-01-02 15:04:05"),
		SoalArray:soalPeserta,
		UserUpdate:userId.(string),
		LastUpdate:time.Now().In(config.TimeZone).Format("2006-01-02 15:04:05"),
	}
	insertStatus,err := service.ProcessInsertPesertaUjian(&pesertaData)
	if err != nil || insertStatus == false {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error detected : "+err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"jadwal_data": jadwalData,
		"peserta_data":pesertaData,
		"sisa_waktu":jadwalData.DurasiUjian,
	})
	return
} 

func SaveJawabanPeserta(c *gin.Context) {
	var postData struct {
		PesertaId string `json:"peserta_id"`
		JawabanArray map[string]string `json:"jawaban_array"`
	}

	if c.Bind(&postData) != nil {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Failed To read Body",
		})
		return
	}

	updateStatus,err := service.SaveJawabanPesertaUjian(postData.PesertaId,postData.JawabanArray)
	if err != nil || updateStatus == false {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error detected : "+err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"updateStatus": updateStatus,
	})
	return
}

func CheckExamTiming(startDate time.Time, examDuration int ,now time.Time,graceMinutes int) (isLate bool, isClosed bool) {
	// fmt.Printf("Now :%v \n",now)
	lateAt := startDate
	// fmt.Printf("Late At :%v .\n",lateAt)
	closeAt := startDate.Add(time.Duration(graceMinutes) * time.Minute)
	// fmt.Printf("Close At :%v .\n",closeAt)


	isLate = now.After(lateAt)
	isClosed = now.After(closeAt)

	// fmt.Printf("Is Closed : %v",isClosed)
	// fmt.Printf("\n Is Late : %v",isLate)

	return isLate, isClosed
}

func CanResumeExam(examStartAt time.Time,examDurationMin int,studentStartAt time.Time,now time.Time, validateType string) (bool, time.Duration) {
	examEndAt := examStartAt.Add(
		time.Duration(examDurationMin) * time.Minute,
	)

	studentEndAt := studentStartAt.Add(
		time.Duration(examDurationMin) * time.Minute,
	)

	actualEndAt := studentEndAt
	if validateType == "mulai_potong" {
		actualEndAt = examEndAt
	}
	// if examEndAt.Before(actualEndAt) {
		// 	actualEndAt = examEndAt
		// }

		if now.After(actualEndAt) {
			return false, 0
		}

		remaining := actualEndAt.Sub(now)
		return true, remaining
	}