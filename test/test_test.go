package test

import (
	"fmt"
	"reflect"
	"testing"
)

type Person struct {
	Name   string
	Age    int
	Gender string
	Sex    *int
}

func (p *Person) SetName(name string) {
	p.Name = name
}

func (p *Person) SetAge(age int) {
	p.Age = age
}

func (p *Person) SetGender(gender string) {
	p.Gender = gender
}

func (p *Person) ToString() {
	fmt.Println("{Name:", p.Name, " Age:", p.Age, " Gender:", p.Gender)
}

func (p Person) Error() string {
	return ""
}
func TestMessage(s *testing.T) {
	var e error
	var a interface{}
	p := &Person{}
	e = p
	a = e
	v := reflect.ValueOf(a)
	fmt.Println(v.Elem().Kind().String())

}
