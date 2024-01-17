package main

func main() {
	// c, err := client.Dial(client.Options{})
	// if err != nil {
	// 	log.Fatalln("Unable to create client", err)
	// }
	// defer c.Close()

	// option := client.StartWorkflowOptions{
	// 	ID:        shared.Workflow + "_" + uuid.New(),
	// 	TaskQueue: shared.TaskQueue,
	// }

	// transferInfo := shared.TransactionInfo{
	// 	TransactionId: uuid.New(),
	// 	FromAccountId: "1",
	// 	ToAccountId:   "2",
	// 	Amount:        100,
	// }

	// _, err = c.ExecuteWorkflow(context.Background(), option, app.MoneyTransferWorkflow, transferInfo)
	// if err != nil {
	// 	log.Fatalf("Unable to execute %s workflow\n", option.ID, err)
	// }
}
