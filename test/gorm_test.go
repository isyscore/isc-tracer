package test

import (
	"github.com/isyscore/isc-gobase/extend/orm"
	_ "github.com/isyscore/isc-tracer"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"testing"
	"time"
)

// 使用环境变量：base.profiles.active=database
func TestGorm(t *testing.T) {
	c := &gorm.Config{
		SkipDefaultTransaction: true,
		//Logger:                 DBlogger.Default.LogMode(DBlogger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}}
	db, _ := orm.NewGormDbWitConfig(c)

	// 删除表
	db.Exec("drop table isc_demo.gobase_demo")
	//
	//// 测试异常
	db.Exec("drop table isc_demo.gobase_demoxxx")
	//
	////新增表
	db.Exec("CREATE TABLE isc_demo.gobase_demo(\n" +
		"  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',\n" +
		"  `name` char(20) NOT NULL COMMENT '名字',\n" +
		"  `age` INT NOT NULL COMMENT '年龄',\n" +
		"  `address` char(20) NOT NULL COMMENT '名字',\n" +
		"  \n" +
		"  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',\n" +
		"  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',\n" +
		"\n" +
		"  PRIMARY KEY (`id`)\n" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='测试表'")

	db.Exec("INSERT INTO `gobase_demo` (`name`,`age`,`address`) VALUES ('xx', 12, '杭州')")
	// 新增x
	db.Create(&GobaseDemo{Name: "zhou", Age: 18, Address: "杭州"})
	db.Create(&GobaseDemo{Name: "zhou", Age: 11, Address: "杭州2"})

	// 查询：一行
	var demo GobaseDemo
	db.Where("name=?", "zhou").First(&demo)
	t.Log(demo)
	// 测试sql
	dd, _ := db.DB()
	query, err := dd.Query("select * from isc_demo.gobase_demo")
	if err == nil {
		var list []GobaseDemo
		for query.Next() {
			var a GobaseDemo
			query.Scan(&a)
			list = append(list, a)
		}
		t.Log(len(list))
	}
	// 测试参数
	query, err = dd.Query("select * from isc_demo.gobase_demo where age = ?", 18)
	if err == nil {
		var list []GobaseDemo
		for query.Next() {
			var a GobaseDemo
			query.Scan(&a)
			list = append(list, a)
		}
		t.Log(len(list))
	}

	// 测试异常
	// 根据目前的mysql驱动, query只会执行一次
	// 由于被丢弃, 这条trace会丢失
	query, err = dd.Query("select * from isc_demo.gobase_demoxxx where id = ?", 23)

	if err == nil {
		var list []GobaseDemo
		for query.Next() {
			var a GobaseDemo
			query.Scan(&a)
			list = append(list, a)
		}
		t.Log(len(list))
	}
	time.Sleep(time.Second * 2)
}

type GobaseDemo struct {
	Id      uint64
	Name    string
	Age     int
	Address string
}
