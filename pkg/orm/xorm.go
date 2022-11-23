package orm

//
//import (
//	"github.com/isyscore/isc-gobase/extend/orm"
//	"github.com/isyscore/isc-gobase/isc"
//	_const "github.com/isyscore/isc-tracer/internal/const"
//	"github.com/isyscore/isc-tracer/internal/trace"
//	"github.com/isyscore/isc-tracer/pkg"
//	"gorm.io/gorm"
//)
//
//package orm
//
//import (
//"context"
//"encoding/json"
//"github.com/isyscore/isc-gobase/extend/orm"
//"github.com/isyscore/isc-gobase/isc"
//_const "github.com/isyscore/isc-tracer/internal/const"
//"github.com/isyscore/isc-tracer/internal/trace"
//"github.com/isyscore/isc-tracer/pkg"
//"gorm.io/gorm"
//)
//type GobaseXormHook struct {
//}
//
//const (
//	traceContextKey = "gobase-gorm-trace-key"
//
//	// 自定义事件名称
//	_eventBeforeCreate = "gobase-gorm-collector-event:before_create"
//	_eventAfterCreate  = "gobase-gorm-collector-event:after_create"
//	_eventBeforeUpdate = "gobase-gorm-collector-event:before_update"
//	_eventAfterUpdate  = "gobase-gorm-collector-event:after_update"
//	_eventBeforeQuery  = "gobase-gorm-collector-event:before_query"
//	_eventAfterQuery   = "gobase-gorm-collector-event:after_query"
//	_eventBeforeDelete = "gobase-gorm-collector-event:before_delete"
//	_eventAfterDelete  = "gobase-gorm-collector-event:after_delete"
//	_eventBeforeRow    = "gobase-gorm-collector-event:before_row"
//	_eventAfterRow     = "gobase-gorm-collector-event:after_row"
//	_eventBeforeRaw    = "gobase-gorm-collector-event:before_raw"
//	_eventAfterRaw     = "gobase-gorm-collector-event:after_raw"
//
//	// 自定义 span 的操作名称
//	_opCreate = "insert"
//	_opUpdate = "update"
//	_opQuery  = "select"
//	_opDelete = "delete"
//	_opRow    = "row"
//	_opRaw    = "execute"
//)
//
//// 实现 gorm 插件所需方法
//func (i *GobaseGormHook) Name() string {
//	return "gobase_gorm_plugin"
//}
//
//// 实现 gorm 插件所需方法
//func (i *GobaseGormHook) Initialize(db *gorm.DB) (err error) {
//	// 在 gorm 中注册各种回调事件
//	for _, e := range []error{
//		db.Callback().Create().Before("gorm:create").Register(_eventBeforeCreate, beforeCreate),
//		db.Callback().Create().After("gorm:create").Register(_eventAfterCreate, after),
//		db.Callback().Update().Before("gorm:update").Register(_eventBeforeUpdate, beforeUpdate),
//		db.Callback().Update().After("gorm:update").Register(_eventAfterUpdate, after),
//		db.Callback().Query().Before("gorm:query").Register(_eventBeforeQuery, beforeQuery),
//		db.Callback().Query().After("gorm:query").Register(_eventAfterQuery, after),
//		db.Callback().Delete().Before("gorm:delete").Register(_eventBeforeDelete, beforeDelete),
//		db.Callback().Delete().After("gorm:delete").Register(_eventAfterDelete, after),
//		db.Callback().Row().Before("gorm:row").Register(_eventBeforeRow, beforeRow),
//		db.Callback().Row().After("gorm:row").Register(_eventAfterRow, after),
//		db.Callback().Raw().Before("gorm:raw").Register(_eventBeforeRaw, beforeRaw),
//		db.Callback().Raw().After("gorm:raw").Register(_eventAfterRaw, after),
//	} {
//		if e != nil {
//			return e
//		}
//	}
//	return
//}
//
//// 注册各种前置事件时，对应的事件方法
//func _injectBefore(db *gorm.DB, op string) {
//	if !pkg.DatabaseTraceSwitch {
//		return
//	}
//
//	if db == nil {
//		return
//	}
//
//	if db.Statement == nil || db.Statement.Context == nil {
//		db.Logger.Error(context.TODO(), "未定义 db.Statement 或 db.Statement.Context")
//		return
//	}
//
//	tracer := pkg.ServerStartTrace(_const.MYSQL, "gorm:"+op)
//	db.InstanceSet(traceContextKey, tracer)
//}
//
//// 注册后置事件时，对应的事件方法
//func after(db *gorm.DB) {
//	if !pkg.DatabaseTraceSwitch {
//		return
//	}
//
//	if db == nil {
//		return
//	}
//
//	if db.Statement == nil || db.Statement.Context == nil {
//		db.Logger.Error(context.TODO(), "未定义 db.Statement 或 db.Statement.Context")
//		return
//	}
//
//	_tracer, isExist := db.InstanceGet(traceContextKey)
//	if !isExist || _tracer == nil {
//		return
//	}
//
//	tracer, ok := _tracer.(*trace.Tracer)
//	if !ok || tracer == nil {
//		return
//	}
//
//	resultMap := map[string]any{}
//	result := _const.OK
//
//	b, err := json.Marshal(db.Statement.Vars)
//	if err != nil {
//		resultMap["err"] = err.Error()
//		result = _const.ERROR
//	}
//
//	if db.Error != nil {
//		resultMap["err"] = db.Error.Error()
//		result = _const.ERROR
//	}
//	resultMap["sql"] = db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
//	resultMap["table"] = db.Statement.Table
//	resultMap["query"] = db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
//	resultMap["parameters"] = string(b)
//
//	// todo 返回大小，暂时设置为0
//	pkg.ServerEndTrace(tracer, 0, result, isc.ToJsonString(resultMap))
//}
//
//func beforeCreate(db *gorm.DB) {
//	_injectBefore(db, _opCreate)
//}
//
//func beforeUpdate(db *gorm.DB) {
//	_injectBefore(db, _opUpdate)
//}
//
//func beforeQuery(db *gorm.DB) {
//	_injectBefore(db, _opQuery)
//}
//
//func beforeDelete(db *gorm.DB) {
//	_injectBefore(db, _opDelete)
//}
//
//func beforeRow(db *gorm.DB) {
//	_injectBefore(db, _opRow)
//}
//
//func beforeRaw(db *gorm.DB) {
//	_injectBefore(db, _opRaw)
//}
