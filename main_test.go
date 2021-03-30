package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var testUpgrader = websocket.Upgrader{}

func testEchoHandler(responseWriter http.ResponseWriter, request *http.Request) {
	conn, err := testUpgrader.Upgrade(responseWriter, request, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			break
		}
	}
}

func TestEcho(t *testing.T) {
	// setup test server
	server := httptest.NewServer(http.HandlerFunc(testEchoHandler))
	defer server.Close()

	// setup latency server
	latencyServerURL := "localhost:9090"
	forwarder := &Forwarder{
		listenAddress:  latencyServerURL,
		forwardAddress: strings.TrimPrefix(server.URL, `http://`),
		delay:          time.Millisecond * 300,
	}
	go func() {
		listenAndForward(forwarder)
	}()

	// setup client
	latencyServerWebsocketURL := fmt.Sprintf("ws://%s", latencyServerURL)
	clientConnection, _, err := websocket.DefaultDialer.Dial(latencyServerWebsocketURL, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer clientConnection.Close()

	// echo through latency server
	for i := 0; i < 10; i++ {
		start := time.Now()
		testMessage := fmt.Sprintf("this is message %d", i)
		if err := clientConnection.WriteMessage(websocket.TextMessage, []byte(testMessage)); err != nil {
			t.Errorf("%v", err)
		}
		_, returnMessage, err := clientConnection.ReadMessage()
		if err != nil {
			t.Errorf("%v", err)
		}
		duration := time.Since(start)
		if string(returnMessage) != testMessage {
			t.Errorf("Expected '%s', got '%s'", testMessage, string(returnMessage))
		}
		minExpectedDuration := time.Millisecond * 600
		maxExpectedDuration := time.Millisecond * 602
		if duration < minExpectedDuration || duration > maxExpectedDuration {
			t.Errorf("Expected '%s' to be greater than '%s' and less than '%s'", duration, minExpectedDuration, maxExpectedDuration)
		}
	}
	err = clientConnection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		t.Errorf("%v", err)
	}
}
