package main

import (
    "fmt"
    "log"
    "net"
    "encoding/gob"
)

type P struct {
    M, N int64
}

func main() {
    fmt.Println("start client");
    conn, err := net.Dial("tcp", "127.0.0.1:8081")
    if err != nil {
        log.Fatal("Connection error", err)
    }
    encoder := gob.NewEncoder(conn)
    p := &P{1, 2}
    encoder.Encode(p)
    conn.Close()
    fmt.Println("done");
}

func check_err(err error, message string) {
    if err != nil {
        panic(err)
    }
    fmt.Printf("%s\n", message)
}
