package main

import (
	"fmt"
	"unsafe"
)

type MyInt int

type MyIntAlias = int

// add method to MyInt
func (m MyInt) Add(n MyInt) MyInt {
	return m + n
}

// add method to MyIntAlias
func (m MyIntAlias) Add(n MyIntAlias) MyIntAlias {
	return m + n
}

// add method to int
func (m int) Add(n int) int {
	return m + n
}


func main() {
	var m MyInt = 1
	var n MyInt = 2
	fmt.Println(m.Add(n))

	var m1 MyIntAlias = 1
	// var n1 MyIntAlias = 2
	// fmt.Println(m1.Add(n1))

	var m2 int = 1
	// var n2 int = 2
	// fmt.Println(m2.Add(n2))

	// size of MyInt and MyIntAlias is same
	fmt.Println("Size of MyInt: ", unsafe.Sizeof(m))
	fmt.Println("Size of MyIntAlias: ", unsafe.Sizeof(m1))
	fmt.Println("Size of int: ", unsafe.Sizeof(m2))

}