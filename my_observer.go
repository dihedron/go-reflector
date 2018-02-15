// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

type MyObserver struct {
	counter *int
	buffer  *bytes.Buffer
}

// String returns the contents of the internal storage buffer.
func (o MyObserver) String() string {
	return o.buffer.String()
}

func (o MyObserver) Reset() {
	*(o.counter) = 0
	o.buffer.Reset()
}

func (o MyObserver) OnNil(path string, name string, tags string, typ reflect.Type) bool {
	fmt.Fprintf(o.buffer, "%s%s: <nil>,\n", tab(*(o.counter)), name)
	//log.Debugf("%-64s", fmt.Sprintf("%s%s: <invalid> \"<invalid>\",", tab(*(o.counter)), name))
	return true
}

func (o MyObserver) OnValue(path string, name string, tags string, object reflect.Value) bool {
	if object.Kind() == reflect.Invalid {
		if len(tags) > 0 {
			fmt.Fprintf(o.buffer, "%s%s: <invalid> \"<invalid>\" `%s`,\n", tab(*(o.counter)), name, tags)
		} else {
			fmt.Fprintf(o.buffer, "%s%s: <invalid> \"<invalid>\",\n", tab(*(o.counter)), name)
		}
		// log.Debugf("%-64s", fmt.Sprintf("%s%s: <invalid> \"<invalid>\",", tab(*(o.counter)), name))
	} else {
		if object.CanInterface() {
			if len(tags) > 0 {
				fmt.Fprintf(o.buffer, "%s%s: %s \"%v\" `%s`,\n", tab(*(o.counter)), name, object.Type(), object.Interface(), tags)
			} else {
				fmt.Fprintf(o.buffer, "%s%s: %s \"%v\",\n", tab(*(o.counter)), name, object.Type(), object.Interface())
			}
		} else {
			if len(tags) > 0 {
				fmt.Fprintf(o.buffer, "%s%s: <unexported> `%s`,\n", tab(*(o.counter)), name, tags)
			} else {
				fmt.Fprintf(o.buffer, "%s%s: <unexported>,\n", tab(*(o.counter)), name)
			}
		}
	}
	return true
}

func (o MyObserver) OnPointer(path string, name string, start bool, tags string, object reflect.Value) bool {
	if start {
		if len(tags) > 0 {
			fmt.Fprintf(o.buffer, "%s%s: %s `%s` -> {\n", tab(*(o.counter)), name, object.Type(), tags)
		} else {
			fmt.Fprintf(o.buffer, "%s%s: %s -> {\n", tab(*(o.counter)), name, object.Type())
		}
		*(o.counter)++
	} else {
		*(o.counter)--
		fmt.Fprintf(o.buffer, "%s},\n", tab(*(o.counter)))
	}
	return true
}

func (o MyObserver) OnList(path string, name string, start bool, tags string, object reflect.Value) bool {
	// has access to object.Len()
	if start {
		if len(tags) > 0 {
			fmt.Fprintf(o.buffer, "%s%s: %s `%s` [\n", tab(*(o.counter)), name, object.Type(), tags)
		} else {
			fmt.Fprintf(o.buffer, "%s%s: %s [\n", tab(*(o.counter)), name, object.Type())
		}
		*(o.counter)++
	} else {
		*(o.counter)--
		fmt.Fprintf(o.buffer, "%s],\n", tab(*(o.counter)))
	}
	return true
}

func (o MyObserver) OnStruct(path string, name string, start bool, tags string, object reflect.Value) bool {
	if start {
		if len(tags) > 0 {
			fmt.Fprintf(o.buffer, "%s%s: %s `%s` {\n", tab(*(o.counter)), name, object.Type(), tags)
		} else {
			fmt.Fprintf(o.buffer, "%s%s: %s {\n", tab(*(o.counter)), name, object.Type())
		}
		*(o.counter)++
	} else {
		*(o.counter)--
		fmt.Fprintf(o.buffer, "%s},\n", tab(*(o.counter)))
	}
	return true
}

// func (o MyObserver) OnStructField(path string, name string, field reflect.StructField, tags string, object reflect.Value) bool {

// 	if object.CanInterface() {
// 		fmt.Fprintf(o.buffer, "%s%s: %s \"%v\" {\n", tab(*(o.counter)), name, object.Type(), object.Interface())
// 		*(o.counter)++
// 		// TODO: loop on headers
// 		*(o.counter)--
// 	} else {
// 		fmt.Fprintf(o.buffer, "%s%s: <unexported> {\n", tab(*(o.counter)), name)
// 	}
// 	fmt.Fprintf(o.buffer, "%s  ... headers ...\n", tab(*(o.counter)))
// 	fmt.Fprintf(o.buffer, "%s},\n", tab(*(o.counter)))
// 	return true
// }

func (o MyObserver) OnMap(path string, name string, start bool, tags string, object reflect.Value) bool {
	if start {
		if len(tags) > 0 {
			fmt.Fprintf(o.buffer, "%s%s: map[%s]%s `%s` {\n", tab(*(o.counter)), name, object.Type().Key().String(), object.Type().Elem().String(), tags)
		} else {
			fmt.Fprintf(o.buffer, "%s%s: map[%s]%s {\n", tab(*(o.counter)), name, object.Type().Key().String(), object.Type().Elem().String())
		}
		*(o.counter)++
	} else {
		*(o.counter)--
		fmt.Fprintf(o.buffer, "%s}\n", tab(*(o.counter)))
	}
	return true
}

func (o MyObserver) OnInterface(path string, name string, start bool, tags string, object reflect.Value) bool {
	if start {
		if len(tags) > 0 {
			fmt.Fprintf(o.buffer, "%s%s: %s `%s` {\n", tab(*(o.counter)), name, object.Type(), tags)
		} else {
			fmt.Fprintf(o.buffer, "%s%s: %s {\n", tab(*(o.counter)), name, object.Type())
		}
		*(o.counter)++
	} else {
		*(o.counter)--
		fmt.Fprintf(o.buffer, "%s},\n", tab(*(o.counter)))
	}
	return true
}

func (o MyObserver) OnChannel(path string, name string, tags string, object reflect.Value) bool {
	if len(tags) > 0 {
		fmt.Fprintf(o.buffer, "%s%s: [%d]%s `%s`,\n", tab(*(o.counter)), name, object.Len(), object.Type(), tags)
	} else {
		fmt.Fprintf(o.buffer, "%s%s: [%d]%s,\n", tab(*(o.counter)), name, object.Len(), object.Type())
	}
	return true
}

func (o MyObserver) OnFunction(path string, name string, tags string, object reflect.Value) bool {
	if len(tags) > 0 {
		fmt.Fprintf(o.buffer, "%s%s: %s `%s`,\n", tab(*(o.counter)), name, object.Type(), tags)
	} else {
		fmt.Fprintf(o.buffer, "%s%s: %s,\n", tab(*(o.counter)), name, object.Type())
	}
	return true
}

func (o MyObserver) OnUnsafePointer(path string, name string, tags string, object reflect.Value) bool {
	if len(tags) > 0 {
		fmt.Fprintf(o.buffer, "%s%s: %s 0x%s `%s`,\n", tab(*(o.counter)), name, object.Type(), strconv.FormatUint(uint64(object.Pointer()), 16), tags)
	} else {
		fmt.Fprintf(o.buffer, "%s%s: %s 0x%s,\n", tab(*(o.counter)), name, object.Type(), strconv.FormatUint(uint64(object.Pointer()), 16))
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
