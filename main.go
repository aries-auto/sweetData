package main

import (
	// "github.com/curt-labs/sweetData/data/category"
	"flag"
	"github.com/curt-labs/sweetData/data/part"
	"log"
)

//https://github.com/curt-labs/GoAPI/issues/17

func main() {
	flag.Parse()
	//update Parts Before Categories
	err := part.GetAndInsertParts()
	log.Print(err)

	//get cats by brand
	// cats, err := category.GetCategories(3)
	// log.Print(cats, err)

	// //insert cats
	// err = category.InsertCategories(cats)
	// log.Print(err)

}
