package controller

import(
	// "fmt"
	"time"

	// IMPORTANT IMPORT FOR ALL CONTROLLER
	"net/http"
	"github.com/gin-gonic/gin"

	"hongde_backend/internal/config"
	"hongde_backend/internal/service"
	"hongde_backend/internal/model"
)

func GetDetailJadwalUjian(c *gin.Context) {
	jadwalId := c.Param("jadwalId")
	if jadwalId == "" {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Jadwal ID Is Not Provided",
		})
		return
	}
	canRecap := "y"
	recapReason := ""
	allPesertaUjian,err := service.GetAllPesertaUjianJadwal(jadwalId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error detected : "+err.Error(),
		})
		return
	}

	doneTotal := 0
	for _, rowData := range allPesertaUjian {
		if rowData.TotalWaiting > 0 {
			canRecap = "n"
			recapReason = "Masih Ada Siswa Yang Membutuhkan Penilaian Manual"
			break
		}
		if rowData.IsRecapped == "y" {
			doneTotal++
		}
	}
	if len(allPesertaUjian) == doneTotal {
		canRecap = "n"
		recapReason = "Seluruh Siswa Telah Selesai Direkap"
	}
	jadwalData, err := service.GetJadwalById(jadwalId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error detected : "+err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"jadwal_data": jadwalData,
		"can_rekap":canRecap,
		"alasan_rekap":recapReason,
		"all_peserta_ujian":allPesertaUjian,
	})
}

func DetailSingleSiswa(c *gin.Context) {
	pesertaId := c.Param("pesertaId")
	if pesertaId == "" {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Peserta ID Is Not Provided",
		})
		return
	}
	allDone := `y`
	currentScore := 0.0
	pesertaData,err := service.GetDetailPesertaId(pesertaId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error detected : "+err.Error(),
		})
		return
	}
	for _, rowData := range pesertaData.SoalArray {
		if rowData.IsRight == "w" {
			allDone = "n"
		}else if rowData.IsRight == "y" {
			currentScore += rowData.BobotSoal 
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"peserta_data": pesertaData,
		"all_done":allDone,
		"current_score":currentScore,
	})
}

func FinishManualPeserta(c *gin.Context) {
	var postData struct {
		PesertaId string `json:"peserta_id"`
	}
	if c.Bind(&postData) != nil {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Failed To read Body",
		})
		return
	}

	allSoalPeserta, err := service.GetJawabanPesertaWithSoalData(postData.PesertaId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error detected : "+err.Error(),
		})
		return
	}
	updateStatus,err := service.FinalizePesertaUjian(postData.PesertaId,"",allSoalPeserta)
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

func SavePenilaianEssay(c *gin.Context) {
	var postData struct {
		PesertaId string `json:"peserta_id"`
		SoalId int `json:"soal_id"`
		BobotJawaban float64 `json:"bobot_jawaban"`
	}
	if c.Bind(&postData) != nil {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Failed To read Body",
		})
		return
	}

	isRight := "y"
	if postData.BobotJawaban == 0 {
		isRight = "n"
	}

	updateData := model.SoalPesertaUjian{
		PesertaId:postData.PesertaId,
		SoalId:postData.SoalId,
		IsRight:isRight,
		BobotSoal:postData.BobotJawaban,
	}
	statusUpdate,err := service.UpdateSinglePenilaianJawaban(&updateData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error detected : "+err.Error(),
		})
		return
	}
	allDone := `y`
	currentScore := 0.0
	soalArray,err := service.GetJawabanPesertaUjian(postData.PesertaId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error detected : "+err.Error(),
		})
		return
	}
	for _, rowData := range soalArray {
		if rowData.IsRight == "w" {
			allDone = "n"
		}else if rowData.IsRight == "y" {
			currentScore += rowData.BobotSoal 
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"status_update": statusUpdate,
		"all_done":allDone,
		"current_score":currentScore,
	})
}

func SavePenilaianOther(c *gin.Context) {
	var postData struct {
		PesertaId string `json:"peserta_id"`
		SoalId int `json:"soal_id"`
		isBenar string `json:"is_benar"`
	}
	if c.Bind(&postData) != nil {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Failed To read Body",
		})
		return
	}


	updateData := model.SoalPesertaUjian{
		PesertaId:postData.PesertaId,
		SoalId:postData.SoalId,
		IsRight:postData.isBenar,
	}
	statusUpdate,err := service.UpdateSinglePenilaianJawaban(&updateData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error detected : "+err.Error(),
		})
		return
	}
	allDone := `y`
	currentScore := 0.0
	soalArray,err := service.GetJawabanPesertaUjian(postData.PesertaId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error detected : "+err.Error(),
		})
		return
	}
	for _, rowData := range soalArray {
		if rowData.IsRight == "w" {
			allDone = "n"
		}else if rowData.IsRight == "y" {
			currentScore += rowData.BobotSoal 
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"status_update": statusUpdate,
		"all_done":allDone,
		"current_score":currentScore,
	})
}

func FinalizeJadwalUjian(c *gin.Context) {
	var postData struct {
		JadwalId string `json:"jadwal_id"`
	}
	if c.Bind(&postData) != nil {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Failed To read Body",
		})
		return
	}
	allPesertaUjian,err := service.GetAllPesertaUjianJadwal(postData.JadwalId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error detected : "+err.Error(),
		})
		return
	}
	userId,_ := c.Get("userID")
	for index,_ := range allPesertaUjian {
		allPesertaUjian[index].UserUpdate = userId.(string)
		allPesertaUjian[index].LastUpdate = time.Now().In(config.TimeZone).Format("2006-01-02 15:04:05")
	}
	statusUpdate,err := service.FinalizeHasilUjianJadwal(allPesertaUjian)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error detected : "+err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"status_update":statusUpdate,
		"all_peserta_ujian":allPesertaUjian,
	})
}