package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func serveHijack(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serveHijack")
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Don't forget to close the connection:
	defer conn.Close()

	log.Println("Upgrade: " + r.Header.Get("Upgrade"))
	log.Println("Connection: " + r.Header.Get("Connection"))
	bufrw.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
	bufrw.WriteString("Upgrade: " + r.Header.Get("Upgrade") + "\r\n")
	bufrw.WriteString("Connection: " + r.Header.Get("Connection") + "\r\n")
	bufrw.WriteString("\r\n")
	bufrw.Flush()
	bufrw.WriteString("Now we're speaking raw TCP!\r\n")
	bufrw.Flush()
}

func main() {
	http.HandleFunc("/hijack", serveHijack)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
