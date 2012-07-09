// Package jsonp helps to handle multiply nested JSON objects.
// With the help of this pkg handling data like that feels more like in a dynamic language (eg. JavaScript), and even better since you get an exception
// in JS if you want to access a member of a nonobject/nonarray.
// + Some utility functions to alleviate the pain of migration from a dynamic language.
package jsonp

import (
	"reflect"
	"strconv"
	"encoding/json"
	"fmt"
)

func explode(str string) []string {
	l := make([]string, 1)
	i := 0
	for _, v := range str {
		s := string(v)
		if s == "." || s == "[" || s == "]" {
			l = append(l, "")
			i++
		} else {
			l[i] += s
		}
	}
	return l
}

// Check if an interface is a map[string]interface{} and contains the element with the key str.
func getFromMap(ob interface{}, str string) (interface{}, bool) { // ok
	if m, k := ob.(map[string]interface{}); k {
		if elem, okay := m[str]; okay {
			return elem, true
		} else {
			return nil, false
		}
	}
	return nil, false
}

// Core of the package
// JSON: {"a":{"b":{"c":{"d":{"e"}}}}}
// val, ok := jsonp.Get(object_name, "a.b.c.d")
func Get(ob interface{}, str string) (interface{}, bool) {
	l := explode(str)
	for _, v := range l {
		if v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				if slice, ok := ob.([]interface{}); ok && len(slice) < n {
					ob = slice[n]
				} else if elem, okay := getFromMap(ob, v); okay {
					ob = elem
				} else {
					return nil, false
				}
			} else if elem, okay := getFromMap(ob, v); okay { // objektumot keresünk
				ob = elem
			} else {
				return nil, false
			}
		}
	}
	return ob, true
}

// Get map. You spare a type assertion with this.
func GetM(ob interface{}, str string) (map[string]interface{}, bool) {
	o, ok := Get(ob, str)
	if ok {
		val, ismap := o.(map[string]interface{})
		if ismap {
			return val, true
		} else {
			return nil, false
		}
	}
	return nil, false
}

// Get (interface{}) slice. You spare a type assertion with this.
func GetS(ob interface{}, str string) ([]interface{}, bool) {
	o, ok := Get(ob, str)
	if ok {
		val, is_slice := o.([]interface{})
		if is_slice {
			return val, true
		} else {
			return nil, false
		}
	}
	return nil, false
}

// Get integer. You spare a type assertion with this.
func GetI(ob interface{}, str string) (int, bool) {
	o, ok := Get(ob, str)
	if ok {
		val, isint := o.(int)
		if isint {
			return val, true
		} else {
			return 0, false
		}
	}
	return 0, false
}

// If ob and str identifies a map[string]interface{} or []interface{}, then this function iterates trough it, and compares every element to val. Returns true of finds equality.
func HasVal(ob interface{}, str string, val interface{}) bool {
	o, ok := Get(ob, str)
	if ok {
		if m, k := o.(map[string]interface{}); k {
			for _, v := range m {
				if reflect.DeepEqual(v, val) {
					return true
				}
			}
		} else if sl, okay := o.([]interface{}); okay {
			for _, v := range sl {
				if reflect.DeepEqual(v, val) {
					return true
				}
			}
		}
	}
	return false
}

// Convenient way of encoding to JSON.
func Encode(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	return string(b), err
}

// Convenient way of decoding from JSON.
func Decode(str string) (interface{}, error) {
	var v interface{}
	err := json.Unmarshal([]byte(str), &v)
	return v, err
}

func DecodeM(str string) (map[string]interface{}, error) {
	v, err := Decode(str)
	if err != nil {
		return nil, err
	}
	if ma, ok := v.(map[string]interface{}); ok {
		return ma, err
	}
	return nil, fmt.Errorf("JSON is not a map.")
}

// Some random type converting method...

// This method converts an []interface{} to []string
func ToStringSlice(ob interface{}) []string {
	n := make([]string, 0)
	sl, ok := ob.([]interface{})
	if ok {
		for _, v := range sl {
			if str, ok := v.(string); ok {
				n = append(n, str)
			}
		}
	}
	return n
}
