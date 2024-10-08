package xsql

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"sort"
// 	"testing"
// 	"time"

// 	"github.com/wfunc/util/attrvalid"
// 	"github.com/wfunc/util/converter"
// 	"github.com/wfunc/util/xmap"
// 	"github.com/wfunc/util/xtime"
// 	"github.com/jackc/pgx/v4/pgxpool"
// 	"github.com/shopspring/decimal"
// )

// var Pool func() *pgxpool.Pool

// func init() {
// 	pool, err := pgxpool.Connect(context.Background(), "postgresql://dev:123@psql.loc:5432/dev")
// 	if err != nil {
// 		panic(err)
// 	}
// 	Pool = func() *pgxpool.Pool {
// 		return pool
// 	}
// }

// type timeObject struct {
// 	CreateTime Time
// }

// func TestTime(t *testing.T) {
// 	obj := &timeObject{
// 		CreateTime: Time(time.Now()),
// 	}
// 	bys, err := json.Marshal(obj)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	obj2 := &timeObject{}
// 	err = json.Unmarshal(bys, obj2)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if obj.CreateTime.Timestamp() != obj2.CreateTime.Timestamp() {
// 		t.Error("error")
// 		return
// 	}
// 	//
// 	t1 := TimeZero()
// 	bys, err = t1.MarshalJSON()
// 	if err != nil || string(bys) != "0" {
// 		t.Errorf("err:%v,bys:%v", err, string(bys))
// 		return
// 	}
// 	//
// 	//
// 	t2 := Time{}
// 	err = t2.UnmarshalJSON([]byte("null"))
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	bys, err = t2.MarshalJSON()
// 	if err != nil || string(bys) != "0" {
// 		t.Errorf("err:%v,bys:%v", err, string(bys))
// 		return
// 	}
// 	err = t2.Scan(time.Now())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	//
// 	//
// 	fmt.Printf("TimeZero--->%v\n", TimeZero())
// 	fmt.Printf("TimeNow--->%v\n", TimeNow())
// 	fmt.Printf("TimeStartOfToday--->%v\n", TimeStartOfToday())
// 	fmt.Printf("TimeStartOfWeek--->%v\n", TimeStartOfWeek())
// 	fmt.Printf("TimeStartOfMonth--->%v\n", TimeStartOfMonth())
// 	fmt.Printf("TimeUnix--->%v\n", TimeUnix(0))
// 	fmt.Printf("AsTime--->%v\n", t2.AsTime())
// }

// func TestMap(t *testing.T) {
// 	var eval M
// 	err := eval.Scan(1)
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	//
// 	_, err = Pool().Exec(context.Background(), `drop table if exists xsql_test_map`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	_, err = Pool().Exec(context.Background(), `create table xsql_test_map(tid int,mval text)`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	{ //normal
// 		var mval = M{"a": 1}
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_map values ($1,$2)`, 1, mval)
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var mval1 M
// 		err = Pool().QueryRow(context.Background(), `select mval from xsql_test_map where tid=$1`, 1).Scan(&mval1)
// 		if err != nil || len(mval1) != 1 {
// 			t.Error(err)
// 			return
// 		}
// 		mv1 := xmap.Wrap(xmap.M(mval1))
// 		if mv1.Int("a") != 1 {
// 			t.Error("error")
// 		}
// 	}
// 	{ //nil
// 		var mval M = nil
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_map values ($1,$2)`, 2, mval)
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var mval1 M
// 		err = Pool().QueryRow(context.Background(), `select mval from xsql_test_map where tid=$1`, 2).Scan(&mval1)
// 		if err != nil || len(mval1) != 0 {
// 			t.Error(err)
// 			return
// 		}
// 		mval1.RawMap()
// 		mval1.AsMap()
// 	}
// }

// func TestMapArray(t *testing.T) {
// 	var eval MArray
// 	err := eval.Scan(1)
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	//
// 	_, err = Pool().Exec(context.Background(), `drop table if exists xsql_test_map`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	_, err = Pool().Exec(context.Background(), `create table xsql_test_map(tid int,mval text)`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	{ //normal
// 		var mval = MArray{{"a": 1}}
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_map values ($1,$2)`, 1, mval)
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var mval1 MArray
// 		err = Pool().QueryRow(context.Background(), `select mval from xsql_test_map where tid=$1`, 1).Scan(&mval1)
// 		if err != nil || len(mval1) != 1 {
// 			t.Error(err)
// 			return
// 		}
// 		mv1 := xmap.Wrap(mval1[0])
// 		if mv1.Int("a") != 1 {
// 			t.Error("error")
// 		}
// 	}
// 	{ //nil
// 		var mval MArray = nil
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_map values ($1,$2)`, 2, mval)
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var mval1 MArray
// 		err = Pool().QueryRow(context.Background(), `select mval from xsql_test_map where tid=$1`, 2).Scan(&mval1)
// 		if err != nil || len(mval1) != 0 {
// 			t.Error(err)
// 			return
// 		}
// 	}
// }

