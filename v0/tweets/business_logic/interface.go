package tweetsservice

import (
	"context"

	model "github.com/okpalaChidiebere/chirper-app-api-tweet/v0/tweets/model"
)

//go:generate mockgen -destination mock.go -source=interface.go -package=tweetsservice
type Service interface {
	SaveTweet(ctx context.Context, tweet *model.Tweet) (*model.Tweet, error)
	BulkSaveTweet(ctx context.Context, tweets []*model.Tweet) error
	ListTweets(ctx context.Context, limit int32, nextKey string) ([]*model.Tweet, string, error)
	SaveLikeToggle(ctx context.Context, tweetID, author, authedUserID string, hasLiked bool) error
}