package controller

import (
	"database/sql"
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

func (DBController *DataBaseController) CheckUser(email string, password string) {
	result := DBController.DB.QueryRow("SELECT id FROM users WHERE email=$1 AND encrypted_password=$2",
		email, password)
	var id int
	err := result.Scan(&id)
	if err == sql.ErrNoRows {
		ShowError(err)
	} else if err != nil {
		ShowError(err)
	} else {
		log.Println("User Found: id-", string(id))
	}
}

func (DBController *DataBaseController) CreateUser(email string, password string) {
	result, err := DBController.DB.Exec("insert into users (email, encrypted_password) values ($1, $2)",
		email, password)
	if err != nil {
		ShowError(err)
		return
	}
	log.Println(result.RowsAffected())
}
