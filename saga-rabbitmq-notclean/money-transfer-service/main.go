package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"saga-rabbitmq-notclean/config"
	model "saga-rabbitmq-notclean/money-transfer-service/shared"

	"github.com/gofiber/fiber/v2"
	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
)

// request, response to client. starts a workflow then wait for the workflow to finish
func Handler(c client.Client) fiber.Handler {
	fn := func(ctx *fiber.Ctx) error {
		var clientReq model.CLientRequest
		err := ctx.BodyParser(&clientReq)
		if err != nil {
			log.Println(err.Error())
			ctx.Status(fiber.StatusBadRequest).JSON(&model.SaferResponse{
				Code:    fiber.StatusInternalServerError,
				Message: err.Error(),
			})

		}

		var workflowInput = &model.WorkflowInput{
			TransactionID: uuid.New(),
			FromAccountID: clientReq.FromAccountID,
			ToAccountID:   clientReq.ToAccountID,
			Amount:        clientReq.Amount,
		}

		options := client.StartWorkflowOptions{
			ID:        config.GetConfig().Temporal.Workflow + "-" + workflowInput.TransactionID,
			TaskQueue: config.GetConfig().Temporal.TaskQueue,
		}

		we, err := c.ExecuteWorkflow(context.Background(), options, config.GetConfig().Temporal.Workflow, workflowInput)
		if err != nil {
			ctx.Status(500).JSON(model.SaferResponse{
				WorkflowID: workflowInput.TransactionID,
				RunID:      we.GetID(),
				Code:       http.StatusInternalServerError,
				Message:    err.Error(),
			})
			return err
		}

		var clientResponse model.ClientResponse
		err = we.Get(context.Background(), &clientResponse)
		if err != nil {
			ctx.Status(500).JSON(model.SaferResponse{
				WorkflowID: workflowInput.TransactionID,
				RunID:      we.GetID(),
				Code:       http.StatusInternalServerError,
				Message:    err.Error(),
			})
			return err
		}

		ctx.Status(200).JSON(clientResponse)
		return nil
	}
	return fn
}

func main() {
	config := config.GetConfig()

	c, err := client.Dial(client.Options{
		HostPort: fmt.Sprintf("%s:%s", config.Temporal.Host, config.Temporal.Port),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	g := fiber.New()
	publicRouter := g.Group("/api/v1")
	publicRouter.Post("/moneytransfer", Handler(c))

	g.Listen(fmt.Sprintf("%s:%s", config.MoneyTransfer.Host, config.MoneyTransfer.Port))
}
