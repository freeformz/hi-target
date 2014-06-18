package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

func main() {

	app := os.Getenv("APP")
	if app == "" {
		log.Fatal("$APP needs to be set to the app name to contact.")
	}

	for {

		//We're even using SSL/TLS!
		conn, err := tls.Dial("tcp", app+".herokuapp.com:443", nil)
		if err != nil {
			log.Println("Dial ", err)
			panic(err)
		}

		cc := httputil.NewClientConn(conn, nil)

		nr, err := http.NewRequest("GET", "https://"+app+".herokuapp.com/hijack", nil)
		if err != nil {
			log.Println("New Request ", err)
			panic(err)
		}

		nr.Header.Add("Connection", "upgrade")
		nr.Header.Add("Upgrade", "tcp")

		resp, err := cc.Do(nr)

		if err != nil {
			log.Println("Do ", err)
			panic(err)
		}

		if resp.StatusCode != 101 {
			log.Fatal("Unable to Upgrade Connection!")
		} else {
			log.Println("Upgraded Connection!")
			hc, rdr := cc.Hijack()

			data, err := ioutil.ReadAll(rdr)
			if err != nil {
				log.Println("ReadAll 1 ", err)
				panic(err)
			}

			log.Println("drain rdr")
			log.Println(string(data))

			data, err = ioutil.ReadAll(hc)
			if err != nil {
				log.Println("ReadAll 2 ", err)
				panic(err)
			}

			log.Println("drain the hijacked connection")
			log.Println(string(data))

			hc.Close()

			time.Sleep(1 * time.Second)
		}
	}

}
