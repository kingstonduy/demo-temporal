package bootstrap

import (
	"fmt"
	"log"

	"github.com/lengocson131002/go-clean-core/config"
	"go.temporal.io/sdk/client"
)

func GetTemporalClient(cfg config.Configure) *client.Client {
	c, err := client.Dial(client.Options{
		HostPort: fmt.Sprintf("%s:%s", GetConfig().Temporal.Host, GetConfig().Temporal.Port),
	})
	if err != nil {
		log.Fatalf("unable to create client, %v", err)
	}

	return &c
}