// func TestSQLScan(t *testing.T) {
// 	err := sqlScan(nil, nil, nil)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	err = sqlScan("xx", nil, nil)
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	err = sqlScan("[xx", nil, nil)
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	err = sqlScan("", nil, nil)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	err = sqlScan("null", nil, nil)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// }

// func TestIntArray(t *testing.T) {
// 	var ary IntArray
// 	err := ary.Scan(1)
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	err = ary.Scan("a")
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	// ary.Value()
// 	//
// 	v0, v1, v2 := int(3), int(2), int(1)
// 	ary = append(ary, v0)
// 	ary = append(ary, v1)
// 	ary = append(ary, v2)
// 	sort.Sort(ary)
// 	//
// 	_, err = Pool().Exec(context.Background(), `drop table if exists xsql_test_int`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	_, err = Pool().Exec(context.Background(), `create table xsql_test_int(tid int,iarry text)`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	{ //normal
// 		_, err = Pool().Exec(context.Background(), `delete from xsql_test_int`)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int values ($1,$2)`, 1, ary)
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 IntArray
// 		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int where tid=$1`, 1).Scan(&ary1)
// 		if err != nil || len(ary1) != 3 {
// 			t.Error(err)
// 			return
// 		}
// 		if !ary1.HavingOne(3) {
// 			t.Error("error")
// 			return
// 		}
// 		if ary1.HavingOne(4) {
// 			t.Error("error")
// 			return
// 		}
// 		if ary1.Join(",") != "1,2,3" {
// 			t.Error("error")
// 		}
// 		if ary1.DbArray() != "{1,2,3}" {
// 			t.Error("error")
// 		}
// 		ary2 := append(ary1, 3).RemoveDuplicate()
// 		if ary2.DbArray() != "{1,2,3}" {
// 			t.Error("error")
// 		}
// 	}
// 	{ //string array
// 		_, err = Pool().Exec(context.Background(), `delete from xsql_test_int`)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int values ($1,$2)`, 1, ary.StrArray())
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 IntArray
// 		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int where tid=$1`, 1).Scan(&ary1)
// 		if err != nil || len(ary1) != 3 {
// 			t.Error(err)
// 			return
// 		}
// 		if ary1.Join(",") != "1,2,3" {
// 			t.Error("error")
// 		}
// 	}
// 	{ //nil
// 		var arynil IntArray = nil
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int values ($1,$2)`, 2, arynil)
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 IntArray
// 		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int where tid=$1`, 2).Scan(&ary1)
// 		if err != nil || len(ary1) != 0 {
// 			t.Error(err)
// 			return
// 		}
// 	}
// 	ary.AsPtrArray()
// }

