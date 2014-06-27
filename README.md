#GoRealtimeWeb

Examples how to write realtime web applications in Golang

This repository contains examples for the following real-time implementations:  

* [Server-sent events](http://en.wikipedia.org/wiki/Server-sent_events)
* [Long Polling](http://en.wikipedia.org/wiki/Push_technology#Long_polling)
* [Websocket](http://en.wikipedia.org/wiki/WebSocket)

##howto

1. install dependencies ```go get code.google.com/p/go.net/websocket```
2. run the go script ```go run realtimeweb.go```
3. open **<http://localhost:8000>** and choose an example  
4. connect to localhost port 8001 via telnet (```telnet localhost 8001```)
5. have fun

##info

for more informations about real-time technologies you can read [this excelent stackoverflow article (AJAX vs Long-Polling vs SSE vs Websockets vs Comet)](http://stackoverflow.com/a/12855533).  

##license

[MIT (see LICENSE file)](https://github.com/SimonWaldherr/GoRealtimeWeb/blob/master/LICENSE)
