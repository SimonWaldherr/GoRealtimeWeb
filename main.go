package main

import (
	rtw "./rtw"
	"fmt"
	"golang.org/x/net/websocket"
	"net"
	"net/http"
	"os"
	"runtime"
	"syscall"
	"time"
)

func main() {
	fmt.Printf("application started at: %s\n", time.Now().Format(time.RFC822))
	var starttime int64 = time.Now().Unix()
	runtime.GOMAXPROCS(8)

	rtw.Hub.Init()

	http.HandleFunc("/send", rtw.HttpHandler)
	http.HandleFunc("/events/sse", rtw.HttpHandler)
	http.HandleFunc("/events/lp", rtw.HttpHandler)
	http.Handle("/events/ws", websocket.Handler(rtw.WebsocketHandler))
	http.Handle("/", http.FileServer(http.Dir("./")))

	go http.ListenAndServe(":8000", nil)

	ln, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: 8001,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go rtw.ListenForTelnet(ln)

	var input string
	for input != "exit" {
		_, _ = fmt.Scanf("%v", &input)
		if input != "exit" {
			switch input {
			case "", "0", "5", "help", "info":
				fmt.Print("you can type \n1: \"exit\" to kill this application")
				fmt.Print("\n2: \"clients\" to show the amount of connected clients")
				fmt.Print("\n3: \"system\" to show info about the server")
				fmt.Print("\n4: \"time\" to show since when this application is running")
				fmt.Print("\n5: \"help\" to show this information")
				fmt.Println()
			case "1", "exit", "kill":
				fmt.Println("application get killed in 5 seconds")
				input = "exit"
				time.Sleep(5 * time.Second)
			case "2", "clients":
				fmt.Printf("connected to %d clients\n", len(rtw.Hub.Clients))
			case "3", "system":
				fmt.Printf("CPU cores: %d\nGo calls: %d\nGo routines: %d\nGo version: %v\nProcess ID: %v\n", runtime.NumCPU(), runtime.NumCgoCall(), runtime.NumGoroutine(), runtime.Version(), syscall.Getpid())
			case "4", "time":
				fmt.Printf("application running since %d minutes\n", (time.Now().Unix()-starttime)/60)
			}
		}
	}
	os.Exit(0)
}
