package tweetsdataaccess

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"

	"github.com/okpalaChidiebere/chirper-app-api-tweet/v0/common"
	model "github.com/okpalaChidiebere/chirper-app-api-tweet/v0/tweets/model"
	"github.com/stretchr/testify/assert"
)

const fakeTable = "fake-table-name"

type DynamodbMockClient struct {
	common.DynamoDBAPI
}

func initializeFakeDynamoDBRepository() (Repository, error) {
	return NewDynamoDbRepo(&DynamodbMockClient{}, fakeTable), nil
}

func (m *DynamodbMockClient) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	result := dynamodb.GetItemOutput{}

    if params.TableName == nil || *params.TableName == "" {
        return &result, errors.New("Missing required field CreateTableInput.TableName")
    }

	item, err := attributevalue.MarshalMap(model.Tweet{
		Id: "r0xu2v1qrxa6ygtvf2rkjw",
		Author: "dan_abramov",
		Text: "This is a great idea.",
		Timestamp: model.ChirperAppUnixTime(time.UnixMilli(1510044395650)),
		Likes: []string{"tylermcginnis"},
		ReplyingTo: "6h5ims9iks66d4m7kqizmv",
		Replies: make([]string, 0),
	})

	if err != nil {
        return &result, err
    }

	result.Item = item

	return &result, nil
}

func (m *DynamodbMockClient) UpdateItem(ctx context.Context, input *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
    result := dynamodb.UpdateItemOutput{}

    if input.TableName == nil || *input.TableName == "" {
        return &result, errors.New("Missing required field UpdateItemInput.TableName")
    }

    return &result, nil
}

func (m *DynamodbMockClient) Query(ctx context.Context, input *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
    if input.TableName == nil || *input.TableName == "" {
        return &dynamodb.QueryOutput{}, errors.New("Missing required field UpdateItemInput.TableName")
    }

	return &dynamodb.QueryOutput{
		Items:[]map[string]types.AttributeValue{
			{
				"id":        &types.AttributeValueMemberS{Value: "r0xu2v1qrxa6ygtvf2rkjw"},
				"author":        &types.AttributeValueMemberS{Value: "dan_abramov"},
				"text":        &types.AttributeValueMemberS{Value: "This is a great idea."},
				"timestamp":  &types.AttributeValueMemberN{Value: "1510044395650"},
				"likes":        &types.AttributeValueMemberSS{Value: []string{"tylermcginnis"}},
				"replyingTo":        &types.AttributeValueMemberS{Value: "6h5ims9iks66d4m7kqizmv"},
				"replies":        &types.AttributeValueMemberSS{Value: make([]string, 0)},
			},
		},
	}, nil
}

func (m *DynamodbMockClient) BatchWriteItem(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error){
	 result := dynamodb.BatchWriteItemOutput{}

	// keys := make([]string, 0, len(params.RequestItems))
	// for k := range params.RequestItems {
	// 	keys = append(keys, k)
	// }

	//Check Map Contains a key "fake-table-name"
	requests, ok := params.RequestItems[fakeTable]
	if !ok {
		return &result, errors.New("table or tables specified in the BatchWriteItem request does not exist")
	}

	if len(requests) > 25 {
		return &result, errors.New("There are more than 25 requests in the batch")
	}

	for _, wr := range requests {
		if wr.DeleteRequest != nil && wr.PutRequest != nil {
			return &result, errors.New("you cannot put and delete the same item in the same BatchWriteItem request")
		}
	}

	return &result, nil
}

func Test_ListTweetsFromDynamoDb_ReturnsWithNoError(t *testing.T) {
	ctx := context.Background()
	repo, err := initializeFakeDynamoDBRepository()
	if err != nil {
		t.Fatalf("error initializing repository: %s", err.Error())
	}

	tweets, _, err := repo.ListTweetsFromDynamoDb(ctx, "dan_abramov", "", 0)
	if err != nil {
        t.Fatal(err)
    }

	assert.Equal(t, 1, len(tweets))
}

func Test_GetTweetFromDynamoDb_ReturnsWithNoError(t *testing.T) {
	ctx := context.Background()
	repo, err := initializeFakeDynamoDBRepository()
	if err != nil {
		t.Fatalf("error initializing repository: %s", err.Error())
	}

	_, err = repo.GetTweetFromDynamoDb(ctx, "someIDWeDidNotUse")
	if err != nil {
        t.Fatal(err)
    }

	// t.Log("Retrieved test user '" + tweet.Id + "' from table ")
}

func Test_BatchSaveTweetToDynamoDb(t *testing.T) {

	testCases := []struct {
		name string

		tweets []*model.Tweet

		expectedError error
	}{
		{
			name: "Should return no error",
			tweets:  randomTweets(10),
			expectedError: nil,
		},
		{
			name: "Should return an error when requests is more than 25",
			tweets:  randomTweets(26),
			expectedError: errors.New("There are more than 25 requests in the batch"),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			repo, err := initializeFakeDynamoDBRepository()
			if err != nil {
				t.Fatalf("error initializing repository: %s", err.Error())
			}

			err = repo.BulkSaveTweetToDynamoDb(ctx, tc.tweets)
			assert.Equal(t, err, tc.expectedError)
		})
	}
}

// func Test_UpsertTweetFromDynamoDb_ReturnsWithNoError(t *testing.T) {
// 	ctx := context.Background()
// 	repo, err := initializeFakeDynamoDBRepository()
// 	if err != nil {
// 		t.Fatalf("error initializing repository: %s", err.Error())
// 	}

// 	err = repo. UpsertTweetFromDynamoDb(ctx, &model.Tweet{})
// 	if err != nil {
//         t.Fatal(err)
//     }
// }

func randomTweets(tweetsCount int) []*model.Tweet {
	tweets := make([]*model.Tweet, 0)
	n := 1
	for n <= tweetsCount {
		tweet := &model.Tweet{Id: uuid.NewString()}
		tweets = append(tweets, tweet)
		n += 1
	}
	return tweets
}
