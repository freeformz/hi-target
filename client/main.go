package main

import (
	"crypto/tls"
	"encoding/binary"
	"io"
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
			log.Println("Hijacked Connection!")

			var messageSizeU uint32
			err := binary.Read(rdr, binary.LittleEndian, &messageSizeU)
			if err != nil {
				log.Fatal("Unable to read message size: ", err)
			}

			log.Println("Read Message Length: ", messageSizeU)

			buff := make([]byte, messageSizeU)
			_, err = io.ReadFull(rdr, buff)
			if err != nil {
				log.Fatal("Unable to read message: ", err)
			}

			log.Println("Message from Server: " + string(buff))

			msg := []byte("Ohai There!. Do you like length prefixed communication?")
			binary.Write(hc, binary.LittleEndian, uint32(len(msg)))
			hc.Write(msg)

			hc.Close()

			time.Sleep(1 * time.Second)
		}
	}

}
