package controller

import (
	"database/sql"
	"fmt"
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
	var id int32
	err := result.Scan(&id)
	if err != nil {
		ShowError(err)
		log.Println(err)
		return ""
	} else {
		log.Println("User Found")
	}
	return DBController.GenerateUUID(id)
}

func (DBController *DataBaseController) CheckAccessToken(token string) (bool, int32) {
	result := DBController.DB.QueryRow("SELECT owner_id FROM access_tokens WHERE token=$1",
		token)
	var uid int32
	err := result.Scan(&uid)
	if err != nil {
		ShowError(err)
		return false, -1
	} else {
		log.Println("Token Found")
	}
	return true, uid
}

func (DBContoller *DataBaseController) InvalidateAccessToken(token string) {
	result, err := DBContoller.DB.Exec("update access_tokens set valid = $1 where token = $2", false, token)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(result.RowsAffected())
}

func (DBController *DataBaseController) RemoveAccessToken(token string) {
	result, err := DBController.DB.Exec("delete from access_tokens where token = $1", token)
	if err != nil {
		log.Println(err)
	}
	result.RowsAffected()
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

func (DBController *DataBaseController) GenerateUUID(id int32) string {
	acces_token := uuid.New()
	result, err := DBController.DB.Exec("insert into Access_tokens (token, owner_id) values ($1, $2)",
		acces_token, id)
	if err != nil {
		ShowError(err)
		return ""
	}
	result.RowsAffected()
	log.Println("Access Token Created")
	return acces_token.String()
}

func (DBController *DataBaseController) GetCharacters(ownerID int32) []string {
	rows, err := DBController.DB.Query("SELECT name,appearance_data,level,id FROM characters WHERE owner_id=$1",
		ownerID)
	var charactersData []string
	for rows.Next() {
		var characterData string
		var appearanceData string
		var name string
		var level string
		var id string
		err = rows.Scan(&name, &appearanceData, &level, &id)
		if err != nil {
			fmt.Println(err)
			continue
		}
		characterData = name + "|" + appearanceData + "|" + level + "|" + id
		charactersData = append(charactersData, characterData)
	}
	if err != nil {
		ShowError(err)
		log.Println(err)
	} else {
		log.Println("Character loads")
	}

	log.Println(charactersData)
	return charactersData
}

func (DBController *DataBaseController) CreateCharacter(name string, appearanceData string, ownerID int32) bool {
	result, err := DBController.DB.Exec("insert into characters (name, appearance_data, owner_id) values ($1, $2, $3)", name, appearanceData, ownerID)
	if err != nil {
		ShowError(err)
	}
	r, _ := result.RowsAffected()
	if r == 1 {
		return true
	} else {
		return false
	}
}
