package attrvalid

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/wfunc/util/attrscan"
	"github.com/wfunc/util/converter"
)

type ZeroChecker interface {
	IsZero() bool
}

// EnumValider is interface to enum valid
type EnumValider interface {
	EnumValid(v interface{}) error
}

// Validable is interface to define object can be valid by valid temple
type Validable interface {
	ValidFormat(format string, args ...interface{}) error
}

// M is an map[string]interface{} which can be valid by valid temple
type M map[string]interface{}

// RawMap will return raw map
func (m M) RawMap() map[string]interface{} {
	return m
}

// Get will return value by key
func (m M) Get(key string) (v interface{}, err error) {
	return m[key], nil
}

// ValidFormat will valid args by format temple
func (m M) ValidFormat(format string, args ...interface{}) error {
	return ValidAttrFormat(format, m, true, args...)
}

// MS is an map[string]string which can be valid by valid temple
type MS map[string]string

// Get will return value by key
func (m MS) Get(key string) (v interface{}, err error) {
	return m[key], nil
}

// ValidFormat will valid args by format temple
func (m MS) ValidFormat(format string, args ...interface{}) error {
	return ValidAttrFormat(format, m, true, args...)
}

// Values is an url.Values which can be valid by valid temple
type Values url.Values

// Get will return value by key
func (v Values) Get(key string) (val interface{}, err error) {
	return url.Values(v).Get(key), nil
}

// ValidFormat will valid args by format temple
func (v Values) ValidFormat(format string, args ...interface{}) error {
	return ValidAttrFormat(format, v, true, args...)
}

// QueryValidFormat will valid args by http request query
func QueryValidFormat(req *http.Request, format string, args ...interface{}) error {
	return Values(req.URL.Query()).ValidFormat(format, args...)
}

// FormValidFormat will valid args by http request form
func FormValidFormat(req *http.Request, format string, args ...interface{}) error {
	return Values(req.Form).ValidFormat(format, args...)
}

// PostFormValidFormat will valid args by http request post form
func PostFormValidFormat(req *http.Request, format string, args ...interface{}) error {
	return Values(req.PostForm).ValidFormat(format, args...)
}

// RequestValidFormat will valid args by http request query/form/postform
func RequestValidFormat(req *http.Request, format string, args ...interface{}) error {
	query := req.URL.Query()
	getter := func(key string) (v interface{}, err error) {
		val := query.Get(key)
		if len(val) < 1 {
			val = req.Form.Get(key)
		}
		if len(val) < 1 {
			val = req.PostForm.Get(key)
		}
		v = val
		return
	}
	return ValidAttrFormat(format, ValueGetterF(getter), true, args...)
}

func checkTemplateRequired(data interface{}, required bool, lts []string) (bool, error) {
	zeroChecker, ok := data.(ZeroChecker)
	if v := reflect.ValueOf(data); (ok && zeroChecker.IsZero()) || v.Kind() == reflect.Invalid || (v.IsZero() || reflect.Indirect(v).IsZero()) {
		if (lts[0] == "R" || lts[0] == "r") && required {
			return true, errors.New("data is empty")
		}
		return true, nil
	}
	// if v, ok := data.(string); data == nil || (ok && len(v) < 1) { //chekc the value required.
	// 	if (lts[0] == "R" || lts[0] == "r") && required {
	// 		return true, errors.New("data is empty")
	// 	}
	// 	return true, nil
	// }
	// if v, ok := data.([]interface{}); data == nil || (ok && len(v) < 1) { //chekc the value required.
	// 	if (lts[0] == "R" || lts[0] == "r") && required {
	// 		return true, errors.New("data is empty")
	// 	}
	// 	return true, nil
	// }
	return false, nil
}

func CompatibleType(typ reflect.Type) (ok bool) {
	if typ == reflect.TypeOf(nil) {
		return
	}
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String:
		ok = true
	case reflect.Ptr, reflect.Slice:
		ok = CompatibleType(typ.Elem())
	}
	return
}

