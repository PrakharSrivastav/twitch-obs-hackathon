package main

import (
	"context"
	"encoding/json"
	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/PrakharSrivastav/twitch-obs-hackathon/azure"
	"github.com/PrakharSrivastav/twitch-obs-hackathon/obs"
	"log"
)

func main() {

	// get the service bus sbClient
	sbClient, err := azure.NewClient()
	if err != nil {
		log.Fatalln("azure connection error : ", err)
	}
	obsClient, err := obs.NewClient()
	if err != nil {
		log.Fatalln("OBS obsClient error : ", err)
	}
	defer closeConnections(sbClient, obsClient)

	// some logic here
	for {
		err = sbClient.Queue.ReceiveOne(context.Background(), MessagePrinter{cc: obsClient})
		if err != nil {
			log.Fatalln("receive error : ", err)
		}
	}

	log.Println("All Good")

}

func closeConnections(client *azure.SBClient, cc *obs.OBSClient) {

	err := client.Close()
	if err != nil {
		log.Println(err)
	}

	err = cc.Client.Conn.Close()
	if err != nil {
		log.Println(err)
	}
	log.Println("connection closed")

}

type MessagePrinter struct {
	cc *obs.OBSClient
}

func (mp MessagePrinter) Handle(ctx context.Context, msg *servicebus.Message) error {

	var mm SbMessage
	err := json.Unmarshal(msg.Data, &mm)
	if err != nil {
		log.Printf("unknown message %s \n", string(msg.Data))
	}

	log.Printf("message %v \n", mm)

	switch mm.Action {
	case "SwitchScene":
		scene := mm.Options["scene"]
		if scene != "" {
			if err := mp.cc.SwitchScene(scene); err != nil {
				log.Printf("error : %v \n", err)
			}
		}

	case "StartStream":
		if err := mp.cc.StartStream(); err != nil {
			log.Printf("error : %v \n", err)
		}

	case "StopStream":
		if err := mp.cc.StopStream(); err != nil {
			log.Printf("error : %v \n", err)
		}

	case "QuestionAnswered":
		// aggregate
		// calculate if highest count is correct
		// if correct
		// switch to next scene
		// go back to start

	default:
		log.Println("unknown action ", mm.Action)
	}

	return msg.Complete(ctx)
}

type SbMessage struct {
	Action  string            `json:"action"`
	Options map[string]string `json:"options"`
}

var currentScene = "BRB"
