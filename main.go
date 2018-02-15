// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"
	"unsafe"

	"github.com/dihedron/go-reflector/log"
	"github.com/dihedron/go-reflector/reflector"
)

type Struct struct {
	MyInterf1 interface{}
	MyInterf2 interface{}
	MyString  string
	//InvalidValue reflect.Value
}

type Embedded struct {
	MyPublic     string
	myPrivate    string
	myPointer    *string
	MyPointer    *string
	MyNilPointer *string
	MyChannel    chan int
}

type Embedder struct {
	Embedded
	StructPlain   Struct
	StructPtr     *Struct
	Array         [6]int
	Slice         []float32
	Map           map[string]int
	UnsafePointer unsafe.Pointer
}

func (e Embedder) String() string {
	return "strin"
}

func main() {

	log.SetLevel(log.DBG)
	log.SetStream(os.Stdout)
	log.SetTimeFormat("15:04:05.000")
	log.SetFlags(log.FlagFunctionInfo)

	s := "string pointer"
	o := Embedder{
		Embedded: Embedded{
			MyPublic:  "public",
			myPrivate: "private",
			myPointer: &s,
			MyPointer: &s,
		},
		StructPlain: Struct{
			MyInterf1: "string as interface in referenced struct",
			MyString:  "string in struct",
		},
		StructPtr: &Struct{
			MyInterf1: "string as interface in pointed struct",
			MyString:  "string in pointed struct",
		},
		Array: [6]int{0, 1, 2, 3, 4, 5},
		Slice: []float32{0, 1, 2, 3, 4, 5, 6},
		Map: map[string]int{
			"name":    1,
			"surname": 2,
			"phone":   3,
		},
		UnsafePointer: unsafe.Pointer(&s),
	}

	observer := MyObserver{
		counter: new(int),
		buffer:  new(bytes.Buffer),
	}
	reflector.Visit("", "o", o, observer)

	c := complex(10.0, 4.0)
	reflector.Visit("", "c", c, observer)

	fmt.Printf("buffer is:\n%s\n", observer)
}
