package reflector

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

func Display(name string, x interface{}) {
	fmt.Printf("Display %s (%T):\n", name, x)
	visit(name, reflect.ValueOf(x))
}

func visit(path string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			visit(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			visit(fieldPath, v.Field(i))
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			visit(fmt.Sprintf("%s[%s]", path, format(key)), v.MapIndex(key))
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			visit(fmt.Sprintf("(*%s)", path), v.Elem())
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			visit(path+".value", v.Elem())
		}
	default:
		// basic types, channels, funcs
		fmt.Printf("%s = %s\n", path, format(v))
	}
}

// formatAtom formats a value without inspecting its internal structure.
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
