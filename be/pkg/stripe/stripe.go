package stripe_client

import (
	"context"

)

type (
	Stripe interface {
		CreatePayment(ctx context.Context) (err error)
	}

	stripeClient struct {
		
	}
)
