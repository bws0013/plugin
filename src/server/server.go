package main

import (
  "net"
  "fmt"
  "io/ioutil"
  "os"
)
import "encoding/gob"

// TODO add permissions as an element passed in the packet

type my_packet struct {
  Message string
  Contains_file bool
  File_name string
  File []byte
}

func main() {
  fmt.Println("start");
  ln, err := net.Listen("tcp", ":8081")
  check_err(err, "Listening!")

  for {
    conn, err := ln.Accept() // this blocks until connection or error
    if err != nil {
        fmt.Println("This connection needs a tissue, skipping!")
        continue
    }
    go listen_packet(conn) // a goroutine handles conn so that the loop can accept other connections
  }
}

func listen_packet(conn net.Conn) {

  // dec := gob.NewDecoder(conn)
  //   p := &P{}
  //   dec.Decode(p)
  //   fmt.Printf("Received : %+v", p);
  //   conn.Close()

  dec := gob.NewDecoder(conn)
  p := &my_packet{}
  err := dec.Decode(p)
  check_err(err, "No problems on read in")

  conn.Close()

  fmt.Printf("Message: %s\n", p.File_name)

  if p.Contains_file && p.File != nil {
    create_file(p.File_name, p.File)
  } else {
    fmt.Println("No file detected!")
  }


  if dec != nil {
    fmt.Printf("Client disconnected.\n")
    return
  }

  // conn.Write([]byte("liftoff" + "\n"))

      // dec := gob.NewDecoder(conn)
      // p := &my_packet{}
      // dec.Decode(p)
      //
      // fmt.Println("Message: %s", p.message)
      //
      // // send new string back to client
      // conn.Write([]byte("Message Recieved: " + p.message + "\n"))

}

func create_file(name string, data []byte) {
  permissions := os.FileMode(0644)
  full_name := "./../../storage/recieved/" + name
  err := ioutil.WriteFile(full_name, data, permissions)
  check_err(err, "File created!")
}

func check_err(err error, message string) {
    if err != nil {
        panic(err)
    }
    fmt.Printf("%s\n", message)
}
