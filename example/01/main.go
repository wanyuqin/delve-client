package main

import "fmt"

func main() {
	a := 1
	s := "a string"
	arr := []string{"a", "b", "c"}

	fmt.Println("hello world")
	fmt.Println("hello world")
	fmt.Printf("var a is %d\n", a)
	fmt.Printf("var s id %s \n", s)
	fmt.Printf("var arr is %v\n", arr)
	for i := 0; i < 100; i++ {
		fmt.Printf("Hello index %d\n", i)
	}
}
