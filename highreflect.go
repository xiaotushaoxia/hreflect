package hreflect

import (
	"fmt"
	"reflect"
)

func GetField(obj any, attr string) (any, error) {
	_, field, err := checkObjAndGetField(obj, attr)
	if err != nil {
		return nil, err
	}
	return field.Interface(), nil
}

func SetField(obj any, attr string, value any) error {
	_, field, err := checkObjAndGetField(obj, attr)
	if err != nil {
		return err
	}
	if !field.CanSet() {
		return fmt.Errorf("cannot set %s field value in obj %T", attr, obj)
	}
	val := reflect.ValueOf(value)

	// 先尝试convert, convert失败就报错
	// 这样会导致字段是int但是set float也会成功，因为float是可以convert到int的。如果偏爱这个行为的话，用这个挺好的
	fuzzError := tryConvertAndSet(field, val)
	if fuzzError == nil {
		return nil
	}
	clearErr := checkKind(obj, attr, field, val)
	if clearErr != nil {
		return clearErr
	}
	return fuzzError
	// 最初的版本 类型别名之间不能Set 这不能接受
	//if field.Type() != val.Type() {
	//	return fmt.Errorf("provided value type %T didn't match obj field type %s", value, field.Type())
	//}
	//field.Set(val)
	//return nil
}

func SetFieldStrict(obj any, attr string, value any) error {
	_, field, err := checkObjAndGetField(obj, attr)
	if err != nil {
		return err
	}
	if !field.CanSet() {
		return fmt.Errorf("cannot set %s field value in obj %T", attr, obj)
	}
	val := reflect.ValueOf(value)
	// 1 先尝试检查Kind, Kind不一样肯定就不能接受Set 这样不会有"隐式"的类型转换
	err = checkKind(obj, attr, field, val)
	if err != nil {
		return err
	}
	return tryConvertAndSet(field, val)
}

func IsStruct(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Struct
}

func IsPointer(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Ptr
}

func checkKind(obj any, attr string, field reflect.Value, val reflect.Value) error {
	vt := val.Type()
	ft := field.Type()
	if ft.Kind() != vt.Kind() {
		return fmt.Errorf("provided value %v(type:%s, kind:%s) didn't match obj(%T) field %s (type:%s, kind:%s)",
			val, vt, vt.Kind(), obj, attr, ft, ft.Kind())
	}
	return nil
}

func tryConvertAndSet(field reflect.Value, val reflect.Value) error {
	if converted, er := tryConvert(field.Type(), val); er == nil {
		return trySet(field, converted)
	}
	return fmt.Errorf("convert %s to %s error", val.Type(), field.Type())
}

func trySet(field reflect.Value, v reflect.Value) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("%v", p)
		}
	}()
	field.Set(v)
	return
}

func tryConvert(t reflect.Type, oldV reflect.Value) (v reflect.Value, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("%v", p)
		}
	}()
	v = oldV.Convert(t)
	return
}

func checkObjAndGetField(obj any, attr string) (ov reflect.Value, field reflect.Value, err error) {
	if obj == nil {
		err = fmt.Errorf("field set/get on nil")
		return
	}
	k := reflect.TypeOf(obj).Kind()
	if k != reflect.Struct && k != reflect.Ptr {
		err = fmt.Errorf("cannot set/get field on a non-struct interface: %T", obj)
		return
	}
	ov = reflect.Indirect(reflect.ValueOf(obj))

	field = ov.FieldByName(attr)
	if !field.IsValid() {
		err = fmt.Errorf("no such field: %s in obj %T", attr, obj)
	}
	return
}
