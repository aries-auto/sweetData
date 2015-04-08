package vehicle

import (
	"github.com/curt-labs/sweetData/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	// "gopkg.in/mgo.v2"

	"database/sql"
	"encoding/json"
	// "errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	ID     int
	AAIAID int
	Year   Year
	Make   Make
	Model  Model
}

type Year struct {
	ID int
}
type Make struct {
	ID     int
	AAIAID int
	Name   string
}
type Model struct {
	ID            int
	AAIAID        int
	Name          string
	VehicleTypeID int
}

type Submodel struct {
	ID     int
	AAIAID int
	Name   string
}

type Configuration struct {
	ID      int
	Configs []ConfigAttribute
}

type ConfigAttribute struct {
	ID     int
	VcdbID int
	Value  string
	Type   ConfigAttributeType
}

type ConfigAttributeType struct {
	ID         int
	Name       string
	AcesTypeID int
	Sort       int
}

//from parts
type PartVehicle struct {
	ID       int      `json:"id,omitempty" xml:"id,omitempty"`
	Year     int      `json:"year,omitempty" xml:"year,omitempty"`
	Make     string   `json:"make,omitempty" xml:"make,omitempty"`
	Model    string   `json:"model,omitempty" xml:"model,omitempty"`
	Submodel string   `json:"submodel,omitempty" xml:"submodel,omitempty"`
	Config   []Config `json:"configuration,omitempty" xml:"configuration,omitempty"`
}

type Config struct {
	Type  string `json:"type,omitempty" xml:"type,omitempty"`
	Value string `json:"value,omitempty" xml:"value,omitempty"`
}

const (
	mongoErrorCollection = "errors"
)

var (
	insertVehicle = `insert into vcdb_Vehicle (BaseVehicleID, SubModelID, ConfigID, AppID, RegionID) values (?,?,?,?,?)`
	getMakeByName = `select ID from vcdb_Make where AAIAMakeID = ? and MakeName = ?`
	insertMake    = `insert into vcdb_Make (AAIAMakeID, MakeName) values (?,?)`
	allMakes      = `select AAIAMakeID, MakeName from vcdb_Make `
	makeIDs       = `select ID, AAIAMakeID from vcdb_Make`

	getModelByName = `select ID from vcdb_Model where AAIAModelID = ? and ModelName = ?`
	insertModel    = `insert into vcdb_Model (AAIAModelID, ModelName, VehicleTypeID) values (?,?,?)`
	allModels      = `select AAIAModelID, ModelName, VehicleTypeID from vcdb_Model`
	modelIDs       = `select ID, AAIAModelID from vcdb_Model`

	allBaseVehicles = `select b.AAIABaseVehicleID, b.YearID, ma.AAIAMakeID, mo.AAIAModelID
		from BaseVehicle as b
		join vcdb_Make as ma on ma.ID = b.MakeID
		join vcdb_Model as mo on mo.ID = b.ModelID`

	getBaseVehicleByYMM = `select b.ID
		from BaseVehicle as b
		join vcdb_Make as ma on ma.ID = b.MakeID
		join vcdb_Model as mo on mo.ID = b.ModelID
		where b.YearID = ?
		and ma.MakeName = ?
		and mo.ModelName = ?`
	checkBaseVehicles = `select ID from BaseVehicle where AAIABaseVehicleID = ?`
	insertBaseVehicle = `insert into BaseVehicle (AAIABaseVehicleID, YearID, MakeID, ModelID) values (?,?,?,?)`

	getSubmodelByName = `select ID from Submodel where SubmodelName = ? and AAIASubmodelID = ?`
	getSubmodel       = `select ID from Submodel where SubmodelName = ?`
	insertSubmodel    = `insert into Submodel (AAIASubmodelID, SubmodelName) values (?,?)`
	allSubmodels      = `select AAIASubmodelID, SubmodelName from Submodel`

	checkAttributeType = `select ID from ConfigAttributeType where name = ? and AcesTypeID = ? and sort = ?`
	checkAttribute     = `select ca.ID from ConfigAttribute as ca join ConfigAttributeType as cat on 
		cat.ID = ca.ConfigAttributeTypeID where ca.vcdbID = ? and ca.value = ? and cat.name = ?`
	insertAttributeType = `insert into ConfigAttributeType (name, AcesTypeID, sort) values (?,?,?)`
	insertAttribute     = `insert into ConfigAttribute (ConfigAttributeTypeID, parentID, vcdbID, value) values (?,0,?,?)`
	allConfigAttributes = `select cat.name, cat.AcesTypeID, cat.sort,
		ca.vcdbID, ca.value
		from ConfigAttributeType as cat
		join ConfigAttribute as ca on ca.ConfigAttributeTypeID = cat.ID`

	getConfigAttribute = `select ca.ID from ConfigAttribute as ca 
		join ConfigAttributeType as cat on cat.ID = ca.ConfigAttributeTypeID
		where cat.name = ?
		and ca.value = ?`
	insertVehicleConfig          = `insert into VehicleConfig (AAIAVehicleConfigID) values (0)`
	insertVehicleConfigAttribute = `insert into VehicleConfigAttribute (AttributeID, VehicleConfigID) values (?,?)`
	insertVehiclePartJoin        = `insert into vcdb_VehiclePart (VehicleID, PartNumber) values (?,?)`
	findVehicle                  = `select v.ID, b.AAIABaseVehicleID, b.YearID, ma.AAIAMakeID, ma.MakeName, mo.AAIAModelID, mo.ModelName, 
		s.AAIASubmodelID, s.SubmodelName, ca.ID, ca.vcdbID, ca.value, cat.AcesTypeID, cat.Name
		from vcdb_Vehicle as v 
		join BaseVehicle as b on b.ID = v.BaseVehicleID
		join vcdb_Make as ma on ma.ID = b.MakeID
		join vcdb_Model as mo on mo.ID = b.ModelID
		left join Submodel as s on s.ID = v.SubmodelID
		left join VehicleConfigAttribute as vca on vca.VehicleConfigID = v.ConfigID
		left join ConfigAttribute as ca on ca.ID = vca.AttributeID
		left join ConfigAttributeType as cat on cat.ID = ca.ConfigAttributeTypeID
		where b.YearID = ?
		and ma.MakeName = ?
		and mo.ModelName = ?`
	submodelAddon      = ` and s.SubmodelName = ?`
	submodelNullAddon  = ` and (v.SubmodelID is null or v.SubmodelID = 0)`
	configNullAddon    = ` and (v.ConfigID = 0 or v.ConfigID is null)`
	configNotNullAddon = ` and (v.ConfigID > 0 and v.ConfigID is not null)`
)

