package xmap

import (
	"fmt"
	"sort"
	"testing"
)

func TestMap(t *testing.T) {
	m := M{}
	m.SetValue("abc", "123")
	m.SetValue("int", 123)
	m.SetValue("x", M{"a": 123})
	m.SetValue("/ary", []interface{}{1, 2, 3})
	m.SetValue("/arm", []interface{}{M{}, M{}, M{}})
	assert := func(v bool) {
		if !v {
			panic("error")
		}
	}
	assert(m.Exist("abc"))
	//
	assert(nil != m.Map("x"))
	assert(nil != m.MapDef(nil, "x"))
	assert(nil == m.MapDef(nil, "not"))
	//
	assert("123" == m.Str("abc"))
	assert("123" == m.StrDef("", "abc"))
	assert("" == m.StrDef("", "not"))
	//
	assert(123 == m.Int("int"))
	assert(123 == m.Int64("int"))
	assert(123 == m.Uint64("int"))
	assert(123 == m.Float64("int"))
	assert(123 == m.Int("x/a"))
	assert(123 == m.Int64("x/a"))
	assert(123 == m.Uint64("x/a"))
	assert(123 == m.Float64("x/a"))
	//
	assert(123 == m.IntDef(0, "int"))
	assert(123 == m.Int64Def(0, "int"))
	assert(123 == m.Uint64Def(0, "int"))
	assert(123 == m.Float64Def(0, "int"))
	//
	assert(0 == m.IntDef(0, "not"))
	assert(0 == m.Int64Def(0, "not"))
	assert(0 == m.Uint64Def(0, "not"))
	assert(0 == m.Float64Def(0, "not"))
	//
	//
	if v, err := m.StrVal("int"); true {
		assert("123" == v && err == nil)
	}
	if v, err := m.IntVal("int"); true {
		assert(123 == v && err == nil)
	}
	if v, err := m.Int64Val("int"); true {
		assert(123 == v && err == nil)
	}
	if v, err := m.Uint64Val("int"); true {
		assert(123 == v && err == nil)
	}
	if v, err := m.Float64Val("int"); true {
		assert(123 == v && err == nil)
	}
	if v, err := m.MapVal("x"); true {
		assert(nil != v && err == nil)
	}
	//
	if v, err := m.ArrayVal("ary"); true {
		assert(3 == len(v) && err == nil)
	}
	if v, err := m.ArrayStrVal("ary"); true {
		assert(3 == len(v) && err == nil)
	}
	if v, err := m.ArrayIntVal("ary"); true {
		assert(3 == len(v) && err == nil)
	}
	if v, err := m.ArrayInt64Val("ary"); true {
		assert(3 == len(v) && err == nil)
	}
	if v, err := m.ArrayUint64Val("ary"); true {
		assert(3 == len(v) && err == nil)
	}
	if v, err := m.ArrayFloat64Val("ary"); true {
		assert(3 == len(v) && err == nil)
	}
	if v, err := m.ArrayMapVal("arm"); true {
		assert(3 == len(v) && err == nil)
	}
	//
	//test remove
	m.SetValue("having", "123")
	assert(m.Value("having") != nil)
	m.Delete("having")
	assert(m.Value("having") == nil)
}

func TestArrayMap(t *testing.T) {
	var m map[string]interface{}
	//test all type
	m = map[string]interface{}{}
	m["arr1"] = []map[string]interface{}{{}, {}}
	m["arr2"] = []interface{}{M{}, M{}}
	m["nil"] = nil
	for key, val := range m {
		_, err := ArrayMapVal(val)
		if err != nil {
			t.Error(key)
			return
		}
	}
	//test error
	m = map[string]interface{}{}
	m["int"] = 1
	m["i1"] = []interface{}{"aaa"}
	m["i2"] = []*testing.T{nil}
	for key, val := range m {
		_, err := ArrayMapVal(val)
		if err == nil {
			t.Error(key)
			return
		}
	}
}

