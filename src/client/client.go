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
    file, err := ioutil.ReadFile(file_path)

    fileInfo, err := os.Stat(file_path)

    var mode uint
    if err == nil { mode = uint(fileInfo.Mode()) }

    file_name := filepath.Base(file_path)

    return my_packet{current_time, message, file_exists, file_name, mode, file}
  } else {
    fmt.Println("No file included")
    return my_packet{current_time, message, file_exists, "", 0000, nil}
  }

}

func check_err(err error, message string) {
    if err != nil {
      panic(err)
    }
    if len(message) != 0 {
      fmt.Printf("%s\n", message)
    }
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
