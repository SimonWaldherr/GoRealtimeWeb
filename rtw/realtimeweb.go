package rtw

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Connections struct {
	Clients      map[chan string]bool
	AddClient    chan chan string
	RemoveClient chan chan string
	Messages     chan string
}

var Hub = &Connections{
	Clients:      make(map[chan string]bool),
	AddClient:    make(chan (chan string)),
	RemoveClient: make(chan (chan string)),
	Messages:     make(chan string),
}

func (Hub *Connections) Init() {
	go func() {
		for {
			select {
			case s := <-Hub.AddClient:
				Hub.Clients[s] = true
				log.Println("Added new client")
			case s := <-Hub.RemoveClient:
				delete(Hub.Clients, s)
				log.Println("Removed client")
			case msg := <-Hub.Messages:
				for s, _ := range Hub.Clients {
					s <- msg
				}
				log.Printf("Broadcast \"%v\" to %d Clients", msg, len(Hub.Clients))
			}
		}
	}()
}

func HttpHandler(w http.ResponseWriter, r *http.Request) {
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
		Hub.Messages <- str
		f.Flush()
		return
	}

	messageChannel := make(chan string)
	Hub.AddClient <- messageChannel
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
			Hub.RemoveClient <- messageChannel
		}
	}
}

func WebsocketHandler(ws *websocket.Conn) {
	var in string
	messageChannel := make(chan string)
	Hub.AddClient <- messageChannel

	for i := 0; i < 1440; {
		select {
		case msg := <-messageChannel:
			jsonData, _ := json.Marshal(msg)
			str := string(jsonData)
			in = fmt.Sprintf("{\"str\": %s, \"time\": \"%v\"}\n\n", str, time.Now())

			if err := websocket.Message.Send(ws, in); err != nil {
				Hub.RemoveClient <- messageChannel
				i = 1440
			}
		case <-time.After(time.Second * 60):
			in = fmt.Sprintf("{\"str\": \"No Data\"}\n\n")
			if err := websocket.Message.Send(ws, in); err != nil {
				Hub.RemoveClient <- messageChannel
				i = 1440
			}
		}
		i++
	}
}

func TelnetHandler(c *net.TCPConn) {
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
		Hub.Messages <- str
		io.WriteString(c, "sent to "+strconv.Itoa(len(Hub.Clients))+" Clients\n")
	}
	time.Sleep(150 * time.Millisecond)
	fmt.Printf("Connection from %v closed.\n", c.RemoteAddr())
	c.Close()
	return
}

func ListenForTelnet(ln *net.TCPListener) {
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go TelnetHandler(conn)
	}
}