func TestPathValue(t *testing.T) {
	//data
	m1 := map[string]interface{}{
		"s":   "str",
		"i":   int64(16),
		"f":   float64(16),
		"ary": []interface{}{1, 3, 4},
	}
	m2 := map[string]interface{}{
		"a":   "abc",
		"m":   m1,
		"ary": []interface{}{"1", "3", "4"},
	}
	m3 := map[string]interface{}{
		"b":   "abcc",
		"m":   m2,
		"ary": []interface{}{m1, m2},
	}
	m4 := M{
		"test": 1,
		"ms":   []interface{}{m1, m2, m3},
		"m3":   m3,
		"ary2": []int{1, 3, 4},
		"me":   map[string]string{"a": "b"},
	}
	var v interface{}
	var err error
	v, err = m4.ValueVal("/path")
	assertError(t, v, err)
	v, err = m4.ValueVal("/test")
	assertNotError(t, v, err)
	v, err = m4.ValueVal("/ms")
	assertNotError(t, v, err)
	v, err = m4.ValueVal("/m3")
	assertNotError(t, v, err)
	//
	v, err = m4.ValueVal("/m3/b")
	assertNotError(t, v, err)
	v, err = m4.ValueVal("/m3/b2")
	assertError(t, v, err)
	v, err = m4.ValueVal("/m3/ary")
	assertNotError(t, v, err)
	v, err = m4.ValueVal("/ms/1")
	assertNotError(t, v, err)
	v, err = m4.ValueVal("/ms/100")
	assertError(t, v, err)
	v, err = m4.ValueVal("/ms/a")
	assertError(t, v, err)
	v, err = m4.ValueVal("/ary2/100")
	assertError(t, v, err)
	v, err = m4.ValueVal("/ms/@len")
	assertNotError(t, v, err)
	v, err = m4.ValueVal("/ary2/@len")
	assertError(t, v, err)
	v, err = m4.ValueVal("/test/abc")
	assertError(t, v, err)
	v, err = m4.ValueVal("/me/a")
	assertError(t, v, err)
	v, err = m4.ValueVal("/mekkkk/a")
	assertError(t, v, err)
}

func TestSetValue(t *testing.T) {
	var v interface{}
	var err error
	m := M{
		"eary":  []string{},
		"ary":   []interface{}{456},
		"emap":  map[string]string{},
		"map":   map[string]interface{}{},
		"ntype": "kkkk",
	}
	m.SetValue("/abc", M{"a": 1})
	v, err = m.ValueVal("/abc/a")
	assertNotError(t, v, err)
	err = m.SetValue("/ary/0", 123)
	assertNotError(t, nil, err)

	err = m.SetValue("/map/a", 123)
	assertNotError(t, nil, err)
	_, err = m.ValueVal("/map/a")
	assertNotError(t, nil, err)
	//
	//error
	err = m.SetValue("/abcd/abc", 123)
	assertError(t, nil, err)
	err = m.SetValue("/eary/1", 123)
	assertError(t, nil, err)
	err = m.SetValue("/ary/5", 123)
	assertError(t, nil, err)
	err = m.SetValue("/ary/a", 123)
	assertError(t, nil, err)
	err = m.SetValue("/emap/a", 123)
	assertError(t, nil, err)
	err = m.SetValue("/ntype/a", 123)
	assertError(t, nil, err)
	err = m.SetValue("", 123)
	assertError(t, nil, err)
	err = m.setPathValue("", 123)
	assertError(t, nil, err)
	//
	mv, err := m.MapVal("/abc")
	assertNotError(t, mv, err)
	v, err = mv.ValueVal("/a")
	assertNotError(t, v, err)
}

func assertNotError(t *testing.T, v interface{}, err error) {
	if err != nil {
		panic(err)
	}
	fmt.Println(v)
}

