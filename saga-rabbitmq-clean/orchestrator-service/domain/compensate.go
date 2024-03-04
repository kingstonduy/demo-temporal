package domain

import (
	"go.temporal.io/sdk/workflow"
)

type Compensations interface {
	AddCompensation(activity any, parameters ...any)
	Compensate(ctx workflow.Context, inParallel bool)
}
