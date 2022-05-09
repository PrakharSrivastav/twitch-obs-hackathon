package obs

import (
	"github.com/andreykaipov/goobs"
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

	return &OBSClient{Client: client}, nil
}
