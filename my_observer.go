// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"reflect"

	"github.com/dihedron/go-reflector/log"
)

type MyObserver struct{}

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

func (o MyObserver) OnList(context interface{}, path string, name string, start bool, kind reflect.Kind, typ reflect.Type, length int) bool {
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

func (o MyObserver) OnStruct(context interface{}, path string, name string, start bool, object reflect.Value) bool {
	counter := context.(*int)
	if start {
		log.Debugf("%-64s%s", fmt.Sprintf("%s%s: {", tab(*counter), name), fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q",
			object.Kind(), object.Type(), path, name))
		*counter++
	} else {
		*counter--
		log.Debugf("%-64s%s", fmt.Sprintf("%s},", tab(*counter)), fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q",
			object.Kind(), object.Type(), path, name))
	}
	return true
}

func (o MyObserver) OnMap(context interface{}, path string, name string, start bool, object reflect.Value) bool {
	counter := context.(*int)
	if start {
		log.Debugf("%-64s%s", fmt.Sprintf("%s%s: map[%s]%s {", tab(*counter), name, object.Type().Key().String(), object.Type().Elem().String()),
			fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", object.Kind(), object.Type(), path, name))
		*counter++
	} else {
		*counter--
		log.Debugf("%-64s%s", fmt.Sprintf("%s},", tab(*counter)), fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", object.Kind(),
			object.Type(), path, name))
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
