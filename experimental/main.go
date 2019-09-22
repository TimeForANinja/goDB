package main

import (
	"fmt"

	"github.com/timeforaninja/goDB/reflector"
)

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

	fmt.Println("-- Testing reflection")
	testReflector()
}

type test struct {
	A int64  `goDB:"AUTO_INCREMENT,UNIQUE"`
	B string `goDB:"NOT_NULL"`
	C uint8
}

func testReflector() {
	x := test{-5, "yey", 9}

	fmt.Print("x(pre): ")
	fmt.Println(x)

	t := reflector.NewTable(x)
	t.AddRow(x)
	x.A = -6

	fmt.Print("x(read): ")
	fmt.Println(t.ReadRow(0))
	fmt.Print("x(org): ")
	fmt.Println(x)
}
