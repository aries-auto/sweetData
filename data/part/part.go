package part

import (
	"github.com/curt-labs/sweetData/data/vehicle"
	"github.com/curt-labs/sweetData/data/video"
	"github.com/curt-labs/sweetData/helpers/database"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"encoding/json"
	"io/ioutil"
	// "log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//attributes,content, pricing, reviews, images, categories

type Part struct {
	ID             int       `json:"id" xml:"id,attr"`
	Status         int       `json:"status" xml:"status,attr"`
	DateModified   time.Time `json:"date_modified" xml:"date_modified,attr"`
	DateAdded      time.Time `json:"date_added" xml:"date_added,attr"`
	ShortDesc      string    `json:"short_description" xml:"short_description,attr"`
	OldPartNumber  string    `json:"oldPartNumber,omitempty" xml:"oldPartNumber,omitempty"`
	PriceCode      int       `json:"price_code" xml:"price_code,attr"`
	ClassID        int       `json:"classID,omitempty" xml:"classID,omitempty,attr"`
	Featured       bool      `json:"featured,omitempty" xml:"featured,omitempty"`
	AcesPartTypeID int       `json:"acesPartTypeId,omitempty" xml:"acesPartTypeId,omitempty"`
	ReplacedBy     int       `json:"replacedBy,omitempty" xml:"replacedBy,omitempty"`
	BrandID        int       `json:"brandId" xml:"brandId,attr"`

	InstallSheet *url.URL    `json:"install_sheet" xml:"install_sheet"` //ignore - handled in content/contentBridge
	Attributes   []Attribute `json:"attributes" xml:"attributes"`
	// VehicleAttributes []string    `json:"vehicle_atttributes" xml:"vehicle_attributes"` //ignore
	Vehicles []Vehicle `json:"vehicles,omitempty" xml:"vehicles,omitempty"`
	Content  []Content `json:"content" xml:"content"`
	Pricing  []Price   `json:"pricing" xml:"pricing"`
	Reviews  []Review  `json:"reviews" xml:"reviews"`
	Images   []Image   `json:"images" xml:"images"`
	// Related       []int         `json:"related" xml:"related"`       //ignore
	// Categories    []Category    `json:"categories" xml:"categories"` //ignore - handled in categories
	Videos   []video.Video `json:"videos" xml:"videos"`
	Packages []Package     `json:"packages" xml:"packages"`
	// Customer      CustomerPart  `json:"customer,omitempty" xml:"customer,omitempty"`         //ignore
	Class         Class         `json:"class,omitempty" xml:"class,omitempty"`               //use for classID
	Installations Installations `json:"installation,omitempty" xml:"installation,omitempty"` //VehiclePart
	// Inventory         PartInventory     `json:"inventory,omitempty" xml:"inventory,omitempty"`       //ignore
	// RelatedCount  int     `json:"related_count" xml:"related_count,attr"`   //ignore
	// AverageReview float64 `json:"average_review" xml:"average_review,attr"` //ignore

}
type Attribute struct {
	Key   string `json:"key" xml:"key,attr"`
	Value string `json:"value" xml:",chardata"`
	Sort  int    `json:"sort,omitempty" xml:"sort,omitempty"`
}
type Content struct {
	ID          int         `json:"id,omitempty" xml:"id,omitempty"`
	Text        string      `json:"text,omitempty" xml:"text,omitempty"`
	ContentType ContentType `json:"contentType,omitempty" xml:"contentType,omitempty"`
	UserID      string      `json:"userId,omitempty" xml:"userId,omitempty"`
	Deleted     bool        `json:"deleted,omitempty" xml:"deleted,omitempty"`
}

type ContentType struct {
	Id        int
	Type      string
	AllowHtml bool
}

type Price struct {
	Id           int       `json:"id,omitempty" xml:"id,omitempty"`
	PartId       int       `json:"partId,omitempty" xml:"partId,omitempty"`
	Type         string    `json:"type,omitempty" xml:"type,omitempty"`
	Price        float64   `json:"price" xml:"price"`
	Enforced     bool      `json:"enforced,omitempty", xml:"enforced, omitempty"`
	DateModified time.Time `json:"dateModified,omitempty" xml:"dateModified,omitempty"`
}

type Image struct {
	ID     int      `json:"id,omitempty" xml:"id,omitempty"`
	Size   string   `json:"size,omitempty" xml:"size,omitempty"`
	Sort   string   `json:"sort,omitempty" xml:"sort,omitempty"`
	Height int      `json:"height,omitempty" xml:"height,omitempty"`
	Width  int      `json:"width,omitempty" xml:"width,omitempty"`
	Path   *url.URL `json:"path,omitempty" xml:"path,omitempty"`
	PartID int      `json:"partId,omitempty" xml:"partId,omitempty"`
}
type Category struct {
	ID              int        `json:"id" xml:"id,attr"`
	ParentID        int        `json:"parent_id" xml:"parent_id,attr"`
	Sort            int        `json:"sort" xml:"sort,attr"`
	DateAdded       time.Time  `json:"date_added" xml:"date_added,attr"`
	Title           string     `json:"title" xml:"title,attr"`
	ShortDesc       string     `json:"short_description" xml:"short_description"`
	LongDesc        string     `json:"long_description" xml:"long_description"`
	ColorCode       string     `json:"color_code" xml:"color_code,attr"`
	FontCode        string     `json:"font_code" xml:"font_code,attr"`
	Image           *url.URL   `json:"image" xml:"image"`
	Icon            *url.URL   `json:"icon" xml:"icon"`
	IsLifestyle     bool       `json:"lifestyle" xml:"lifestyle,attr"`
	VehicleSpecific bool       `json:"vehicle_specific" xml:"vehicle_specific,attr"`
	VehicleRequired bool       `json:"vehicle_required" xml:"vehicle_required,attr"`
	Content         []Content  `json:"content,omitempty" xml:"content,omitempty"`
	SubCategories   []Category `json:"sub_categories,omitempty" xml:"sub_categories,omitempty"`
	// ProductListing  *PaginatedProductListing `json:"product_listing,omitempty" xml:"product_listing,omitempty"`
	Filter          interface{} `json:"filter,omitempty" xml:"filter,omitempty"`
	MetaTitle       string      `json:"metaTitle,omitempty" xml:"v,omitempty"`
	MetaDescription string      `json:"metaDescription,omitempty" xml:"metaDescription,omitempty"`
	MetaKeywords    string      `json:"metaKeywords,omitempty" xml:"metaKeywords,omitempty"`
	BrandID         int         `json:"categoryId,omitempty" xml:"categoryId,omitempty"`
}

type Package struct {
	ID                 int         `json:"id,omitempty" xml:"id,omitempty"`
	PartID             int         `json:"partId,omitempty" xml:"partId,omitempty"`
	Height             float64     `json:"height,omitempty" xml:"height,omitempty"`
	Width              float64     `json:"width,omitempty" xml:"width,omitempty"`
	Length             float64     `json:"length,omitempty" xml:"length,omitempty"`
	Weight             float64     `json:"weight,omitempty" xml:"weight,omitempty"`
	DimensionUnit      string      `json:"dimensionUnit,omitempty" xml:"dimensionUnit,omitempty"`
	DimensionUnitLabel string      `json:"dimensionUnitLabel,omitempty" xml:"dimensionUnitLabel,omitempty"`
	WeightUnit         string      `json:"weightUnit,omitempty" xml:"weightUnit,omitempty"`
	WeightUnitLabel    string      `json:"weightUnitLabel,omitempty" xml:"weightUnitLabel,omitempty"`
	PackageUnit        string      `json:"packageUnit,omitempty" xml:"packageUnit,omitempty"`
	PackageUnitLabel   string      `json:"packageUnitLabel,omitempty" xml:"packageUnitLabel,omitempty"`
	Quantity           int         `json:"quantity,omitempty" xml:"quantity,omitempty"`
	PackageType        PackageType `json:"packageType,omitempty" xml:"packageType,omitempty"`
}

type PackageType struct {
	ID   int    `json:"id,omitempty" xml:"id,omitempty"`
	Name string `json:"name,omitempty" xml:"name,omitempty"`
}

type CustomerPart struct {
	Price         float64 `json:"price" xml:"price,attr"`
	CartReference int     `json:"cart_reference" xml:"cart_reference,attr"`
}

type Class struct {
	ID    int    `json:"id,omitempty" xml:"id,omitempty"`
	Name  string `json:"name,omitempty" xml:"name,omitempty"`
	Image string `json:"image,omitempty" xml:"image,omitempty"`
}

type Installation struct { //aka VehiclePart Table
	ID          int     `json:"id,omitempty" xml:"id,omitempty"`
	Vehicle     Vehicle `json:"vehicle,omitempty" xml:"vehicle,omitempty"`
	Part        Part    `json:"part,omitempty" xml:"part,omitempty"`
	Drilling    string  `json:"drilling,omitempty" xml:"v,omitempty"`
	Exposed     string  `json:"exposed,omitempty" xml:"exposed,omitempty"`
	InstallTime int     `json:"installTime,omitempty" xml:"installTime,omitempty"`
}
type Installations []Installation

type Vehicle struct {
	ID            int      `json:"id,omitempty" xml:"id,omitempty"`
	Year          int      `json:"year,omitempty" xml:"year,omitempty"`
	Make          string   `json:"make,omitempty" xml:"make,omitempty"`
	Model         string   `json:"model,omitempty" xml:"model,omitempty"`
	Submodel      string   `json:"submodel,omitempty" xml:"submodel,omitempty"`
	Configuration []Config `json:"configuration,omitempty" xml:"configuration,omitempty"`
}

type Config struct {
	Type  string `json:"type,omitempty" xml:"type,omitempty"`
	Value string `json:"value,omitempty" xml:"value,omitempty"`
}

type Review struct {
	Id          int       `json:"id,omitempty" xml:"id,omitempty"`
	PartID      int       `json:"partId,omitempty" xml:"partId,omitempty"`
	Rating      int       `json:"rating,omitempty" xml:"rating,omitempty"`
	Subject     string    `json:"subject,omitempty" xml:"subject,omitempty"`
	ReviewText  string    `json:"reviewText,omitempty" xml:"reviewText,omitempty"`
	Name        string    `json:"name,omitempty" xml:"name,omitempty"`
	Email       string    `json:"email,omitempty" xml:"email,omitempty"`
	Active      bool      `json:"active,omitempty" xml:"active,omitempty"`
	Approved    bool      `json:"approved,omitempty" xml:"approved,omitempty"`
	CreatedDate time.Time `json:"createdDate,omitempty" xml:"createdDate,omitempty"`
	CustomerID  int       `json:"customerId,omitempty" xml:"customerId,omitempty"`
}

var pageCounter = 1

var (
	checkPart  = `select partID from Part where partID = ?`
	insertPart = `insert into Part (partID, status, dateModified, dateAdded, shortDesc,oldPartNumber,priceCode,classID, featured,
		ACESPartTypeID, replacedBy, brandID) values(?,?,?,?,?,?,?,?,?,?,?,?)`
	checkAttribute  = `select pAttrID from PartAttribute where partID = ? and value = ? and field = ? and sort = ? `
	insertAttribute = `insert into PartAttribute (partID, value, field, sort, canFilter) values (?,?,?,?,?)`
	checkPrice      = `select priceID from Price where partID = ? and priceType = ? and price = ? and enforced = ? `
	insertPrice     = `insert into Price (partID, priceType, price, enforced) values (?,?,?,?)`
	checkReview     = `select reviewID from Review where partID = ? and rating = ? and subject = ? and review _text = ? and 
		name = ? and email = ? and active = ? and approved = ? and createdDate = ? and cust_id = ? `
	insertReview = `insert into Review (partID,rating, subject, review _text,
		name, email, active, approved, createdDate, cust_id) values (?,?,?,?,?,?,?,?,?,?,?)`
	checkImage      = `select imageID from PartImages where sizeID = ? and sort = ? and path = ? and height = ? and width = ? and partID = ?`
	insertImage     = `insert into PartImages (sizeID, sort, path, height, width, partID) values (?,?,?,?,?,?)`
	getImageSizeIds = `select size, sizeID from PartImageSizes`
	checkPackage    = `select ID from PartPackage where partID = ? and height = ? and width = ? and length = ? and weight = ? and dimensionUOM = ? 
		and weightUOM = ? and packageUOM = ? and quantity = ? and typeID = ?`
	insertPackage = ` insert into PartPackage (partID, height, width, length, weight, dimensionUOM, weightUOM, packageUOM, quantity, typeID)
		values (?,?,?,?,?,?,?,?,?,?)`
	getUOMs           = `select code, ID from UnitOfMeasure`
	checkContent      = `select contentID from Content where text = ? and cTypeID = ? and userID = ? and deleted = ?`
	insertContent     = `insert into Content (text, cTypeID, userID, deleted) values (?,?,?,?)`
	checkPartContent  = `select cBridgeID from ContentBridge where partID = ? and contentID = ?`
	insertPartContent = `insert into ContentBridge (partID, contentID) values (?,?)`
	checkContentType  = `select cTypeID from ContentType where type = ? and allowHTML = ?`
	insertContentType = `insert into ContentType (type, allowHTML, isPrivate) values (?,?,0)`
)

func GetAndInsertParts() error {
	var parts []Part
	var err error
	parts, err = getPartsByPage(pageCounter)
	if err != nil {
		return err
	}
	pageCounter++
	err = InsertParts(parts)
	if err != nil {
		return err
	}
	if len(parts) > 0 {
		//recursion
		err = GetAndInsertParts()
		if err != nil {
			return err
		}
	}
	return err
}

func getPartsByPage(page int) ([]Part, error) {
	var ps []Part
	res, err := http.Get(database.Api + "part?key=" + database.ApiKey + "&count=5&page=" + strconv.Itoa(page))
	if err != nil {
		return ps, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ps, err
	}

	err = json.Unmarshal(body, &ps)
	return ps, err
}

//Check for part's existence in New DB
func (p *Part) Check() (int, error) {
	var err error
	var id int
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return id, err
	}
	defer db.Close()

	stmt, err := db.Prepare(checkPart)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(p.ID).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, err
}

