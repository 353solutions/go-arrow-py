package main

import (
	"fmt"
	"time"

	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/memory"
)

func recExample() {
	pool := memory.NewGoAllocator()

	schema := arrow.NewSchema(
		[]arrow.Field{
			arrow.Field{Name: "f1-i32", Type: arrow.PrimitiveTypes.Int32},
			arrow.Field{Name: "f2-f64", Type: arrow.PrimitiveTypes.Float64},
		},
		nil,
	)
	b := array.NewRecordBuilder(pool, schema)
	defer b.Release()

	ib := b.Field(0).(*array.Int32Builder)
	fb := b.Field(1).(*array.Float64Builder)
	for i := 0; i < 10; i++ {
		ib.Append(int32(i * 7))
		fb.Append(float64(i * 10))
	}

	r := b.NewRecord()
	defer r.Release()

	tbl := array.NewTableFromRecords(schema, []array.Record{r})
	printTable(tbl)
}

func printTable(tbl array.Table) {
	tr := array.NewTableReader(tbl, 5)

	for n := 0; tr.Next(); n++ {
		rec := tr.Record()
		for i, col := range rec.Columns() {
			fmt.Printf("rec %d (%q): %+v\n", n, rec.ColumnName(i), col)
		}
	}

}

func colExample(pool memory.Allocator) {

	f1 := arrow.Field{Name: "f1-i64", Type: arrow.PrimitiveTypes.Int64}
	f2 := arrow.Field{Name: "f2-f64", Type: arrow.PrimitiveTypes.Float64}
	schema := arrow.NewSchema(
		[]arrow.Field{f1, f2},
		nil,
	)
	ib := array.NewInt64Builder(pool)
	fb := array.NewFloat64Builder(pool)

	for i := 0; i < 10; i++ {
		ib.Append(int64(i * 3))
		fb.Append(float64(i * 7))
	}

	ic := array.NewChunked(f1.Type, []array.Interface{ib.NewArray()})
	icol := array.NewColumn(f1, ic)
	fc := array.NewChunked(f2.Type, []array.Interface{fb.NewArray()})
	fcol := array.NewColumn(f2, fc)

	tbl := array.NewTable(schema, []array.Column{*icol, *fcol}, -1)
	printTable(tbl)
}

func shmExample() error {
	sh, err := NewSharedMemory("lassie", 1024)
	if err != nil {
		return err
	}
	defer sh.Close(false)

	copy(sh.Data(), []byte("hello\n"))

	return nil
}

func main() {
	// recExample()
	// shmExample()
	// pool := memory.NewGoAllocator()

	pool, err := NewShmAlloactor("shm-allocator", 1<<20)
	if err != nil {
		panic(err)
	}
	defer pool.Close(true)

	colExample(pool)
	fmt.Println("OK")
	for {
		time.Sleep(time.Second)
	}
}
