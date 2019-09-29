package main

import (
	"fmt"

	"github.com/timeforaninja/goDB"
)

func main() {
	newDB()
	newEncDB()
}

func newDB() {
	db := goDB.NewDB()
	err := db.Open("test.db")
	if err != nil {
		fmt.Print("Failed: ")
		fmt.Println(err)
	} else {
		fmt.Print("Success!")
		fmt.Println(db)
	}
}

func newEncDB() {
	db := goDB.NewEncDB("myPW")
	err := db.Open("test.enc.db")
	if err != nil {
		fmt.Print("Failed: ")
		fmt.Println(err)
	} else {
		fmt.Print("Success!")
		fmt.Println(db)
	}
}
