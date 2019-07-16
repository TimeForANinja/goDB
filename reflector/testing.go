package main

import (
	"fmt"
	"reflect"
	"strings"
)

var SupportetTypes = []string{"uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "_date", "string", "bool", "_float_number_real", "_blob"}

type Table struct {
	master      interface{}
	Name        string
	ColumnCount int
	Columns     []Column
	Rows        []Row
}

type Column struct {
	Name  string
	Type  string
	Flags Flags
}

type Row struct {
	Values []reflect.Value
}

func NewTable(i interface{}) *Table {
	v := reflect.ValueOf(i)
	t := Table{
		master:      reflect.New(reflect.TypeOf(i)).Elem().Interface(),
		Name:        v.Type().Name(),
		ColumnCount: v.NumField(),
		Columns:     make([]Column, v.NumField()),
		Rows:        make([]Row, 0),
	}
	for i := 0; i < t.ColumnCount; i++ {
		structField := v.Type().Field(i)
		t.Columns[i] = Column{
			Name:  structField.Name,
			Type:  structField.Type.Name(),
			Flags: FlagsFromStrings(strings.Split(structField.Tag.Get("goDB"), ",")),
		}
	}
	return &t
}

func (t *Table) AddRow(item interface{}) {
	values := make([]reflect.Value, t.ColumnCount)
	for i := 0; i < t.ColumnCount; i++ {
		values[i] = reflect.ValueOf(item).Field(i)
	}
	t.Rows = append(t.Rows, Row{values})
}

func (t *Table) ReadRow(index int) interface{} {
	row := t.Rows[index]
	response := reflect.New(reflect.TypeOf(t.master))
	a := reflect.ValueOf(response.Interface()).Elem()
	for i := 0; i < t.ColumnCount; i++ {
		col := t.Columns[i]
		a.FieldByName(col.Name).Set(row.Values[i])
		/*
			might also filter by type and use .SetString()|.SetBool()|...
			Row#Values would habe to be of type []interface{}
			and we would have to add another .Interface() in NewRow>for
		*/
	}
	return response.Elem()
}

/*
 * Data for testing purpose
 */
type test struct {
	A int64  `goDB:"AUTO_INCREMENT,UNIQUE"`
	B string `goDB:"NOT_NULL"`
	C uint8
}

func main() {
	x := test{-5, "yey", 9}

	fmt.Print("x: ")
	fmt.Println(x)

	t := NewTable(x)
	t.AddRow(x)
	x.A = -6

	fmt.Print("t: ")
	fmt.Println(t)
	fmt.Print("x(read): ")
	fmt.Println(t.ReadRow(0))
	fmt.Print("x(org): ")
	fmt.Println(x)
}
