package main

import "fmt"


func test(mylist []int){
	mylist = append(mylist,2)
}

func main(){
	mylist := []int{1,3}
	fmt.Println("before", mylist)
	test(&mylist)
	fmt.Println("after", mylist)
	fmt.Println("hi")
}
