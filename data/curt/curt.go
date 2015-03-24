package curt

import (
	"database/sql"
	"github.com/curt-labs/sweetData/helpers/database"
)

var (
	everytable = `select distinct TABLE_NAME from INFORMATION_SCHEMA.COLUMNS where TABLE_SCHEMA='CurtDev'`
)

func GetEverything() {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(everytable)
	if err != nil {
		return err
	}
	defer stmt.Close()
}
