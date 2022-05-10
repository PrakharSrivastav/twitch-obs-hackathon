package main

import (
	"context"
	"encoding/json"
	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/PrakharSrivastav/twitch-obs-hackathon/azure"
	"github.com/PrakharSrivastav/twitch-obs-hackathon/obs"
	"log"
	"time"
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

	go func() {
		log.Println(currentScene)

		for {

			if currentScene != "StartStream" && currentScene != "SceneBRB" {
				time.Sleep(time.Second * 30)

				log.Println("calculating")
				log.Println("currentScene ", currentScene, currentSceneRight, currentSceneWrong)
				if currentScene == "scene-15" && currentSceneRight > currentSceneWrong {
					obsClient.SwitchScene("scene-complete")
				}

				switch {
				case currentScene == "scene-15" && changeScene():
					obsClient.SwitchScene("scene-complete")

				case currentScene == "scene-1" && changeScene():
					resetCount("scene-2", obsClient, sbClient)

				case currentScene == "scene-2" && changeScene():
					resetCount("scene-3", obsClient, sbClient)

				case currentScene == "scene-3" && changeScene():
					resetCount("scene-4", obsClient, sbClient)

				case currentScene == "scene-4" && changeScene():
					resetCount("scene-5", obsClient, sbClient)

				case currentScene == "scene-5" && changeScene():
					resetCount("scene-6", obsClient, sbClient)

				case currentScene == "scene-6" && changeScene():
					resetCount("scene-7", obsClient, sbClient)

				case currentScene == "scene-7" && changeScene():
					resetCount("scene-8", obsClient, sbClient)

				case currentScene == "scene-8" && changeScene():
					resetCount("scene-9", obsClient, sbClient)

				case currentScene == "scene-9" && changeScene():
					resetCount("scene-10", obsClient, sbClient)

				case currentScene == "scene-10" && changeScene():
					resetCount("scene-11", obsClient, sbClient)

				case currentScene == "scene-11" && changeScene():
					resetCount("scene-12", obsClient, sbClient)

				case currentScene == "scene-12" && changeScene():
					resetCount("scene-13", obsClient, sbClient)

				case currentScene == "scene-13" && changeScene():
					resetCount("scene-14", obsClient, sbClient)

				case currentScene == "scene-14" && changeScene():
					resetCount("scene-15", obsClient, sbClient)

				default:
					resetCount("SceneBRB", obsClient, sbClient)
				}
			}
		}
	}()

	// some logic here
	go func() {
		for {
			err = sbClient.Queue.ReceiveOne(context.Background(), MessagePrinter{cc: obsClient, aa: sbClient})
			if err != nil {
				log.Fatalln("receive error : ", err)
			}
		}
	}()

	for {

	}
	log.Println("All Good")

}
func resetCount(scene string, client *obs.OBSClient, q *azure.SBClient) {
	log.Println("changing scene ", scene)
	currentScene = scene
	currentSceneRight = 0
	currentSceneWrong = 0
	if scene == "scene-1" {
		err := q.SendQueue.Send(context.Background(), servicebus.NewMessageFromString(scene))
		if err != nil {
			log.Println("error with sending")
		}
		client.SwitchScene("SceneStarting")
		time.Sleep(time.Second * 10)
		client.SwitchScene(scene)
	}
	if scene == "SceneBRB" {
		err := q.SendQueue.Send(context.Background(), servicebus.NewMessageFromString(scene))
		if err != nil {
			log.Println("error with sending")
		}
		client.SwitchScene("SceneNextQuestion")
		time.Sleep(time.Second * 10)
		client.SwitchScene(scene)
	}
	if scene != "scene-complete" && scene != "scene-1" {
		err := q.SendQueue.Send(context.Background(), servicebus.NewMessageFromString(scene))
		if err != nil {
			log.Println("error with sending")
		}
		client.SwitchScene("SceneNextQuestion")
		time.Sleep(time.Second * 10)
		client.SwitchScene(scene)
	}
	log.Println("current Scene ", currentScene)
}
func changeScene() bool {
	return currentSceneRight > currentSceneWrong
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
	aa *azure.SBClient
}

func (mp MessagePrinter) Handle(ctx context.Context, msg *servicebus.Message) error {

	var mm Answers
	err := json.Unmarshal(msg.Data, &mm)
	if err != nil {
		log.Printf("unknown message %s \n", string(msg.Data))
	}

	log.Printf("message %v \n", mm)

	// start the stream
	if mm.Scene == "StartStream" {
		if err := mp.cc.StartStream(); err != nil {
			log.Printf("error : %v \n", err)
		}
		log.Println("... waiting for 15 seconds")
		time.Sleep(time.Second * 15)

		if err := mp.cc.SwitchScene("SceneBRB"); err != nil {
			log.Printf("error : %v \n", err)
		}
	}

	// start the trivia
	if mm.Scene == "StartTrivia" {
		resetCount("scene-1", mp.cc, mp.aa)
	}

	// stop the stream
	if mm.Scene == "StopStream" {

		if err := mp.cc.SwitchScene("SceneBRB"); err != nil {
			log.Printf("error : %v \n", err)
		}

		if err := mp.cc.StopStream(); err != nil {
			log.Printf("error : %v \n", err)
		}
	}
	mm.Scene = currentScene
	if dict[currentScene] == mm.Answer {
		currentSceneRight = currentSceneRight + 1
	} else {
		currentSceneWrong = currentSceneWrong + 1
	}

	return msg.Complete(ctx)
}

var currentScene = "SceneBRB"
var currentSceneRight = 0
var currentSceneWrong = 0

var dict = map[string]string{
	"scene-1":  "A",
	"scene-2":  "B",
	"scene-3":  "A",
	"scene-4":  "A",
	"scene-5":  "B",
	"scene-6":  "B",
	"scene-7":  "B",
	"scene-8":  "A",
	"scene-9":  "A",
	"scene-10": "A",
	"scene-11": "B",
	"scene-12": "B",
	"scene-13": "B",
	"scene-14": "B",
	"scene-15": "A",
}

type Answers struct {
	Scene  string `json:"scene"`
	Answer string `json:"answer"`
}
