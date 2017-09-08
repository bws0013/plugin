package main

import (
  "net"
  "fmt"
  "os"
  "io/ioutil"
  "path/filepath"
  "encoding/gob"
  "bufio"
  "time"
)

// *****Attention, see directory creation requirements*****

// TODO compare current method of getting file to just using os.Open

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

  reader := bufio.NewReader(os.Stdin)
  fmt.Print("Enter text: ")
  text, _ := reader.ReadString('\n')

  // fmt.Println(get_all_files_from_dir("./../../storage/sent/"))

  p := form_packet("hello", "sleep 10 && echo done", "./../../storage/sent/numbers.in")
  p2 := form_packet("hello_1", "ls -alh && ls","./../../storage/sent/")

  if true != true {
    fmt.Println(p)
    fmt.Println(p2)
  }

  // connect to this socket
  if text == "1\n" {
    dial_server_packet(p)
  }
  if text == "2\n" {
    dial_server_packet(p2)
  }

}

// Send a packet to the server
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

// Create a packet to be sent to the server
func form_packet(message, commands, file_path string) my_packet {
  current_time := time.Now().Format(time.RFC3339)

  file_exists := check_for_file(file_path)
  dir_exists := check_for_dir(file_path)

  something_exists := file_exists || dir_exists

  var all_names []string
  var all_permissions []uint
  var all_files_text [][]byte

  if file_exists {

    file_name := filepath.Base(file_path)
    local_permissions, local_text := get_file_and_permissions(file_path)

    all_names = append(all_names, file_name)
    all_permissions = append(all_permissions, local_permissions)
    all_files_text = append(all_files_text, local_text)
  } else if dir_exists {
    files, err := ioutil.ReadDir(file_path)
    check_err(err, "")
    for _, f := range files {
      full_file_path := file_path + f.Name()
      local_permissions, local_text := get_file_and_permissions(full_file_path)

      all_names = append(all_names, f.Name())
      all_permissions = append(all_permissions, local_permissions)
      all_files_text = append(all_files_text, local_text)
    }
  } else {
    fmt.Println("Error Reading File/Dir, Returning empty packet")
  }


  return my_packet{current_time, message, commands, something_exists, all_names, all_permissions, all_files_text}

}

// Get a file and the file permissions
func get_file_and_permissions(file_path string) (uint, []byte) {
  if !check_for_file(file_path) { return 0, nil }

  file_text, err := ioutil.ReadFile(file_path)
  check_err(err, "")
  fileInfo, err := os.Stat(file_path)
  check_err(err, "")

  var mode uint
  if err == nil { mode = uint(fileInfo.Mode()) }

  return mode, file_text
}


// Used for error checking/logging
func check_err(err error, message string) {
    if err != nil {
      panic(err)
    }
    if len(message) != 0 {
      fmt.Printf("%s\n", message)
    }
}

// Check if a directory exists
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

// Check if a file exists
func check_for_file(file_path string) bool {
  finfo, err := os.Stat(file_path)
  if err != nil {
    return false
  }
  if !finfo.IsDir() {
    return true
  }
  return false
}
