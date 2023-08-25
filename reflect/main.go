package main

import (
	"fmt"
	"reflect"
)

type SSS struct {
	A string
	B int
}

func (s SSS) String() string {
	return "sss" + s.A
}

func main() {
	var num float64 = 1.2345
	fmt.Println("type: ", reflect.TypeOf(num))   // type:  float64
	fmt.Println("value: ", reflect.ValueOf(num)) // value:  1.2345

	s := SSS{A: "ggg"}
	// get concrete type
	sType := reflect.TypeOf(s)
	// get value
	sValue := reflect.ValueOf(s)
	fmt.Println("type: ", sType, sType.Kind())    // type:  main.SSS struct
	fmt.Println("value: ", sValue, sValue.Kind()) // value:  sssggg struct
	// parse struct use field
	for i := 0; i < sType.NumField(); i++ {
		field := sType.Field(i)
		value := sValue.Field(i).Interface() // get value as interface{}
		fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
		// A: string = ggg
		// B: int = 0
	}
	// use value call method
	println(sValue.MethodByName("String").Call([]reflect.Value{})[0].Interface().(string)) // sssggg
	for i := 0; i < sType.NumMethod(); i++ {
		method := sType.Method(i)
		fmt.Printf("%s: %v\n", method.Name, method.Type)                           // String: func(main.SSS) string
		println(method.Func.Call([]reflect.Value{sValue})[0].Interface().(string)) // sssggg
	}
	// Elem get value pointer pointed
	sptr := &SSS{A: "ptr"}
	println(reflect.ValueOf(sptr).Elem().Kind().String())                                                        // struct
	println(reflect.ValueOf(sptr).Elem().MethodByName("String").Call([]reflect.Value{})[0].Interface().(string)) // sssptr
	// Elem get value interface contained
	println(reflect.ValueOf(struct{ i interface{} }{i: "imstring"}).Field(0).Kind().String()) // interface
	println(reflect.ValueOf(struct{ i interface{} }{i: "imstring"}).Field(0).Elem().String()) // imstring
}
