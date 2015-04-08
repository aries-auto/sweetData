package testimonial

import (
	"github.com/curt-labs/sweetData/helpers/database"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"time"
)

type Testimonial struct {
	ID          int
	Rating      float64
	Title       string
	Testimonial string
	DateAdded   *time.Time
	Approved    bool
	Active      bool
	FirstName   string
	LastName    string
	Location    string
	BrandID     int
}

var (
	getAll = `select testimonialID, rating, title, testimonial, dateAdded, approved, active, first_name, last_name, location, brandID from Testimonial where brandID = 3`
	check  = `select testimonialID from Testimonial where rating = ? and title = ? and  testimonial = ? and  
		dateAdded = ? and  approved = ? and  active = ? and  first_name = ? and  last_name = ? and  location = ? and brandID = ?`
	insert = `insert into Testimonial (rating, title, testimonial, dateAdded, approved, active, first_name, last_name, location, brandID) values (?,?,?,?,?,?,?,?,?,?)`
)

func ImportTestimonials() error {
	ts, err := GetAllTestimonials()
	if err != nil {
		return err
	}
	for _, t := range ts {
		err = t.Check()
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if err == sql.ErrNoRows {
			err = t.Insert()
			if err != nil {
				return err
			}
		}
	}
	return err
}

func GetAllTestimonials() ([]Testimonial, error) {
	var err error
	var ts []Testimonial
	db, err := sql.Open("mysql", database.OldDBConnectionString())
	if err != nil {
		return ts, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getAll)
	if err != nil {
		return ts, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	if err != nil {
		return ts, err
	}
	var t Testimonial
	var title, test, first, last, location *string
	for res.Next() {
		err = res.Scan(
			&t.ID,
			&t.Rating,
			&title,
			&test,
			&t.DateAdded,
			&t.Approved,
			&t.Active,
			&first,
			&last,
			&location,
			&t.BrandID,
		)
		if err != nil {
			return ts, err
		}
		if title != nil {
			t.Title = *title
		}
		if test != nil {
			t.Testimonial = *test
		}
		if first != nil {
			t.FirstName = *first
		}
		if last != nil {
			t.LastName = *last
		}
		if location != nil {
			t.Location = *location
		}
		ts = append(ts, t)
	}
	return ts, err
}

func (t *Testimonial) Check() error {
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(check)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		t.Rating,
		t.Title,
		t.Testimonial,
		t.DateAdded,
		t.Approved,
		t.Active,
		t.FirstName,
		t.LastName,
		t.Location,
		t.BrandID,
	).Scan(&t.ID)
	return err
}

func (t *Testimonial) Insert() error {
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insert)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(
		t.Rating,
		t.Title,
		t.Testimonial,
		t.DateAdded,
		t.Approved,
		t.Active,
		t.FirstName,
		t.LastName,
		t.Location,
		t.BrandID,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = int(id)
	return err

}
