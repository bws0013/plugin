package main

import (
  "net"
  "fmt"
  "io/ioutil"
  "os"
  "os/exec"
)
import "encoding/gob"

// TODO add logging feature
// TODO create log directory and ability to export it
// TODO create time based files for input
// TODO check for redundant/unused methods

type my_packet struct {
  Current_time string
  Message string
  Contains_file bool
  File_name []string
  Permissions []uint
  File [][]byte
}

func main() {
  // execute_commands("ls && ls -alh")

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

  // fmt.Println(p)

  conn.Write([]byte("liftoff"))

  conn.Close()

  // fmt.Printf("Message: %s\n", p.File_name)

  if p.Contains_file && p.File != nil {
    create_file(p.File_name, p.Permissions, p.File)
    execute_commands(p.Message, p.Current_time)
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
func execute_commands(message, time string) {
  path := "./../../storage/scripts/"
  name := path + "bs_" + time + ".sh"
  permissions := uint(511)
  message_in_byte_form := []byte(message)

  err := ioutil.WriteFile(name, message_in_byte_form, os.FileMode(permissions))
  check_err(err, "Creating Script")

  out, err := exec.Command("/bin/sh", name).Output()
  check_err(err, "Executing Script")
  fmt.Printf("%s\n", out)
}

// Clear the directory where executables will be stored
func clear_execution_dir() {
  to_delete := "./../../storage/recieved/"

  if !check_for_dir(to_delete) {
    return
  }

  // Delete files

}

func create_file(name []string, permissions []uint, data [][]byte) {
  current_path := "./../../storage/recieved/"

  for i, _ := range name {

    current_name := current_path + name[i]
    current_permissions := permissions[i]
    current_data := data[i]

    err := ioutil.WriteFile(current_name, current_data, os.FileMode(current_permissions))
    check_err(err, "File created!")
  }
}

func check_for_dir(file_path string) bool {
  finfo, err := os.Stat(file_path)
  if err != nil {
    return false
  }
  if finfo.IsDir() {
    return true
  }
  return false
}

func check_err(err error, message string) {
    if err != nil {
        panic(err)
    }
    fmt.Printf("%s\n", message)
}