//Insert Part into New DB
func InsertParts(parts []Part) error {
	var err error
	//you'll want these maps, friend
	imageSizeMap, err := getImageSizeMap()
	if err != nil {
		return err
	}
	uomMap, err := getUOMmap()
	if err != nil {
		return err
	}

	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertPart)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, p := range parts {
		id, err := p.Check()
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		//insert part
		if id < 1 || err == sql.ErrNoRows {
			_, err = stmt.Exec(
				p.ID,
				p.Status,
				p.DateModified,
				p.DateAdded,
				p.ShortDesc,
				p.OldPartNumber,
				p.PriceCode,
				p.ClassID,
				p.Featured,
				p.AcesPartTypeID,
				p.ReplacedBy,
				p.BrandID,
			)
			if err != nil {
				return err
			}
		}

		//vehicles
		err = vehicle.InsertPartVehicles(p.ID)
		if err != nil {
			return err
		}

		// break //TOOD
		//TODO -insert VEHICLE PART !!!!

		//videos
		err = video.InsertVideos(p.Videos)
		if err != nil {
			return err
		}

		//attributes
		for _, a := range p.Attributes {
			attrID, err := a.Check(p)
			if attrID > 0 {
				continue
			}
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			err = a.Insert(p)
			if err != nil {
				return err
			}
		}

		//prices
		for _, price := range p.Pricing {
			priceID, err := price.Check(p)
			if priceID > 0 {
				continue
			}
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			err = price.Insert(p)
			if err != nil {
				return err
			}
		}

		//reviews
		for _, r := range p.Reviews {
			reviewID, err := r.Check(p)
			if reviewID > 0 {
				continue
			}
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			err = r.Insert(p)
			if err != nil {
				return err
			}
		}

		//images
		for _, image := range p.Images {
			imageSizeID := imageSizeMap[image.Size] //need sizeIDs from map
			imageID, err := image.Check(p, imageSizeID)
			if imageID > 0 {
				continue
			}
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			err = image.Insert(p, imageSizeID)
			if err != nil {
				return err
			}
		}

		//packages
		for _, pack := range p.Packages {

			dimUOMID := uomMap[pack.DimensionUnit]
			weiUOMID := uomMap[pack.WeightUnit]
			packUOMID := uomMap[pack.PackageUnit]

			packID, err := pack.Check(p, dimUOMID, weiUOMID, packUOMID)
			if packID > 0 {
				continue
			}
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			err = pack.Insert(p, dimUOMID, weiUOMID, packUOMID)
			if err != nil {
				return err
			}
		}

		//content
		for _, c := range p.Content {
			//check contentType
			c.ContentType.Id, err = c.ContentType.Check()
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			if err == sql.ErrNoRows || c.ContentType.Id == 0 {
				err = c.ContentType.Insert()
				if err != nil {
					return err
				}
			}

			//then, actual content
			c.ID, err = c.Check(p)
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			if err == sql.ErrNoRows || c.ID == 0 {
				err = c.Insert(p)
				if err != nil {
					return err
				}
			}
			partContentID, err := c.CheckPartContent(p)
			if partContentID > 0 {
				continue
			}
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			err = c.InsertPartContent(p)
			if err != nil {
				return err
			}
		}

		//install sheet
		if p.InstallSheet.Path != "" {
			var c Content
			c.ContentType.Id = 43
			c.Text = p.InstallSheet.Path
			c.ID, err = c.Check(p)
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			if err == sql.ErrNoRows || c.ID == 0 {
				err = c.Insert(p)
				if err != nil {
					return err
				}
			}
			partContentID, err := c.CheckPartContent(p)
			if partContentID > 0 {
				continue
			}
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			err = c.InsertPartContent(p)
			if err != nil {
				return err
			}
		}

	}

	return err
}