func assertError(t *testing.T, v interface{}, err error) {
	fmt.Println(err)
	if err == nil {
		panic("not error")
	}
}

// func TestArray2(t *testing.T) {
// 	fmt.Println([]int{1, 3, 5}[:3])
// }

// func TestNewMap(t *testing.T) {
// 	fmt.Println(NewMap("map.json"))
// 	fmt.Println(NewMap("map.jsn"))
// 	fmt.Println(NewMaps("maps.json"))
// 	fmt.Println(NewMaps("maps.jsn"))
// }

func TestValidFormat(t *testing.T) {
	m := M(map[string]interface{}{
		"ab1": 1,
		"ab2": "xxx",
		"map": map[string]interface{}{
			"x1": 100,
		},
	})
	var v1 int64
	var v2 string
	var v3 int
	var v4 int
	err := m.ValidFormat(`
		ab1,R|I,R:0;
		ab2,R|S,L:0;
		/map/x1,R|I,R:0;
		not|ab1,R|I,R:0;
		`, &v1, &v2, &v3, &v4)
	if v1 != 1 || v2 != "xxx" || v3 != 100 || v4 != 1 {
		t.Error("error")
		return
	}
	fmt.Println(v1, v2, v3)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestMSorter(t *testing.T) {
	newall := func() []M {
		return []M{
			{
				"s": "a",
				"i": 1,
				"f": 1.0,
			},
			{
				"s": "c",
				"i": 3,
				"f": 3.0,
			},
			{
				"s": "b",
				"i": 2,
				"f": 2.0,
			},
		}
	}
	sort.Sort(NewMSorter(newall(), 0, false, "i"))
	sort.Sort(NewMSorter(newall(), 0, true, "i"))
	sort.Sort(NewMSorter(newall(), 1, false, "f"))
	sort.Sort(NewMSorter(newall(), 1, true, "f"))
	sort.Sort(NewMSorter(newall(), 2, false, "s"))
	sort.Sort(NewMSorter(newall(), 2, true, "s"))
}

// func TestAryMapVal(t *testing.T) {
// 	mv, err := NewMap("map2.json")
// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	var ary_1 = mv.AryVal("ary_1")
// 	if ary_1 == nil || len(ary_1) < 1 {
// 		t.Error("error")
// 		return
// 	}
// 	fmt.Println(ary_1)
// 	var ary_2 = mv.AryVal("ary_2")
// 	if ary_2 == nil || len(ary_2) < 1 {
// 		t.Error("error")
// 		return
// 	}
// 	fmt.Println(ary_2)
// 	var ary_3 = mv.AryMapVal("ary_3")
// 	if ary_3 == nil || len(ary_3) < 1 {
// 		t.Error("error")
// 		return
// 	}
// 	fmt.Println(ary_3)
// 	var ary_4 = mv.AryStrVal("ary_1")
// 	if ary_4 == nil || len(ary_4) < 1 {
// 		t.Error("error")
// 		return
// 	}
// 	fmt.Println(ary_4)
// 	var ary_5 = mv.AryStrVal("ary_2")
// 	if ary_5 == nil || len(ary_5) < 1 {
// 		t.Error("error")
// 		return
// 	}
// 	fmt.Println(ary_5)
// 	mv.AryValP("/ary_1")
// 	mv.AryValP("/ary_2")
// 	mv.AryMapValP("/ary_3")
// 	mv.AryStrValP("/ary_2")
// }

// func TestAryMapVal2(t *testing.T) {
// 	var data = Map{}
// 	data.SetValue("abc", []Map{
// 		Map{
// 			"a": 1,
// 		},
// 	})
// 	fmt.Println(reflect.ValueOf([]int{1, 2}).Kind())
// 	var vals = data.AryMapVal("abc")
// 	if len(vals) < 1 {
// 		t.Error("error")
// 		return
// 	}
// 	if vals[0].IntVal("a") != 1 {
// 		t.Error("error")
// 		return
// 	}
// }
