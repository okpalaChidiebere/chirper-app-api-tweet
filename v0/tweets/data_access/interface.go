package tweetsdataaccess

import (
	"context"

	model "github.com/okpalaChidiebere/chirper-app-api-tweet/v0/tweets/model"
)

//go:generate mockgen -destination mock.go -source=interface.go -package=tweetsdataaccess
type Repository interface {
	//Creates a new tweet or replaces a an old tweet with a new tweet(if the tweet id exists) in the in tweets table.
	SaveTweetToDynamoDb(ctx context.Context, replyingToAuthor string, tweet *model.Tweet) (*model.Tweet, error)
	//returns the list of tweets
	ListTweetsFromDynamoDb(ctx context.Context, authedUserID, nextKey string, limit int32) (results []*model.Tweet, nextCursor string, err error)
	//get a tweet by ID
	GetTweetFromDynamoDb(ctx context.Context, tweetID string) (*model.Tweet, error)
	//Update
	// UpsertTweetFromDynamoDb(ctx context.Context, tweet *model.Tweet) error
	//toogle likes
	SaveLikeToggleInDynamoDb(ctx context.Context, tweetID, author, authedUserID string, hasLiked bool) error
	//scan
	ScanTweetsFromDynamoDb(ctx context.Context, limit int32, nextKey string) ([]*model.Tweet, string, error)
	//Multi Create or replace
	BulkSaveTweetToDynamoDb(ctx context.Context, tweets []*model.Tweet) error
}