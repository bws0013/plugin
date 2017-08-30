package main

import (
  "net"
  "fmt"
  "sync"
)
import "encoding/gob"

// TODO incorporate this https://medium.com/@lhartikk/a-blockchain-in-200-lines-of-code-963cc1cc0e54

type my_packet struct {
  Key int32
  Bet float32
  Res int
}

var bet_map sync.Map

func main() {
  // execute_commands("ls && ls -alh")

  fmt.Println("start");
  ln, err := net.Listen("tcp", ":8081")
  // check_err(err, "Listening!")
  if err != nil {
    fmt.Println("Error at listen")
    return
  }

  total_connections := 0

  for {
    if total_connections > 2 {
      break
    }
    conn, err := ln.Accept() // this blocks until connection or error
    if err != nil {
        fmt.Println("This connection needs a tissue, skipping!")
        continue
    }
    go listen_packet(conn) // a goroutine handles conn so that the loop can accept other connections
    total_connections++
  }

  // TODO: Figure out this range stuff
  // for x, _ := bet_map.Range() {
  //
  //
  // }

}



func listen_packet(conn net.Conn) {

  dec := gob.NewDecoder(conn)
  p := &my_packet{}
  err := dec.Decode(p)

  if err != nil { fmt.Println("Tell me about it") }

  bet_map.Store(p.Key, p)

  conn.Write([]byte("liftoff"))

  conn.Close()

  if dec != nil {
    fmt.Printf("Client disconnected.\n")
    return
  }

}