func init() {
	database.InitMongo()
}

func ImportVehicles() error {
	var err error
	err = ImportMakes()
	if err != nil {
		log.Print("makes", err)
		return err
	}
	err = ImportModels()
	if err != nil {
		log.Print("models", err)
		return err
	}
	err = ImportBaseVehicles()
	if err != nil {
		log.Print("base", err)
		return err
	}
	err = ImportSubmodels()
	if err != nil {
		log.Print("subs", err)
		return err
	}
	err = ImportConfigs()
	if err != nil {
		log.Print("configs", err)
		return err
	}
	log.Print("V import done")
	return err
}

func InsertPartVehicles(partID int) error {
	var err error
	vs, err := getPartVehicles(partID)
	if err != nil {
		return err
	}

	for _, vehicle := range vs {
		err = InsertPartVehicle(vehicle, partID)
	}
	return err
}

//api call to get PartVehicle
func getPartVehicles(partID int) ([]PartVehicle, error) {
	var vs []PartVehicle
	res, err := http.Get(database.Api + "part/" + strconv.Itoa(partID) + "/vehicles?key=" + database.ApiKey)
	if err != nil {
		return vs, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return vs, err
	}

	err = json.Unmarshal(body, &vs)
	return vs, err
}

func InsertPartVehicle(ve PartVehicle, partID int) error {
	var err error
	var v Vehicle
	//find vehicle
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	var res *sql.Rows

	//Get V.ID, based on whether there is submodel/configs or not
	if ve.Submodel == "" {
		stmt, err := db.Prepare(findVehicle + submodelNullAddon)
		if err != nil {
			return err
		}
		defer stmt.Close()
		res, err = stmt.Query(ve.Year, ve.Make, ve.Model)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
	} else if len(ve.Config) < 1 {
		stmt, err := db.Prepare(findVehicle + submodelAddon + configNullAddon)
		if err != nil {
			return err
		}
		defer stmt.Close()
		res, err = stmt.Query(ve.Year, ve.Make, ve.Model, ve.Submodel)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
	} else {
		stmt, err := db.Prepare(findVehicle + submodelAddon + configNotNullAddon)
		if err != nil {
			return err
		}
		defer stmt.Close()
		res, err = stmt.Query(ve.Year, ve.Make, ve.Model, ve.Submodel)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
	}

	var c ConfigAttribute
	var cs []ConfigAttribute
	var subai, convcid, contypeid, caid *int
	var subname, convalue, contypename *string
	for res.Next() {
		err = res.Scan(
			&v.ID,
			&v.BaseVehicle.AAIAID,
			&v.BaseVehicle.Year.ID,
			&v.BaseVehicle.Make.AAIAID,
			&v.BaseVehicle.Make.Name,
			&v.BaseVehicle.Model.AAIAID,
			&v.BaseVehicle.Model.Name,
			&subai,
			&subname,
			&caid,
			&convcid,
			&convalue,
			&contypeid,
			&contypename,
		)
		if err != nil {
			return err
		}

		if subai != nil {
			v.Submodel.AAIAID = *subai
		}
		if subname != nil {
			v.Submodel.Name = *subname
		}
		if convcid != nil {
			c.VcdbID = *convcid
		}
		if convalue != nil {
			c.Value = *convalue
		}
		if contypeid != nil {
			c.Type.AcesTypeID = *contypeid
		}
		if contypename != nil {
			c.Type.Name = *contypename
		}
		if caid != nil {
			c.ID = *caid
		}
		cs = append(cs, c)
	}
	v.Configuration.Configs = cs

	//there is a vehicle, check for matching configs
	if v.ID > 0 {
		//check a match of config arrays
		veConfigMap := make(map[string]string)
		for _, con := range ve.Config {
			veConfigMap[con.Type] = con.Value
		}
		for _, config := range v.Configuration.Configs {
			if val, ok := veConfigMap[config.Type.Name]; ok {
				if val != config.Value {
					//value wrong
					v.ID, err = ve.Insert()
					if err != nil {
						return err
					}
				}
			}
		}

	}

	//Vehicle doesn't exist (id == 0), create
	sess := database.MongoSession.Copy() //to save errors
	if v.ID == 0 {
		if v.BaseVehicle.ID == 0 {
			err = sess.DB("DataMigration").C(mongoErrorCollection).Insert(v)
		}
		err = v.Insert()
		if err != nil {
			return err
		}
	}
	//join part
	stmt, err := db.Prepare(insertVehiclePartJoin)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(v.ID, partID)

	return err
}

