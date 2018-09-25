package gqlapi

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mercari/mtc2018-web/server/domains"
)

var _ SessionResolver = (*sessionResolver)(nil)

func newSessionResolver(rootResolver *rootResolver) SessionResolver {
	return &sessionResolver{
		rootResolver: rootResolver,

		likeCountCache: make(map[int]sessionLikeCache),
		now:            time.Now,
	}
}

type sessionResolver struct {
	*rootResolver

	lock           sync.Mutex
	likeCountCache map[int]sessionLikeCache
	now            func() time.Time
}

type sessionLikeCache struct {
	expireAt time.Time
	value    int
}

func (r *sessionResolver) ID(ctx context.Context, obj *domains.Session) (string, error) {
	return fmt.Sprintf("Session:%d", obj.ID), nil
}

func (r *sessionResolver) Liked(ctx context.Context, obj *domains.Session) (int, error) {
	if obj == nil {
		return 0, nil
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	cache, ok := r.likeCountCache[obj.ID]
	if ok && r.now().Before(cache.expireAt) {
		return cache.value, nil
	}

	count, err := r.likeRepo.CountBySessionID(ctx, int64(obj.ID))
	if err != nil {
		return 0, err
	}

	r.likeCountCache[obj.ID] = sessionLikeCache{
		expireAt: r.now().Add(r.cfg.SessionLikedCacheExpired),
		value:    count,
	}

	return count, nil
}

func (r *sessionResolver) Speakers(ctx context.Context, obj *domains.Session) ([]domains.Speaker, error) {
	if obj == nil {
		return nil, nil
	}

	speakerIDs := make([]string, 0, len(obj.SpeakerIDs))
	for _, id := range obj.SpeakerIDs {
		speakerIDs = append(speakerIDs, id)
	}

	sessionList, err := r.speakerRepo.Get(ctx, speakerIDs...)
	if err != nil {
		return nil, err
	}

	var resp []domains.Speaker
	for _, session := range sessionList {
		resp = append(resp, *session)
	}

	return resp, nil
}
