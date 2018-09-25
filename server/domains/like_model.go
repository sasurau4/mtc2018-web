package domains

import (
	"context"
	"sync"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
)

// LikeRepo is basic operation unit for Like.
type LikeRepo interface {
	Insert(ctx context.Context, like *Like) (*Like, error)
	BulkInsert(ctx context.Context, like []Like) ([]Like, error)
	CountBySessionID(ctx context.Context, sessionID int64) (int, error)
}

// NewFakeLikeRepo returns new LikeRepo.
func NewFakeLikeRepo() (LikeRepo, error) {
	return &fakeLikeRepo{}, nil
}

type fakeLikeRepo struct {
	list []*Like

	mu sync.Mutex
}

func (repo *fakeLikeRepo) Insert(ctx context.Context, like *Like) (*Like, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	like.UUID = uuid.New().String()
	repo.list = append(repo.list, like)
	return like, nil
}

func (repo *fakeLikeRepo) BulkInsert(ctx context.Context, likes []Like) ([]Like, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, like := range likes {
		like.UUID = uuid.New().String()
		repo.list = append(repo.list, &like)
	}
	return likes, nil
}

func (repo *fakeLikeRepo) CountBySessionID(ctx context.Context, sessionID int64) (int, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var count int
	for _, like := range repo.list {
		if like.SessionID == sessionID {
			count++
		}
	}

	return count, nil
}

// NewLikeRepo returns new LikeRepo.
func NewLikeRepo(spannerClient *spanner.Client) (LikeRepo, error) {
	return &likeRepo{spanner: spannerClient}, nil
}

type likeRepo struct {
	spanner *spanner.Client
}

func (repo *likeRepo) Insert(ctx context.Context, like *Like) (*Like, error) {
	like.UUID = uuid.New().String()
	_, err := repo.spanner.Apply(ctx, []*spanner.Mutation{like.Insert(ctx)})
	if err != nil {
		return nil, err
	}
	return like, nil
}

func (repo *likeRepo) BulkInsert(ctx context.Context, likes []Like) ([]Like, error) {
	muts := make([]*spanner.Mutation, 0, len(likes))
	for _, like := range likes {
		like.UUID = uuid.New().String()
		muts = append(muts, like.Insert(ctx))
	}

	_, err := repo.spanner.Apply(ctx, muts)
	if err != nil {
		return nil, err
	}
	return likes, nil
}

func (repo *likeRepo) CountBySessionID(ctx context.Context, sessionID int64) (int, error) {
	stmt := spanner.Statement{
		SQL: `SELECT COUNT(*) AS CNT FROM Likes WHERE SessionId = @sessionID`,
		Params: map[string]interface{}{
			"sessionID": sessionID,
		},
	}

	iter := repo.spanner.Single().Query(ctx, stmt)
	defer iter.Stop()

	row, err := iter.Next()
	if err != nil {
		return 0, err
	}

	var count int
	if err := row.ColumnByName("CNT", &count); err != nil {
		return 0, err
	}
	return count, nil
}
