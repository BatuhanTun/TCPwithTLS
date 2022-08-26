package game

import (
	"fmt"
	math "math/rand"
	"net"
	"server/models"

	"github.com/astaxie/beego/orm"
)

type mape struct {
	harita    string
	sayaç     int
	tempsayaç int
}

func HandleClient(conn net.Conn) {
	for {
		defer conn.Close()
		my := &models.Veri{}
		orm.NewOrm().QueryTable("Veri").One(my)
		//my.Kapkonum = 0
		//my.Tavkonum = 0
		//my.Map = "."
		data := make([]byte, 1024)
		var dön mape
		if my.Seçim == "false" {
			my.Sayaç1 = math.Intn(2)
			if my.Sayaç1 == 1 {
				my.Seçim = "T"
				conn.Write([]byte("T"))
			} else {
				my.Seçim = "K"
				conn.Write([]byte("K"))
			}
		} else if my.Seçim == "T" {
			conn.Write([]byte("K"))
			my.Seçim = "false"
		} else if my.Seçim == "K" {
			conn.Write([]byte("T"))
			my.Seçim = "false"
		}
		if _, err := orm.NewOrm().Update(my); err != nil {
			panic(err)
		}

		for my.Tavkonum != 30 || my.Kapkonum != 30 {
			orm.NewOrm().QueryTable("Veri").One(my)
			if my.Sayaç1 == 1 {
				if my.Finish != "1" {
					conn.Write([]byte("1"))
				} else {
					orm.NewOrm().QueryTable("Veri").One(my)
					my.Finish = "0"
					if _, err := orm.NewOrm().Update(my); err != nil {
						panic(err)
					}
					conn.Write([]byte("finish"))
					conn.Close()
					break
				}
				dataname, err := conn.Read(data)
				if err != nil {
					fmt.Println(err)
					return
				}
				flag := string(data[:dataname])
				if flag != "true" {
					continue
				}
			} else if my.Sayaç1 == 0 {
				if my.Finish != "1" {
					conn.Write([]byte("0"))
				} else {
					orm.NewOrm().QueryTable("Veri").One(my)
					my.Finish = "0"
					if _, err := orm.NewOrm().Update(my); err != nil {
						panic(err)
					}
					conn.Write([]byte("finish"))
					conn.Close()
					break
				}
				dataname, err := conn.Read(data)
				if err != nil {
					fmt.Println(err)
					return
				}
				flag := string(data[:dataname])
				if flag != "true" {
					continue
				}
			}
			if my.Sayaç1 == 2 {
				my.Seçim = "false"
				if _, err := orm.NewOrm().Update(my); err != nil {
					panic(err)
				}
				conn.Write([]byte("A"))
				conn.Close()
				return
			}
			if my.Sayaç1 == 1 {
				conn.Write([]byte(my.Map))
				conn.Write([]byte("Sıra Tavşanda Lütfen T ye basıp zar atınız"))
			} else {
				conn.Write([]byte(my.Map))
				conn.Write([]byte("Sıra Kaplumbağada Lütfen K ye basıp zar atınız"))
			}
			dataname, err := conn.Read(data)
			if err != nil {
				fmt.Println(err)
				return
			}
			username := string(data[:dataname])
			if username == "A" {
				my.Tavkonum = 0
				my.Kapkonum = 0
				my.Map = "."
				my.Finish = "0"
				my.Sayaç1 = 2
				my.Seçim = "false"
				if _, err := orm.NewOrm().Update(my); err != nil {
					panic(err)
				}
				conn.Write([]byte("A"))
				conn.Close()
				return
			}
			if username != "T" && username != "K" && username != "A" {
				conn.Write([]byte("Lütfen geçerli bir komut giriniz!"))
				continue
			}
			if (my.Sayaç1 == 1 && username != "T") || (my.Sayaç1 == 0 && username != "K") {
				conn.Write([]byte("Lütfen sıradaki oyuncu zar atsın!"))
				continue
			}
			zar = math.Intn(5) + 1
			dön = Oyun(username, my.Sayaç1, my)
			my.Sayaç1 = dön.sayaç
			my.Map = dön.harita
			if _, err := orm.NewOrm().Update(my); err != nil {
				panic(err)
			}
			conn.Write([]byte(dön.harita))
			if dön.tempsayaç == 2 {
				my.Tavkonum = 0
				my.Kapkonum = 0
				my.Map = "."
				my.Finish = "1"
				if _, err := orm.NewOrm().Update(my); err != nil {
					panic(err)
				}
				conn.Write([]byte("A"))
				conn.Close()
				break
			}
		}
		my.Finish = "0"
		if _, err := orm.NewOrm().Update(my); err != nil {
			panic(err)
		}
		break
	}
}
