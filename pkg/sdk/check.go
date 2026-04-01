package sdk

import "context"

type Check interface {
	ID() string
	Name() string
	Run(ctx context.Context) (*Finding, error)
}
