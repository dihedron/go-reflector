// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"

	"github.com/dihedron/go-reflector/log"
)

type MyObserver struct {
	counter *int
	buffer  *bytes.Buffer
}

// String returns the contents of the internal storage buffer.
func (o MyObserver) String() string {
	return o.buffer.String()
}

func (o MyObserver) OnNil(path string, name string, typ reflect.Type) bool {
	fmt.Fprintf(o.buffer, "%s%s: <nil>,\n", tab(*(o.counter)), name)
	//log.Debugf("%-64s", fmt.Sprintf("%s%s: <invalid> \"<invalid>\",", tab(*(o.counter)), name))
	return true
}

func (o MyObserver) OnValue(path string, name string, object reflect.Value) bool {
	if object.Kind() == reflect.Invalid {
		fmt.Fprintf(o.buffer, "%s%s: <invalid> \"<invalid>\",\n", tab(*(o.counter)), name)
		log.Debugf("%-64s", fmt.Sprintf("%s%s: <invalid> \"<invalid>\",", tab(*(o.counter)), name))
	} else if object.CanInterface() {
		fmt.Fprintf(o.buffer, "%s%s: %s \"%v\",\n", tab(*(o.counter)), name, object.Type(), object.Interface())
		log.Debugf("%-64s%s", fmt.Sprintf("%s%s: %s \"%v\",", tab(*(o.counter)), name, object.Type(), object.Interface()),
			fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", object.Kind(), object.Type(), path, name))
	} else {
		fmt.Fprintf(o.buffer, "%s%s: <unexported>,\n", tab(*(o.counter)), name)
		log.Debugf("%-64s%s", fmt.Sprintf("%s%s: <unexported>,", tab(*(o.counter)), name),
			fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", object.Kind(), object.Type(), path, name))
	}
	return true
}

func (o MyObserver) OnPointer(path string, name string, start bool, object reflect.Value) bool {
	if start {
		fmt.Fprintf(o.buffer, "%s%s: %s -> {\n", tab(*(o.counter)), name, object.Type())
		log.Debugf("%-64s%s", fmt.Sprintf("%s%s: %s -> {", tab(*(o.counter)), name, object.Type()),
			fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", object.Kind(), object.Type(), path, name))
		*(o.counter)++
	} else {
		*(o.counter)--
		fmt.Fprintf(o.buffer, "%s},\n", tab(*(o.counter)))
		log.Debugf("%-64s%s", fmt.Sprintf("%s},", tab(*(o.counter))),
			fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", object.Kind(), object.Type(), path, name))
	}
	return true
}

func (o MyObserver) OnList(path string, name string, start bool, object reflect.Value) bool {
	// has access to object.Len()
	if start {
		fmt.Fprintf(o.buffer, "%s%s: %s [\n", tab(*(o.counter)), name, object.Type())
		log.Debugf("%-64s%s", fmt.Sprintf("%s%s: %s [", tab(*(o.counter)), name, object.Type()),
			fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", object.Kind(), object.Type(), path, name))
		*(o.counter)++
	} else {
		*(o.counter)--
		fmt.Fprintf(o.buffer, "%s],", tab(*(o.counter)))
		log.Debugf("%-64s%s", fmt.Sprintf("%s],", tab(*(o.counter))), fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q",
			object.Kind(), object.Type(), path, name))
	}
	return true
}

func (o MyObserver) OnStruct(path string, name string, start bool, object reflect.Value) bool {
	if start {
		fmt.Fprintf(o.buffer, "%s%s: %s {\n", tab(*(o.counter)), name, object.Type())
		log.Debugf("%-64s%s", fmt.Sprintf("%s%s: %s {", tab(*(o.counter)), name, object.Type()), fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q",
			object.Kind(), object.Type(), path, name))
		*(o.counter)++
	} else {
		*(o.counter)--
		fmt.Fprintf(o.buffer, "%s}\n", tab(*(o.counter)))
		log.Debugf("%-64s%s", fmt.Sprintf("%s},", tab(*(o.counter))), fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q",
			object.Kind(), object.Type(), path, name))
	}
	return true
}

func (o MyObserver) OnMap(path string, name string, start bool, object reflect.Value) bool {
	if start {
		fmt.Fprintf(o.buffer, "%s%s: map[%s]%s {\n", tab(*(o.counter)), name, object.Type().Key().String(), object.Type().Elem().String())
		log.Debugf("%-64s%s", fmt.Sprintf("%s%s: map[%s]%s {", tab(*(o.counter)), name, object.Type().Key().String(), object.Type().Elem().String()),
			fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", object.Kind(), object.Type(), path, name))
		*(o.counter)++
	} else {
		*(o.counter)--
		fmt.Fprintf(o.buffer, "%s}\n", tab(*(o.counter)))
		log.Debugf("%-64s%s", fmt.Sprintf("%s},", tab(*(o.counter))), fmt.Sprintf("// kind: %-12s type: %-20s path: %-16q name: %q", object.Kind(),
			object.Type(), path, name))
	}
	return true
}

func (o MyObserver) OnInterface(path string, name string, start bool, object reflect.Value) bool {
	if start {
		fmt.Fprintf(o.buffer, "%s%s: %s {\n", tab(*(o.counter)), name, object.Type())
		*(o.counter)++
	} else {
		*(o.counter)--
		fmt.Fprintf(o.buffer, "%s},\n", tab(*(o.counter)))
	}
	return true
}

func (o MyObserver) OnChannel(path string, name string, object reflect.Value) bool {
	fmt.Fprintf(o.buffer, "%s%s: [%d]%s,\n", tab(*(o.counter)), name, object.Len(), object.Type())
	return true
}

func (o MyObserver) OnFunction(path string, name string, object reflect.Value) bool {
	fmt.Fprintf(o.buffer, "%s%s: %s,\n", tab(*(o.counter)), name, object.Type())
	return true
}

func (o MyObserver) OnUnsafePointer(path string, name string, object reflect.Value) bool {
	fmt.Fprintf(o.buffer, "%s%s: %s 0x%s,\n", tab(*(o.counter)), name, object.Type(), strconv.FormatUint(uint64(object.Pointer()), 16))
	return true
}

func tab(counter int) string {
	s := ""
	for i := 0; i < counter; i++ {
		s += "  "
	}
	return s
}
