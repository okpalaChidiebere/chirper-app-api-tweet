package tweetsdataaccess

import (
	"context"
	"encoding/json"
	"log"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/okpalaChidiebere/chirper-app-api-tweet/v0/common"
	model "github.com/okpalaChidiebere/chirper-app-api-tweet/v0/tweets/model"
)

//We can call this an Adapter! It connects to external service
type DynamoDbRepository struct {
	client common.DynamoDBAPI
	tableName string
}

type NextKey struct {
	Id string `json:"id"`
	Author string `json:"author"`
}

func NewDynamoDbRepo(client common.DynamoDBAPI, tableName string) *DynamoDbRepository{
	return &DynamoDbRepository{
		client: client,
		tableName: tableName,
	}
}

func (r *DynamoDbRepository) SaveTweetToDynamoDb(ctx context.Context, replyingToAuthor string, tweet *model.Tweet) (*model.Tweet, error) {
	item, _ := attributevalue.MarshalMap(tweet)
	if len(tweet.Likes) == 0{
			delete(item, "likes")
	}
	if len(tweet.Replies) == 0{
		delete(item, "replies")
	}

	ti := []types.TransactWriteItem{
            {
                Put: &types.Put{
                    Item: item,
                    TableName: aws.String(r.tableName),
                },
            },
            {
                Update: &types.Update{
                    TableName:  aws.String("chirper-app-users-dev"),
					Key: map[string]types.AttributeValue{
						"id": &types.AttributeValueMemberS{Value: tweet.Author},
					},
					UpdateExpression: aws.String("ADD tweets :tweets"),
					ExpressionAttributeValues: map[string]types.AttributeValue{
						//the fact that the tweets array is a Set type, we are sure that the tweetIds will be unique
						":tweets": &types.AttributeValueMemberSS{ Value: []string{ tweet.Id } },
					},
					ConditionExpression: aws.String("attribute_exists(id)"), //we want to make sure the user exists in the database
                },
            },
        }


	if tweet.ReplyingTo != "" {
		
		ti = append(ti, types.TransactWriteItem{       
			Update: &types.Update{
				TableName:  aws.String(r.tableName),
				Key: map[string]types.AttributeValue{
					"id": &types.AttributeValueMemberS{Value: tweet.ReplyingTo},
					"author": &types.AttributeValueMemberS{Value: replyingToAuthor},
				},
				UpdateExpression: aws.String("ADD replies :replies"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":replies": &types.AttributeValueMemberSS{ Value: []string{tweet.Id} },
				},
				ConditionExpression: aws.String("attribute_exists(id)"),
			},
        })
	}

	input := &dynamodb.TransactWriteItemsInput{
        TransactItems: ti,
    }

	if _, err := r.client.TransactWriteItems(ctx,input); err != nil {
		return nil,  err
	}
	return tweet, nil
}

func (r *DynamoDbRepository) BulkSaveTweetToDynamoDb(ctx context.Context, tweets []*model.Tweet) error {
 	batch := make(map[string][]types.WriteRequest)
 	var requests []types.WriteRequest

	for _, tweet := range tweets {
		item, _ := attributevalue.MarshalMap(tweet)

		if len(tweet.Likes) == 0{
			delete(item, "likes")
		}
		if len(tweet.Replies) == 0{
			delete(item, "replies")
		}

		requests = append(requests, types.WriteRequest{PutRequest: &types.PutRequest{Item: item}})
	}

	batch[r.tableName] = requests

	_, err := r.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
		RequestItems: batch,
	})
	if err != nil {
		return err
	}

	//for the sake of this demo, we don't want to deal with unprocessed items if any. We could let the user know the exact ids of unprocessed items if we wanted :)
	// if len(op.UnprocessedItems) != 0 {
	// 	log.Println("there were", len(op.UnprocessedItems), "unprocessed records")
	// }

	return nil
	//for more of a migration type algorithm :)
	// @see https://towardsdatascience.com/dynamodb-go-sdk-how-to-use-the-scan-and-batch-operations-efficiently-5b41988b4988
}

