package utils

import (
	"reflect"

	"errors"
)

// dst should be a pointer to struct, src should be a struct
func CopyStruct(dst interface{}, src interface{}) (err error) {
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		err = errors.New("dst isn't a pointer to struct")
		return
	}
	dstElem := dstValue.Elem()
	if dstElem.Kind() != reflect.Struct {
		err = errors.New("pointer doesn't point to struct")
		return
	}

	srcValue := reflect.ValueOf(src)
	srcType := reflect.TypeOf(src)
	if srcType.Kind() != reflect.Struct {
		err = errors.New("src isn't struct")
		return
	}

	for i := 0; i < srcType.NumField(); i++ {
		sf := srcType.Field(i)
		sv := srcValue.FieldByName(sf.Name)
		// make sure the value which in dst is valid and can set
		if dv := dstElem.FieldByName(sf.Name); dv.IsValid() && dv.CanSet() {
			if sf.Type == dv.Type() {
				dv.Set(sv)
			}
		}
	}
	return
}

func CopyValue(dst interface{}, src interface{}) (err error) {
	if !IsNil(src) {
		dstValue := reflect.ValueOf(dst)
		if dstValue.Kind() != reflect.Ptr {
			err = errors.New("dst isn't a pointer to struct")
			return
		}

		srcValue := reflect.ValueOf(src)
		if srcValue.Kind() != reflect.Ptr {
			err = errors.New("src isn't a pointer to struct")
			return
		}
		srcV := srcValue.Elem()
		dstV := dstValue.Elem()

		dstType := reflect.TypeOf(dst)
		srcType := reflect.TypeOf(src)
		if srcType == dstType && srcV.CanSet() {
			dstV.Set(srcV)
		}
	}

	return
}

func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}
