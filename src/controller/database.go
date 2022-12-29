package controller

import (
	"database/sql"
	"github.com/google/uuid"
	"log"
)

type DataBaseController struct {
	DB *sql.DB
}

func DataBaseConnect(connStr string) *DataBaseController {
	DBController := new(DataBaseController)
	DBController.DB, _ = sql.Open("postgres", connStr)
	return DBController
}

func (DBController *DataBaseController) CheckUser(email string, password string) string {
	result := DBController.DB.QueryRow("SELECT id FROM users WHERE email=$1 AND encrypted_password=$2",
		email, password)
	var id int
	err := result.Scan(&id)
	if err != nil {
		ShowError(err)
		log.Println(err)
		return ""
	} else {
		log.Println("User Found")
	}
	return DBController.GenerateUUID()
}

func (DBController *DataBaseController) CheckAccessToken(token string) bool {
	result := DBController.DB.QueryRow("SELECT id FROM access_tokens WHERE token=$1",
		token)
	var id int
	err := result.Scan(&id)
	if err != nil {
		ShowError(err)
		return false
	} else {
		log.Println("Token Found")
	}
	return true
}

func (DBController *DataBaseController) CreateUser(email string, password string) bool {
	result, err := DBController.DB.Exec("insert into users (email, encrypted_password) values ($1, $2)",
		email, password)
	if err != nil {
		ShowError(err)
		return false
	}
	log.Println(result.RowsAffected())
	return true
}

func (DBController *DataBaseController) GenerateUUID() string {
	acces_token := uuid.New()
	result, err := DBController.DB.Exec("insert into Access_tokens (token) values ($1)",
		acces_token)
	if err != nil {
		ShowError(err)
		return ""
	}
	result.RowsAffected()
	log.Println("Access Token Created")
	return acces_token.String()
}
