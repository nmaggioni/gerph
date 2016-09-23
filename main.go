package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/boltdb/bolt"
)

const DEBUG bool = true

var (
	DB *bolt.DB
	portPtr *string
	dbPathPtr *string
)

func catchSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		<- c
		fmt.Println("\n* Interrupt received, stopping...")
		fmt.Print("[..] Closing DB\r")
		CloseDB()
		fmt.Println("[OK] Closing DB")

		fmt.Println()
		os.Exit(0)
	}()
}

func main() {
	portPtr = flag.String("port", "3000", "The port to listen on.")
	dbPathPtr = flag.String("dbpath", "./gerph.db", "The path to the file to save the keystore in.")
	flag.Parse()

	fmt.Print("[..] Setting up DB in \"" + *dbPathPtr + "\"\r")
	db, err := SetupDB(*dbPathPtr)
	DB = db
	if err != nil {
		log.Fatal(err)
	}
	defer CloseDB()
	fmt.Println("[OK] Setting up DB in \"" + *dbPathPtr + "\"")

	catchSignals()

	fmt.Print("[..] Listening on port " + *portPtr + "\r")
	go func() {
		Listen(*portPtr)
	}()
	fmt.Println("[OK] Listening on port " + *portPtr)
	fmt.Println("\n* Ready!")
	select{}
}
