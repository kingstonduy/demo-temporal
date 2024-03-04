package temporal1

import (
	"log"
	"orchestrator-service/domain"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type Compensations struct {
	compensations []any
	arguments     [][]any
}

func NewCompensations() domain.Compensations {
	return &Compensations{
		compensations: make([]any, 0),
		arguments:     make([][]any, 0),
	}
}

func (s *Compensations) AddCompensation(activity any, parameters ...any) {
	s.compensations = append(s.compensations, activity)
	s.arguments = append(s.arguments, parameters)
}

func (s Compensations) Compensate(ctx workflow.Context, inParallel bool) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 10,
		RetryPolicy:         &temporal.RetryPolicy{MaximumAttempts: 0},
	}
	ctx = workflow.WithActivityOptions(ctx, options)
	if !inParallel {
		// Compensate in Last-In-First-Out order, to undo in the reverse order that activies were applied.
		for i := len(s.compensations) - 1; i >= 0; i-- {
			errCompensation := workflow.ExecuteActivity(ctx, s.compensations[i], s.arguments[i]...).Get(ctx, nil)
			if errCompensation != nil {
				workflow.GetLogger(ctx).Error("Executing compensation failed", "Error", errCompensation)
			}
		}
	} else {
		log.Println("🔥🔥🔥🔥🔥🔥")
		selector := workflow.NewSelector(ctx)
		for i := 0; i < len(s.compensations); i++ {
			execution := workflow.ExecuteActivity(ctx, s.compensations[i], s.arguments[i]...)
			selector.AddFuture(execution, func(f workflow.Future) {
				if errCompensation := f.Get(ctx, nil); errCompensation != nil {
					workflow.GetLogger(ctx).Error("Executing compensation failed", "Error", errCompensation)
				}
			})
		}
		for range s.compensations {
			selector.Select(ctx)
		}

	}
}
