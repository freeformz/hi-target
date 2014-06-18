package main

import (
	"encoding/binary"
	"io"
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

	data := []byte("Hello")
	binary.Write(conn, binary.LittleEndian, uint32(len(data)))
	conn.Write(data)

	var messageSizeU uint32
	err = binary.Read(conn, binary.LittleEndian, &messageSizeU)
	if err != nil {
		log.Fatalf("Error while reading messsage size %d: \n", messageSizeU, err)
	}

	buff := make([]byte, messageSizeU)
	_, err = io.ReadFull(conn, buff)
	if err != nil {
		log.Fatal("Error while reading message: ", err)
	}

	log.Println("Msg from client: " + string(buff))
}

func main() {
	log.SetFlags(0)
	http.HandleFunc("/hijack", serveHijack)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
