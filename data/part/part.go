package part

import (
	"github.com/curt-labs/sweetData/data/video"
	"github.com/curt-labs/sweetData/helpers/database"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
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

	// InstallSheet      *url.URL    `json:"install_sheet" xml:"install_sheet"` //ignore - handled in content/contentBridge
	Attributes []Attribute `json:"attributes" xml:"attributes"`
	// VehicleAttributes []string    `json:"vehicle_atttributes" xml:"vehicle_attributes"` //ignore
	Vehicles []Vehicle `json:"vehicles,omitempty" xml:"vehicles,omitempty"`
	// Content           []Content   `json:"content" xml:"content"` //ignore - handled in content/contentbridge
	Pricing []Price `json:"pricing" xml:"pricing"`
	// Reviews           []Review      `json:"reviews" xml:"reviews"` //ignore
	Images []Image `json:"images" xml:"images"`
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

// type Review struct {
// 	Id          int       `json:"id,omitempty" xml:"id,omitempty"`
// 	PartID      int       `json:"partId,omitempty" xml:"partId,omitempty"`
// 	Rating      int       `json:"rating,omitempty" xml:"rating,omitempty"`
// 	Subject     string    `json:"subject,omitempty" xml:"subject,omitempty"`
// 	ReviewText  string    `json:"reviewText,omitempty" xml:"reviewText,omitempty"`
// 	Name        string    `json:"name,omitempty" xml:"name,omitempty"`
// 	Email       string    `json:"email,omitempty" xml:"email,omitempty"`
// 	Active      bool      `json:"active,omitempty" xml:"active,omitempty"`
// 	Approved    bool      `json:"approved,omitempty" xml:"approved,omitempty"`
// 	CreatedDate time.Time `json:"createdDate,omitempty" xml:"createdDate,omitempty"`
// 	// Customer    customer.Customer `json:"customer,omitempty" xml:"customer,omitempty"`
// }
// type Reviews []Review

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

var pageCounter = 1

var (
	checkPart  = `select partID from Part partID = ?`
	insertPart = `insert into Part partID, status, dateModified, dateAdded, shortDesc,oldPartNumber,priceCode,classID, featured,
	 ACESPartTypeID, replacedBy, brandID values(?,?,?,?,?,?,?,?,?,?,?,?)`
)

func GetParts() error {
	var parts []Part
	ps, err := getPartsByPage(pageCounter)
	if err != nil {
		return err
	}
	parts = append(parts, ps...)
	if len(ps) > 0 {
		pageCounter++
		err = GetParts()
		if err != nil {
			return err
		}
	} else {
		pageCounter = 1
	}

	err = InsertParts(parts)
	return err
}

func getPartsByPage(page int) ([]Part, error) {
	var ps []Part
	res, err := http.Get(database.Api + "&count=50&page=" + strconv.Itoa(page))
	if err != nil {
		return ps, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ps, err
	}

	err = json.Unmarshal(body, &ps)
	log.Print(len(ps))
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
		if id > 0 {
			continue
		}
		if err != nil && err != sql.ErrNoRows {
			return err
		}
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
		//TODO insert attributes, content, pricing, reviews, images, videos, packages, vehicleJoin/installations,

		//videos
		err = video.InsertVideos(p.Videos)
		if err != nil {
			return err
		}
	}

	return err
}
