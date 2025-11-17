package service

import(
	"database/sql"

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