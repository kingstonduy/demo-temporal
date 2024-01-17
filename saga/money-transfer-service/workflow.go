package app

// func MoneyTransferWorkflow(ctx workflow.Context, info TransactionInfo) (err error) {
// 	options := workflow.ActivityOptions{
// 		StartToCloseTimeout: time.Second * 5,
// 		RetryPolicy:         &temporal.RetryPolicy{MaximumAttempts: 2},
// 	}
// 	ctx = workflow.WithActivityOptions(ctx, options)

// 	// var compensations Compensations

// 	defer func() {
// 		if err != nil {
// 			// activity failed, and workflow context is canceled
// 			// disconnectedCtx, _ := workflow.NewDisconnectedContext(ctx)
// 			// compensations.Compensate(disconnectedCtx, true)
// 		}
// 	}()

// 	err = workflow.ExecuteActivity(ctx, ValidateRecipientActivity, info).Get(ctx, nil)
// 	if err != nil {
// 		return err
// 	}

// 	return err
// }
