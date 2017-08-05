package main

import "net"
import "fmt"
import "bufio"
import "strings" // only needed below for sample processing

func main() {

  listen()
  // for i := 0; i < 2; i++ {
  //   go listen()
  // }

}

func listen() {
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
