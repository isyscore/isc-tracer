package test

import (
	"database/sql"
	"github.com/isyscore/isc-gobase/extend/orm"
	_ "github.com/isyscore/isc-tracer"
	"github.com/isyscore/isc-tracer/internal/trace"
	"testing"
)

func TestGorm(t *testing.T) {
	trace.OsTraceSwitch = true
	trace.DatabaseTraceSwitch = true

	gormDb, err := orm.NewGormDb()
	if err != nil {
		t.Fatal(err)
	}
	db, err := gormDb.DB()
	if err != nil {
		t.Fatal(err)
	}
	rows, err := db.Query("select * from dmc_device limit 1")
	if err != nil {
		t.Fatal(err)
	}

	var results []string
	var tmp string

	// 获取字段名称
	tmp = ""
	cols, _ := rows.Columns()
	for i := range cols {
		tmp += cols[i] + ","
	}
	results = append(results, tmp)

	// 根据字段数量，指定查询scan的参数
	values := make([]sql.RawBytes, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(scans...); err != nil {
			t.Fatal(err)
		}
		// 组装，暂定以逗号隔开
		tmp = ""
		for j := range values {
			tmp += string(values[j]) + ","
		}
		results = append(results, tmp)
	}
	t.Log(results)

}
