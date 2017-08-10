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
//

type my_packet struct {
  Current_time string
  Message string
  Contains_file bool
  File_name string
  Permissions uint
  File []byte
}


func main() {
  //dial_server_message()

  fmt.Println(get_all_files_from_dir("./../../storage/sent/"))

  p := form_packet("hello", "./../../storage/sent/numbers.in")

  // connect to this socket
  dial_server_packet(p)
}

func dial_server_packet(packet my_packet) {
  conn, err := net.Dial("tcp", "192.168.1.5:8081")

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

func form_packet(message, file_path string) my_packet {
  file_exists := check_for_file(file_path)

  current_time := time.Now().Format(time.RFC3339)
  if file_exists {
    file_text, err := ioutil.ReadFile(file_path)
    fileInfo, err := os.Stat(file_path)

    var mode uint
    if err == nil { mode = uint(fileInfo.Mode()) }

    file_name := filepath.Base(file_path)

    return my_packet{current_time, message, file_exists, file_name, mode, file_text}
  } else {
    fmt.Println("No file included")
    return my_packet{current_time, message, file_exists, "", 0000, nil}
  }

}

func get_all_files_from_dir(file_path string) ([]uint, [][]byte) {
  if !check_for_dir(file_path) { return nil, nil }

  var all_permissions []uint
  var all_files_text [][]byte

  files, err := ioutil.ReadDir(file_path)
  check_err(err, "")

  for _, f := range files {
    full_file_path := file_path + f.Name()
    curr_permissions, curr_text := get_file_and_permissions(full_file_path)

    all_permissions = append(all_permissions, curr_permissions)
    all_files_text = append(all_files_text, curr_text)
  }

  return all_permissions, all_files_text
}

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

func check_err(err error, message string) {
    if err != nil {
      panic(err)
    }
    if len(message) != 0 {
      fmt.Printf("%s\n", message)
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
