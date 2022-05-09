package main

import (
	"fmt"
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/events"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"log"
	"os"
)

func main() {
	client, err := goobs.New(
		os.Getenv("WSL_HOST")+":4444",
		goobs.WithPassword("hello"),                   // optional
		goobs.WithDebug(os.Getenv("OBS_DEBUG") != ""), // optional
	)
	if err != nil {
		panic(err)
	}

	go func() {
		for event := range client.IncomingEvents {
			switch e := event.(type) {
			case *events.SourceVolumeChanged:
				fmt.Printf("Volume changed for %-25q: %f\n", e.SourceName, e.Volume)
			default:
				log.Printf("Unhandled event: %#v", e.GetUpdateType())
			}
		}
	}()

	list, err := client.Scenes.GetSceneList()
	if err != nil {
		log.Fatal("XXXX", err)
	}

	log.Printf("%v", list)

	scene, err := client.Scenes.SetCurrentScene(&scenes.SetCurrentSceneParams{SceneName: "test"})
	if err != nil {
		log.Fatal("YYYY", err)
	}
	log.Printf("%v", scene)

	for {
	}

	/*fmt.Println("Setting random volumes for each source...")

	rand.Seed(time.Now().UnixNano())
	list, _ := client.Sources.GetSourcesList()

	for _, v := range list.Sources {
		if _, err := client.Sources.SetVolume(&sources.SetVolumeParams{
			Source: v.Name,
			Volume: rand.Float64(),
		}); err != nil {
			panic(err)
		}
	}

	if len(list.Sources) == 0 {
		fmt.Println("No sources!")
		os.Exit(0)
	}

	fmt.Println("Test toggling the mute status of the first source...")

	name := list.Sources[0].Name
	resp, _ := client.Sources.GetVolume(&sources.GetVolumeParams{Source: name})
	fmt.Printf("%s is muted? %t\n", name, resp.Muted)

	_, _ = client.Sources.ToggleMute(&sources.ToggleMuteParams{Source: name})
	resp, _ = client.Sources.GetVolume(&sources.GetVolumeParams{Source: name})
	fmt.Printf("%s is muted? %t\n", name, resp.Muted)*/
}
