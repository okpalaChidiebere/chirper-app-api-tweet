package tweetmodel

import (
	"strconv"
	"strings"
	"time"
)

type Tweet struct {
	Author string  `json:"author" dynamodbav:"author"`
	Id string  `json:"id" dynamodbav:"id"`
	Likes []string  `json:"likes,omitempty" dynamodbav:"likes,omitempty,omitemptyelem,stringset"`
  	Replies []string  `json:"replies,omitempty" dynamodbav:"replies,omitempty,omitemptyelem,stringset"`
  	Text string  `json:"text" dynamodbav:"text_blob"`
	Timestamp ChirperAppUnixTime `json:"timestamp,omitempty" dynamodbav:"created_at,unixtime"`
  	ReplyingTo string  `json:"replyingTo" dynamodbav:"replyingTo"` //if empty then we know its a new tweet
}


type ChirperAppUnixTime time.Time

//custom serialization
func (c ChirperAppUnixTime) MarshalJSON() ([]byte, error) {
	if c.IsZero(){
		return []byte("null"), nil
	}

	return []byte(strconv.Itoa(int(time.Time(c).UnixMilli()))), nil
}

func (c *ChirperAppUnixTime) UnmarshalJSON(b []byte) (err error) {
	// r := strings.Replace(string(s), `"`, ``, -1)
	r := strings.Trim(string(b), `"`) //get rid of double quotes
    if r == "" || r == "null" {
        return nil
    }

	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}

	t := time.UnixMilli(q)

	*c = ChirperAppUnixTime(t) //set result using the pointer
	// fmt.Printf("v ==== %v \n", c)
    return nil
}

func (c ChirperAppUnixTime) IsZero() bool {
	return time.Time(c).IsZero()
}

func (t ChirperAppUnixTime) String() string { return time.Time(t).String() }
