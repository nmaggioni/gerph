package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const DEBUG bool = true

var DB *bolt.DB

func main() {
	fmt.Println("Setting up DB...")
	db, err := SetupDB("gerph.db")
	DB = db
	if err != nil {
		log.Fatal(err)
	}
	defer CloseDB()

	fmt.Println("Now listening...")
	Listen("3000")
}
