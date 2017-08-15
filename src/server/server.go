package main

import (
  "net"
  "fmt"
  "io/ioutil"
  "os"
  "os/exec"
  "strings"
)
import "encoding/gob"

// TODO add logging feature
// TODO create log directory and ability to export it
// TODO check for redundant/unused methods

type my_packet struct {
  Current_time string
  Message string
  Commands string
  Contains_file bool
  File_name []string
  Permissions []uint
  File [][]byte
}

func main() {
  // execute_commands("ls && ls -alh")

  fmt.Println("start");
  ln, err := net.Listen("tcp", ":8081")
  // check_err(err, "Listening!")
  if err != nil {
    fmt.Println("Error at listen")
    return
  }

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

  dec := gob.NewDecoder(conn)
  p := &my_packet{}
  err := dec.Decode(p)

  // fmt.Println(p)

  conn.Write([]byte("liftoff"))

  conn.Close()

  check_err(err, "No problems on read in", p.Current_time)

  if p.Contains_file && p.File != nil {
    create_file(p.Current_time, p.File_name, p.Permissions, p.File)
    execute_commands(p.Commands, p.Current_time)
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
func execute_commands(commands, time string) {
  path := "./../../storage/scripts/"
  name := path + "bs_" + time + ".sh"
  permissions := uint(511)
  message_in_byte_form := []byte(commands)

  err := ioutil.WriteFile(name, message_in_byte_form, os.FileMode(permissions))
  check_err(err, "Creating Script", time)

  out, err := exec.Command("/bin/sh", name).Output()
  check_err(err, "Executing Script", time)
  fmt.Printf("%s\n", out)
}

// Clear the directory where executables will be stored
func clear_scripts_dir(time string) {
  file_path := "./../../storage/scripts/"

  if !check_for_dir(file_path) { return }

  files, err := ioutil.ReadDir(file_path)
  // check_err(err, "Getting recieved file names")

  for _, file_name := range files {
    fmt.Println(file_path + file_name.Name())
    if strings.HasPrefix(file_name.Name(), "bs_") {
      err = os.Remove(file_path + file_name.Name())
      check_err(err, "Removing file: " + file_name.Name(), time)
    }

  }

  // Delete files

}

func create_file(time string, name []string, permissions []uint, data [][]byte) {
  current_path := "./../../storage/recieved/"

  for i, _ := range name {

    current_name := current_path + name[i]
    current_permissions := permissions[i]
    current_data := data[i]

    err := ioutil.WriteFile(current_name, current_data, os.FileMode(current_permissions))
    check_err(err, "File created!", time)
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

func check_err(err error, message, time string) {
  log_file_name := "bs_" + time + ".txt"
  message += "\n"
  f, err := os.OpenFile(log_file_name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
  if err != nil {
      panic(err)
  }
  defer f.Close()
  if _, err = f.WriteString(message); err != nil {
    panic(err)
  }
}
