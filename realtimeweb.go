package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Connections struct {
	clients      map[chan string]bool
	addClient    chan chan string
	removeClient chan chan string
	messages     chan string
}

var hub = &Connections{
	clients:      make(map[chan string]bool),
	addClient:    make(chan (chan string)),
	removeClient: make(chan (chan string)),
	messages:     make(chan string),
}

func (hub *Connections) Init() {
	go func() {
		for {
			select {
			case s := <-hub.addClient:
				hub.clients[s] = true
				log.Println("Added new client")
			case s := <-hub.removeClient:
				delete(hub.clients, s)
				log.Println("Removed client")
			case msg := <-hub.messages:
				for s, _ := range hub.clients {
					s <- msg
				}
				log.Printf("Broadcast \"%v\" to %d clients", msg, len(hub.clients))
			}
		}
	}()
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported!", http.StatusInternalServerError)
		return
	}

	if r.URL.Path == "/send" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		str := string(body)
		hub.messages <- str
		f.Flush()
		return
	}

	messageChannel := make(chan string)
	hub.addClient <- messageChannel
	notify := w.(http.CloseNotifier).CloseNotify()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for i := 0; i < 1440; {
		select {
		case msg := <-messageChannel:
			jsonData, _ := json.Marshal(msg)
			str := string(jsonData)
			if r.URL.Path == "/events/sse" {
				fmt.Fprintf(w, "data: {\"str\": %s, \"time\": \"%v\"}\n\n", str, time.Now())
			} else if r.URL.Path == "/events/lp" {
				fmt.Fprintf(w, "{\"str\": %s, \"time\": \"%v\"}", str, time.Now())
			}
			f.Flush()
		case <-time.After(time.Second * 60):
			if r.URL.Path == "/events/sse" {
				fmt.Fprintf(w, "data: {\"str\": \"No Data\"}\n\n")
			} else if r.URL.Path == "/events/lp" {
				fmt.Fprintf(w, "{\"str\": \"No Data\"}")
			}
			f.Flush()
			i++
		case <-notify:
			f.Flush()
			i = 1440
			hub.removeClient <- messageChannel
		}
	}
}

func websocketHandler(ws *websocket.Conn) {
	var in string
	messageChannel := make(chan string)
	hub.addClient <- messageChannel

	for i := 0; i < 1440; {
		select {
		case msg := <-messageChannel:
			jsonData, _ := json.Marshal(msg)
			str := string(jsonData)
			in = fmt.Sprintf("{\"str\": %s, \"time\": \"%v\"}\n\n", str, time.Now())
			websocket.Message.Send(ws, in)
		case <-time.After(time.Second * 60):
			in = fmt.Sprintf("{\"str\": \"No Data\"}\n\n")
			websocket.Message.Send(ws, in)
			i++
		}
	}
}

func telnetHandler(c *net.TCPConn) {
	defer c.Close()
	fmt.Printf("Connection from %s to %s established.\n", c.RemoteAddr(), c.LocalAddr())
	io.WriteString(c, fmt.Sprintf("Welcome on %s\n", c.LocalAddr()))
	buf := make([]byte, 4096)
	for {
		n, err := c.Read(buf)
		if (err != nil) || (n == 0) {
			c.Close()
			break
		}
		str := strings.TrimSpace(string(buf[0:n]))
		hub.messages <- str
		io.WriteString(c, "sent to "+strconv.Itoa(len(hub.clients))+" clients\n")
	}
	time.Sleep(150 * time.Millisecond)
	fmt.Printf("Connection from %v closed.\n", c.RemoteAddr())
	c.Close()
	return
}

func listenForTelnet(ln *net.TCPListener) {
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go telnetHandler(conn)
	}
}

func main() {
	fmt.Printf("application started at: %s\n", time.Now().Format(time.RFC822))
	var starttime int64 = time.Now().Unix()
	runtime.GOMAXPROCS(8)

	hub.Init()

	http.HandleFunc("/send", httpHandler)
	http.HandleFunc("/events/sse", httpHandler)
	http.HandleFunc("/events/lp", httpHandler)
	http.Handle("/events/ws", websocket.Handler(websocketHandler))
	http.Handle("/", http.FileServer(http.Dir("./")))

	go http.ListenAndServe(":8000", nil)

	ln, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: 8001,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go listenForTelnet(ln)

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
				fmt.Printf("connected to %d clients\n", len(hub.clients))
			case "3", "system":
				fmt.Printf("CPU cores: %d\nGo calls: %d\nGo routines: %d\nGo version: %v\nProcess ID: %v\n", runtime.NumCPU(), runtime.NumCgoCall(), runtime.NumGoroutine(), runtime.Version(), syscall.Getpid())
			case "4", "time":
				fmt.Printf("application running since %d minutes\n", (time.Now().Unix()-starttime)/60)
			}
		}
	}
	os.Exit(0)
}
