package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	tweetsservice "github.com/okpalaChidiebere/chirper-app-api-tweet/v0/tweets/business_logic"
	model "github.com/okpalaChidiebere/chirper-app-api-tweet/v0/tweets/model"
	tweet_v1 "github.com/okpalaChidiebere/chirper-app-gen-protos/tweet/v1"
	"github.com/stretchr/testify/assert"
)

func TestTweetsSever_SaveTweet(t *testing.T){
	testCases := []struct {
		name          string
		inputReq      *tweet_v1.SaveTweetRequest

		expectedNote  *model.Tweet
		note  *model.Tweet
		saveTweetError  error

		expectedResponse  *tweet_v1.SaveTweetResponse
		expectedError error
	}{
		{
			name: "service error returns error",
			inputReq: &tweet_v1.SaveTweetRequest{},

			note: &model.Tweet{
				Timestamp: model.ChirperAppUnixTime(time.UnixMilli(0)),
			},
			expectedNote: nil,

			saveTweetError:  errors.New("error"),
			expectedError:  errors.New("error"),
			expectedResponse: nil,
		},
		{
			name: "OK request",
			inputReq: &tweet_v1.SaveTweetRequest{},
			
			note: &model.Tweet{
				Timestamp: model.ChirperAppUnixTime(time.UnixMilli(0)),
			},
			expectedNote: &model.Tweet{
				Timestamp: model.ChirperAppUnixTime(time.UnixMilli(1518122597860)),
			},
			saveTweetError:  nil,

			expectedError:  nil,
			expectedResponse: &tweet_v1.SaveTweetResponse{
				Tweet: &tweet_v1.Tweet{
					Timestamp: 1518122597860,
				},
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			tweetsServiceMock := tweetsservice.NewMockService(ctrl)
			tweetsServiceMock.EXPECT().SaveTweet(gomock.Any(), tc.note).Return(tc.expectedNote, tc.saveTweetError).Times(1)

			s := NewTweetServer(tweetsServiceMock)

			got, err := s.SaveTweet(ctx, tc.inputReq)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResponse, got)
		})
	}
}

func TestTweetsSever_SaveLikeToggle(t *testing.T){
	testCases := []struct {
		name          string
		inputReq      *tweet_v1.SaveLikeToggleRequest

		tweetID string 
		author string
		authedUserID string
		hasLiked bool
		saveLikeToggleError error

		expectedError error
	}{
		{
			name: "service error returns error",
			inputReq: &tweet_v1.SaveLikeToggleRequest{},
			saveLikeToggleError:  errors.New("error"),
			expectedError:  errors.New("error"),
		},
		{
			name: "returns nil on success",
			inputReq: &tweet_v1.SaveLikeToggleRequest{},
			saveLikeToggleError: nil,
			expectedError: nil,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			tweetsServiceMock := tweetsservice.NewMockService(ctrl)
			tweetsServiceMock.EXPECT().SaveLikeToggle(gomock.Any(), tc.tweetID, tc.author, tc.authedUserID, tc.hasLiked).Return(tc.saveLikeToggleError)

			s := NewTweetServer(tweetsServiceMock)

			_, err := s.SaveLikeToggle(ctx, tc.inputReq)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestTweetsSever_ListTweets(t *testing.T){
	testCases := []struct {
		name          string
		inputReq      *tweet_v1.ListTweetsRequest

		limit   int32
		nextKey  string

		expectedTweets  []*model.Tweet
		expectedNextKey  string
		listError  error

		expectedResponse  *tweet_v1.ListTweetsResponse
		expectedError error
	}{
		{
			name: "returns error when service error",
			inputReq: &tweet_v1.ListTweetsRequest{},
			listError: errors.New("error"),
			expectedError: errors.New("error"),
		},
		{
			name: "properly converst Notes; OK request",
			inputReq:  &tweet_v1.ListTweetsRequest{},
			listError: nil,
			expectedError: nil,
			expectedTweets: []*model.Tweet{
				{ 
						Id: "8xf0y6ziyjabvozdd253nd",
						Text: "Shoutout to all the speakers I know for whom English is not a first language, but can STILL explain a concept well. It's hard enough to give a good talk in your mother tongue!",
						Author: "sarah_edo",
						Timestamp: model.ChirperAppUnixTime(time.UnixMilli(1518122597860)),
						Likes: []string{"tylermcginnis"},
						Replies: []string{"fap8sdxppna8oabnxljzcv", "3km0v4hf1ps92ajf4z2ytg"},
						ReplyingTo: "",
					},
			},
			expectedResponse: &tweet_v1.ListTweetsResponse{
				Items: []*tweet_v1.Tweet{
					{
						Id: "8xf0y6ziyjabvozdd253nd",
						Text: "Shoutout to all the speakers I know for whom English is not a first language, but can STILL explain a concept well. It's hard enough to give a good talk in your mother tongue!",
						Author: "sarah_edo",
						Timestamp: 1518122597860,
						Likes: []string{"tylermcginnis"},
						Replies: []string{"fap8sdxppna8oabnxljzcv", "3km0v4hf1ps92ajf4z2ytg"},
						ReplyingTo: "",
					},
				},
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			tweetsServiceMock := tweetsservice.NewMockService(ctrl)
			tweetsServiceMock.EXPECT().ListTweets(gomock.Any(), tc.limit, tc.nextKey).Return(tc.expectedTweets, tc.expectedNextKey, tc.listError).Times(1)

			s := NewTweetServer(tweetsServiceMock)

			got, err := s.ListTweets(ctx, tc.inputReq)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResponse, got)
		})
	}
}
