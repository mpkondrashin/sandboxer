package main

import "fmt"

type Struct struct {
	s string
}

type modifier func(*Struct)

type X struct {
	m modifier
}

func (x *X) AddModifier(m modifier) {
	if x.m == nil {
		x.m = m
		return
	}
	original := x.m
	x.m = func(s *Struct) {
		original(s)
		m(s)
	}
}

func (x *X) Run() {
	s := &Struct{
		s: "ABC",
	}
	if x.m != nil {
		x.m(s)
	}
	fmt.Println(s)
}

func main() {
	x := &X{}
	fmt.Println("1")
	x.Run()

	x.AddModifier(func(s *Struct) {
		s.s += "A"
	})
	fmt.Println("2")
	x.Run()

	x.AddModifier(func(s *Struct) {
		s.s += "B"
	})
	x.Run()
}
