package main

import (
	"log"
	"net/http"
	"os"
)

func serveHijack(w http.ResponseWriter, r *http.Request) {
	upgradeHeader := r.Header.Get("Upgrade")
	if upgradeHeader == "" {
		http.Error(w, "Upgrade header not provided so not hijacking", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Connection") != "Upgrade" {
		http.Error(w, "Connection header must be 'Upgrade' in order to upgrade the connection", http.StatusBadRequest)
		return
	}

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

	log.Println("Trying to upgrade hijacked connection to " + upgradeHeader + " !")

	defer conn.Close()

	// Send the right response...
	bufrw.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
	bufrw.WriteString("Upgrade: " + r.Header.Get("Upgrade") + "\r\n")
	bufrw.WriteString("Connection: " + r.Header.Get("Connection") + "\r\n")
	bufrw.WriteString("\r\n")
	bufrw.Flush()

	// After this we're just a TCP connection
	bufrw.WriteString("Now we're speaking raw TCP!\r\n")
	bufrw.Flush()
}

func main() {
	log.SetFlags(0)
	http.HandleFunc("/hijack", serveHijack)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
