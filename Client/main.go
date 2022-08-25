package main

import (
	"crypto/tls"
	"fmt"
	"log"
)

func main() {
	cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", "localhost:3000", &config)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var username string
	data := make([]byte, 1024)
	key := "0"
	datastring, err := conn.Read(data)
	if err != nil {
		panic(err)
	}
	veri := string(data[:datastring])
	for {
		if veri == "T" {
			fmt.Println("Tavşan sizsiniz.")
			key = "1"
			break
		} else if veri == "K" {
			fmt.Println("Kaplumbağa sizsiniz.")
			key = "0"
			break
		} else if veri == "A" {
			fmt.Println("Çıkış yapıldı.")
			break
		} else {
			fmt.Println("Lütfen geçerli bir seçim yapınız veya çıkış için A' ya basınız!")
		}
	}

	flag := true
	var flaggame string
	for {
		data = make([]byte, 1024)
		datastring, err := conn.Read(data)
		if err != nil {
			panic(err)
		}
		veri := string(data[:datastring])
		if key == veri {
			flaggame = "true"
		} else {
			flaggame = "false"
		}
		conn.Write([]byte(flaggame))
		if key == veri {
			flag = true

			data = make([]byte, 1024)
			datastring, err = conn.Read(data)
			if err != nil {
				panic(err)
			}
			veri = string(data[:datastring])

			if veri == "Tebrikler tavşan kazandı!" || veri == "Tebrikler kaplumbağa kazandı!" || veri == "finish" {
				fmt.Println("Üzgünüz kaybettiniz!")
				break
			}
			fmt.Println(veri)

			data = make([]byte, 1024)
			datastring, err = conn.Read(data)
			if err != nil {
				panic(err)
			}
			veri = string(data[:datastring])
			if veri == "A" {
				fmt.Println("Çıkış yapıldı.")
				break
			}
			fmt.Println(veri)

			fmt.Scanf("%s\n", &username)
			conn.Write([]byte(username))

			data = make([]byte, 1024)
			datastring, err = conn.Read(data)
			if err != nil {
				panic(err)
			}
			veri = string(data[:datastring])
			if veri == "A" {
				fmt.Println("Çıkış yapıldı.")
				break
			}
			fmt.Println(veri)

			if veri == "Lütfen geçerli bir komut giriniz!" || veri == "Lütfen sıradaki oyuncu zar atsın!" {
				flag = false
			}
		}
		if veri == "Tebrikler tavşan kazandı!" || veri == "Tebrikler kaplumbağa kazandı!" || veri == "finish" {
			break
		}
		if veri == "A" {
			fmt.Println("Diğer oyuncu oyundan çıktı kazanan sizsiniz tebrikler.")
			break
		}
		if flag {
			fmt.Println("Lütfen sıradaki oyuncuyu bekleyiniz")
			flag = false
		}
	}
}
