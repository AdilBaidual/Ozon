package graph

import (
	"Service/internal/core/model"
	"Service/internal/core/usecase"
	"context"
)

func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }

func (s subscriptionResolver) Notification(ctx context.Context, postID int) (<-chan *model.Comment, error) {
	ch := make(chan *model.Comment)

	client := usecase.NewSubscriptionClient(ch)
	s.coreUC.AddClient(client)

	go func() {
		defer close(ch)
		s.coreUC.Listen(ctx, client)
	}()

	client.Subscribe(postID)

	return ch, nil
}
