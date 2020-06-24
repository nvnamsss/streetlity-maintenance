package server

import (
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

var server *socketio.Server

const Tag string = "[Realtime-Server]"

//Open new namespace for common user and maintenance user. The create namespace include chat room and location room
func OpenOrderSpace(nsp string) {
	log.Println(Tag, "Creating namespace", nsp)

	server.OnConnect(nsp, func(s socketio.Conn) (e error) {
		log.Println(Tag, nsp, "new connection:", s.RemoteAddr().String())
		s.Join("location")
		s.Join("chat")
		s.Join("information")

		s.Emit("joined")
		return
	})

	server.OnEvent(nsp, "join", func(s socketio.Conn, msg string) {
		log.Println(Tag, "Join", msg)
		s.Join("location")
		s.Join("chat")
		s.Join("information")
	})

	server.OnEvent(nsp, "update-location", func(s socketio.Conn, data string) {
		log.Println(Tag, "update-location", data)
		server.ForEach(nsp, "location", func(c socketio.Conn) {
			if c.ID() != s.ID() {
				c.Emit("location-update", data)
			}
		})
	})

	server.OnEvent(nsp, "pull-location", func(s socketio.Conn) {
		log.Println(Tag, "pull-location")

		server.ForEach(nsp, "location", func(c socketio.Conn) {
			if c.ID() != s.ID() {
				c.Emit("pull-location")
			}
		})
	})

	server.OnEvent(nsp, "update-information", func(s socketio.Conn, data string) {
		log.Println(Tag, "update-information", data)
		server.ForEach(nsp, "information", func(c socketio.Conn) {
			if c.ID() != s.ID() {
				c.Emit("update-information", data)
			}
		})
	})

	server.OnEvent(nsp, "pull-information", func(s socketio.Conn, msg string) {
		log.Println(Tag, "pull-information", msg)
		server.ForEach(nsp, "information", func(c socketio.Conn) {
			if c.ID() != s.ID() {
				c.Emit("pull-information", msg)
			}
		})
	})

	server.OnEvent(nsp, "chat", func(s socketio.Conn, msg string, timestamp string) {
		// s.SetContext(msg)
		log.Println(Tag, "chat", msg, timestamp)
		server.ForEach(nsp, "chat", func(c socketio.Conn) {
			if c.ID() != s.ID() {
				c.Emit("chat", msg, timestamp)
			}
		})
	})

	server.OnEvent(nsp, "typing-chat", func(s socketio.Conn, typing_user string) {
		log.Println(Tag, "typing-chat", "from", typing_user)
		server.ForEach(nsp, "chat", func(c socketio.Conn) {
			if c.ID() != s.ID() {
				c.Emit("typing-chat", typing_user)
			}
		})
	})

	server.OnEvent(nsp, "typed-chat", func(s socketio.Conn, typed_user string) {
		log.Println(Tag, "typed-chat", "from", typed_user)
		server.ForEach(nsp, "chat", func(c socketio.Conn) {
			if c.ID() != s.ID() {
				c.Emit("typed-chat", typed_user)
			}
		})
	})

	server.OnEvent(nsp, "decline", func(s socketio.Conn, msg string) {
		server.ClearRoom(nsp, "join")
		server.ClearRoom(nsp, "chat")
		server.ClearRoom(nsp, "information")
		s.Close()
	})

	server.OnDisconnect(nsp, func(s socketio.Conn, msg string) {
		s.Leave("location")
		s.Leave("chat")
		s.Leave("information")
		s.Close()
	})

	rooms := server.Rooms(nsp)

	for _, room := range rooms {
		server.ForEach(nsp, room, func(s socketio.Conn) {
			log.Println(Tag, "Some one still here, bye bitch")
			s.Close()
		})
	}
}

func OpenOrderSpaceByRoom(room string) {
	server.OnEvent("/order", "update-location", func(s socketio.Conn, msg string) {
		server.ForEach("order", room, func(c socketio.Conn) {
			if c.ID() != s.ID() {
				c.Emit("update-location", msg)
			}
		})
	})

	server.OnEvent("/order", "pull-location", func(s socketio.Conn, msg string) {
		server.ForEach("order", room, func(c socketio.Conn) {
			if c.ID() != s.ID() {
				c.Emit("pull-location", msg)
			}
		})
	})

	server.OnEvent("/order", "information", func(s socketio.Conn, msg string) {
		server.ForEach("order", room, func(c socketio.Conn) {
			if c.ID() != s.ID() {
				c.Emit("information", msg)
			}
		})
	})

	server.OnEvent("/order", "join", func(s socketio.Conn, msg string) {
		log.Println(Tag, "Join room 2", room)
		s.Join(room)
	})

	server.OnEvent("/order", "pull-information", func(s socketio.Conn, msg string) {
		server.ForEach("/order", room, func(c socketio.Conn) {
			if c.ID() != s.ID() {
				c.Emit("pull-information", msg)
			}
		})
	})

	server.OnEvent("/order", "chat", func(s socketio.Conn, msg string) {
		s.SetContext(msg)
		server.ForEach("/order", room, func(c socketio.Conn) {
			if c.ID() != s.ID() {
				c.Emit("chat", msg)
			}
		})
	})

	server.OnEvent("/order", "decline", func(s socketio.Conn, msg string) {
		server.ClearRoom("order", room)
		server.BroadcastToRoom("order", room, "decline")
	})
}

func Create() {
	var err error
	server, err = socketio.NewServer(nil)
	if err != nil {
		log.Fatal(Tag, err)
	}

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println(Tag, "meet error:", e.Error())
	})

	go server.Serve()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))

	go func() {
		if e := http.ListenAndServe(":6182", nil); e != nil {
			log.Println(Tag, e.Error())
		} else {
			log.Println(Tag, "Serving at localhost:6182...")

		}

	}()

	OpenOrderSpace("/himom")
	// server.Close()
	// log.Println("Hi mom")
	// log.Fatal()
}
