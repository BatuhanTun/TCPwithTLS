package mysql

import (
	"server/lib"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func Mysql() {
	my := lib.ConfigGet("mysqluser") + ":" + lib.ConfigGet("mysqlpass") + "@tcp(" + lib.ConfigGet("mysqlhost") + ":3306)/" + lib.ConfigGet("mysqldb") + "?charset=utf8&loc=Europe%2FIstanbul"
	if err := orm.RegisterDriver("mysql", orm.DRMySQL); err != nil {
		beego.Error(err)
	}
	if err := orm.RegisterDataBase("default", "mysql", my); err != nil {
		beego.Error(err)
	}
	if err := orm.SetDataBaseTZ("default", time.Now().Location()); err != nil {
		beego.Error(err)
	}
}