// func TestIntPtrArray(t *testing.T) {
// 	var ary IntPtrArray
// 	err := ary.Scan(1)
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	err = ary.Scan("a")
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	// ary.Value()
// 	//
// 	v0, v1, v2 := int(3), int(2), int(1)
// 	ary = append(ary, &v0)
// 	ary = append(ary, &v1)
// 	ary = append(ary, &v2)
// 	ary = append(ary, nil)
// 	sort.Sort(ary)
// 	//
// 	_, err = Pool().Exec(context.Background(), `drop table if exists xsql_test_int`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	_, err = Pool().Exec(context.Background(), `create table xsql_test_int(tid int,iarry text)`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	{ //normal
// 		_, err = Pool().Exec(context.Background(), `delete from xsql_test_int`)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int values ($1,$2)`, 1, ary)
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 IntPtrArray
// 		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int where tid=$1`, 1).Scan(&ary1)
// 		if err != nil || len(ary1) != 4 {
// 			t.Errorf("%v,%v", err, ary1)
// 			return
// 		}
// 		if !ary1.HavingOne(3) {
// 			t.Error("error")
// 			return
// 		}
// 		if ary1.HavingOne(4) {
// 			t.Error("error")
// 			return
// 		}
// 		if ary1.Join(",") != "1,2,3" {
// 			t.Errorf("%v", ary1)
// 			return
// 		}
// 		if ary1.DbArray() != "{1,2,3}" {
// 			t.Error("error")
// 		}
// 		ary2 := append(ary1, converter.IntPtr(3)).RemoveDuplicate()
// 		if ary2.DbArray() != "{1,2,3}" {
// 			t.Error("error")
// 		}
// 	}
// 	{ //string array
// 		_, err = Pool().Exec(context.Background(), `delete from xsql_test_int`)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int values ($1,$2)`, 1, ary.StrArray())
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 IntPtrArray
// 		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int where tid=$1`, 1).Scan(&ary1)
// 		if err != nil || len(ary1) != 3 {
// 			t.Error(err)
// 			return
// 		}
// 		if ary1.Join(",") != "1,2,3" {
// 			t.Error("error")
// 		}
// 	}
// 	{ //nil
// 		var arynil IntPtrArray = nil
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int values ($1,$2)`, 2, arynil)
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 IntPtrArray
// 		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int where tid=$1`, 2).Scan(&ary1)
// 		if err != nil || len(ary1) != 0 {
// 			t.Error(err)
// 			return
// 		}
// 		{ //join
// 			var ary1 = IntPtrArray{}
// 			ary1 = append(ary1, &v0)
// 			ary1 = append(ary1, nil)
// 			ary1 = append(ary1, &v2)
// 			sort.Sort(ary1)
// 			if ary1.Join(",") != "1,3" {
// 				t.Error(ary1.Join(","))
// 				return
// 			}
// 		}
// 	}
// 	ary.AsArray()
// 	ary[0] = nil
// 	ary.AsArray()
// }

// func TestInt64Array(t *testing.T) {
// 	var ary Int64Array
// 	err := ary.Scan(1)
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	err = ary.Scan("a")
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	// ary.Value()
// 	//
// 	v0, v1, v2 := int64(3), int64(2), int64(1)
// 	ary = append(ary, v0)
// 	ary = append(ary, v1)
// 	ary = append(ary, v2)
// 	sort.Sort(ary)
// 	//
// 	_, err = Pool().Exec(context.Background(), `drop table if exists xsql_test_int64`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	_, err = Pool().Exec(context.Background(), `create table xsql_test_int64(tid int,iarry text)`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	{ //normal
// 		_, err = Pool().Exec(context.Background(), `delete from xsql_test_int64`)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int64 values ($1,$2)`, 1, ary)
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 Int64Array
// 		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int64 where tid=$1`, 1).Scan(&ary1)
// 		if err != nil || len(ary1) != 3 {
// 			t.Error(err)
// 			return
// 		}
// 		if !ary1.HavingOne(3) {
// 			t.Error("error")
// 			return
// 		}
// 		if ary1.HavingOne(4) {
// 			t.Error("error")
// 			return
// 		}
// 		if ary1.Join(",") != "1,2,3" {
// 			t.Error("error")
// 		}
// 		if ary1.DbArray() != "{1,2,3}" {
// 			t.Error("error")
// 		}
// 		ary2 := append(ary1, 3).RemoveDuplicate()
// 		if ary2.DbArray() != "{1,2,3}" {
// 			t.Error("error")
// 		}
// 	}
// 	{ //string array
// 		_, err = Pool().Exec(context.Background(), `delete from xsql_test_int64`)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int64 values ($1,$2)`, 1, ary.StrArray())
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 Int64Array
// 		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int64 where tid=$1`, 1).Scan(&ary1)
// 		if err != nil || len(ary1) != 3 {
// 			t.Error(err)
// 			return
// 		}
// 		if ary1.Join(",") != "1,2,3" {
// 			t.Error("error")
// 		}
// 	}
// 	{ //nil
// 		var arynil Int64Array = nil
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int64 values ($1,$2)`, 2, arynil)
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 Int64Array
// 		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int64 where tid=$1`, 2).Scan(&ary1)
// 		if err != nil || len(ary1) != 0 {
// 			t.Error(err)
// 			return
// 		}
// 	}
// 	ary.AsPtrArray()
// }

