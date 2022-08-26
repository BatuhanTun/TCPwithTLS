package game

import (
	"encoding/json"
	"io/ioutil"
	"server/models"
	"strconv"

	"github.com/astaxie/beego/orm"
)

var (
	zar int
)

func Oyun(tür string, sayaç int, my *models.Veri) mape {
	var i int
	var dön mape
	dön.tempsayaç = 1
	if tür == "T" {
		dön.sayaç = 0
	} else if tür == "K" {
		dön.sayaç = 1
	}
	if tür == "T" {
		tavşan()
	} else if tür == "K" {
		kaplumbağa()
	}
	orm.NewOrm().QueryTable("Veri").One(my)
	dön.harita = ""
	for i = 0; i < 30; i++ {
		if i == my.Tavkonum && i == my.Kapkonum {
			dön.harita += "KT"
		} else if i == my.Tavkonum {
			dön.harita += " T "
		} else if i == my.Kapkonum {
			dön.harita += " K "
		} else {
			dön.harita += " - "
		}
	}
	if tür == "T" {
		my.Konum = my.Tavkonum
	} else if tür == "K" {
		my.Konum = my.Kapkonum
	}
	if my.Konum == 29 {
		if tür == "T" {
			dön.harita = "Tebrikler tavşan kazandı!"
			dön.tempsayaç = 2
		} else if tür == "K" {
			dön.harita = "Tebrikler kaplumbağa kazandı!"
			dön.tempsayaç = 2
		}
	}
	if _, err := orm.NewOrm().Update(my); err != nil {
		panic(err)
	}
	return dön
}

func tavşan() {
	my := &models.Veri{}
	orm.NewOrm().QueryTable("Veri").One(my)
	stringzar := strconv.Itoa(zar)
	adım := Json("tavşan" + stringzar)
	ekle, _ := strconv.Atoi(adım)
	my.Tavkonum = my.Tavkonum + ekle
	if my.Tavkonum < 0 {
		my.Tavkonum = 0
	}
	if my.Tavkonum > 29 {
		my.Tavkonum = 29
	}
	if _, err := orm.NewOrm().Update(my); err != nil {
		panic(err)
	}

}

func kaplumbağa() {
	my := &models.Veri{}
	orm.NewOrm().QueryTable("Veri").One(my)
	stringzar := strconv.Itoa(zar)
	adım := Json("kaplumbağa" + stringzar)
	ekle, _ := strconv.Atoi(adım)
	my.Kapkonum = my.Kapkonum + ekle
	if my.Kapkonum < 0 {
		my.Kapkonum = 0
	}
	if my.Kapkonum > 29 {
		my.Kapkonum = 29
	}
	if _, err := orm.NewOrm().Update(my); err != nil {
		panic(err)
	}
}

func Json(name string) string {
	file, err := ioutil.ReadFile("bilgi.json")
	if err != nil {
		return "Bulunamadı!"
	}
	var data map[string]string = make(map[string]string)
	if err := json.Unmarshal(file, &data); err != nil {
		panic(err)
	}
	result, ok := data[name]
	if !ok {
		return "Bulunamadı!"
	}
	return result
}
