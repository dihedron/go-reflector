package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/dihedron/go-openstack/log"
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

type MyObserver struct{}

func (o MyObserver) OnNil(context interface{}, path string, name string) bool {
	counter := context.(*int)
	log.Debugf("%s%s: null,", tab(*counter), name)
	log.Debugf("%-64s%s", fmt.Sprintf("%s%s: null,", tab(*counter), name), fmt.Sprintf("// kind: nil          type: nil path: %-16q name: %q", path, name))
	return true
}

func (o MyObserver) OnValue(context interface{}, path string, name string, kind reflect.Kind, typ reflect.Type, object reflect.Value) bool {
	counter := context.(*int)
	if object.CanInterface() {
		log.Debugf("%-64s%s", fmt.Sprintf("%s%s: \"%v\",", tab(*counter), name, object.Interface()), fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", kind, typ, path, name))
	} else {
		log.Debugf("%-64s%s", fmt.Sprintf("%s%s: <unexported>,", tab(*counter), name), fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", kind, typ, path, name))
	}
	return true
}

func (o MyObserver) OnPointer(context interface{}, path string, name string, start bool, kind reflect.Kind, typ reflect.Type) bool {
	counter := context.(*int)
	if start {
		log.Debugf("%-64s%s", fmt.Sprintf("%s%s: -> {", tab(*counter), name), fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", kind, typ, path, name))
		*counter++
	} else {
		*counter--
		log.Debugf("%-64s%s", fmt.Sprintf("%s},", tab(*counter)), fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", kind, typ, path, name))
	}
	return true
}

func (o MyObserver) OnArray(context interface{}, path string, name string, start bool, kind reflect.Kind, typ reflect.Type, length int) bool {
	counter := context.(*int)
	if start {
		log.Debugf("%-64s%s", fmt.Sprintf("%s%s: (size %d) [", tab(*counter), name, length), fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", kind, typ, path, name))
		*counter++
	} else {
		*counter--
		log.Debugf("%-64s%s", fmt.Sprintf("%s],", tab(*counter)), fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", kind, typ, path, name))
	}
	return true
}

func (o MyObserver) OnStruct(context interface{}, path string, name string, start bool, kind reflect.Kind, typ reflect.Type) bool {
	counter := context.(*int)
	if start {
		log.Debugf("%-64s%s", fmt.Sprintf("%s%s: {", tab(*counter), name), fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", kind, typ, path, name))
		*counter++
	} else {
		*counter--
		log.Debugf("%-64s%s", fmt.Sprintf("%s},", tab(*counter)), fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", kind, typ, path, name))
	}
	return true
}

func tab(counter int) string {
	s := ""
	for i := 0; i < counter; i++ {
		s += "  "
	}
	return s
}

func main() {

	log.SetLevel(log.DBG)
	log.SetStream(os.Stdout)
	log.SetTimeFormat("15:04:05.000")

	log.Debugf("---------------------------------------------------------------------")

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
