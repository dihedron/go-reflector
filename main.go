package main

import (
	"fmt"

	"github.com/dihedron/go-reflector/reflector"
)

func main() {
	fmt.Println("hallo, world!")

	/*
		req, err := http.NewRequest("GET", "http://example.com", nil)
		if err != nil {
			return
		}
		reflector.Display("req", req)
	*/

	c := complex(10.0, 4.0)
	reflector.Display("c", c)
}
