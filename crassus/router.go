package crassus

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

type Router struct {
	Channels []Channel
	Clients  []Client
}

func (r *Router) Run() {
	log.Println("Starting Crassus router...")
	http.HandleFunc("/ws", r.connectionHandler)
	go r.handleMessages()
	waitForSigterm()
}

func (r *Router) handleMessages() {
	for {
		for i, _ := range r.Channels {
			go r.pushMessages(r.Channels[i])
		}
	}
}

func (rtr *Router) connectionHandler(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return // stop processing or else we throw a runtime error below
	}

	client := Client{
		ID:        uuid.Must(uuid.NewV4(), err).String(),
		Conn:      ws,
		Connected: true,
	}
	privChan := Channel{
		ID:            client.ID,
		Channel:       make(chan Message),
		Subscriptions: []Client{client},
	}

	// Register client
	rtr.Clients = append(rtr.Clients, client)
	rtr.Channels = append(rtr.Channels, privChan)

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %+v\n", err)
			break
		}
		msg.From = client.ID
		for i, v := range rtr.Channels {
			if msg.Channel == v.ID {
				rtr.Channels[i].Channel <- msg // Can't use v because it's a copy
				break                          // Channels should be unique
			}
		}
	}
}

func (r *Router) pushMessages(ch Channel) {
	// Pull message from the channel
	msg := <-ch.Channel

	switch msg.Action {
	case SUBSCRIBE:
		client, err := r.getClient(msg.From)
		if err != nil {
			log.Println("Unable to lookup client for subscription: " + err.Error())
		} else {
			ch.Subscribe(client)
		}
	case TELL:
		// Send message to all clients subscribed to channel
		for i := range ch.Subscriptions {
			if ch.Subscriptions[i].Connected {
				err := ch.Subscriptions[i].Conn.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %+v\n", err)
					ch.Subscriptions[i].Conn.Close()
					ch.Subscriptions[i].Connected = false
				}
			}
		}
	}
}

func (r *Router) getClient(id string) (Client, error) {
	var err error
	for i := range r.Clients {
		if r.Clients[i].ID == id {
			return r.Clients[i], err
		}
	}
	err = fmt.Errorf("Client not found")
	return Client{}, err
}
