package main

import (
  "net"
  "fmt"
  "bufio"
  "strings"
)
import "encoding/gob"

// TODO make a struct for handling file transfers

type my_packet struct {
  Message string
  Contains_file bool
  File_name string
  File []byte
}

func main() {
  //listen_message()
  //listen_packet()
  // for i := 0; i < 2; i++ {
  //   go listen()
  // }

  fmt.Println("Launching server...")
  ln, err := net.Listen("tcp", ":8081")
  check(err, "Server is ready.")
  for {
    conn, err := ln.Accept()
    check(err, "Accepted connection.")
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
  check(err, "No problems on read in")

  fmt.Println("Message: %s", p.File_name)


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
  check(err, "Server is ready.")

  // run loop forever (or until ctrl-c)
  for {
    conn, _ := ln.Accept()
    check(err, "Accepted connection.")

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

func check(err error, message string) {
    if err != nil {
        panic(err)
    }
    fmt.Printf("%s\n", message)
}