func (a *Attribute) Check(p Part) (int, error) {
	var id int
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return id, err
	}
	defer db.Close()

	stmt, err := db.Prepare(checkAttribute)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(p.ID, a.Value, a.Key, a.Sort).Scan(&id)
	return id, err
}

func (a *Attribute) Insert(p Part) error {
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertAttribute)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.ID, a.Value, a.Key, a.Sort, 0)
	return err
}

func (price *Price) Check(p Part) (int, error) {
	var id int
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return id, err
	}
	defer db.Close()

	stmt, err := db.Prepare(checkPrice)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(p.ID, price.Type, price.Price, price.Enforced).Scan(&id)
	return id, err
}

func (price *Price) Insert(p Part) error {
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertPrice)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.ID, price.Type, price.Price, price.Enforced)
	return err
}

func (r *Review) Check(p Part) (int, error) {
	var id int
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return id, err
	}
	defer db.Close()

	stmt, err := db.Prepare(checkReview)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		p.ID,
		r.Rating,
		r.Subject,
		r.ReviewText,
		r.Name,
		r.Name,
		r.Email,
		r.Active,
		r.Approved,
		r.CreatedDate,
		r.CustomerID,
	).Scan(&id)
	return id, err
}

func (r *Review) Insert(p Part) error {
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertReview)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		p.ID,
		r.Rating,
		r.Subject,
		r.ReviewText,
		r.Name,
		r.Name,
		r.Email,
		r.Active,
		r.Approved,
		r.CreatedDate,
		r.CustomerID,
	)
	return err
}

