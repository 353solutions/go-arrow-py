package main

import (
	"fmt"

	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/memory"
)

func main() {
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
	tr := array.NewTableReader(tbl, 5)

	for n := 0; tr.Next(); n++ {
		rec := tr.Record()
		for i, col := range rec.Columns() {
			fmt.Printf("rec %d (%q): %+v\n", n, rec.ColumnName(i), col)
		}
	}

	fmt.Printf("%+v\n", tbl)
	fmt.Println("OK")

}