// ValidAttrTemple will valid the data to specified value by limit
//
// data: target value for valding
//
// valueType: target value type limit by <R or O>|<S or I or F>
//
//	R is value must be having and valid
//
//	O is vlaue can be empty or nil, but must be valid if it having value
//
//	S:string value,I:integet value,F:float value
//
//	example "R|F" is required float value
//
//	example "O|F" is optional float value
//
// valueRange: taget value range limit by <O or R or P>:<limit pattern>
//
//	O is value must be in options, all optional is seperated by -
//
//	R is value must be in range by "start-end", or "start" to positive infinite or negative infinite  to "-end"
//
//	P is value must be matched by regex pattern
//
//	example "O:1-2-3-4" is valid by value is in options 1-2-3-4)
//	example "P:^.*\@.*$" is valid by string having "@"
//
// required: if true, ValidAttrTemple will return fail when require value is empty or nil,
// if false, ValidAttrTemple will return success although setting required for emppty/nil value
func ValidAttrTemple(data interface{}, valueType string, valueRange string, required bool, enum EnumValider) (interface{}, error) {
	valueRange = strings.Replace(valueRange, "%N", ",", -1)
	valueRange = strings.Replace(valueRange, "%%", "%", -1)
	lts := strings.SplitN(valueType, "|", 2) //valid required type
	if len(lts) < 2 {
		return nil, fmt.Errorf("invalid type limit:%s", valueType)
	}
	lrs := strings.SplitN(valueRange, ":", 2) //valid value range.
	if len(lrs) < 2 {
		return nil, fmt.Errorf("invalid range limit:%s", valueRange)
	}
	if ret, err := checkTemplateRequired(data, required, lts); ret { //check required
		return nil, err
	}
	// required = required && (lts[0] == "R" || lts[0] == "r")
	//define the valid string function.
	validStr := func(ds string) (interface{}, error) {
		//check range limit.
		switch lrs[0] {
		case "o", "O": //option limit.
			options := strings.Split(lrs[1], "~")
			if converter.ArrayHaving(options, ds) {
				return ds, nil
			}
			return nil, fmt.Errorf("invalid value(%s) for options(%s)", ds, lrs[1])
		case "l", "L": //length limit
			slen := int64(len(ds))
			rgs := strings.Split(lrs[1], "~")
			var beg, end int64 = 0, 0
			var err error = nil
			if len(rgs) > 0 && len(rgs[0]) > 0 {
				beg, err = strconv.ParseInt(rgs[0], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid range begin number(%s)", rgs[0])
				}
			} else {
				beg = 0
			}
			if len(rgs) > 1 && len(rgs[1]) > 0 {
				end, err = strconv.ParseInt(rgs[1], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid range end number option(%s)", rgs[1])
				}
			} else {
				end = math.MaxInt64
			}
			if beg < slen && end > slen {
				return ds, nil
			}
			return nil, fmt.Errorf("string length must match %d<len<%d, but %d", beg, end, slen)
		case "p", "P": //regex pattern limit
			mched, err := regexp.MatchString(lrs[1], ds)
			if err != nil {
				return nil, err
			}
			if mched {
				return ds, nil
			}
			return nil, fmt.Errorf("value(%s) not match regex(%s)", ds, lrs[1])
		case "e", "E":
			if enum == nil {
				return nil, fmt.Errorf("target is not enum able")
			}
			return ds, enum.EnumValid(ds)
		case "n", "N":
			return ds, nil
		}
		//unknow range limit type.
		return nil, fmt.Errorf("invalid range limit %s for string", lrs[0])
	}
	//define valid number function.
	validNum := func(ds float64) (interface{}, error) {
		//check range limit.
		switch lrs[0] {
		case "r", "R":
			var beg, end float64 = 0, 0
			var err error = nil
			rgs := strings.Split(lrs[1], "~")
			if len(rgs) > 0 && len(rgs[0]) > 0 {
				beg, err = strconv.ParseFloat(rgs[0], 64)
				if err != nil {
					return nil, fmt.Errorf("invalid range begin number(%s)", rgs[0])
				}
			} else {
				beg = 0
			}
			endStr := ""
			if len(rgs) > 1 && len(rgs[1]) > 0 {
				end, err = strconv.ParseFloat(rgs[1], 64)
				if err != nil {
					return nil, fmt.Errorf("invalid range end number option(%s)", rgs[1])
				}
				endStr = rgs[1]
			} else {
				end = math.MaxFloat64
				endStr = "MaxFloat64"
			}
			if beg < ds && end > ds {
				return ds, nil
			}
			return nil, fmt.Errorf("value must match %f<val<%v, but %v", beg, endStr, ds)
		case "o", "O":
			options := strings.Split(lrs[1], "~")
			var oary []float64
			for _, o := range options { //covert to float array.
				v, err := strconv.ParseFloat(o, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid number option(%s)", lrs[1])
				}
				oary = append(oary, v)
			}
			if converter.ArrayHaving(oary, ds) {
				return ds, nil
			}
			return nil, fmt.Errorf("invalid value(%f) for options(%s)", ds, lrs[1])
		case "e", "E":
			if enum == nil {
				return nil, fmt.Errorf("target is not enum able")
			}
			return ds, enum.EnumValid(ds)
		case "n", "N":
			return ds, nil
		}
		//unknow range limit type.
		return nil, fmt.Errorf("invalid range limit %s for float", lrs[0])
	}
	//define valid number function.
	validInt := func(ds int64) (interface{}, error) {
		//check range limit.
		switch lrs[0] {
		case "r", "R":
			var beg, end int64 = 0, 0
			var err error = nil
			rgs := strings.Split(lrs[1], "~")
			if len(rgs) > 0 && len(rgs[0]) > 0 {
				beg, err = strconv.ParseInt(rgs[0], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid range begin number(%s)", rgs[0])
				}
			} else {
				beg = 0
			}
			endStr := ""
			if len(rgs) > 1 && len(rgs[1]) > 0 {
				end, err = strconv.ParseInt(rgs[1], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid range end number option(%s)", rgs[1])
				}
				endStr = rgs[1]
			} else {
				end = math.MaxInt64
				endStr = "MaxInt64"
			}
			if beg < ds && end > ds {
				return ds, nil
			}
			return nil, fmt.Errorf("value must match %v<val<%v, but %v", beg, endStr, ds)
		case "o", "O":
			options := strings.Split(lrs[1], "~")
			var oary []int64
			for _, o := range options { //covert to float array.
				v, err := strconv.ParseInt(o, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid number option(%s)", lrs[1])
				}
				oary = append(oary, v)
			}
			if converter.ArrayHaving(oary, ds) {
				return ds, nil
			}
			return nil, fmt.Errorf("invalid value(%v) for options(%s)", ds, lrs[1])
		case "e", "E":
			if enum == nil {
				return nil, fmt.Errorf("target is not enum able")
			}
			return ds, enum.EnumValid(ds)
		case "n", "N":
			return ds, nil
		}
		//unknow range limit type.
		return nil, fmt.Errorf("invalid range limit %s for float", lrs[0])
	}
	//define value type function
	validValeuType := func(ds interface{}) (interface{}, error) {
		switch lts[1] {
		case "s", "S":
			sval, _ := converter.StringVal(ds)
			return validStr(sval)
		case "i", "I":
			ids, err := converter.Int64Val(ds)
			if err != nil {
				return nil, fmt.Errorf("invalid value(%s) for type(%s):%v", ds, lts[1], err)
			}
			return validInt(ids)
		case "f", "F":
			fds, err := converter.Float64Val(ds)
			if err != nil {
				return nil, fmt.Errorf("invalid value(%s) for type(%s):%v", ds, lts[1], err)
			}
			return validNum(fds)
		}
		return nil, fmt.Errorf("invalid value type:%s", lts[1])
	}
	return validValeuType(data)
}

