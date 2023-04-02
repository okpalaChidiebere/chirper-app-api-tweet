package tweetsservice

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	repo "github.com/okpalaChidiebere/chirper-app-api-tweet/v0/tweets/data_access"
	model "github.com/okpalaChidiebere/chirper-app-api-tweet/v0/tweets/model"
)

type ServiceImpl struct {
	repo repo.Repository
}

func New(repo repo.Repository) *ServiceImpl {
	return &ServiceImpl{repo}
}

func (s *ServiceImpl) SaveTweet(ctx context.Context, tweet *model.Tweet) (*model.Tweet, error){
	var replyingToAuthor string
	if tweet.Author == "" {
		return nil, errors.New("author is required")
	}

	if tweet.ReplyingTo != "" {		
		tokens := strings.Split(tweet.ReplyingTo, ":")

		if len(tokens) != 2{
			return nil, errors.New("invalid format for replyingTo. It should be eg {reply_tweet_id}:{reply_tweet_author}")
		}
		tweet.ReplyingTo = tokens[0]
		replyingToAuthor = tokens[1]
	}

	if tweet.Timestamp.IsZero() {
			tweet.Timestamp =  model.ChirperAppUnixTime(time.Now())
	}

	if tweet.Id == "" {
		tweet.Id = uuid.NewString()
	}

	newTweet, err := s.repo.SaveTweetToDynamoDb(ctx, replyingToAuthor,tweet)
	if err != nil {
		return nil, err
	}
	return newTweet, nil
}

func (s *ServiceImpl) BulkSaveTweet(ctx context.Context, tweets []*model.Tweet) error {
	if len(tweets) == 0 {
		return errors.New("cannot perform action on an empty list")
	}

	for _, tweet := range tweets {
		if tweet.Author == "" {
			return fmt.Errorf("author is required for tweetID: %s", tweet.Id)
		}
	}

	for i := range tweets {
		if tweets[i].Timestamp.IsZero() {
			tweets[i].Timestamp =  model.ChirperAppUnixTime(time.Now())
		}

		if tweets[i].Id == "" {
			tweets[i].Id = uuid.NewString()
		}
			
	}

	if err := s.repo.BulkSaveTweetToDynamoDb(ctx, tweets); err != nil {
		return err
	}
	return nil
}

func (s *ServiceImpl) ListTweets(ctx context.Context, limit int32, nextKey string) ([]*model.Tweet, string, error) {
	if (limit <= 0){
		limit = 10
	} else if limit > 30 {
		return nil, "", errors.New("limit cannot be more than 30")
	}

	return s.repo.ScanTweetsFromDynamoDb(ctx, limit, nextKey)
}

func (s *ServiceImpl) SaveLikeToggle(ctx context.Context, tweetID, author, authedUserID string, hasLiked bool) error {
	if tweetID == "" {
		return errors.New("id is required")
	}
	if author == "" {
		return errors.New("author is required")
	}
	if authedUserID == "" {
		return errors.New("authedUserID is required")
	}
	return s.repo.SaveLikeToggleInDynamoDb(ctx, tweetID, author, authedUserID, hasLiked)
}