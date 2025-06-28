package consumer_interfaces

import "context"

type OnboardingConsumer interface {
	VerifyOnboardingToken(ctx context.Context) error
	UserPointUpdate(ctx context.Context) error
}
