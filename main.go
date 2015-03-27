package main

import (
	// "github.com/curt-labs/sweetData/data/category"
	"flag"
	"github.com/curt-labs/sweetData/data/testimonial"
	// "github.com/curt-labs/sweetData/data/part"
	// "github.com/curt-labs/sweetData/data/vehicle"
	"log"
)

//https://github.com/curt-labs/GoAPI/issues/17

func main() {
	flag.Parse()

	err := testimonial.ImportTestimonials()
	log.Print(err)
	// err := vehicle.ImportVehicles()
	// log.Print(err)

	// update Parts Before Categories
	// err := part.GetAndInsertParts()
	// log.Print(err)

	//get cats by brand
	// cats, err := category.GetCategories(3)
	// log.Print(cats, err)

	// //insert cats
	// err = category.InsertCategories(cats)
	// log.Print(err)

}
