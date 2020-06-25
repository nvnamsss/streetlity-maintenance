package server

import (
	"log"
	"net/http"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

var server *socketio.Server
var chat_stack map[string][]string

const Tag string = "[Realtime-Server]"

//Open new namespace for common user and maintenance user. The create namespace include chat room and location room
func OpenOrderSpace(nsp string) {
	log.Println(Tag, "Creating namespace", nsp)
	chat_stack[nsp] = []string{}

	server.OnConnect(nsp, func(s socketio.Conn) (e error) {
		log.Println(Tag, nsp, "new connection:", s.RemoteAddr().String())
		s.Join("location")
		s.Join("chat")
		s.Join("information")
		s.Emit("joined")

		return
	})

	// server.OnEvent(nsp, "join", func(s socketio.Conn, msg string) {
	// 	log.Println(Tag, "Join", msg)
	// 	s.Join("location")
	// 	s.Join("chat")
	// 	s.Join("information")
	// })

	server.OnEvent(nsp, "update-location", func(s socketio.Conn, data string) {
		log.Println(Tag, "update-location", data)
		address := s.RemoteAddr()

		server.ForEach(nsp, "location", func(c socketio.Conn) {
			if c.RemoteAddr() != address {
				c.Emit("location-update", data)
			}
		})
	})

	server.OnEvent(nsp, "pull-location", func(s socketio.Conn) {
		log.Println(Tag, "pull-location")
		address := s.RemoteAddr()
		server.ForEach(nsp, "location", func(c socketio.Conn) {
			if c.RemoteAddr() != address {
				c.Emit("pull-location")
			}
		})
	})

	server.OnEvent(nsp, "update-information", func(s socketio.Conn, data string) {
		log.Println(Tag, "update-information", data)
		address := s.RemoteAddr()

		server.ForEach(nsp, "information", func(c socketio.Conn) {
			if c.RemoteAddr() != address {
				c.Emit("update-information", data)
			}
		})
	})

	server.OnEvent(nsp, "pull-information", func(s socketio.Conn, msg string) {
		log.Println(Tag, "pull-information", msg)
		address := s.RemoteAddr()

		server.ForEach(nsp, "information", func(c socketio.Conn) {
			if c.RemoteAddr() != address {
				c.Emit("pull-information", msg)
			}
		})
	})

	server.OnEvent(nsp, "chat", func(s socketio.Conn, msg string) {
		// s.SetContext(msg)
		log.Println(Tag, "chat", msg)
		log.Println(Tag, "send message to ", server.RoomLen(nsp, "chat"))
		chat_stack[nsp] = append(chat_stack[nsp], msg)
		address := s.RemoteAddr()
		server.ForEach(nsp, "chat", func(c socketio.Conn) {
			if c.RemoteAddr() != address {
				c.Emit("chat", msg)
			}
		})
	})

	server.OnEvent(nsp, "pull-chat", func(s socketio.Conn) {
		stack := chat_stack[nsp]
		for _, msg := range stack {
			s.Emit("chat", msg)
			time.Sleep(50 * time.Millisecond)
		}
	})

	server.OnEvent(nsp, "typing-chat", func(s socketio.Conn, typing_user string) {
		log.Println(Tag, "typing-chat", "from", typing_user)
		address := s.RemoteAddr()

		server.ForEach(nsp, "chat", func(c socketio.Conn) {
			if c.RemoteAddr() != address {
				c.Emit("typing-chat", typing_user)
			}
		})
	})

	server.OnEvent(nsp, "typed-chat", func(s socketio.Conn, typed_user string) {
		log.Println(Tag, "typed-chat", "from", typed_user)
		address := s.RemoteAddr()

		server.ForEach(nsp, "chat", func(c socketio.Conn) {
			if c.RemoteAddr() != address {
				c.Emit("typed-chat", typed_user)
			}
		})
	})

	server.OnEvent(nsp, "decline", func(s socketio.Conn, msg string) {
		server.ClearRoom(nsp, "join")
		server.ClearRoom(nsp, "chat")
		server.ClearRoom(nsp, "information")
		delete(chat_stack, nsp)

		server.ForEach(nsp, "chat", func(c socketio.Conn) {
			c.Close()
		})
	})

	server.OnDisconnect(nsp, func(s socketio.Conn, msg string) {
		log.Println(Tag, nsp, "connection is closed:", s.RemoteAddr().String())
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
	chat_stack = make(map[string][]string)
	if err != nil {
		log.Fatal(Tag, err)
	}

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println(Tag, "meet error:", e.Error())
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
