package utils

import (
	"fmt"
	"reflect"

	"github.com/labstack/gommon/log"
)

// Error prints an error and sends it to elastic
func Error(s string, i interface{}) {
	log.Error(s, i)
	// go sendToElastic("errors", "error", i)
}

// Pretty prints types and values
func Pretty(t interface{}) {
	fmt.Println("*-*-*-*utils.Pretty:*-*-*-*-*")

	// Dont remove links !
	// https://github.com/shiena/ansicolor
	//\x1b[0m	All attributes off(color at startup)

	colorOff := "\x1b[0m"
	color1 := "\x1b[31;1m"
	color2 := "\x1b[36;1m"

	// https://jimmyfrasche.github.io/go-reflection-codex/#indirect

	// ValueOf returns a new Value initialized to the concrete value
	// stored in the interface i. ValueOf(nil) returns the zero Value.
	s := reflect.ValueOf(t)

	// Type is the representation of a Go type.
	typeOfT := s.Type()
	fmt.Println( s.Kind())

	if s.Kind() == reflect.Slice {
		s := reflect.ValueOf(t)
		fmt.Println("reflect.Slice")
		for i := 0; i < s.Len(); i++ {
			fmt.Println(s.Index(i))
		}
	}
	if s.Kind() == reflect.Struct {
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)


			if f.Kind() == reflect.Ptr && f.Elem().Kind() == reflect.Struct {

				//switch f.Elem().Type().String() {
				//
				//case "nullable.NullString":
				//	nlField := f.Elem().Interface().(nl.NullString)
				//	fmt.Println(">>>>>>>", nlField.String)
				//case "nullable.NullTime":
				//	nlField := f.Elem().Interface().(nl.NullTime)
				//	fmt.Println(">>>>>>>", nlField.Time)
				//case "nullable.ZeroOmitInt64":
				//	nlField := f.Elem().Interface().(nl.ZeroOmitInt64)
				//	fmt.Println(">>>>>>>", nlField.Int64)
				//case "nullable.NullInt64":
				//	nlField := f.Elem().Interface().(nl.NullInt64)
				//	fmt.Println(">>>>>>>", nlField.Int64)
				//}

			}

			// https://golang.org/pkg/fmt/
			// %d	base 10
			// %s	the uninterpreted bytes of the string or slice
			// %v	the value in a default format

			if f.Kind() == reflect.Struct {
				fmt.Printf("%d: struct -> "+color1+"%s %s"+colorOff+"\n", i, typeOfT.Field(i).Name, f.Type())
				Pretty(f.Interface())
			} else {
				fmt.Printf("..%d: "+color1+"%s %s"+colorOff+" ->"+color2+"%v\n"+colorOff, i, typeOfT.Field(i).Name, f.Type(), f.Interface())
			}
			if f.Type().String() == "*int64" {
				fmt.Printf(".. %s %s\n", f.Type().String(), f.Elem())
				// fmt.Printf(f.Type().String(), f.Int())
			}
		}
	} else {
		fmt.Printf(color1+"%s"+colorOff+" -> "+color2+"%v\n"+colorOff, typeOfT.Name(), t)
	}
	fmt.Println("*-*-*-*-*-*-*-*-*-*-*-*-*-*-*")
}