func (i *Image) Check(p Part, imageSizeID int) (int, error) {
	var id int
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return id, err
	}
	defer db.Close()

	stmt, err := db.Prepare(checkImage)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	path := i.Path.Path
	err = stmt.QueryRow(
		imageSizeID,
		i.Sort,
		path,
		i.Height,
		i.Width,
		p.ID,
	).Scan(&id)
	return id, err
}

func (i *Image) Insert(p Part, imageSizeID int) error {
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertImage)
	if err != nil {
		return err
	}
	defer stmt.Close()
	path := i.Path.Path
	_, err = stmt.Exec(
		imageSizeID,
		i.Sort,
		path,
		i.Height,
		i.Width,
		p.ID,
	)
	return err
}

func (pack *Package) Check(p Part, dimUOMID, weiUOMID, packUOMID int) (int, error) {
	var id int
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return id, err
	}
	defer db.Close()

	stmt, err := db.Prepare(checkPackage)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		p.ID,
		pack.Height,
		pack.Width,
		pack.Length,
		pack.Weight,
		dimUOMID,
		weiUOMID,
		packUOMID,
		pack.Quantity,
		pack.PackageType.ID,
	).Scan(&id)
	return id, err
}

func (pack *Package) Insert(p Part, dimUOMID, weiUOMID, packUOMID int) error {
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertPackage)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		p.ID,
		pack.Height,
		pack.Width,
		pack.Length,
		pack.Weight,
		dimUOMID,
		weiUOMID,
		packUOMID,
		pack.Quantity,
		pack.PackageType.ID,
	)
	return err
}

