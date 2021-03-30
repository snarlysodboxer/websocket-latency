package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	listenAddress  = flag.String("listenAddress", "localhost:9090", "Listen Address")
	forwardAddress = flag.String("forwardAddress", "localhost:8080", "Forward Address")
	delay          = flag.Duration("delay", time.Millisecond*300, `Delay, units are "ns", "us" "ms", "s", "m", "h"`)
)

var upgrader = websocket.Upgrader{}

// Forwarder defines listen and forward addresses
type Forwarder struct {
	listenAddress  string
	forwardAddress string
	delay          time.Duration
}

func (forwarder *Forwarder) forward(w http.ResponseWriter, r *http.Request) {
	listenConnection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer listenConnection.Close()

	forwardURL := fmt.Sprintf("ws://%s", forwarder.forwardAddress)
	forwardConnection, _, err := websocket.DefaultDialer.Dial(forwardURL, nil)
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	defer forwardConnection.Close()

	go func() {
		for {
			messageType, message, err := forwardConnection.ReadMessage()
			if err != nil {
				log.Println("Forward connection read error:", err)
				break
			}
			time.Sleep(forwarder.delay)
			err = listenConnection.WriteMessage(messageType, message)
			if err != nil {
				log.Println("Listen connection write error:", err)
				break
			}
		}
	}()

	for {
		messageType, message, err := listenConnection.ReadMessage()
		if err != nil {
			log.Println("Listen connection read error:", err)
			break
		}
		time.Sleep(forwarder.delay)
		err = forwardConnection.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Forward connection write error:", err)
			break
		}
	}
}

func listenAndForward(forwarder *Forwarder) error {
	http.HandleFunc("/", forwarder.forward)
	return http.ListenAndServe(forwarder.listenAddress, nil)
}

func main() {
	flag.Parse()
	forwarder := &Forwarder{
		listenAddress:  *listenAddress,
		forwardAddress: *forwardAddress,
		delay:          *delay,
	}
	log.Fatal(listenAndForward(forwarder))
}
