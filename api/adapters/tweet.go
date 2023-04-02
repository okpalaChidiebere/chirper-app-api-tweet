package api_adapters

import (
	"time"

	model "github.com/okpalaChidiebere/chirper-app-api-tweet/v0/tweets/model"
	pb "github.com/okpalaChidiebere/chirper-app-gen-protos/tweet/v1"
)



func TweetToProto (t *model.Tweet) *pb.Tweet{
	return &pb.Tweet{
		Id: t.Id,
		Author: t.Author,
		Likes: t.Likes,
		Replies: t.Replies,
		Text: t.Text,
		/*
		FYI: 
		Proto3 to JSON Mapping by design:
		int64, uint64 ---> String
		float, double ----> number

		This means the timestamp will be marshalled as string in http response :)
		@see https://protobuf.dev/programming-guides/proto3/#json
		*/
		Timestamp: time.Time(t.Timestamp).UnixMilli(),
		ReplyingTo: t.ReplyingTo,
	}
}

func TweetsToProto (ts []*model.Tweet) []*pb.Tweet{
	var tweets []*pb.Tweet
	for _, t := range ts {
		tweets = append(tweets, TweetToProto(t))
	}
	return tweets
}