func validAttrTemple(data interface{}, temple string, parts []string, required bool, enum EnumValider) (val interface{}, err error) {
	val, err = ValidAttrTemple(data, parts[1], parts[2], required, enum)
	if err != nil {
		err = fmt.Errorf("limit(%s),%s", temple, err.Error())
		if len(parts) > 3 {
			err = errors.New(parts[3])
		}
	}
	return
}

// ValueSetter is interface to set value to target arg
type ValueSetter interface {
	//Set the value
	Set(value interface{}) error
}

// ValueSetterF is func for implment ValueSetter
type ValueSetterF func(interface{}) error

// Set will set value to arg
func (v ValueSetterF) Set(value interface{}) error { return v(value) }

// ValueGetter is inteface for using get the value by key
type ValueGetter interface {
	//Get the value by key
	Get(key string) (interface{}, error)
}

// ValueGetterF is func for implment ValueGetter
type ValueGetterF func(key string) (interface{}, error)

// Get will call the func
func (v ValueGetterF) Get(key string) (interface{}, error) { return v(key) }

// ValidAttrFormat will valid multi value by foramt template, return error if fail
//
// format is temple set is seperated by ";", general it is one line one temple end with ";"
//
//	arg1,R|I,R:0;
//	arg2,O|F,R:0;
//	...
//
// valueGetter is value getter by key
//
// required if true, ValidAttrTemple will return fail when require value is empty or nil,
// if false, ValidAttrTemple will return success although setting required for emppty/nil value
//
// args is variable list for store value, it must be go pointer
//
//	var arg1 int
//	var arg2 float64
//	ValidAttrFormat(format,getter,&arg1,&arg2)
func ValidAttrFormat(format string, valueGetter ValueGetter, required bool, args ...interface{}) error {
	format = regexp.MustCompile(`\/\/.*`).ReplaceAllString(format, "")
	format = strings.Replace(format, "\n", "", -1)
	format = strings.Trim(format, " \t;")
	if len(format) < 1 {
		return errors.New("format not found")
	}
	temples := strings.Split(format, ";")
	if len(args) < 1 {
		args = make([]interface{}, len(temples))
	}
	if len(temples) != len(args) {
		return errors.New("args count is not equal format count")
	}
	for idx, temple := range temples {
		temple = strings.TrimSpace(temple)
		parts := strings.SplitN(temple, ",", 4)
		if len(parts) < 3 {
			return fmt.Errorf("temple error:%s", temple)
		}
		sval, err := valueGetter.Get(parts[0])
		if err != nil {
			return fmt.Errorf("get value by key %v fail with %v", parts[0], err)
		}
		var target interface{}
		var targetValue reflect.Value
		var targetKind reflect.Kind
		var enum EnumValider
		if args[idx] != nil {
			target = args[idx]
			targetValue = reflect.Indirect(reflect.ValueOf(target))
			targetKind = targetValue.Kind()
			enum, _ = args[idx].(EnumValider)
		}
		checkValue := reflect.ValueOf(sval)
		if checkValue.Kind() == reflect.Ptr && !checkValue.IsZero() {
			sval = reflect.Indirect(reflect.ValueOf(sval)).Interface()
		}
		if targetKind != reflect.Slice {
			rval, err := validAttrTemple(sval, temple, parts, required, enum)
			if err != nil {
				return err
			}
			if rval == nil {
				continue
			}
			var setted bool
			if target != nil && !CompatibleType(targetValue.Type()) {
				if setter, ok := target.(ValueSetter); ok {
					setted = true
					err = setter.Set(rval)
				} else if sc, ok := target.(sql.Scanner); ok {
					setted = true
					err = sc.Scan(rval)
				}
			}
			if target != nil && !setted {
				err = ValidSetValue(target, rval)
			}
			if err != nil {
				return fmt.Errorf("set value to %v by key %v,%v fail with %v", reflect.TypeOf(target), parts[0], reflect.TypeOf(sval), err)
			}
			continue
		}
		if target != nil && !CompatibleType(targetValue.Type()) {
			if setter, ok := target.(ValueSetter); ok {
				_, err = validAttrTemple(sval, temple, parts, required, enum)
				if err != nil {
					return err
				}
				err = setter.Set(sval)
				if err != nil {
					return fmt.Errorf("set value to %v by key %v,%v fail with %v", reflect.TypeOf(target), parts[0], reflect.TypeOf(sval), err)
				}
				continue
			}
			if sc, ok := target.(sql.Scanner); ok {
				_, err = validAttrTemple(sval, temple, parts, required, enum)
				if err != nil {
					return err
				}
				err = sc.Scan(sval)
				if err != nil {
					return fmt.Errorf("set value to %v by key %v,%v fail with %v", reflect.TypeOf(target), parts[0], reflect.TypeOf(sval), err)
				}
				continue
			}
		}
		svals, _ := converter.ArrayValAll(sval, true) //ignore error
		// if err != nil && err != converter.ErrNil {
		// 	return err
		// }
		if ret, err := checkTemplateRequired(svals, required, strings.SplitN(parts[1], "|", 2)); ret { //check required
			if err != nil {
				return err
			}
			continue
		}
		var targetArray reflect.Value
		if target != nil {
			targetArray = targetValue
		}
		for _, sval = range svals {
			rval, err := validAttrTemple(sval, temple, parts, required, enum)
			if err != nil {
				return err
			}
			if rval == nil {
				continue
			}
			if targetArray.IsValid() {
				tval, err := ValidValue(targetValue.Type().Elem(), rval)
				if err != nil {
					return err
				}
				targetArray = reflect.Append(targetArray, reflect.ValueOf(tval))
			}
		}
		if targetArray.IsValid() {
			targetValue.Set(targetArray)
		}
	}
	return nil
}

