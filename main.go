package main

import (
	"log"
	"time"

	"github.com/christopher-dG/go-obs-websocket"
)

func main() {
	// Connect a client.
	c := obsws.Client{Host: "localhost", Port: 4444}
	if err := c.Connect(); err != nil {
		log.Fatal(err)
	}
	defer c.Disconnect()
	log.Println("1")

	// Send and receive a request asynchronously.
	req := obsws.NewGetStreamingStatusRequest()
	if err := req.Send(c); err != nil {
		log.Println("i am here")
		log.Fatal(err)
	}
	log.Println("2")

	// This will block until the response comes (potentially forever).
	resp, err := req.Receive()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("3")

	log.Println("streaming:", resp.Streaming)

	// Set the amount of time we can wait for a response.
	obsws.SetReceiveTimeout(time.Second * 2)

	// Send and receive a request synchronously.
	req = obsws.NewGetStreamingStatusRequest()
	// Note that we create a new request,
	// because requests have IDs that must be unique.
	// This will block for up to two seconds, since we set a timeout.
	resp, err = req.SendReceive(c)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("streaming:", resp.Streaming)

	// Respond to events by registering handlers.
	err = c.AddEventHandler("SwitchScenes", func(e obsws.Event) {
		// Make sure to assert the actual event type.
		log.Println("new scene:", e.(obsws.SwitchScenesEvent).SceneName)
	})

	if err != nil {
		log.Println(err)
		return
	}


	c.

	for {
	}
}