// func TestInt64PtrArray(t *testing.T) {
// 	var ary Int64PtrArray
// 	err := ary.Scan(1)
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	err = ary.Scan("a")
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	// ary.Value()
// 	//
// 	v0, v1, v2 := int64(3), int64(2), int64(1)
// 	ary = append(ary, &v0)
// 	ary = append(ary, &v1)
// 	ary = append(ary, &v2)
// 	sort.Sort(ary)
// 	//
// 	_, err = Pool().Exec(context.Background(), `drop table if exists xsql_test_int64`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	_, err = Pool().Exec(context.Background(), `create table xsql_test_int64(tid int,iarry text)`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	{ //normal
// 		_, err = Pool().Exec(context.Background(), `delete from xsql_test_int64`)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int64 values ($1,$2)`, 1, ary)
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 Int64PtrArray
// 		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int64 where tid=$1`, 1).Scan(&ary1)
// 		if err != nil || len(ary1) != 3 {
// 			t.Error(err)
// 			return
// 		}
// 		if !ary1.HavingOne(3) {
// 			t.Error("error")
// 			return
// 		}
// 		if ary1.HavingOne(4) {
// 			t.Error("error")
// 			return
// 		}
// 		if ary1.Join(",") != "1,2,3" {
// 			t.Error("error")
// 		}
// 		if ary1.DbArray() != "{1,2,3}" {
// 			t.Error("error")
// 		}
// 		ary2 := append(ary1, converter.Int64Ptr(3)).RemoveDuplicate()
// 		if ary2.DbArray() != "{1,2,3}" {
// 			t.Error("error")
// 		}
// 	}
// 	{ //string array
// 		_, err = Pool().Exec(context.Background(), `delete from xsql_test_int64`)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int64 values ($1,$2)`, 1, ary.StrArray())
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 Int64PtrArray
// 		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int64 where tid=$1`, 1).Scan(&ary1)
// 		if err != nil || len(ary1) != 3 {
// 			t.Error(err)
// 			return
// 		}
// 		if ary1.Join(",") != "1,2,3" {
// 			t.Error("error")
// 		}
// 	}
// 	{ //nil
// 		var arynil Int64PtrArray = nil
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int64 values ($1,$2)`, 2, arynil)
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 Int64PtrArray
// 		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int64 where tid=$1`, 2).Scan(&ary1)
// 		if err != nil || len(ary1) != 0 {
// 			t.Error(err)
// 			return
// 		}
// 		{ //join
// 			var ary1 = Int64PtrArray{}
// 			ary1 = append(ary1, &v0)
// 			ary1 = append(ary1, nil)
// 			ary1 = append(ary1, &v2)
// 			sort.Sort(ary1)
// 			if ary1.Join(",") != "1,3" {
// 				t.Error(ary1.Join(","))
// 				return
// 			}
// 		}
// 	}
// 	ary.AsArray()
// 	ary[0] = nil
// 	ary.AsArray()
// }

// func TestFloat64Array(t *testing.T) {
// 	var ary Float64Array
// 	err := ary.Scan(1)
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	err = ary.Scan("a")
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	// ary.Value()
// 	//
// 	v0, v1, v2 := float64(3), float64(2), float64(1)
// 	ary = append(ary, v0)
// 	ary = append(ary, v1)
// 	ary = append(ary, v2)
// 	sort.Sort(ary)
// 	//
// 	_, err = Pool().Exec(context.Background(), `drop table if exists xsql_test_float64`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	_, err = Pool().Exec(context.Background(), `create table xsql_test_float64(tid int,iarry text)`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	{ //normal
// 		_, err = Pool().Exec(context.Background(), `delete from xsql_test_float64`)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_float64 values ($1,$2)`, 1, ary)
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 Float64Array
// 		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_float64 where tid=$1`, 1).Scan(&ary1)
// 		if err != nil || len(ary1) != 3 {
// 			t.Error(err)
// 			return
// 		}
// 		if !ary1.HavingOne(3) {
// 			t.Error("error")
// 			return
// 		}
// 		if ary1.HavingOne(4) {
// 			t.Error("error")
// 			return
// 		}
// 		if ary1.Join(",") != "1,2,3" {
// 			t.Error("error")
// 		}
// 		if ary1.DbArray() != "{1,2,3}" {
// 			t.Error("error")
// 		}
// 		ary2 := append(ary1, 3).RemoveDuplicate()
// 		if ary2.DbArray() != "{1,2,3}" {
// 			t.Error("error")
// 		}
// 	}
// 	{ //string array
// 		_, err = Pool().Exec(context.Background(), `delete from xsql_test_float64`)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_float64 values ($1,$2)`, 1, ary.StrArray())
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 Float64Array
// 		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_float64 where tid=$1`, 1).Scan(&ary1)
// 		if err != nil || len(ary1) != 3 {
// 			t.Error(err)
// 			return
// 		}
// 		if ary1.Join(",") != "1,2,3" {
// 			t.Error("error")
// 		}
// 	}
// 	{ //nil
// 		var arynil Float64Array = nil
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_float64 values ($1,$2)`, 2, arynil)
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 Float64Array
// 		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_float64 where tid=$1`, 2).Scan(&ary1)
// 		if err != nil || len(ary1) != 0 {
// 			t.Error(err)
// 			return
// 		}
// 	}
// 	ary.AsPtrArray()
// }