// ValidSetValue will convert src value to dst type and set it
func ValidSetValue(dst, src interface{}) error {
	if sc, ok := dst.(sql.Scanner); ok {
		return sc.Scan(src)
	}
	if setter, ok := dst.(ValueSetter); ok {
		return setter.Set(src)
	}
	dstValue := reflect.Indirect(reflect.ValueOf(dst))
	val, err := ValidValue(dstValue.Type(), src)
	if err == nil {
		targetValue := reflect.ValueOf(val)
		if targetValue.Type() == dstValue.Type() {
			dstValue.Set(targetValue)
		} else {
			dstValue.Set(targetValue.Convert(dstValue.Type()))
		}
	}
	return err
}

// ValidValue will convert src value to dst type and return it
func ValidValue(dst reflect.Type, src interface{}) (val interface{}, err error) {
	srcType := reflect.TypeOf(src)
	srcValue := reflect.ValueOf(src)
	if srcType.Kind() == dst.Kind() {
		return srcValue.Convert(dst).Interface(), nil
	}
	if dst.Kind() != reflect.String && srcValue.CanConvert(dst) { //skip string for int=>string is invalid
		return srcValue.Convert(dst).Interface(), nil
	}
	var isptr = false
	var kind = dst.Kind()
	if kind == reflect.Ptr {
		kind = dst.Elem().Kind()
		isptr = true
	}
	var tiv int64
	var tfv float64
	var tsv string
	switch kind {
	case reflect.Int:
		tiv, err = converter.Int64Val(src)
		if err == nil {
			target := int(tiv)
			if isptr {
				val = &target
			} else {
				val = target
			}
		}
	case reflect.Int16:
		tiv, err = converter.Int64Val(src)
		if err == nil {
			target := int16(tiv)
			if isptr {
				val = &target
			} else {
				val = target
			}
		}
	case reflect.Int32:
		tiv, err = converter.Int64Val(src)
		if err == nil {
			target := int32(tiv)
			if isptr {
				val = &target
			} else {
				val = target
			}
		}
	case reflect.Int64:
		tiv, err = converter.Int64Val(src)
		if err == nil {
			target := int64(tiv)
			if isptr {
				val = &target
			} else {
				val = target
			}
		}
	case reflect.Uint:
		tiv, err = converter.Int64Val(src)
		if err == nil {
			target := uint(tiv)
			if isptr {
				val = &target
			} else {
				val = target
			}
		}
	case reflect.Uint16:
		tiv, err = converter.Int64Val(src)
		if err == nil {
			target := uint16(tiv)
			if isptr {
				val = &target
			} else {
				val = target
			}
		}
	case reflect.Uint32:
		tiv, err = converter.Int64Val(src)
		if err == nil {
			target := uint32(tiv)
			if isptr {
				val = &target
			} else {
				val = target
			}
		}
	case reflect.Uint64:
		tiv, err = converter.Int64Val(src)
		if err == nil {
			target := uint64(tiv)
			if isptr {
				val = &target
			} else {
				val = target
			}
		}
	case reflect.Float32:
		tfv, err = converter.Float64Val(src)
		if err == nil {
			target := float32(tfv)
			if isptr {
				val = &target
			} else {
				val = target
			}
		}
	case reflect.Float64:
		tfv, err = converter.Float64Val(src)
		if err == nil {
			target := float64(tfv)
			if isptr {
				val = &target
			} else {
				val = target
			}
		}
	case reflect.String:
		tsv, err = converter.StringVal(src)
		if isptr {
			val = &tsv
		} else {
			val = tsv
		}
	default:
		err = fmt.Errorf("not supported")
	}
	if err == nil {
		return val, err
	}
	return nil, fmt.Errorf("parse kind(%v) value to kind(%v) value->%v", srcType.Kind(), dst, err)
}

