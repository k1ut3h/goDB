package main

import (
  "fmt"
  "net"
  "os"
  "strings"
  "regexp"
)

const (
  SERVER_HOST = "localhost"
  SERVER_PORT = "3000"
  SERVER_TYPE = "tcp"
)

func main(){
  server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

  store := make(map[string]string)

  if err!=nil{
    fmt.Println("Error Listening: ", err.Error())
    os.Exit(1)
  }

  defer server.Close()

  fmt.Println("Welcome to the in-memory key-value store")
  fmt.Println("Listening on "+SERVER_HOST+":"+SERVER_PORT)
  fmt.Println("Waiting for client... ")

  for{
    connection, err := server.Accept()
    if err!=nil{
      fmt.Println("Error accepting: ", err.Error())
    }
    fmt.Println("Client Connected")
    go processClient(connection, store)
  }
}

func processClient(connection net.Conn, store map[string]string){
  buffer := make([]byte, 1024)
  mLen, err := connection.Read(buffer)

  if err!=nil{
    fmt.Println("Error reading: ", err.Error())
  }

  fmt.Println()
  content := strings.Split(string(buffer[:mLen]), "\n")

  header := string(content[0])
  data := string(content[len(content)-1])

  set, seterr := regexp.MatchString("SET", header)
  get, geterr := regexp.MatchString("GET", header)

  if seterr!=nil{
    fmt.Println(seterr.Error())
  }

  if geterr!=nil{
    fmt.Println(geterr.Error())
  }

  if get {
    fmt.Println("Get request made")
    ret_val := store[data]
    if ret_val==""{
      _, err = connection.Write([]byte("No value\n"))
    } else {
      _, err = connection.Write([]byte(ret_val+"\n"))
    }
  }

  if set{
    fmt.Println("Set request made")
    key_value := strings.Split(data, ",")
    store[strings.TrimSpace(key_value[0])] = strings.TrimSpace(key_value[1])
    _, err = connection.Write([]byte("Data written successfully\n"))
    if err!=nil{
      fmt.Println(err.Error())
    }
  }

  defer connection.Close()
}
