package main

import "fmt"

func main() {
	var s = []byte("3a2g")
	fmt.Println(int(s[0]) - int(byte('0')))
}
