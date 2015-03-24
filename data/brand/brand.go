package brand

import (
	"github.com/curt-labs/sweetData/helpers/database"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
)

var (
	checkVideoBrand  = `select ID from VideoNewToBrand where videoID = ? and brandID = ?`
	insertVideoBrand = `insert into VideoNewToBrand (videoID,brandID) values (?,?)`
)

type Brand struct {
	ID int
}

func (b Brand) InsertVideoBrand(vid int) error {
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertVideoBrand)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(vid, b.ID)
	return err
}

func (b Brand) CheckVideoBrand(vid int) (int, error) {
	var err error
	var id int
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return id, err
	}
	defer db.Close()

	stmt, err := db.Prepare(checkVideoBrand)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(vid, b.ID).Scan(&id)
	return id, err
}
