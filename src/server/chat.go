package server

import (
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

var server *socketio.Server

//Open new namespace for common user and maintenance user. The create namespace include chat room and location room
func OpenOrderSpace(nsp string) {
	server.OnEvent(nsp, "location-update", func(s socketio.Conn, msg string) {
		server.ForEach(nsp, "msg", func(c socketio.Conn) {
			if c.ID() != s.ID() {
				c.Emit("location-update", msg)
			}
		})
	})

	server.OnEvent(nsp, "join", func(s socketio.Conn, msg string) {
		log.Println("[Server]", "Join", msg)
		s.Join("location")
		s.Join("chat")
	})

	server.OnEvent(nsp, "chat", func(s socketio.Conn, msg string) {
		s.SetContext(msg)
		server.ForEach(nsp, "chat", func(c socketio.Conn) {
			if c.ID() != s.ID() {
				c.Emit("msg", msg)
			}
		})
	})

	server.OnDisconnect(nsp, func(s socketio.Conn, msg string) {
		s.Leave("location")
		s.Leave("msg")
		s.Close()
	})
}

func Create() {
	var err error
	server, err = socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/location", "update", func(s socketio.Conn, msg string) {
		log.Println("[SocketIO]", "update location", msg)
		// server.ForEach("/location", room, func(c socketio.Conn) {
		// 	if c.ID() != s.ID() {
		// 		c.Emit("")
		// 	}
		// })
	})

	server.OnEvent("/chat", "join", func(s socketio.Conn, msg string) string {
		log.Println("[SocketIO]", "join chat", msg)
		s.Join(msg)
		return "join " + msg
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		log.Println("[SocketIO]", "message", msg)
		s.SetContext(msg)

		// rooms := s.Rooms()
		server.ForEach("/chat", "", func(c socketio.Conn) {
			if c.ID() != s.ID() {
				c.Emit("msg", msg)
			}
		})
		// for _, room := range rooms {

		// 	// server.BroadcastToRoom("/chat", room, "msg", msg)
		// }

		return "recv " + msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})
	go server.Serve()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	go func() {
		if e := http.ListenAndServe(":6182", nil); e != nil {
			log.Println("[Server]", e.Error())
		} else {
			log.Println("Serving at localhost:6182...")

		}

	}()
	// server.Close()
	// log.Println("Hi mom")
	// log.Fatal()
}