// func TestFloat64PtrArray(t *testing.T) {
// 	var ary Float64PtrArray
// 	err := ary.Scan(1)
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	err = ary.Scan("a")
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	// ary.Value()
// 	//
// 	v0, v1, v2 := float64(3), float64(2), float64(1)
// 	ary = append(ary, &v0)
// 	ary = append(ary, &v1)
// 	ary = append(ary, &v2)
// 	sort.Sort(ary)
// 	//
// 	_, err = Pool().Exec(context.Background(), `drop table if exists xsql_test_float64`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	_, err = Pool().Exec(context.Background(), `create table xsql_test_float64(tid int,iarry text)`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	{ //normal
// 		_, err = Pool().Exec(context.Background(), `delete from xsql_test_float64`)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_float64 values ($1,$2)`, 1, ary)
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 Float64PtrArray
// 		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_float64 where tid=$1`, 1).Scan(&ary1)
// 		if err != nil || len(ary1) != 3 {
// 			t.Error(err)
// 			return
// 		}
// 		if !ary1.HavingOne(3) {
// 			t.Error("error")
// 			return
// 		}
// 		if ary1.HavingOne(4) {
// 			t.Error("error")
// 			return
// 		}
// 		if ary1.Join(",") != "1,2,3" {
// 			t.Error(ary1.Join(","))
// 		}
// 		if ary1.DbArray() != "{1,2,3}" {
// 			t.Error("error")
// 		}
// 		ary2 := append(ary1, converter.Float64Ptr(3)).RemoveDuplicate()
// 		if ary2.DbArray() != "{1,2,3}" {
// 			t.Error("error")
// 		}
// 	}
// 	{ //string array
// 		_, err = Pool().Exec(context.Background(), `delete from xsql_test_float64`)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_float64 values ($1,$2)`, 1, ary.StrArray())
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 Float64PtrArray
// 		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_float64 where tid=$1`, 1).Scan(&ary1)
// 		if err != nil || len(ary1) != 3 {
// 			t.Error(err)
// 			return
// 		}
// 		if ary1.Join(",") != "1,2,3" {
// 			t.Error("error")
// 		}
// 	}
// 	{ //nil
// 		var arynil Float64PtrArray = nil
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_float64 values ($1,$2)`, 2, arynil)
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 Float64PtrArray
// 		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_float64 where tid=$1`, 2).Scan(&ary1)
// 		if err != nil || len(ary1) != 0 {
// 			t.Error(err)
// 			return
// 		}
// 	}
// 	{ //join
// 		var ary1 = Float64PtrArray{}
// 		ary1 = append(ary1, &v0)
// 		ary1 = append(ary1, nil)
// 		ary1 = append(ary1, &v2)
// 		sort.Sort(ary1)
// 		if ary1.Join(",") != "1,3" {
// 			t.Error(ary1.Join(","))
// 			return
// 		}
// 	}
// 	ary.AsArray()
// 	ary[0] = nil
// 	ary.AsArray()
// }