func (r *DynamoDbRepository) ListTweetsFromDynamoDb(ctx context.Context, authedUserID, nextKey string, limit int32) (results []*model.Tweet, nextCursor string, err error) {
	items := []*model.Tweet{}

	p := &dynamodb.QueryInput{
		TableName: aws.String(r.tableName),
		KeyConditionExpression: aws.String("author = :author"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":author":        &types.AttributeValueMemberS{Value: authedUserID},
		},
		ScanIndexForward: aws.Bool(false), //it reverses the order of the list. The latest images will be first
	}

	if nextKey != "" {
		nk := &NextKey{}

		//We decode the key
		k, _ := url.QueryUnescape(nextKey)

		//parse the key
		json.Unmarshal([]byte(k), nk)

		st, _ := attributevalue.MarshalMap(nk)

		p.ExclusiveStartKey = st
	}

	out, err := r.client.Query(ctx,p)

	if err != nil {
		return items, "", err
	}

	err = attributevalue.UnmarshalListOfMaps(out.Items, &items)
	if err != nil {
		return items, "", err
	}

	var mNextKey NextKey
	if err := attributevalue.UnmarshalMap(out.LastEvaluatedKey, &mNextKey); err != nil {
		return items, "", err
	}

	
	log.Printf("%+v\n", nextKey)

	var finalKeyValue string
	if mNextKey.Id == "" {
		//when the next key is null it means there is no more items ot return
		finalKeyValue = string("null")
	} else {
		out, _ := json.Marshal(mNextKey.Id)
		finalKeyValue = string(out)
	}

	return items, finalKeyValue, nil
}

func (r *DynamoDbRepository) GetTweetFromDynamoDb(ctx context.Context, tweetID string) (*model.Tweet, error){
	item := model.Tweet{}

	p := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: tweetID},
        },
	}

	out, err := r.client.GetItem(ctx,p)

	if err != nil {
		return &item, err
	}

	err = attributevalue.UnmarshalMap(out.Item, &item)
	if err != nil {
		return &item, err
	}

	return &item, nil
}

func (r *DynamoDbRepository) SaveLikeToggleInDynamoDb(ctx context.Context, tweetID, author, authedUserID string, hasLiked bool) error {
	updateString := "ADD likes :likes"

	if hasLiked {
		updateString = "DELETE likes :likes"
		//		updateString = "SET likes = list_append(likes, :likes)"
	}

	input := &dynamodb.UpdateItemInput{
        TableName:  aws.String(r.tableName),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: tweetID},
			"author": &types.AttributeValueMemberS{Value: author},
        },
        UpdateExpression: aws.String(updateString),
        ExpressionAttributeValues: map[string]types.AttributeValue{
			":likes": &types.AttributeValueMemberSS{ Value: []string{authedUserID} },
        },
		//https://www.alexdebrie.com/posts/dynamodb-condition-expressions/
		//https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.ConditionExpressions.html
		ConditionExpression: aws.String("attribute_exists(id)"), //we want to update ONLY if it exists
		ReturnValues: types.ReturnValueAllNew,
    }

	if _, err := r.client.UpdateItem(ctx,input); err != nil {
		return  err
	}
	return nil
}

func (r *DynamoDbRepository) ScanTweetsFromDynamoDb(ctx context.Context, limit int32, nextKey string) ([]*model.Tweet, string, error) {
	/*
	we expect the next key to be an object like { "id": "", "<range_key>": "" } . 
	
	Range key will be required for pagination only if the table has a range key was set during the configuration of the table. 
	For the Tweets table, the range key is `authors` that is of type string
	*/
	var mNextKey map[string]string
	pItems := []*model.Tweet{}

	input := &dynamodb.ScanInput{
		TableName:  aws.String(r.tableName),
		Limit:      aws.Int32(limit),
	}

	if nextKey != "" {
		//We decode the key
		k, _ := url.QueryUnescape(nextKey)

		//parse the key
		json.Unmarshal([]byte(k), &mNextKey)

		st, _ := attributevalue.MarshalMap(mNextKey)

		input.ExclusiveStartKey = st
	}

	out, err := r.client.Scan(ctx, input)
	if err != nil {
		return pItems, "", err
	}

	err = attributevalue.UnmarshalListOfMaps(out.Items, &pItems)
	if err != nil {
		return pItems, "", err
	}

	var finalKeyValue string
	//if the an empty key is returned we know we have reached the end of the page
	if out.LastEvaluatedKey == nil {
		finalKeyValue = ""
	} else {
		if err := attributevalue.UnmarshalMap(out.LastEvaluatedKey, &mNextKey); err != nil {
			return pItems, "", err
		}
		out, _ := json.Marshal(mNextKey)
		finalKeyValue = url.QueryEscape(string(out))
	}
	// log.Printf("%+v\n", out.LastEvaluatedKey)
	return pItems, finalKeyValue, nil

	//for more on how to improve your scanning speed see the link below. Ideally you may want to use dynamoDB Query operation which is faster
	// @see https://towardsdatascience.com/dynamodb-go-sdk-how-to-use-the-scan-and-batch-operations-efficiently-5b41988b4988
}