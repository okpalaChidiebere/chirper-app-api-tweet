package tweetsservice

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	tweetsrepo "github.com/okpalaChidiebere/chirper-app-api-tweet/v0/tweets/data_access"
	model "github.com/okpalaChidiebere/chirper-app-api-tweet/v0/tweets/model"
)

func Test_SaveTweet(t *testing.T) {
		testCases := []struct {
		name          string
		tweet *model.Tweet
		replyingToAuthor string

		expectedRepoResp *model.Tweet
		repoError error
		expectedError error

		expectedRepoCallTimes int
	}{
		{
			name: "should return error with repo error",
			tweet: &model.Tweet{Id: "SomeID", Author: "some_handle", ReplyingTo: "tweetID:another_author",},
			repoError: errors.New("repo error"),

			expectedError: errors.New("repo error"),
			replyingToAuthor: "another_author",

			expectedRepoCallTimes: 1,

		},
		{
			name: "should return error when replyingToAuthor tweet params format is invalid",
			tweet: &model.Tweet{Id: "SomeID", Author: "some_handle", ReplyingTo: "tweetID-I-am-ReplyingTo"},
			repoError:  errors.New("invalid format for replyingTo. It should be eg {reply_tweet_id}:{reply_tweet_author}"),

			expectedError:  errors.New("invalid format for replyingTo. It should be eg {reply_tweet_id}:{reply_tweet_author}"),
			expectedRepoCallTimes: 0,
		},
		{
			name: "should return error with author ID not provided",
			tweet: &model.Tweet{Id: "SomeID"},

			repoError: errors.New("author is required"),
			expectedError: errors.New("author is required"),
			
			expectedRepoCallTimes: 0,

		},
		{
			name: "should return no error with repo doesn't error",
			tweet: &model.Tweet{Id: "SomeID", Author: "some_handle", ReplyingTo: "tweetID:another_author"},
			replyingToAuthor: "another_author",
		 	expectedRepoResp:  &model.Tweet{Id: "SomeID", Author: "some_handle", ReplyingTo: "tweetID"},

			repoError: nil,
			expectedError: nil,

			expectedRepoCallTimes: 1,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)

			repoMock := tweetsrepo.NewMockRepository(ctrl)
			
			repoMock.EXPECT().SaveTweetToDynamoDb(ctx, tc.replyingToAuthor, tc.tweet).Times(tc.expectedRepoCallTimes).Return(tc.expectedRepoResp, tc.repoError)

			service := New(repoMock)
			_, err := service.SaveTweet(ctx, tc.tweet)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func Test_BatchSaveTweet(t *testing.T) {
		testCases := []struct {
		name          string
		tweets []*model.Tweet

		buildStubs func(ctx context.Context, tweets []*model.Tweet, repoMock *tweetsrepo.MockRepository)

		expectedError error
	}{
		{
			name: "should return error with repo error",
			tweets: []*model.Tweet{
				{Id: "SomeID1", Author: "some_handle1"},
				{Id: "SomeID2", Author: "some_handle2"},
			},
			buildStubs: func(ctx context.Context, tweets []*model.Tweet, repoMock *tweetsrepo.MockRepository) {
				repoMock.EXPECT().BulkSaveTweetToDynamoDb(ctx, tweets).Times(1).Return(errors.New("repo error"))
			},
			expectedError: errors.New("repo error"),
		},
		{
			name: "should return error with author ID not provided",
			tweets: []*model.Tweet{
				{Id: "SomeID1", Author: "some_handle1"},
				{Id: "SomeID2"},
			},
			buildStubs: func(ctx context.Context, tweets []*model.Tweet, repoMock *tweetsrepo.MockRepository) {
				repoMock.EXPECT().BulkSaveTweetToDynamoDb(ctx, tweets).Times(0)
			},
			expectedError: errors.New("author is required for tweetID: SomeID2"),
		},
		{
			name: "should return no error with repo doesn't error",
			tweets: []*model.Tweet{
				{Id: "SomeID1", Author: "some_handle1"},
				{Id: "SomeID2", Author: "some_handle2"},
			},
			buildStubs: func(ctx context.Context, tweets []*model.Tweet, repoMock *tweetsrepo.MockRepository) {
				repoMock.EXPECT().BulkSaveTweetToDynamoDb(ctx, tweets).Times(1).Return(nil)
			},
			expectedError: nil,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)

			repoMock := tweetsrepo.NewMockRepository(ctrl)

			tc.buildStubs(ctx, tc.tweets, repoMock)

			service := New(repoMock)
			err := service.BulkSaveTweet(ctx, tc.tweets)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}