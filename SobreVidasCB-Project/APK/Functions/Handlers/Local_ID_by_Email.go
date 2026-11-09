package Handlers

import (
	Database "SobreVidasCB-layout/Functions/DB"
)

func LocalizaIdPorEmail(email string) int {
	var id int
	db := Database.OpenConnection()
	db.QueryRow("SELECT id FROM ACS WHERE email=$1", email).Scan(&id)
	db.Close()

	return id
}