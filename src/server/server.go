package main

import (
  "net"
  "fmt"
  "bufio"
  "strings"
  "io/ioutil"
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

  fmt.Printf("Message: %s\n", p.File_name)

  conn.Close()

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

func listen_message() {
  fmt.Println("Launching server...")

  // listen on all interfaces
  ln, err := net.Listen("tcp", ":8081")
  check_err(err, "Server is ready.")

  // run loop forever (or until ctrl-c)
  for {
    conn, _ := ln.Accept()
    check_err(err, "Accepted connection.")

    go func() {
      // will listen for message to process ending in newline (\n)
      buf := bufio.NewReader(conn)

      for {
        message, err := buf.ReadString('\n')
        if err != nil {
          fmt.Printf("Client disconnected.\n")
          break
        }

        // output message received
        fmt.Print("Message Received:", string(message))
        // sample process for string received
        newmessage := strings.ToUpper(message)
        // send new string back to client
        conn.Write([]byte(newmessage + "\n"))
      }
    }()
  }
}

func create_file(name string, data []byte) {
  //permissions := int(0644)
  full_name := "./../../storage/recieved/" + name
  err := ioutil.WriteFile(full_name, data, 0644)
  check_err(err, "File created!")
}

func check_err(err error, message string) {
    if err != nil {
        panic(err)
    }
    fmt.Printf("%s\n", message)
}
