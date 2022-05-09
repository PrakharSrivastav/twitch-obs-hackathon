package main

import (
	"context"
	"fmt"
	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/PrakharSrivastav/twitch-obs-hackathon/azure"
	"github.com/PrakharSrivastav/twitch-obs-hackathon/obs"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"log"
)

func main() {

	// get the service bus sbClient
	sbClient, err := azure.NewClient()
	if err != nil {
		log.Fatalln("azure connection error : ", err)
	}
	defer closeConnections(err, sbClient)()

	obsClient, err := obs.NewClient()
	if err != nil {
		log.Fatalln("OBS obsClient error : ", err)
	}

	mp := MessagePrinter{cc: obsClient}

	// some logic here

	for {
		err = sbClient.Queue.ReceiveOne(context.Background(), mp)
		if err != nil {
			log.Fatalln("receive error : ", err)
		}
	}

	log.Println("All Good")

}

func closeConnections(err error, client *azure.SBClient) func() {
	return func() {
		err = client.Close()
		if err != nil {
			log.Println(err)
		}
		log.Println("connection closed")
	}
}

type MessagePrinter struct {
	cc *obs.OBSClient
}

func (mp MessagePrinter) Handle(ctx context.Context, msg *servicebus.Message) error {
	fmt.Println(string(msg.Data))
	
	scene := string(msg.Data)

	_, err := mp.cc.Client.Scenes.SetCurrentScene(&scenes.SetCurrentSceneParams{SceneName: scene})
	if err != nil {
		log.Println("scene change error :", err)
	}

	return msg.Complete(ctx)
}
