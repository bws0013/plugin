package main

import (
  "net"
  "fmt"
  "io/ioutil"
  "os"
  "os/exec"
)
import "encoding/gob"

// TODO add multiread capability

type my_packet struct {
  Current_time string
  Message string
  Contains_file bool
  File_name string
  Permissions uint
  File []byte
}

func main() {
  execute_commands("ls -alh && ls")

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

  fmt.Println(p)

  conn.Write([]byte("liftoff"))

  conn.Close()

  fmt.Printf("Message: %s\n", p.File_name)

  if p.Contains_file && p.File != nil {
    create_file(p.File_name, p.File, p.Permissions)
  } else {
    fmt.Println("No file detected!")
  }

  if dec != nil {
    fmt.Printf("Client disconnected.\n")
    return
  }

}

// Bad way to execute any command passed in.
// TODO figure out a better way to pass arbitrary commands
func execute_commands(message string) {
  path := "./../../storage/scripts/"
  name := path + "exec1.sh"
  permissions := uint(511)
  message_in_byte_form := []byte(message)

  err := ioutil.WriteFile(name, message_in_byte_form, os.FileMode(permissions))
  check_err(err, "Creating Script")

  out, err := exec.Command("/bin/sh", name).Output()
  check_err(err, "Executing Script")
  fmt.Printf("%s\n", out)
}

// Clear the directory where executables will be stored
func clear_execution_dir() { return }

func create_file(name string, data []byte, permissions uint) {
  full_name := "./../../storage/recieved/" + name
  err := ioutil.WriteFile(full_name, data, os.FileMode(permissions))
  check_err(err, "File created!")
}

func check_err(err error, message string) {
    if err != nil {
        panic(err)
    }
    fmt.Printf("%s\n", message)
}
