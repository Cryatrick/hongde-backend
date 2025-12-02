package service

import(
	"database/sql"

	"hongde_backend/internal/model"
	"hongde_backend/internal/database"
	"hongde_backend/internal/middleware"
)


func LoginAdmin(userName, userPassword string) (*model.UserLogin, error) {
	row := database.DbMain.QueryRow("SELECT usr_username, usr_nama, usr_role FROM hd_useradmin WHERE usr_username = ? AND usr_password = ? AND usr_isdelete = 'n'",userName, userPassword)
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

func GetAdminData(userName string) (*model.UserLogin, error) {
	row := database.DbMain.QueryRow("SELECT usr_username, usr_nama, usr_role FROM hd_useradmin WHERE usr_username = ? AND usr_isdelete = 'n'",userName)
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