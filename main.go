// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/dihedron/go-reflector/log"
	"github.com/dihedron/go-reflector/reflector"
)

type Struct struct {
	MyInterf interface{}
}

type Embedded struct {
	MyPublic  string
	myPrivate string
	myPointer *string
	MyPointer *string
}

type Embedder struct {
	Embedded
	StructPlain Struct
	StructPtr   *Struct
	Array       [6]int
	Slice       []float32
}

func main() {

	log.SetLevel(log.DBG)
	log.SetStream(os.Stdout)
	log.SetTimeFormat("15:04:05.000")

	log.SetFlags(0)
	log.Debugf("---------------------------------------------------------------------")
	log.Debugf("ciao %s", "mondo")

	log.SetFlags(log.FlagFunctionInfo)
	log.Debugf("---------------------------------------------------------------------")
	log.Debugf("ciao %s", "mondo")

	log.SetFlags(log.FlagFunctionInfo | log.FlagSourceInfo)
	log.Debugf("---------------------------------------------------------------------")
	log.Debugf("ciao %s", "mondo")

	os.Exit(0)

	s := "string pointer"
	o := Embedder{
		Embedded: Embedded{
			MyPublic:  "public",
			myPrivate: "private",
			myPointer: &s,
			MyPointer: &s,
		},
		StructPlain: Struct{
			MyInterf: "string as interface in referenced struct",
		},
		StructPtr: &Struct{
			MyInterf: "string as interface in pointed struct",
		},
		Array: [6]int{0, 1, 2, 3, 4, 5},
		Slice: []float32{0, 1, 2, 3, 4, 5, 6},
	}

	observer := MyObserver{}
	counter := 0
	reflector.Visit(&counter, "", "o", o, observer)

	c := complex(10.0, 4.0)
	reflector.Visit(&counter, "", "c", c, observer)
}
