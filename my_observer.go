// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"reflect"

	"github.com/dihedron/go-reflector/log"
)

type MyObserver struct {
	counter *int
}

func (o MyObserver) OnValue(path string, name string, object reflect.Value) bool {
	if object.Kind() == reflect.Invalid {
		log.Debugf("%-64s", fmt.Sprintf("%s%s: <invalid> \"<invalid>\",", tab(*(o.counter)), name))
		// fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", object.Kind(), object.Type(), path, name))
	} else if object.CanInterface() {
		log.Debugf("%-64s%s", fmt.Sprintf("%s%s: %s \"%v\",", tab(*(o.counter)), name, object.Type(), object.Interface()),
			fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", object.Kind(), object.Type(), path, name))
	} else {
		log.Debugf("%-64s%s", fmt.Sprintf("%s%s: <unexported>,", tab(*(o.counter)), name),
			fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", object.Kind(), object.Type(), path, name))
	}
	return true
}

func (o MyObserver) OnPointer(path string, name string, start bool, object reflect.Value) bool {
	if start {
		log.Debugf("%-64s%s", fmt.Sprintf("%s%s: %s -> {", tab(*(o.counter)), name, object.Type()),
			fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", object.Kind(), object.Type(), path, name))
		*(o.counter)++
	} else {
		*(o.counter)--
		log.Debugf("%-64s%s", fmt.Sprintf("%s},", tab(*(o.counter))),
			fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", object.Kind(), object.Type(), path, name))
	}
	return true
}

func (o MyObserver) OnList(path string, name string, start bool, object reflect.Value) bool {
	// has access to object.Len()
	if start {
		log.Debugf("%-64s%s", fmt.Sprintf("%s%s: %s [", tab(*(o.counter)), name, object.Type()),
			fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", object.Kind(), object.Type(), path, name))
		*(o.counter)++
	} else {
		*(o.counter)--
		log.Debugf("%-64s%s", fmt.Sprintf("%s],", tab(*(o.counter))), fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q",
			object.Kind(), object.Type(), path, name))
	}
	return true
}

func (o MyObserver) OnStruct(path string, name string, start bool, object reflect.Value) bool {
	if start {
		log.Debugf("%-64s%s", fmt.Sprintf("%s%s: %s {", tab(*(o.counter)), name, object.Type()), fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q",
			object.Kind(), object.Type(), path, name))
		*(o.counter)++
	} else {
		*(o.counter)--
		log.Debugf("%-64s%s", fmt.Sprintf("%s},", tab(*(o.counter))), fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q",
			object.Kind(), object.Type(), path, name))
	}
	return true
}

func (o MyObserver) OnMap(path string, name string, start bool, object reflect.Value) bool {
	if start {
		log.Debugf("%-64s%s", fmt.Sprintf("%s%s: map[%s]%s {", tab(*(o.counter)), name, object.Type().Key().String(), object.Type().Elem().String()),
			fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", object.Kind(), object.Type(), path, name))
		*(o.counter)++
	} else {
		*(o.counter)--
		log.Debugf("%-64s%s", fmt.Sprintf("%s},", tab(*(o.counter))), fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", object.Kind(),
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
