package obs

import (
	"fmt"
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/events"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/andreykaipov/goobs/api/requests/streaming"
	"log"
	"os"
)

type OBSClient struct {
	Client *goobs.Client
}

func NewClient() (*OBSClient, error) {
	client, err := goobs.New(
		os.Getenv("WSL_HOST")+":4444",
		goobs.WithPassword("hello"),                   // optional
		goobs.WithDebug(os.Getenv("OBS_DEBUG") != ""), // optional
	)

	if err != nil {
		log.Println("connection error")
		return nil, err
	}
	go func() {
		for event := range client.IncomingEvents {
			switch e := event.(type) {
			case *events.SourceVolumeChanged:
				fmt.Printf("Volume changed for %-25q: %f\n", e.SourceName, e.Volume)
			case *events.StreamStarted:
				fmt.Println("Stream started successfully")
			case *events.StreamStopped:
				fmt.Println("Stream stopped successfully")
			case *events.ScenesChanged:
				fmt.Println("Scene changed successfully")

				//default:
				//log.Println("--", e.GetUpdateType())
				//log.Println("--")
			}
		}
	}()

	return &OBSClient{Client: client}, nil
}

func (o *OBSClient) SwitchScene(sceneName string) error {
	resp, err := o.Client.Scenes.SetCurrentScene(&scenes.SetCurrentSceneParams{SceneName: sceneName})
	if err != nil {
		log.Println("scene change error :", err)
		return err
	}

	log.Printf("%v \n", resp)
	return nil
}

func (o *OBSClient) StartStream() error {
	resp, err := o.Client.Streaming.StartStreaming(&streaming.StartStreamingParams{})
	if err != nil {
		log.Println("start stream error :", err)
		return err
	}

	log.Printf("%v \n", resp)
	return nil
}

func (o *OBSClient) StopStream() error {
	resp, err := o.Client.Streaming.StopStreaming(&streaming.StopStreamingParams{})
	if err != nil {
		log.Println("stop stream error :", err)
		return err
	}

	log.Printf("%v \n", resp)
	return nil
}