type rawMapable interface {
	RawMap() map[string]interface{}
}

// Struct is validable struct impl
type Struct struct {
	Target   interface{}
	Required bool
	Tag      string
	loaded   map[string]interface{}
}

// NewStruct will return new struct
func NewStruct(target interface{}) (s *Struct) {
	if reflect.TypeOf(target).Kind() != reflect.Ptr {
		panic("target must be pointer")
	}
	s = &Struct{
		Target:   target,
		Required: true,
		Tag:      "json",
		loaded:   map[string]interface{}{},
	}
	return
}

// Get will return field value by key
func (s *Struct) Get(key string) (value interface{}, err error) {
	if len(s.loaded) < 1 {
		value := reflect.ValueOf(s.Target).Elem()
		vtype := reflect.TypeOf(s.Target).Elem()
		for i := 0; i < vtype.NumField(); i++ {
			valueField := value.Field(i)
			typeField := vtype.Field(i)
			if !typeField.IsExported() {
				continue
			}
			tag := strings.SplitN(typeField.Tag.Get(s.Tag), ",", 2)[0]
			targetValue := valueField.Interface()
			if mv, ok := targetValue.(map[string]interface{}); ok {
				targetValue = M(mv)
			} else if mv, ok := targetValue.(rawMapable); ok {
				targetValue = M(mv.RawMap())
			} else {
				if typeField.Type.Kind() == reflect.Struct {
					targetValue = NewStruct(valueField.Addr().Interface())
				} else if typeField.Type.Kind() == reflect.Ptr && typeField.Type.Elem().Kind() == reflect.Struct {
					targetValue = NewStruct(valueField.Interface())
				}
			}
			s.loaded[typeField.Name] = targetValue
			s.loaded[tag] = targetValue
		}
	}
	key = strings.Trim(key, "/")
	parts := strings.SplitN(key, "/", 2)
	value = s.loaded[parts[0]]
	if len(parts) < 2 || value == nil {
		return
	}
	if getter, ok := value.(ValueGetter); ok {
		value, err = getter.Get(parts[1])
	}
	return
}

