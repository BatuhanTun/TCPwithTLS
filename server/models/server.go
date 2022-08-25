package models

import (
	"github.com/astaxie/beego/orm"
)

type Veri struct {
	Id       int64  `orm:"auto"`
	Sayaç1   int    `orm:"int"`
	Tavkonum int    `orm:"int"`
	Kapkonum int    `orm:"int"`
	Konum    int    `orm:"int"`
	Map      string `orm:"size(255)"`
	Finish   string `orm:"size(255)"`
	Seçim    string `orm:"size(255)"`
}

func init() {
	orm.RegisterModel(
		new(Veri),
	)
}
