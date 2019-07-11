package main

import "fmt"

func main() {
	_, err := NewDBEnc("./test.db.enc", "asdfMovie")
	fmt.Println("-- newDBEnc returned")
	fmt.Print("err: ")
	fmt.Println(err)
	_, err2 := OpenEnc("./test.db.enc", "asdfMovie")
	fmt.Println("-- OpenEnc returned")
	fmt.Print("err: ")
	fmt.Println(err2)
	_, err3 := NewDB("./test.db")
	fmt.Println("-- newDB returned")
	fmt.Print("err: ")
	fmt.Println(err3)
	_, err4 := Open("./test.db")
	fmt.Println("-- Open returned")
	fmt.Print("err: ")
	fmt.Println(err4)
}