// ValidFormat will valid format to struct filed
func (s *Struct) ValidFormat(format string, args ...interface{}) error {
	return ValidAttrFormat(format, s, s.Required, args...)
}

// ValidStructAttrFormat will valid struct by filed
func ValidStructAttrFormat(format string, target interface{}, required bool, args ...interface{}) error {
	return ValidAttrFormat(format, NewStruct(target), required, args...)
}

// ValidFormat will check all supported type and run valid format
func ValidFormat(format string, target interface{}, args ...interface{}) error {
	if getter, ok := target.(ValueGetter); ok {
		return ValidAttrFormat(format, getter, true, args...)
	}
	if req, ok := target.(*http.Request); ok {
		return QueryValidFormat(req, format, args...)
	}
	if val, ok := target.(url.Values); ok {
		return Values(val).ValidFormat(format, args...)
	}
	if ms, ok := target.(map[string]string); ok {
		return MS(ms).ValidFormat(format, args...)
	}
	if mv, ok := target.(map[string]interface{}); ok {
		return M(mv).ValidFormat(format, args...)
	}
	if mv, ok := target.(rawMapable); ok {
		return M(mv.RawMap()).ValidFormat(format, args...)
	}
	return NewStruct(target).ValidFormat(format, args...)
}

type Valider struct {
	attrscan.Scanner
}

