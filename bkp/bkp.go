package bkp

import (
	"log"
	"time"

	"github.com/christopher-dG/go-obs-websocket"
)

func unmain() {
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
		log.Fatal("req.Send", err)
	}

	// This will block until the response comes (potentially forever).
	resp, err := req.Receive()
	if err != nil {
		log.Fatal("req.Receive", err)
	}
	log.Println("streaming:", resp.Streaming)

	// Set the amount of time we can wait for a response.
	obsws.SetReceiveTimeout(time.Second * 10)

	// Send and receive a request synchronously.
	req = obsws.NewGetStreamingStatusRequest()
	resp, err = req.SendReceive(c)
	if err != nil {
		log.Fatal("req.SendReceive", err)
	}

	// Respond to events by registering handlers.
	if err = c.AddEventHandler("SwitchScenes", func(e obsws.Event) {
		log.Println("new scene:", e.(obsws.SwitchScenesEvent).SceneName)
	}); err != nil {
		log.Println("c.AddEventHandler", err)
		return
	}

	log.Println("sending the request to get scene")
	/*request := obsws.GetSceneListRequest{}
	receive, err := request.SendReceive(c)
	if err != nil {
		log.Fatal("req.SendReceive :: ", err)
	}

	log.Printf("%v \n", receive)*/

	request := obsws.GetCurrentSceneRequest{}
	rx, err := request.SendReceive(c)
	if err != nil {
		log.Fatal("gcs error :: ", err)
	}

	log.Printf("%v \n", rx.Name)

	for {
	}
}
