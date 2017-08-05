package main

import (
  "net"
  "fmt"
  "bufio"
  "os"
  "io/ioutil"
  "path/filepath"
  "encoding/gob"
)
//

type my_packet struct {
  Message string
  Contains_file bool
  File_name string
  File []byte
}


func main() {
  //dial_server_message()

  p := form_packet("hello", "./../../storage/sent/numbers.in")

  // connect to this socket
  dial_server_packet(p)
  fmt.Println(p.Message)
}

func dial_server_packet(packet my_packet) {
  conn, err := net.Dial("tcp", "127.0.0.1:8081")

  if err != nil {
    fmt.Println("Unable to send!")
    return
  }

  encoder := gob.NewEncoder(conn)
  err = encoder.Encode(packet)
  check_err(err, "everything is fine")
  conn.Close()

  message, _ := bufio.NewReader(conn).ReadString('\n')
  fmt.Println("Message from server: " + message)


}

func dial_server_message() {
  conn, _ := net.Dial("tcp", "127.0.0.1:8081")

  for {
    // read in input from stdin
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Text to send: ")
    text, _ := reader.ReadString('\n')
    // send to socket
    fmt.Fprintf(conn, text + "\n")
    // listen for reply
    message, _ := bufio.NewReader(conn).ReadString('\n')
    fmt.Print("Message from server: "+message)
  }
}

func form_packet(message, file_path string) my_packet {
  file_exists := check_for_file(file_path)
  if file_exists {
    file, err := ioutil.ReadFile(file_path)
    file_name := filepath.Base(file_path)
    check_err(err, "")
    return my_packet{message, file_exists, file_name, file}
  } else {
    return my_packet{message, file_exists, "", nil}
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
  check_err(err, "")
  if !finfo.IsDir() {
    //fmt.Println(file_path)
    return true
  }
  return false
}
