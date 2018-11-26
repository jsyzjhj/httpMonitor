package models

import (
	"github.com/astaxie/beego/orm"
)

func GetTableList(tableName string, filter map[string]interface{}, order string, limit, start int, list interface{}) (cnt int64) {
	o := orm.NewOrm().QueryTable(tableName)
	for k, v := range filter {
		o = o.Filter(k, v)
	}
	if order != "" {
		o = o.OrderBy(order)
	}
	o.Limit(limit, start).All(list)
	cnt, _ = o.Count()
	return
}

func GetNum(tableName string, filter map[string]interface{}) (cnt int64) {
	o := orm.NewOrm().QueryTable(tableName)
	for k, v := range filter {
		o = o.Filter(k, v)
	}
	cnt, _ = o.Count()
	return
}

func Delete(model interface{}) (int64, error) {
	return orm.NewOrm().Delete(model)
}

func Read(model interface{}) error {
	return orm.NewOrm().Read(model)
}

func Update(model interface{}) (int64, error) {
	return orm.NewOrm().Update(model)
}

func Insert(model interface{}) (int64, error) {
	return orm.NewOrm().Insert(model)
}
