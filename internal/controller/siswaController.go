package controller


import(
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"net/http"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"

	"hongde_backend/internal/middleware"
	"hongde_backend/internal/service"
	"hongde_backend/internal/model"
	"hongde_backend/internal/thirdparty"
)
 
func GetSiswa(c *gin.Context) {
	siswaid := c.Param("siswaId")
	if siswaid != "" {
		siswa, err := service.GetSiswaId(siswaid)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Error detected : "+err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":    http.StatusOK,
			"siswa_data": siswa,
		})
	} else {
		allSiswa, err := service.GetAllSiswa()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Error detected : "+err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":    http.StatusOK,
			"siswa_data": allSiswa,
		})
	}
}

func InsertSiswa(c *gin.Context) {
	validatedInput, _ := c.Get("validatedInput")
	siswaInput := validatedInput.(*model.SiswaList)

	GeneratedId, err := service.GenerateSiswaId()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error detected : "+err.Error(),
		})
		return
	}
	siswaInput.SiswaId = GeneratedId
	hash := md5.Sum([]byte(siswaInput.SiswaId))
	passwordHash := hex.EncodeToString(hash[:])
	siswaInput.PasswordSiswa = passwordHash
	userId,_ := c.Get("userID")
	siswaInput.UserUpdate = userId.(string)
	siswaInput.LastUpdate = time.Now().Format("2006-01-02 15:04:05")
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"check_data":  siswaInput,
	})
	return
}

func UpdateSiswa(c *gin.Context) {
	validatedInput, _ := c.Get("validatedInput")
	siswaInput := validatedInput.(*model.SiswaList)
	userId,_ := c.Get("userID")
	siswaInput.UserUpdate,_ = userId.(string)
	siswaInput.LastUpdate = time.Now().Format("2006-01-02 15:04:05")
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"check_data":  siswaInput,
	})
	return
}

func DeleteSiswa(c *gin.Context) {
	siswaid := c.Param("siswaId")

	if siswaid == "" {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Siswa ID Is Not Provided",
		})
		return
	}
	statusDelete, err := service.DeleteSiswa(siswaid)

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

func ResetPasswordSiswa(c *gin.Context) {
	var PasswordData struct {
		SiswaId string `json:"siswa_id"`
	}
	if c.Bind(&PasswordData) != nil {
		c.JSON(http.StatusOK, gin.H{
			"status" :http.StatusInternalServerError,
			"error": "Failed To read Body",
		})
		return
	}
	newPassword := service.RandSeq(15)
	hash := md5.Sum([]byte(newPassword))
	passwordHash := hex.EncodeToString(hash[:])
	userId,_ := c.Get("userID")
	siswaData := model.SiswaList{
		SiswaId:PasswordData.SiswaId,
		PasswordSiswa:passwordHash,
		UserUpdate:userId.(string),
		LastUpdate:time.Now().Format("2006-01-02 15:04:05"),
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"check_data":  siswaData,
	})
	return
}

func ReadExcelSiswa(c *gin.Context) {
	sheetName := c.PostForm("sheetName")
	if sheetName == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"error":  "sheet name is requred",
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
	excelHeader, excelData, err := thirdparty.ReadExcelFile(sheetName, localFilePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed To Read Excel : " + err.Error(),
		})
		return
	}

	startId, err := service.GenerateSiswaId()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Error Generating Siswa ID : "+err.Error(),
		})
		return
	}
	intStartId,_ := strconv.Atoi(startId)
	var insertSiswaList []model.SiswaList
	userId,_ := c.Get("userID")

	for index, excelRow := range excelData {
		if excelRow["Nama Siswa"].(string) == "" || excelRow["Email Siswa"].(string) == "" || excelRow["Jenis Identitas"].(string) == "" || excelRow["Nomor Identitas"].(string) == "" || excelRow["Tempat Lahir"].(string) == "" || excelRow["Tanggal Lahir"].(string) == "" || excelRow["Tempat Tinggal"].(string) == "" || excelRow["No Kontak"].(string) == "" || excelRow["Jenis Siswa"].(string) == "" {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Detected Null For Required Field On Row - "+strconv.Itoa(index),
			})
			return
		}
		var siswaRow model.SiswaList
		siswaRow.SiswaId = strconv.Itoa(intStartId)

		hash := md5.Sum([]byte(siswaRow.SiswaId))
		passwordHash := hex.EncodeToString(hash[:])
		siswaRow.PasswordSiswa = passwordHash


		siswaRow.NamaSiswa = excelRow["Nama Siswa"].(string)
		siswaRow.EmailSiswa = excelRow["Email Siswa"].(string)
		if excelRow["Nama Mandarin"].(string) != "" {
			siswaRow.NamaMandarin.Valid = true
			siswaRow.NamaMandarin.Str = excelRow["Nama Mandarin"].(string)
		}
		siswaRow.JenisIdentitas = excelRow["Jenis Identitas"].(string)
		siswaRow.NoIdentitas = excelRow["No Identitas"].(string)
		siswaRow.TempatLahir = excelRow["Tempat Lahir"].(string)
		siswaRow.TanggalLahir = excelRow["Tanggal Lahir"].(string)
		siswaRow.NoKontak = excelRow["No Kontak"].(string)
		siswaRow.JenisSiswa = excelRow["Jenis Siswa"].(string)

		siswaRow.UserUpdate,_ = userId.(string)
		siswaRow.LastUpdate = time.Now().Format("2006-01-02 15:04:05")

		insertSiswaList = append(insertSiswaList,siswaRow)

		intStartId += 1
	}


	// Delete the local file after successful upload
	if err := os.Remove(localFilePath); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "failed to clean up local file",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Check If Excel Data Is Right",
		"excelData": gin.H{
			"header": excelHeader,
			"data":   insertSiswaList,
		},
	})
}

