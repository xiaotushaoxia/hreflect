package hreflect

import (
	"fmt"
	"reflect"
	"testing"
)

type TaskState int
type Task struct {
	State TaskState
}

type TT struct {
	Mp map[int]int
}

func TestSetGetField(t *testing.T) {
	a := &TT{}
	mp := map[int]int{}
	err := SetField(a, "Mp", mp)
	mp[1] = 1
	fmt.Println(err, a)

	field, err := GetField(a, "Mp")

	fmt.Println(err)

	fmt.Println(field)
	printAddr(mp)

	printAddr(field)

	err = SetField(a, "Mp", 1)
	fmt.Println(err)
}

func printAddr(a any) {
	fmt.Printf("%p\n", a)
}

func TestSetFieldMap(t *testing.T) {
	//a := &Task{}
	//err := SetFieldOld(a, "State", 3)
	//fmt.Println(err) // provided value type int didn't match obj field type TaskState

	a := &TT{}
	err := SetField(a, "Mp", map[string]string{})
	fmt.Println(err) // provided value type int didn't match obj field type TaskState
}

func TestSetFieldStrict(t *testing.T) {
	a := &Task{}

	err := SetFieldStrict(a, "State", 4.5)
	if err == nil {
		t.Fatalf("err = SetFieldStrict(a, State, 4.5)  no err !")
	}
	if a.State != 0 {
		t.Fatalf("a.State != 0 after set failed")
	}

	err = SetFieldStrict(a, "State", 4)
	if err != nil {
		t.Fatalf("SetFieldStrict(a, State, 4) error")
	}

	if a.State != 4 {
		t.Fatalf("a.State != 4 after set 4")
	}
}

func TestSetFieldNotStrict(t *testing.T) {
	a := &Task{}
	err := SetField(a, "State", 4.5)
	if err != nil {
		t.Fatalf(" SetField(a, \"State\", 4.5) err !")
	}
	if a.State != 4 {
		t.Fatalf("a.State != 4 after set 4.5")
	}

	err = SetField(a, "State", func() {})
	if err == nil {
		t.Fatalf("err = SetField(a, State, func() {})  no err !")
	}

	err = SetField(a, "State", 3)
	if err != nil {
		t.Fatalf("SetField(a, State, 4) error")
	}

	if a.State != 3 {
		t.Fatalf("a.State != 3 after set 3")
	}
}

func setFieldOld(obj any, attr string, value any) error {
	_, field, err := checkObjAndGetField(obj, attr)
	if err != nil {
		return err
	}
	if !field.CanSet() {
		return fmt.Errorf("cannot set %s field value in obj %T", attr, obj)
	}

	val := reflect.ValueOf(value)
	if field.Type() != val.Type() {
		return fmt.Errorf("provided value type %T didn't match obj field type %s", value, field.Type())
	}

	field.Set(val)
	return nil
}
