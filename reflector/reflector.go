// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package reflector

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/dihedron/go-reflector/log"
)

type Observer interface {
	OnNil(path string, name string, typ reflect.Type) bool
	OnValue(path string, name string, object reflect.Value) bool
	OnPointer(path string, name string, start bool, object reflect.Value) bool
	OnList(path string, name string, start bool, object reflect.Value) bool
	OnStruct(path string, name string, start bool, object reflect.Value) bool
	OnMap(path string, name string, start bool, object reflect.Value) bool
	OnInterface(path string, name string, start bool, object reflect.Value) bool
	OnChannel(path string, name string, object reflect.Value) bool
}

func Visit(path string, name string, object interface{}, observer Observer) {

	switch object := object.(type) {
	case reflect.Value:
		switch object.Kind() {

		case reflect.Invalid:
			observer.OnValue(path, "?", object)

		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:

			observer.OnValue(path, name, object)

		case reflect.Chan:
			observer.OnChannel(path, name, object)

		case reflect.Func:
		case reflect.UnsafePointer:

		case reflect.Slice, reflect.Array:
			observer.OnList(path, name, true, object)
			for i := 0; i < object.Len(); i++ {
				Visit(chain(path, name), fmt.Sprintf("[%d]", i), object.Index(i), observer)
			}
			observer.OnList(path, name, false, object)

		case reflect.Struct:
			observer.OnStruct(path, name, true, object)
			for i := 0; i < object.NumField(); i++ {
				Visit(chain(path, name), object.Type().Field(i).Name, object.Field(i), observer)
			}
			observer.OnStruct(path, name, false, object)

		case reflect.Map:
			observer.OnMap(path, name, true, object)
			for _, key := range object.MapKeys() {
				Visit(path, format(key), object.MapIndex(key), observer)
			}
			observer.OnMap(path, name, false, object)

		case reflect.Ptr:
			observer.OnPointer(path, name, true, object)
			if object.IsNil() {
				observer.OnNil(path, ".value", object.Type())
			} else {
				Visit(path, ".value", object.Elem(), observer)
			}
			observer.OnPointer(path, name, false, object)

		case reflect.Interface:
			observer.OnInterface(path, name, true, object)
			if object.IsNil() {
				observer.OnNil(path, ".value", object.Type())
			} else {
				Visit(path, ".value", object.Elem(), observer)
			}
			observer.OnInterface(path, name, false, object)

		default:
			// basic types, channels, funcs
			fmt.Printf("***%s = %s\n", path, format(object))
		}
	default:
		log.Debugf("Starting visit of: %s%s (type: %T):\n", path, name, object)
		// if observer != nil {
		// 	observer.OnValue(path, reflect.TypeOf(object).Kind(), reflect.TypeOf(object), nil)
		// }
		Visit(path, name, reflect.ValueOf(object), observer)
	}
}

func chain(path, name string) string {
	if path == "" {
		return name
	} else if name == "" {
		return path
	}
	return path + "." + name
}

// format formats a value without inspecting its internal structure.
func format(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'E', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'E', -1, 64)
	case reflect.Complex64:
		result := strconv.FormatFloat(real(v.Complex()), 'E', -1, 32)
		if imag(v.Complex()) > 0 {
			result += "+"
		} else {
			result += "-"
		}
		result += strconv.FormatFloat(math.Abs(imag(v.Complex())), 'E', -1, 32)
		result += "j"
		return result
	case reflect.Complex128:
		result := "("
		result += strconv.FormatFloat(real(v.Complex()), 'E', -1, 64)
		result += ")"
		if imag(v.Complex()) > 0 {
			result += "+i("
		} else {
			result += "-i("
		}
		result += strconv.FormatFloat(math.Abs(imag(v.Complex())), 'E', -1, 64)
		result += ")"
		return result
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" + strconv.FormatUint(uint64(v.Pointer()), 16)
	default:
		// reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}
