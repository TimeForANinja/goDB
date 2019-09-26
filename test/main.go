package main

import (
	"fmt"

	"github.com/timeforaninja/goDB"
)

func main() {
	db := goDB.NewDB()
	err := db.Open("test.db")
	if err != nil {
		fmt.Print("Failed: ")
		fmt.Println(err)
	} else {
		fmt.Println("Success!")
	}
}