// func TestStringArray(t *testing.T) {
// 	var ary StringArray
// 	err := ary.Scan(1)
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	err = ary.Scan("a")
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	// ary.Value()
// 	//
// 	v0, v1, v2 := "3", "2", "1"
// 	ary = append(ary, v0)
// 	ary = append(ary, v1)
// 	ary = append(ary, v2)
// 	sort.Sort(ary)
// 	//
// 	_, err = Pool().Exec(context.Background(), `drop table if exists xsql_test_string`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	_, err = Pool().Exec(context.Background(), `create table xsql_test_string(tid int,sarry text)`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	{ //normal
// 		_, err = Pool().Exec(context.Background(), `delete from xsql_test_string`)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_string values ($1,$2)`, 1, ary)
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 StringArray
// 		err = Pool().QueryRow(context.Background(), `select sarry from xsql_test_string where tid=$1`, 1).Scan(&ary1)
// 		if err != nil || len(ary1) != 3 {
// 			t.Error(err)
// 			return
// 		}
// 		if !ary1.HavingOne("3") {
// 			t.Error("error")
// 			return
// 		}
// 		if ary1.HavingOne("4") {
// 			t.Error("error")
// 			return
// 		}
// 		if ary1.Join(",") != "1,2,3" {
// 			t.Error("error")
// 		}
// 		if ary1.DbArray() != "{1,2,3}" {
// 			t.Error("error")
// 		}
// 		ary2 := append(ary1, "3", "", " ").RemoveDuplicate(true, true)
// 		if ary2.DbArray() != "{1,2,3}" {
// 			t.Error("error")
// 		}
// 		ary3 := append(ary1, "", " ").RemoveEmpty(true)
// 		if ary3.DbArray() != "{1,2,3}" {
// 			t.Error("error")
// 		}
// 	}
// 	{ //string array
// 		_, err = Pool().Exec(context.Background(), `delete from xsql_test_string`)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_string values ($1,$2)`, 1, ary.StrArray())
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 StringArray
// 		err = Pool().QueryRow(context.Background(), `select sarry from xsql_test_string where tid=$1`, 1).Scan(&ary1)
// 		if err != nil || len(ary1) != 3 {
// 			t.Error(err)
// 			return
// 		}
// 		if ary1.Join(",") != "1,2,3" {
// 			t.Error("error")
// 		}
// 	}
// 	{ //nil
// 		var arynil StringArray = nil
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_string values ($1,$2)`, 2, arynil)
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 StringArray
// 		err = Pool().QueryRow(context.Background(), `select sarry from xsql_test_string where tid=$1`, 2).Scan(&ary1)
// 		if err != nil || len(ary1) != 0 {
// 			t.Error(err)
// 			return
// 		}
// 	}
// 	ary.AsPtrArray()
// }

// func TestStringPtrArray(t *testing.T) {
// 	var ary StringPtrArray
// 	err := ary.Scan(1)
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	err = ary.Scan("a")
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	// ary.Value()
// 	//
// 	v0, v1, v2 := "3", "2", "1"
// 	ary = append(ary, &v0)
// 	ary = append(ary, &v1)
// 	ary = append(ary, &v2)
// 	sort.Sort(ary)
// 	//
// 	_, err = Pool().Exec(context.Background(), `drop table if exists xsql_test_string`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	_, err = Pool().Exec(context.Background(), `create table xsql_test_string(tid int,sarry text)`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	{ //normal
// 		_, err = Pool().Exec(context.Background(), `delete from xsql_test_string`)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_string values ($1,$2)`, 1, ary)
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 StringPtrArray
// 		err = Pool().QueryRow(context.Background(), `select sarry from xsql_test_string where tid=$1`, 1).Scan(&ary1)
// 		if err != nil || len(ary1) != 3 {
// 			t.Error(err)
// 			return
// 		}
// 		if !ary1.HavingOne("3") {
// 			t.Error("error")
// 			return
// 		}
// 		if ary1.HavingOne("4") {
// 			t.Error("error")
// 			return
// 		}
// 		if ary1.Join(",") != "1,2,3" {
// 			t.Error("error")
// 		}
// 		if ary1.DbArray() != "{1,2,3}" {
// 			t.Error("error")
// 		}
// 		ary2 := append(ary1, nil, converter.StringPtr("3"), converter.StringPtr(""), converter.StringPtr(" ")).RemoveDuplicate(true, true)
// 		if ary2.DbArray() != "{1,2,3}" {
// 			t.Error("error")
// 		}
// 		ary3 := append(ary1, nil, converter.StringPtr(""), converter.StringPtr(" ")).RemoveEmpty(true)
// 		if ary3.DbArray() != "{1,2,3}" {
// 			t.Error("error")
// 		}
// 	}
// 	{ //string array
// 		_, err = Pool().Exec(context.Background(), `delete from xsql_test_string`)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_string values ($1,$2)`, 1, ary.StrArray())
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 StringPtrArray
// 		err = Pool().QueryRow(context.Background(), `select sarry from xsql_test_string where tid=$1`, 1).Scan(&ary1)
// 		if err != nil || len(ary1) != 3 {
// 			t.Error(err)
// 			return
// 		}
// 		if ary1.Join(",") != "1,2,3" {
// 			t.Error("error")
// 		}
// 	}
// 	{ //nil
// 		var arynil StringPtrArray = nil
// 		res, err := Pool().Exec(context.Background(), `insert into xsql_test_string values ($1,$2)`, 2, arynil)
// 		if err != nil || !res.Insert() {
// 			t.Error(err)
// 			return
// 		}
// 		var ary1 StringPtrArray
// 		err = Pool().QueryRow(context.Background(), `select sarry from xsql_test_string where tid=$1`, 2).Scan(&ary1)
// 		if err != nil || len(ary1) != 0 {
// 			t.Error(err)
// 			return
// 		}
// 	}
// 	{ //join
// 		var ary1 = StringPtrArray{}
// 		ary1 = append(ary1, &v0)
// 		ary1 = append(ary1, nil)
// 		ary1 = append(ary1, &v2)
// 		sort.Sort(ary1)
// 		if ary1.Join(",") != "1,3" {
// 			t.Error(ary1.Join(","))
// 			return
// 		}
// 	}
// 	ary.AsArray()
// 	ary[0] = nil
// 	ary.AsArray()
// }

