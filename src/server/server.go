package main

import (
    "fmt"
    "net"
    "encoding/gob"
)

type P struct {
    M, N int64
}
func handleConnection(conn net.Conn) {
    dec := gob.NewDecoder(conn)
    p := &P{}
    err := dec.Decode(p)
    check_err(err, "Decoding")
    fmt.Printf("Received : %+v\n", p);
    fmt.Printf("M : %d\n", p.M);
    conn.Close()
}

func main() {
    fmt.Println("start");
   ln, err := net.Listen("tcp", ":8081")
    if err != nil {
      check_err(err, "Listening!")
        // handle error
    }
    for {
        conn, err := ln.Accept() // this blocks until connection or error
        if err != nil {
            check_err(err, "Accepting!")
            continue
        }
        go handleConnection(conn) // a goroutine handles conn so that the loop can accept other connections
    }
}

func check_err(err error, message string) {
    if err != nil {
        panic(err)
    }
    fmt.Printf("%s\n", message)
}
