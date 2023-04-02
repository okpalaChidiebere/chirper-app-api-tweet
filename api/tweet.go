package api

import (
	"context"
	"log"
	"time"

	apiadapters "github.com/okpalaChidiebere/chirper-app-api-tweet/api/adapters"
	tweetsservice "github.com/okpalaChidiebere/chirper-app-api-tweet/v0/tweets/business_logic"
	model "github.com/okpalaChidiebere/chirper-app-api-tweet/v0/tweets/model"
	pb "github.com/okpalaChidiebere/chirper-app-gen-protos/tweet/v1"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type TweetServer struct {
	// pb.UnimplementedTweetServiceServer
	TweetService tweetsservice.Service
}

func NewTweetServer(tweetsService tweetsservice.Service) pb.TweetServiceServer {
	return &TweetServer{ 
		TweetService: tweetsService,
	}
}

func (s *TweetServer) SaveTweet(ctx context.Context, req *pb.SaveTweetRequest) (*pb.SaveTweetResponse, error) {
	t := &model.Tweet{
		Id: req.GetId(),
		Author: req.GetAuthor(),
		Likes: req.GetReplies(),
		Replies: req.GetReplies(),
		Text: req.GetText(),
		Timestamp: model.ChirperAppUnixTime(time.UnixMilli(req.GetTimestamp())),
		ReplyingTo: req.GetReplyingTo(),
	}

	tweet, err := s.TweetService.SaveTweet(ctx, t)
	if err != nil {
		log.Printf("SaveTweet Err: %v", err.Error())
		return nil, err
	}

	return &pb.SaveTweetResponse{Tweet: apiadapters.TweetToProto(tweet) } , nil
}

func (s *TweetServer) ListTweets(ctx context.Context, req *pb.ListTweetsRequest) (*pb.ListTweetsResponse, error){
	tweets, nk, err := s.TweetService.ListTweets(ctx, req.GetLimit(), req.GetNextKey())
	if err != nil {
		log.Printf("ListTweets Err: %v", err.Error())
		return nil, err
	}

	todos := apiadapters.TweetsToProto(tweets)
	return &pb.ListTweetsResponse{ Items: todos, NextKey: nk } , nil
}

func (s *TweetServer) SaveLikeToggle(ctx context.Context, req *pb.SaveLikeToggleRequest) (*emptypb.Empty, error){
	err := s.TweetService.SaveLikeToggle(ctx, req.GetId(), req.GetAuthor(), req.GetAuthedUserId(), req.GetHasLiked())
	if err != nil {
		log.Printf("SaveLikeToggle Err: %v", err.Error())
		return nil, err
	}
	
	return &emptypb.Empty{}, nil
}