// func TestIsNilZero(t *testing.T) {
// 	var smap M
// 	if !smap.IsNil() || !smap.IsZero() {
// 		t.Error("error")
// 		return
// 	}
// 	var smapArray MArray
// 	if !smapArray.IsNil() || !smapArray.IsZero() {
// 		t.Error("error")
// 		return
// 	}
// 	var stime Time
// 	if stime.IsNil() || !stime.IsZero() {
// 		t.Error("error")
// 		return
// 	}
// 	var sint IntArray
// 	if !sint.IsNil() || !sint.IsZero() {
// 		t.Error("error")
// 		return
// 	}
// 	var sintPtr IntPtrArray
// 	if !sintPtr.IsNil() || !sintPtr.IsZero() {
// 		t.Error("error")
// 		return
// 	}
// 	var sint64 Int64Array
// 	if !sint64.IsNil() || !sint64.IsZero() {
// 		t.Error("error")
// 		return
// 	}
// 	var sint64Ptr Int64PtrArray
// 	if !sint64Ptr.IsNil() || !sint64Ptr.IsZero() {
// 		t.Error("error")
// 		return
// 	}
// 	var sfloat64 Float64Array
// 	if !sfloat64.IsNil() || !sfloat64.IsZero() {
// 		t.Error("error")
// 		return
// 	}
// 	var sfloat64Ptr Float64PtrArray
// 	if !sfloat64Ptr.IsNil() || !sfloat64Ptr.IsZero() {
// 		t.Error("error")
// 		return
// 	}
// 	var sstr StringArray
// 	if !sstr.IsNil() || !sstr.IsZero() {
// 		t.Error("error")
// 		return
// 	}
// 	var sstrPtr StringPtrArray
// 	if !sstrPtr.IsNil() || !sstrPtr.IsZero() {
// 		t.Error("error")
// 		return
// 	}
// }

// func TestValidFormat(t *testing.T) {
// 	data := attrvalid.M{
// 		"json":      converter.JSON(M{"abc": 123}),
// 		"json_list": converter.JSON([]M{{"abc": 123}}),
// 		"time":      xtime.TimeNow(),
// 		"empty":     "",
// 	}
// 	var smap M
// 	var smapArray MArray
// 	var stime Time
// 	var sint IntArray
// 	var sintPtr IntPtrArray
// 	var sint64 Int64Array
// 	var sint64Ptr Int64PtrArray
// 	var sfloat64 Float64Array
// 	var sfloat64Ptr Float64PtrArray
// 	var sstr StringArray
// 	var sstrPtr StringPtrArray
// 	var etime Time
// 	err := data.ValidFormat(`
// 		json,R|S,L:0;json_list,R|S,L:0;
// 		time,R|I,R:0;
// 		time,R|I,R:0;time,R|I,R:0;
// 		time,R|I,R:0;time,R|I,R:0;
// 		time,R|I,R:0;time,R|I,R:0;
// 		time,R|I,R:0;time,R|I,R:0;
// 		empty,O|I,R:0;
// 		`,
// 		&smap, &smapArray,
// 		&stime,
// 		&sint, &sintPtr,
// 		&sint64, &sint64Ptr,
// 		&sfloat64, &sfloat64Ptr,
// 		&sstr, &sstrPtr,
// 		&etime,
// 	)
// 	if err != nil || stime.Timestamp() < 1 || len(sint) < 1 || len(sintPtr) < 1 {
// 		t.Error(err)
// 		return
// 	}
// 	fmt.Println("-->", stime, sint, sintPtr)
// 	if err = stime.Set(int64(0)); err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if err = stime.Set(Time{}); err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if err = stime.Set(&stime); err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if err = stime.Set(time.Now()); err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	ntime := time.Now()
// 	if err = stime.Set(&ntime); err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if err = stime.Set("xxx"); err == nil {
// 		t.Error(err)
// 		return
// 	}
// }

