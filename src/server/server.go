package main

import (
  "net"
  "fmt"
  "bufio"
  "strings"
)
// import "encoding/gob"

// TODO make a struct for handling file transfers

type My_Packet struct {
  message string
  contains_file bool
  file_name string
  file []byte
}

func main() {
  listen_message()
  //listen_packet()
  // for i := 0; i < 2; i++ {
  //   go listen()
  // }

}

func listen_packet() {
  fmt.Println("Launching server...")

  // listen on all interfaces
  ln, err := net.Listen("tcp", ":8081")
  check(err, "Server is ready.")

  for {
    conn, _ := ln.Accept()
    check(err, "Accepted connection.")

    go func() {
      // dec := gob.NewDecoder(conn)
      // p := &My_Packet{}
      // dec.Decode(p)
      //
      // fmt.Println("Message: %s", p.message)
      //
      // // send new string back to client
      // conn.Write([]byte("Message Recieved: " + p.message + "\n"))
      conn.Write([]byte("liftoff" + "\n"))
      return
    }()
  }

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
