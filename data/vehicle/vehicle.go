package vehicle

import (
	"github.com/curt-labs/sweetData/helpers/database"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
)

type Vehicle struct {
	ID            int
	BaseVehicle   BaseVehicle
	Submodel      Submodel
	Configuration Configuration
	AppID         int
	RegionID      int
}

type BaseVehicle struct {
	ID   int
	Name string
}

type Submodel struct {
	ID   int
	Name string
}

type Configuration struct {
	ID      int
	Configs []ConfigAttribute
}

type ConfigAttribute struct {
	ID     int
	VcdbID int
	Value  string
	Type   AttributeType
}

type ConfigAttributeType struct {
	ID         int
	Name       string
	AcesTypeID int
	Sort       int
}

var (
	checkVehicle   = `select ID from vcdb_Vehicle where BaseVehicleID = ? and SubModelID = ? and ConfigID = ? and AppID = ? and RegionID = ?`
	getVehicle     = `select ID, BaseVehicleID, SubModelID, ConfigID, AppID, RegionID from vcdb_Vehicle where ID = ?`
	insertVehicle  = `insert into vcdb_Vehicle (BaseVehicleID, SubModelID, ConfigID, AppID, RegionID) values (?,?,?,?,?)`
	getMakeByName  = `select ID from vcdb_Make where MakeName = ?`
	getModelByName = `select ID from vcdb_Model where ModelName = ?`
)

//TODO  -crap - vehicles differenced go all the way back to the vcdb
func InsertVehicles(vs []Vehicle) error {
	for _, v := range vs {
		err = v.Get()
		if err != nil {
			return err
		}
		err = v.Insert()
		if err != nil {
			return err
		}
	}
}

//from Old DB, by ID
func (v *Vehicle) Get() error {
	var err error
	db, err := sql.Open("mysql", database.OldDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getVehicle)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var sub, con, app *int
	err = stmt.QueryRow(v.ID).Scan(&v.ID, &v.BaseVehicle.ID, &sub, &con, &app, &v.RegionID)
	if err != nil {
		return err
	}
	if sub != nil {
		v.Submodel.ID = *sub
	}
	if con != nil {
		v.Configuration.ID = *con
	}
	if app != nil {
		v.AppID = *app
	}
}

func (v *Vehicle) Check() (int, error) {
	var err error
	var i int
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return i, err
	}
	defer db.Close()

	stmt, err := db.Prepare(checkVehicle)
	if err != nil {
		return i, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(v.BaseVehicle.ID, v.Submodel.ID, v.Configuration.ID, v.AppID, v.RegionID).Scan(&i)
	return i, err
}

//TODO - make, model, sb and config FIRST
func (v *Vehicle) Insert() error {
	var err error
	//check for vehicle in new db
	v.ID, err = v.Check()
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if v.ID > 0 {
		return nil
	}

	//TODO - check/insert  Config Existence

	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertVehicle)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(v.BaseVehicle.ID, v.Submodel.ID, v.Configuration.ID, v.AppID, v.RegionID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	v.ID = int(id)
	return err
}