//From a partVehicle, insert a REAL vehicle
func (ve *PartVehicle) Insert() (int, error) {
	var err error
	var vid int
	var BaseVehicleID, SubModelID, ConfigID int

	sess := database.MongoSession.Copy()

	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return 0, err
	}
	defer db.Close()

	//get baseID, subId, ConfigID
	stmt, err := db.Prepare(getBaseVehicleByYMM)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(ve.Year, ve.Make, ve.Model).Scan(&BaseVehicleID)
	if err != nil {
		// log.Panic("You missed importing a base vehicle")
		err = sess.DB("DataMigration").C(mongoErrorCollection).Insert(err.Error())
		return 0, err

	}

	stmt, err = db.Prepare(getSubmodel)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(ve.Submodel).Scan(&SubModelID)
	if err != nil {
		// log.Panic("You missed importing a submodel")

		err = sess.DB("DataMigration").C(mongoErrorCollection).Insert(err.Error())
		return 0, err
	}

	//config
	var attIDs []int
	for _, partVehicleConfig := range ve.Config {
		stmt, err = db.Prepare(getConfigAttribute)
		if err != nil {
			return 0, err
		}
		defer stmt.Close()
		var attID int
		err = stmt.QueryRow(partVehicleConfig.Type, partVehicleConfig.Value).Scan(&attID)
		if err != nil {
			// log.Panic("You missed importing a config")

			err = sess.DB("DataMigration").C(mongoErrorCollection).Insert(err.Error())
			return 0, err
		}
		attIDs = append(attIDs, attID)
	}
	//we have attriIDs - insert on VehicleConfig, then VCA's
	//Vehicle Config
	stmt, err = db.Prepare(insertVehicleConfig)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec()
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	ConfigID = int(id)

	//VehicleConfigAttributeID
	stmt, err = db.Prepare(insertVehicleConfigAttribute)
	if err != nil {
		return 0, err
	}
	for _, attID := range attIDs {
		_, err = stmt.Exec(attID, ConfigID)
		if err != nil {
			return 0, err
		}
	}

	//insert vehicle
	stmt, err = db.Prepare(insertVehicle)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err = stmt.Exec(BaseVehicleID, SubModelID, ConfigID, 0, 0)
	if err != nil {
		return 0, err
	}
	id, err = res.LastInsertId()
	vid = int(id)
	return vid, err
}

