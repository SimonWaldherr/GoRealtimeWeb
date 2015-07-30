#GoRealtimeWeb

[![Flattr donate button](https://raw.github.com/balupton/flattr-buttons/master/badge-89x18.gif)](https://flattr.com/submit/auto?user_id=SimonWaldherr&url=http%3A%2F%2Fgithub.com%2FSimonWaldherr%2FGoRealtimeWeb "Donate monthly to this project using Flattr")


Examples how to write realtime web applications in Golang

This repository contains examples for the following real-time implementations:  

* [Server-sent events](http://en.wikipedia.org/wiki/Server-sent_events)
* [Long Polling](http://en.wikipedia.org/wiki/Push_technology#Long_polling)
* [Websocket](http://en.wikipedia.org/wiki/WebSocket)

Most of the **long polling** systems close the connection after each transmission from the server, with the help of [oboe.js](https://github.com/jimhigson/oboe.js) this example can handle multiple messages from the server in realtime without closing/reconnecting.  

##howto

1. install dependencies ```go get golang.org/x/net/websocket```
2. run the go script ```go run realtimeweb.go```
3. open **<http://localhost:8000>** and choose an example  
4. connect to localhost port 8001 via telnet (```telnet localhost 8001```)
5. have fun

##info

for more informations about real-time technologies you can read [this excelent stackoverflow article (AJAX vs Long-Polling vs SSE vs Websockets vs Comet)](http://stackoverflow.com/a/12855533).  

##license

[MIT (see LICENSE file)](https://github.com/SimonWaldherr/GoRealtimeWeb/blob/master/LICENSE)
