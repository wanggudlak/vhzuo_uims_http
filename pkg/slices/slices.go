package slices

import (
	"errors"
	"reflect"
)

// 判断是否为slice数据
func IsSlice(arg interface{}) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)

	if val.Kind() == reflect.Slice {
		ok = true
	}

	return
}

// interface{}转为 []interface{}
func CreateAnyTypeSlice(slice interface{}) ([]interface{}, bool) {
	val, ok := IsSlice(slice)

	if !ok {
		return nil, false
	}

	sliceLen := val.Len()

	out := make([]interface{}, sliceLen)

	for i := 0; i < sliceLen; i++ {
		out[i] = val.Index(i).Interface()
	}

	return out, true
}

// 删除切片的某个元素
func RemoveSlice(slice []interface{}, elem interface{}) []interface{} {
	if len(slice) == 0 {
		return slice
	}
	for i, v := range slice {
		if v == elem {
			slice = append(slice[:i], slice[i+1:]...)
			return RemoveSlice(slice, elem)
		}
	}
	return slice
}

// 删除切片为INT类型的某个元素
func RemoveIntSlice(slice []int, elem interface{}) []int {
	if len(slice) == 0 {
		return slice
	}
	for i, v := range slice {
		if v == elem {
			slice = append(slice[:i], slice[i+1:]...)
			return RemoveIntSlice(slice, elem)
		}
	}
	return slice
}

// 判断切片的值是否存在 或 map类型的key是否存在
func IsExistValue(value interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)

	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == value {
				return true, nil
			}
		}
	case reflect.Map:
		// isExist map key
		if targetValue.MapIndex(reflect.ValueOf(value)).IsValid() {
			return true, nil
		}

		// todo isExist map value
	}

	return false, errors.New("not in array")
}