func (v *Vehicle) Insert() error {
	var err error
	newdb, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer newdb.Close()
	//insert vehicle
	stmt, err := newdb.Prepare(insertVehicle)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(v.BaseVehicle.ID, v.Submodel.ID, v.Configuration.ID, 0, 0)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	v.ID = int(id)
	return err
}

//Imports Functions - import base,make,model,sub,config tables
func ImportBaseVehicles() error {
	var err error
	var bs []BaseVehicle
	var b BaseVehicle
	makemap, err := makeMap()
	if err != nil {
		return err
	}
	modelmap, err := modelMap()
	if err != nil {
		return err
	}
	db, err := sql.Open("mysql", database.OldDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	newdb, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer newdb.Close()

	stmt, err := db.Prepare(allBaseVehicles)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	if err != nil {
		return err
	}
	for res.Next() {
		err = res.Scan(
			&b.AAIAID,
			&b.Year.ID,
			&b.Make.AAIAID,
			&b.Model.AAIAID,
		)
		if err != nil {
			return err
		}
		bs = append(bs, b)
	}

	//check & insert
	bStmt, err := newdb.Prepare(checkBaseVehicles)
	if err != nil {
		return err
	}
	defer bStmt.Close()
	insStmt, err := newdb.Prepare(insertBaseVehicle)
	if err != nil {
		return err
	}
	defer insStmt.Close()
	for _, b := range bs {
		var bid int
		err = bStmt.QueryRow(b.AAIAID).Scan(&bid)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if err == sql.ErrNoRows {
			err = nil
			//insert
			b.Make.ID = makemap[b.Make.AAIAID]
			b.Model.ID = modelmap[b.Model.AAIAID]
			res, err := insStmt.Exec(b.AAIAID, b.Year.ID, b.Make.ID, b.Model.ID)
			if err != nil {
				return err
			}
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			b.ID = int(id)
		}
	}
	return err
}

func makeMap() (map[int]int, error) {
	outmap := make(map[int]int)
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return outmap, err
	}
	defer db.Close()

	stmt, err := db.Prepare(makeIDs)
	if err != nil {
		return outmap, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	if err != nil {
		return outmap, err
	}
	var i, a int
	for res.Next() {
		err = res.Scan(&i, &a)
		if err != nil {
			return outmap, err
		}
		outmap[a] = i
	}
	return outmap, err
}

func modelMap() (map[int]int, error) {
	outmap := make(map[int]int)
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return outmap, err
	}
	defer db.Close()

	stmt, err := db.Prepare(modelIDs)
	if err != nil {
		return outmap, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	if err != nil {
		return outmap, err
	}
	var i, a int
	for res.Next() {
		err = res.Scan(&i, &a)
		if err != nil {
			return outmap, err
		}
		outmap[a] = i
	}
	return outmap, err
}

func ImportMakes() error {
	var err error
	var ms []Make
	var m Make
	db, err := sql.Open("mysql", database.OldDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	newdb, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer newdb.Close()

	stmt, err := db.Prepare(allMakes)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	if err != nil {
		return err
	}
	for res.Next() {
		err = res.Scan(&m.AAIAID, &m.Name)
		if err != nil {
			return err
		}
		ms = append(ms, m)
	}

	sStmt, err := newdb.Prepare(getMakeByName)
	if err != nil {
		return err
	}
	defer stmt.Close()

	insStmt, err := newdb.Prepare(insertMake)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, m := range ms {
		var subID int
		err = sStmt.QueryRow(m.AAIAID, m.Name).Scan(&subID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if err == sql.ErrNoRows {
			err = nil
			res, err := insStmt.Exec(m.AAIAID, m.Name)
			if err != nil {
				return err
			}
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			m.ID = int(id)
		}
	}
	return err
}

func ImportModels() error {
	var err error
	var ms []Model
	var m Model
	db, err := sql.Open("mysql", database.OldDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	newdb, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer newdb.Close()

	stmt, err := db.Prepare(allModels)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	if err != nil {
		return err
	}
	for res.Next() {
		err = res.Scan(&m.AAIAID, &m.Name, &m.VehicleTypeID)
		if err != nil {
			return err
		}
		ms = append(ms, m)
	}

	sStmt, err := newdb.Prepare(getModelByName)
	if err != nil {
		return err
	}
	defer stmt.Close()

	insStmt, err := newdb.Prepare(insertModel)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, m := range ms {
		var subID int
		err = sStmt.QueryRow(m.AAIAID, m.Name).Scan(&subID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if err == sql.ErrNoRows {
			err = nil
			res, err := insStmt.Exec(m.AAIAID, m.Name, m.VehicleTypeID)
			if err != nil {
				return err
			}
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			m.ID = int(id)
		}
	}
	return err
}

func ImportSubmodels() error {
	var err error
	var ss []Submodel
	var s Submodel
	db, err := sql.Open("mysql", database.OldDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	newdb, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer newdb.Close()

	stmt, err := db.Prepare(allSubmodels)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	if err != nil {
		return err
	}

	for res.Next() {
		err = res.Scan(&s.AAIAID, &s.Name)
		if err != nil {

			return err
		}
		ss = append(ss, s)
	}
	sStmt, err := newdb.Prepare(getSubmodelByName)
	if err != nil {
		return err
	}
	defer stmt.Close()

	insStmt, err := newdb.Prepare(insertSubmodel)
	if err != nil {

		return err
	}
	defer stmt.Close()
	for _, s := range ss {
		var subID int
		err = sStmt.QueryRow(s.Name, s.AAIAID).Scan(&subID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if err == sql.ErrNoRows {
			err = nil
			res, err := insStmt.Exec(s.AAIAID, s.Name)
			if err != nil {
				return err
			}
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			s.ID = int(id)
		}
	}
	return err
}

func ImportConfigs() error {
	var err error
	var configs []ConfigAttribute
	var c ConfigAttribute
	db, err := sql.Open("mysql", database.OldDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	newdb, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer newdb.Close()

	stmt, err := db.Prepare(allConfigAttributes)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	if err != nil {
		return err
	}
	//get all configs that need to be in new db
	for res.Next() {
		err = res.Scan(&c.Type.Name, &c.Type.AcesTypeID, &c.Type.Sort, &c.VcdbID, &c.Value)
		if err != nil {
			return err
		}
		configs = append(configs, c)
	}
	//check and insert
	atStmt, err := newdb.Prepare(checkAttributeType)
	if err != nil {
		return err
	}
	defer atStmt.Close()
	aStmt, err := newdb.Prepare(checkAttribute)
	if err != nil {
		return err
	}
	defer aStmt.Close()

	atInsStmt, err := newdb.Prepare(insertAttributeType)
	if err != nil {
		return err
	}
	defer atInsStmt.Close()
	aInsStmt, err := newdb.Prepare(insertAttribute)
	if err != nil {
		return err
	}
	defer aInsStmt.Close()

	for _, c := range configs {

		//check attr type
		err = atStmt.QueryRow(c.Type.Name, c.Type.AcesTypeID, c.Type.Sort).Scan(&c.Type.ID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if err == sql.ErrNoRows {
			err = nil
			//insert type
			res, err := atInsStmt.Exec(c.Type.Name, c.Type.AcesTypeID, c.Type.Sort)
			if err != nil {
				return err
			}
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			c.Type.ID = int(id)
		}

		//check attr
		var attID int
		err = aStmt.QueryRow(c.VcdbID, c.Value, c.Type.Name).Scan(&attID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if err == sql.ErrNoRows {
			err = nil
			res, err := aInsStmt.Exec(c.Type.ID, c.VcdbID, c.Value)
			if err != nil {
				return err
			}
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			c.ID = int(id)
		}
	}
	return err
}
