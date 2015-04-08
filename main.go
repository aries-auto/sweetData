package main

import (
	"flag"
	// "github.com/curt-labs/sweetData/data/category"
	// "github.com/curt-labs/sweetData/data/testimonial"
	"github.com/curt-labs/sweetData/data/part"
	// "github.com/curt-labs/sweetData/data/vehicle"
	"log"
)

//https://github.com/curt-labs/GoAPI/issues/17

func main() {
	flag.Parse()

	// err := testimonial.ImportTestimonials()
	// log.Print(err)
	// err := vehicle.ImportVehicles()
	// log.Print(err)

	// update Parts Before Categories
	err := part.GetAndInsertParts()
	log.Print(err)

	//categories
	// err = category.GetAndInsertCategories(3)
	// log.Print(err)

}
