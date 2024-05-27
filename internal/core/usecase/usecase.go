package usecase

import (
	"Service/internal/core"
	"Service/internal/core/model"
	"context"
	"github.com/pkg/errors"
	"sync"

	"go.uber.org/zap"
)

const MaxContentLen = 2000

type SubscriptionManager struct {
	clients map[*Client]map[int]struct{}
	mu      *sync.Mutex
}

func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		clients: make(map[*Client]map[int]struct{}),
		mu:      &sync.Mutex{},
	}
}

type UC struct {
	lg                  *zap.Logger
	repo                core.Repo
	SubscriptionManager *SubscriptionManager
}

func NewUseCase(logger *zap.Logger, repo core.Repo) *UC {
	return &UC{
		lg:                  logger,
		repo:                repo,
		SubscriptionManager: NewSubscriptionManager(),
	}
}

func (u *UC) CreatePostUC(ctx context.Context, params model.NewPost, authorUUID string) (model.Post, error) {
	logger, ok := ctx.Value("logger").(*zap.Logger)
	if !ok {
		logger = u.lg
	}
	var post model.Post

	if len(params.Content) > MaxContentLen {
		err := errors.New("line length is too long")
		logger.Error("validating params", zap.Error(err))
		return post, err
	}

	postId, err := u.repo.CreatePost(ctx, params, authorUUID)
	if err != nil {
		logger.Error("creating post", zap.Error(err))
		return post, err
	}

	post, err = u.repo.GetPostById(ctx, postId)
	if err != nil {
		logger.Error("get post by id", zap.Error(err))
		return post, err
	}

	return post, nil
}

func (u *UC) CreateCommentUC(ctx context.Context, params model.NewComment, authorUUID string) (model.Comment, error) {
	logger, ok := ctx.Value("logger").(*zap.Logger)
	if !ok {
		logger = u.lg
	}
	var comment model.Comment

	if len(params.Content) > MaxContentLen {
		err := errors.New("line length is too long")
		logger.Error("validating params", zap.Error(err))
		return comment, err
	}

	post, err := u.repo.GetPostById(ctx, params.PostID)
	if err != nil {
		logger.Error("getting post", zap.Error(err))
		return comment, err
	}

	if !post.CommentsEnabled {
		logger.Error("comment for post not enabled", zap.Error(err))
		return comment, errors.New("comment for post not enabled")
	}

	commentId, err := u.repo.CreateComment(ctx, params, authorUUID)
	if err != nil {
		logger.Error("creating post", zap.Error(err))
		return comment, err
	}

	comment, err = u.repo.GetCommentById(ctx, commentId)
	if err != nil {
		logger.Error("get post by id", zap.Error(err))
		return comment, err
	}
	go u.SubscriptionManager.broadcast(comment.PostID, comment)
	return comment, nil
}

func (u *UC) GetPostsUC(ctx context.Context) ([]*model.Post, error) {
	logger, ok := ctx.Value("logger").(*zap.Logger)
	if !ok {
		logger = u.lg
	}
	posts, err := u.repo.GetPosts(ctx)
	if err != nil {
		logger.Error("get posts", zap.Error(err))
		return nil, err
	}
	return posts, err
}

func (u *UC) GetCommentsByPostIdUC(ctx context.Context, postId int) ([]*model.Comment, error) {
	logger, ok := ctx.Value("logger").(*zap.Logger)
	if !ok {
		logger = u.lg
	}
	comments, err := u.repo.GetCommentsByPostId(ctx, postId)
	if err != nil {
		logger.Error("get comments", zap.Error(err))
		return nil, err
	}
	return comments, err
}

func (u *UC) GetCommentsByParentIdUC(ctx context.Context, parentId int) ([]*model.Comment, error) {
	logger, ok := ctx.Value("logger").(*zap.Logger)
	if !ok {
		logger = u.lg
	}
	comments, err := u.repo.GetCommentsByParentId(ctx, parentId)
	if err != nil {
		logger.Error("get comments", zap.Error(err))
		return nil, err
	}
	return comments, err
}

func (u *UC) GetCommentByIdUC(ctx context.Context, id int) (model.Comment, error) {
	logger, ok := ctx.Value("logger").(*zap.Logger)
	if !ok {
		logger = u.lg
	}
	comment, err := u.repo.GetCommentById(ctx, id)
	if err != nil {
		logger.Error("get comment by id", zap.Error(err))
		return comment, err
	}
	return comment, err
}

func (u *UC) GetPostByIdUC(ctx context.Context, id int) (model.Post, error) {
	logger, ok := ctx.Value("logger").(*zap.Logger)
	if !ok {
		logger = u.lg
	}
	post, err := u.repo.GetPostById(ctx, id)
	if err != nil {
		logger.Error("get post by id", zap.Error(err))
		return post, err
	}
	return post, err
}

type Client struct {
	ID          string
	subscribe   chan int
	unsubscribe chan int
	events      chan *model.Comment
	done        chan struct{}
}

func (c *Client) Subscribe(event int) {
	c.subscribe <- event
}

func (u *UC) Listen(ctx context.Context, c *Client) {
	for {
		select {
		case eventType := <-c.subscribe:
			u.SubscriptionManager.addSubscription(c, eventType)
		case eventType := <-c.unsubscribe:
			u.SubscriptionManager.removeSubscription(c, eventType)
		case <-c.done:
			u.SubscriptionManager.removeClient(c)
			return
		case <-ctx.Done():
			u.SubscriptionManager.removeClient(c)
			return
		}
	}
}

func (u *UC) AddClient(client *Client) {
	u.SubscriptionManager.mu.Lock()
	defer u.SubscriptionManager.mu.Unlock()
	u.SubscriptionManager.clients[client] = make(map[int]struct{})
}

func (manager *SubscriptionManager) removeClient(client *Client) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	delete(manager.clients, client)
}

func (manager *SubscriptionManager) addSubscription(client *Client, eventType int) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.clients[client][eventType] = struct{}{}
}

func (manager *SubscriptionManager) removeSubscription(client *Client, eventType int) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	delete(manager.clients[client], eventType)
}

func (manager *SubscriptionManager) broadcast(eventType int, message model.Comment) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	for client, subscriptions := range manager.clients {
		if _, subscribed := subscriptions[eventType]; subscribed {
			select {
			case client.events <- &message:
			}
		}
	}
}

func NewSubscriptionClient(events chan *model.Comment) *Client {
	return &Client{
		subscribe:   make(chan int),
		unsubscribe: make(chan int),
		events:      events,
		done:        make(chan struct{}),
	}
}
