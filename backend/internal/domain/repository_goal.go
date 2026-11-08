package domain

import "context"

type GoalRepository interface {
	List(ctx context.Context) ([]Goal, error)
	Get(ctx context.Context, id string) (*Goal, error)
	Create(ctx context.Context, input CreateGoalInput) (*Goal, error)
	Update(ctx context.Context, id string, input UpdateGoalInput) (*Goal, error)
	Delete(ctx context.Context, id string) error
	Contribute(ctx context.Context, goalID string, input ContributeGoalInput) (*Goal, error)
	Reach(ctx context.Context, id string) (*Goal, error)
	Contributions(ctx context.Context, goalID string) ([]GoalContribution, error)
}
