package category

import (
	"database/sql"
	"github.com/curt-labs/sweetData/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/url"
	"time"
)

var (
	allCats    = `select * from categories where brandID = ?`
	insertCats = `insert into Categories (dateAdded, parentID, catTitle, shortDesc, longDesc, image, isLifestyle, 
		codeId, sort, vehicleSpecific, vehicleRequired, metaTitle, metaDesc, metaKeywords, icon, path, brandID)
		values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	insertCatPart = `insert into CatPart (catID, partID) values (?,?)`
	checkCategory = `select catID from Categories where dateAdded = ? and parentID = ? and catTitle = ? and shortDesc = ? and 
			longDesc = ? and image = ? and isLifestyle = ? and codeId = ? and sort = ? and vehicleSpecific = ? and vehicleRequired = ? and
			metaTitle = ? and metaDesc = ? and metaKeywords = ? and icon = ? and path = ? and brandID = ?`
	checkCatPart = `select catPartID from CatPart where catID = ? and partID = ?`
	getCatParts  = `select partID from CatPart where catID = ?`
)

type Category struct {
	ID              int       `json:"id" xml:"id,attr"`
	DateAdded       time.Time `json:"date_added" xml:"date_added,attr"`
	ParentID        int       `json:"parent_id" xml:"parent_id,attr"`
	Title           string    `json:"title" xml:"title,attr"`
	ShortDesc       string    `json:"short_description" xml:"short_description"`
	LongDesc        string    `json:"long_description" xml:"long_description"`
	Image           *url.URL  `json:"image" xml:"image"`
	IsLifestyle     bool      `json:"lifestyle" xml:"lifestyle,attr"`
	CodeID          int
	Sort            int      `json:"sort" xml:"sort,attr"`
	VehicleSpecific bool     `json:"vehicle_specific" xml:"vehicle_specific,attr"`
	VehicleRequired bool     `json:"vehicle_required" xml:"vehicle_required,attr"`
	MetaTitle       string   `json:"metaTitle,omitempty" xml:"v,omitempty"`
	MetaDescription string   `json:"metaDescription,omitempty" xml:"metaDescription,omitempty"`
	MetaKeywords    string   `json:"metaKeywords,omitempty" xml:"metaKeywords,omitempty"`
	Icon            *url.URL `json:"icon" xml:"icon"`
	Path            string
	BrandID         int `json:"categoryId,omitempty" xml:"categoryId,omitempty"`

	PartIDs []int
}

type Scanner interface {
	Scan(...interface{}) error
}

func GetAndInsertCategories(brandId int) error {
	var err error
	cats, err := GetCategories(brandId)
	if err != nil {
		return err
	}
	err = InsertCategories(cats)
	return err
}

//Get Cats from old DB
func GetCategories(brandId int) ([]Category, error) {
	var err error
	var cs []Category
	db, err := sql.Open("mysql", database.OldDBConnectionString())
	if err != nil {
		return cs, err
	}
	defer db.Close()

	stmt, err := db.Prepare(allCats)
	if err != nil {
		return cs, err
	}
	defer stmt.Close()
	res, err := stmt.Query(brandId)
	if err != nil {
		return cs, err
	}
	for res.Next() {
		log.Print(res)
		c, err := PopulateCategory(res)
		if err != nil {
			return cs, err
		}
		err = c.GetParts()
		if err != nil {
			return cs, err
		}
		cs = append(cs, c)
	}
	return cs, err
}

//Get Category Parts from OldDb
func (c *Category) GetParts() error {
	var err error
	db, err := sql.Open("mysql", database.OldDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCatParts)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Query(c.ID)
	if err != nil {
		return err
	}
	var i int
	for res.Next() {
		err = res.Scan(&i)
		if err != nil {
			return err
		}
		c.PartIDs = append(c.PartIDs, i)
	}
	return nil
}

//Check to see it cat exists in NEW CurtDev
func (c *Category) Check() (int, error) {
	var err error
	var id int
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return 1, err
	}
	defer db.Close()

	stmt, err := db.Prepare(checkCategory)
	if err != nil {
		return 1, err
	}
	defer stmt.Close()
	var image, icon string
	if c.Image != nil {
		image = c.Image.Path
	}
	if c.Icon != nil {
		icon = c.Icon.Path
	}
	err = stmt.QueryRow(
		c.DateAdded,
		c.ParentID,
		c.Title,
		c.ShortDesc,
		c.LongDesc,
		image,
		c.IsLifestyle,
		c.CodeID,
		c.Sort,
		c.VehicleSpecific,
		c.VehicleRequired,
		c.MetaTitle,
		c.MetaDescription,
		c.MetaKeywords,
		icon,
		c.Path,
		c.BrandID,
	).Scan(&id)
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		return 1, err
	}

	return id, err
}

//Insert Cat in New DB
func InsertCategories(cs []Category) error {
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertCats)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, c := range cs {
		id, err := c.Check()
		if err != nil {
			return err
		}
		if id > 0 {
			continue
		}
		var image, icon string
		if c.Image != nil {
			image = c.Image.Path
		}
		if c.Icon != nil {
			icon = c.Icon.Path
		}
		res, err := stmt.Exec(
			c.DateAdded,
			c.ParentID,
			c.Title,
			c.ShortDesc,
			c.LongDesc,
			image,
			c.IsLifestyle,
			c.CodeID,
			c.Sort,
			c.VehicleSpecific,
			c.VehicleRequired,
			c.MetaTitle,
			c.MetaDescription,
			c.MetaKeywords,
			icon,
			c.Path,
			c.BrandID,
		)
		if err != nil {
			return err
		}
		cid, err := res.LastInsertId()
		if err != nil {
			return err
		}
		c.ID = int(cid)
		for _, part := range c.PartIDs {
			cpid, err := c.CheckCatPart(part)
			if err != nil {
				return err
			}
			if cpid > 0 {
				continue
			}
			err = c.InsertCatPart(part)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

//Check existence of CatPart in new DB
func (c *Category) CheckCatPart(partID int) (int, error) {
	var err error
	var i int
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return 1, err
	}
	defer db.Close()

	stmt, err := db.Prepare(checkCatPart)
	if err != nil {
		return 1, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(c.ID, partID).Scan(&i)
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		return 1, err
	}
	return i, err
}

//Insert catPart in new Db
func (c *Category) InsertCatPart(partID int) error {
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertCatPart)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.ID, partID)
	return err
}

func PopulateCategory(row Scanner) (Category, error) {
	var initCat Category
	var added *time.Time
	var life *bool
	var title, short, long, catIcon, catImg, mt, md, mk, path *string
	var codeID, brand *int
	err := row.Scan(
		&initCat.ID,
		&added,
		&initCat.ParentID,
		&title,
		&short,
		&long,
		&catImg,
		&life,
		&codeID,
		&initCat.Sort,
		&initCat.VehicleSpecific,
		&initCat.VehicleRequired,
		&mt,
		&md,
		&mk,
		&catIcon,
		&path,
		&brand,
	)
	if err != nil {
		return initCat, err
	}

	if catImg != nil {
		initCat.Image, _ = url.Parse(*catImg)
	}
	if catIcon != nil {
		initCat.Icon, _ = url.Parse(*catIcon)
	}
	if added != nil {
		initCat.DateAdded = *added
	}
	if title != nil {
		initCat.Title = *title
	}
	if short != nil {
		initCat.ShortDesc = *short
	}
	if long != nil {
		initCat.LongDesc = *long
	}
	if life != nil {
		initCat.IsLifestyle = *life
	}
	if codeID != nil {
		initCat.CodeID = *codeID
	}
	if mt != nil {
		initCat.MetaTitle = *mt
	}
	if md != nil {
		initCat.MetaDescription = *md
	}
	if mt != nil {
		initCat.MetaTitle = *mt
	}
	if path != nil {
		initCat.Path = *path
	}
	if brand != nil {
		initCat.BrandID = *brand
	}

	return initCat, err
}