// func TestValidDecimal(t *testing.T) {
// 	data := attrvalid.M{
// 		"int":    100,
// 		"float":  100.0,
// 		"string": "100.0",
// 	}
// 	var val0, val1, val2 decimal.Decimal
// 	err := data.ValidFormat(`
// 		int,R|I,R:0;
// 		float,R|F,R:0;
// 		string,R|F,R:0;
// 		`,
// 		&val0, &val1, &val2,
// 	)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	fmt.Println("-->", val0, val1, val2)

// 	var args = struct {
// 		A decimal.Decimal `json:"a" valid:"a,r|f,r:0"`
// 	}{}
// 	err = attrvalid.Valid(&args, "#all", "")
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	args.A = decimal.NewFromFloat(-1)
// 	err = attrvalid.Valid(&args, "#all", "")
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	args.A = decimal.NewFromFloat(10)
// 	err = attrvalid.Valid(&args, "#all", "")
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// }

// func TestValid(t *testing.T) {
// 	var err error
// 	errObject := struct {
// 		Map      M      `json:"map" valid:"map,r|s,l:0;"`
// 		MapArray MArray `json:"map_array" valid:"map_array,r|s,l:0;"`
// 	}{}
// 	err = attrvalid.Valid(&errObject, "#all", "")
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	ok1Object := struct {
// 		Map      M      `json:"map" valid:"map,r|s,l:0;"`
// 		MapArray MArray `json:"map_array" valid:"map_array,r|s,l:0;"`
// 		Time     Time   `json:"time" valid:"time,r|i,r:0;"`
// 	}{
// 		Map:      M{},
// 		MapArray: MArray{M{}},
// 		Time:     TimeNow(),
// 	}
// 	err = attrvalid.Valid(&ok1Object, "#all", "")
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	ok2Object := struct {
// 		TID  int64 `json:"tid"  valid:"tid,r|i,r:0;"`
// 		Time Time  `json:"time" valid:"time,r|i,r:0;"`
// 	}{
// 		Time: TimeZero(),
// 	}
// 	err = attrvalid.Valid(&ok2Object, "", "")
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// }

// func TestAs(t *testing.T) {
// 	if len(AsIntArray([]int{1})) != 1 {
// 		t.Error("eror")
// 	}
// 	if len(AsIntPtrArray([]int{1})) != 1 {
// 		t.Error("eror")
// 	}
// 	func() {
// 		defer func() {
// 			recover()
// 		}()
// 		AsIntArray("xxx")
// 	}()
// 	if len(AsInt64Array([]int{1})) != 1 {
// 		t.Error("eror")
// 	}
// 	if len(AsInt64PtrArray([]int{1})) != 1 {
// 		t.Error("eror")
// 	}
// 	func() {
// 		defer func() {
// 			recover()
// 		}()
// 		AsInt64Array("xxx")
// 	}()
// 	if len(AsFloat64Array([]int{1})) != 1 {
// 		t.Error("eror")
// 	}
// 	if len(AsFloat64PtrArray([]int{1})) != 1 {
// 		t.Error("eror")
// 	}
// 	func() {
// 		defer func() {
// 			recover()
// 		}()
// 		AsFloat64Array("xxx")
// 	}()
// 	if len(AsStringArray([]int{1})) != 1 {
// 		t.Error("eror")
// 	}
// 	if len(AsStringPtrArray([]int{1})) != 1 {
// 		t.Error("eror")
// 	}
// 	func() {
// 		defer func() {
// 			recover()
// 		}()
// 		AsStringArray(nil)
// 	}()
// }