var Default = &Valider{
	Scanner: attrscan.Scanner{
		Tag: "json",
		NameConv: func(on, name string, field reflect.StructField) string {
			return name
		},
	},
}

func ValidArgs(target interface{}, filter string, args ...interface{}) (format string, args_ []interface{}) {
	format, args_ = Default.ValidArgs(target, filter, args...)
	return
}

func (v *Valider) addValidArgs(target interface{}, filter string, format string, args []interface{}) (format_ string, args_ []interface{}) {
	format_, args_ = format, args
	v.FilterFieldCall("valid", target, filter, func(fieldName, fieldFunc string, field reflect.StructField, value interface{}) {
		valid := field.Tag.Get("valid")
		if field.Type.Kind() == reflect.Struct && valid == "inline" {
			format_, args_ = v.addValidArgs(value, filter, format_, args_)
			return
		}
		if len(valid) < 1 || valid == "-" {
			return
		}
		valid = strings.TrimSpace(valid)
		if strings.HasSuffix(valid, ";") {
			format_ += valid + "\n"
		} else {
			format_ += valid + ";\n"
		}
		args_ = append(args_, value)
	})
	return
}

func (v *Valider) ValidArgs(target interface{}, filter string, args ...interface{}) (format string, args_ []interface{}) {
	format, args_ = v.addValidArgs(target, filter, "", nil)
	for _, arg := range args {
		if v, ok := arg.(string); ok {
			arg = strings.TrimSpace(v)
			if strings.HasSuffix(v, ";") {
				format += v + "\n"
			} else {
				format += v + ";\n"
			}
		} else {
			args_ = append(args_, arg)
		}
	}
	format = strings.TrimSpace(format)
	return
}

func Valid(target interface{}, filter, optional string) (err error) {
	err = Default.Valid(target, filter, optional)
	return
}

func (v *Valider) Valid(target interface{}, filter, optional string) (err error) {
	errList := []string{}
	optional = strings.TrimSpace(optional)
	isExc := strings.HasPrefix(optional, "^")
	optional = strings.TrimPrefix(optional, "^")
	optional = strings.Trim(optional, ",")
	optional = "," + optional + ","
	isRequired := func(fieldName string) bool {
		if isExc {
			return strings.Contains(optional, ","+fieldName+",")
		} else {
			return !strings.Contains(optional, ","+fieldName+",")
		}
	}
	v.FilterFieldCall("valid", target, filter, func(fieldName, fieldFunc string, field reflect.StructField, value interface{}) {
		valid := field.Tag.Get("valid")
		if len(valid) < 1 {
			return
		}
		valid = strings.TrimSpace(valid)
		valid = strings.TrimSuffix(valid, ";")
		parts := strings.SplitN(valid, ",", 4)
		if len(parts) < 3 {
			errList = append(errList, fmt.Sprintf("valid error:%s", valid))
			return
		}
		var xerr error
		enum, _ := value.(EnumValider)
		targetValue := reflect.Indirect(reflect.ValueOf(value))
		if targetValue.Kind() == reflect.Slice {
			n := targetValue.Len()
			for i := 0; i < n; i++ {
				targetItem := targetValue.Index(i)
				_, xerr = validAttrTemple(targetItem.Interface(), valid, parts, isRequired(fieldName), enum)
				if xerr != nil {
					break
				}
			}
			if xerr == nil && n < 1 {
				_, xerr = validAttrTemple(nil, valid, parts, isRequired(fieldName), enum)
			}
		} else {
			_, xerr = validAttrTemple(targetValue.Interface(), valid, parts, isRequired(fieldName), enum)
		}
		if xerr != nil {
			errList = append(errList, xerr.Error())
		}
	})
	if len(errList) > 0 {
		err = fmt.Errorf("%v", strings.Join(errList, "\n"))
	}
	return
}
