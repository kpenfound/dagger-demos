package main

import "fmt"

func main() {
	say := greeting()
	fmt.Println(say)
}

func greeting() string {
	return "Hello"
}
