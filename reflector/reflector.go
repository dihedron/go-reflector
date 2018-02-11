// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package reflector

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/dihedron/go-openstack/log"
)

type Observer interface {
	OnValue(context interface{}, path string, name string, kind reflect.Kind, typ reflect.Type, object reflect.Value) bool
	OnPointer(context interface{}, path string, name string, start bool, kind reflect.Kind, typ reflect.Type) bool
	OnArray(context interface{}, path string, name string, start bool, kind reflect.Kind, typ reflect.Type, length int) bool
	OnStruct(context interface{}, path string, name string, start bool, kind reflect.Kind, typ reflect.Type) bool
}

func Visit(context interface{}, path string, name string, object interface{}, observer Observer) {

	switch object := object.(type) {
	case reflect.Value:
		switch object.Kind() {
		case reflect.Invalid:
			observer.OnValue(context, path, "?", object.Kind(), object.Type(), object.Elem())

		case
			reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Uintptr,
			reflect.Float32, reflect.Float64,
			reflect.Complex64, reflect.Complex128,
			reflect.String:

			observer.OnValue(context, path, name, object.Kind(), object.Type(), object)

		case reflect.Chan:
		case reflect.Func:
		case reflect.UnsafePointer:
		case reflect.Slice, reflect.Array:
			observer.OnArray(context, path, name, true, object.Kind(), object.Type(), object.Len())
			for i := 0; i < object.Len(); i++ {
				Visit(context, chain(path, name), fmt.Sprintf("[%d]", i), object.Index(i), observer)
			}
			observer.OnArray(context, path, name, false, object.Kind(), object.Type(), object.Len())
		case reflect.Struct:
			observer.OnStruct(context, path, name, true, reflect.Struct, object.Type())

			for i := 0; i < object.NumField(); i++ {
				Visit(context, chain(path, name), object.Type().Field(i).Name, object.Field(i), observer)
			}
			observer.OnStruct(context, path, name, false, reflect.Struct, object.Type())
		case reflect.Map:
			for _, key := range object.MapKeys() {
				//Visit(context, fmt.Sprintf("%s[%s]", path, format(key)), object.MapIndex(key), observer)
				Visit(context, path, format(key), object.MapIndex(key), observer)
			}
		case reflect.Ptr:
			observer.OnPointer(context, path, name, true, reflect.Ptr, object.Type())
			if object.IsNil() {
				Visit(context, path, fmt.Sprintf("*(%s)", name), nil, observer)
			} else {
				Visit(context, path, fmt.Sprintf("*(%s)", name), object.Elem(), observer)
			}
			observer.OnPointer(context, path, name, false, reflect.Ptr, object.Type())
		case reflect.Interface:
			if object.IsNil() {
				fmt.Printf("%s = nil\n", path)
			} else {
				//fmt.Printf("%s.type = %s\n", path, object.Elem().Type())
				Visit(context, path, "value", object.Elem(), observer)
			}
		default:
			// basic types, channels, funcs
			fmt.Printf("***%s = %s\n", path, format(object))
		}
	default:
		log.Debugf("Starting visit of: %s%s (type: %T):\n", path, name, object)
		// if observer != nil {
		// 	observer.OnValue(context, path, reflect.TypeOf(object).Kind(), reflect.TypeOf(object), nil)
		// }
		Visit(context, path, name, reflect.ValueOf(object), observer)
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
