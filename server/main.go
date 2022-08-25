package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	math "math/rand"
	"net"
	lib "server/lib"
	"server/models"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var (
	zar int
)

type mape struct {
	harita    string
	sayaç     int
	tempsayaç int
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

func mysql() {
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

func main() {
	mysql()
	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader
	service := ":3000"
	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	log.Print("server: listening")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		defer conn.Close()
		log.Printf("server: accepted from %s", conn.RemoteAddr())
		tlscon, ok := conn.(*tls.Conn)
		if ok {
			log.Print("ok=true")
			state := tlscon.ConnectionState()
			for _, v := range state.PeerCertificates {
				log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
			}
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
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
			dön = oyun(username, my.Sayaç1, my)
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

func oyun(tür string, sayaç int, my *models.Veri) mape {
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
