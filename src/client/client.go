package main

import (
  "net"
  "fmt"
  "os"
  "io/ioutil"
  "encoding/gob"
  "bufio"
  "math/rand"
  "time"
)

// TODO incorporate this https://medium.com/@lhartikk/a-blockchain-in-200-lines-of-code-963cc1cc0e54

type my_packet struct {
  Key int32
  Bet float32
  Res int
}

func main() {
  rand.Seed(time.Now().UTC().UnixNano())

  k := rand.Int31()
  b := float32(10)
  r := rand.Intn(2)

  fmt.Println(my_packet{k, b, r})
}

func dial_server_packet(packet my_packet) {
  conn, err := net.Dial("tcp", "127.0.0.1:8081")

  if err != nil {
    fmt.Println("Unable to send!")
    return
  }
  // fmt.Println(packet)
  encoder := gob.NewEncoder(conn)
  err = encoder.Encode(&packet)

  message, _ := bufio.NewReader(conn).ReadString('\n')
  fmt.Println("Message from server: " + message)
  check_err(err, "everything is fine")
  conn.Close()

}

func check_err(err error, message string) {
    if err != nil {
      panic(err)
    }
    if len(message) != 0 {
      fmt.Printf("%s\n", message)
    }
}
