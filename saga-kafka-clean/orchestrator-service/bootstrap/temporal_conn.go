package bootstrap

import (
	"fmt"
	"log"

	"go.temporal.io/sdk/client"
)

func (c *Config) GetTemporalClient() *client.Client {
	var temporalClient client.Client = nil
	if temporalClient != nil {
		return &temporalClient
	}
	temporalClient, err := client.Dial(client.Options{
		HostPort: fmt.Sprintf("%s:%s", c.Temporal.Host, c.Temporal.Port),
	})

	if err != nil {
		log.Fatalf("unable to create client, %v", err)
	}

	return &temporalClient
}
