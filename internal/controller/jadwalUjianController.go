package controller

import(
	_ "fmt"
	"time"
	// "strconv"

	// IMPORTANT IMPORT FOR ALL CONTROLLER
	"net/http"
	"github.com/gin-gonic/gin"

	"hongde_backend/internal/service"
	"hongde_backend/internal/model"
	"hongde_backend/internal/config"
	// "hongde_backend/internal/middleware"
	// "hongde_backend/internal/thirdparty"
)

// Function For Jadwal Ujian
func GetJadwalUjian(c *gin.Context) {
	jadwalId := c.Param("jadwalId")
	if jadwalId != "" {
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
		})
	} else {
		jadwalData, err := service.GetJadwalUjian()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Error detected : "+err.Error(),
			})
			return
		}
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
			"all_bank_soal":allBank,
			"jadwal_data": jadwalData,
		})
	}
}
func InsertJadwalUjian(c *gin.Context) {
	validatedInput, _ := c.Get("validatedInput")
	jadwalInput := validatedInput.(*model.JadwalUjian)

	userId,_ := c.Get("userID")
	jadwalInput.UserUpdate = userId.(string)
	jadwalInput.LastUpdate = time.Now().In(config.TimeZone).Format("2006-01-02 15:04:05")

	jadwalInput.TokenUjian = service.RandSeq(5)

	statusInsert, err := service.InsertJadwalUjian(jadwalInput) 
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
func UpdateJadwalUjian(c *gin.Context) {
	validatedInput, _ := c.Get("validatedInput")
	jadwalInput := validatedInput.(*model.JadwalUjian)
	userId,_ := c.Get("userID")
	jadwalInput.UserUpdate,_ = userId.(string)
	jadwalInput.LastUpdate = time.Now().In(config.TimeZone).Format("2006-01-02 15:04:05")
	statusUpdate, err := service.UpdateJadwalUjian(jadwalInput) 
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

func ResetTokenJadwalUjian(c *gin.Context) {
	var TokenData struct {
		JadwalId int `json:"jadwal_id"`
	}
	if c.Bind(&TokenData) != nil {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Failed To read Body",
		})
		return
	}
	newToken := service.RandSeq(5)
	// intId,_ := strconv.Atoi(TokenData.JadwalId)
	userId,_ := c.Get("userID")
	jadwalData := model.JadwalUjian{
		JadwalId:TokenData.JadwalId,
		TokenUjian:newToken,
		UserUpdate:userId.(string),
		LastUpdate:time.Now().In(config.TimeZone).Format("2006-01-02 15:04:05"),
	}
	statusUpdate, err := service.UpdateJadwalUjian(&jadwalData) 
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
		"newToken":newToken,
	})
	return
}

func DeleteJadwalUjian(c *gin.Context) {
	jadwalId := c.Param("jadwalId")

	if jadwalId == "" {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "JadwalUjian ID Is Not Provided",
		})
		return
	}
	statusDelete, err := service.DeleteJadwalUjian(jadwalId)

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
// End Function For Jadwal Ujian