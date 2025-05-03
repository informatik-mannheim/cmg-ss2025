package main

import "fmt"




func myArray(){

	b := []int(1,2,3,4)
	b = append(b,5)


	var arr [10] int

	for i := range arr{
	fmt.Println(i)
	}
}