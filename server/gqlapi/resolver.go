//go:generate gqlgen

package gqlapi

import (
	"sync"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/mercari/mtc2018-web/server/domains"
	"go.uber.org/zap"
)

// ResolverConfig provides configuration for ResolverRoot.
type ResolverConfig struct {
	Logger                   *zap.Logger
	SpannerClient            *spanner.Client
	SessionLikedCacheExpired time.Duration
}

// NewResolver returns GraphQL root resolver.
func NewResolver(cfg *ResolverConfig) (ResolverRoot, error) {

	sessionRepo, err := domains.NewSessionRepo()
	if err != nil {
		return nil, err
	}
	speakerRepo, err := domains.NewSpeakerRepo()
	if err != nil {
		return nil, err
	}

	newsRepo, err := domains.NewNewsRepo()
	if err != nil {
		return nil, err
	}

	var (
		likeRepo    domains.LikeRepo
		likeSumRepo domains.LikeSummaryRepo
	)
	if cfg.SpannerClient != nil {
		var err error
		likeRepo, err = domains.NewLikeRepo(cfg.SpannerClient)
		if err != nil {
			return nil, err
		}
		likeSumRepo, err = domains.NewLikeSummaryRepo(cfg.SpannerClient)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		likeRepo, err = domains.NewFakeLikeRepo()
		if err != nil {
			return nil, err
		}
		likeSumRepo, err = domains.NewFakeLikeSummaryRepo()
		if err != nil {
			return nil, err
		}
	}

	storer := newStorer(cfg.Logger, likeRepo, likeSumRepo)
	go storer.Run()

	eventCh := make(chan likeEvent, 128)

	observer := newObserver(cfg.Logger, eventCh, likeSumRepo)
	go observer.Run()

	listener := newListener(cfg.Logger, eventCh, sessionRepo)
	go listener.Run()

	r := &rootResolver{
		cfg: cfg,

		sessionRepo: sessionRepo,
		speakerRepo: speakerRepo,
		likeRepo:    likeRepo,
		newsRepo:    newsRepo,
		likeSumRepo: likeSumRepo,

		Logger:   cfg.Logger,
		storer:   storer,
		observer: observer,
		listener: listener,
	}
	r.sessionResolver = newSessionResolver(r)

	return r, nil
}

type rootResolver struct {
	cfg *ResolverConfig

	sessionRepo domains.SessionRepo
	speakerRepo domains.SpeakerRepo
	likeRepo    domains.LikeRepo
	newsRepo    domains.NewsRepo
	likeSumRepo domains.LikeSummaryRepo

	Logger *zap.Logger
	mu     sync.Mutex

	storer   *storer
	observer *observer
	listener *listener

	sessionResolver SessionResolver
}

func (r *rootResolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *rootResolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *rootResolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

func (r *rootResolver) Session() SessionResolver {
	return r.sessionResolver
}

func (r *rootResolver) Speaker() SpeakerResolver {
	return &speakerResolver{r}
}

func (r *rootResolver) Like() LikeResolver {
	return &likeResolver{r}
}

func (r *rootResolver) News() NewsResolver {
	return &newsResolver{r}
}