func (c *Content) Check(p Part) (int, error) {
	var id int
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return id, err
	}
	defer db.Close()

	stmt, err := db.Prepare(checkContent)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		c.Text,
		c.ContentType.Id,
		c.UserID,
		c.Deleted,
	).Scan(&id)
	return id, err
}

func (c *Content) Insert(p Part) error {
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertContent)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(
		c.Text,
		c.ContentType.Id,
		c.UserID,
		c.Deleted,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = int(id)
	return err
}

func (c *Content) CheckPartContent(p Part) (int, error) {
	var id int
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return id, err
	}
	defer db.Close()

	stmt, err := db.Prepare(checkPartContent)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		p.ID,
		c.ID,
	).Scan(&id)
	return id, err
}

func (c *Content) InsertPartContent(p Part) error {
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertPartContent)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		p.ID,
		c.ID,
	)
	return err
}

func (ct *ContentType) Check() (int, error) {
	var id int
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return id, err
	}
	defer db.Close()

	stmt, err := db.Prepare(checkContentType)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(ct.Type, ct.AllowHtml).Scan(&id)
	return id, err
}

func (ct *ContentType) Insert() error {
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertContentType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(ct.Type, ct.AllowHtml)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	ct.Id = int(id)
	return err
}

//get image sizes
func getImageSizeMap() (map[string]int, error) {
	var err error
	imageSizeMap := make(map[string]int)
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return imageSizeMap, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getImageSizeIds)
	if err != nil {
		return imageSizeMap, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	if err != nil {
		return imageSizeMap, err
	}
	var i int
	var s string
	for res.Next() {
		err = res.Scan(&s, &i)
		if err != nil {
			return imageSizeMap, err
		}
		imageSizeMap[s] = i
	}
	return imageSizeMap, err
}

//get maps of UOM [code]id
func getUOMmap() (map[string]int, error) {
	var err error
	uomMap := make(map[string]int)
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return uomMap, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getUOMs)
	if err != nil {
		return uomMap, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	if err != nil {
		return uomMap, err
	}
	var i int
	var s string
	for res.Next() {
		err = res.Scan(&s, &i)
		if err != nil {
			return uomMap, err
		}
		uomMap[s] = i
	}
	return uomMap, err
}
