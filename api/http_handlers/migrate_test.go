package api_http_handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	tweetsservice "github.com/okpalaChidiebere/chirper-app-api-tweet/v0/tweets/business_logic"
	model "github.com/okpalaChidiebere/chirper-app-api-tweet/v0/tweets/model"
	"github.com/stretchr/testify/require"
)


func Test_MigrateTweetsHandler(t *testing.T){
	testCases := []struct {
		name          string
		body          []byte
		buildStubs    func(tweetsService *tweetsservice.MockService)
		expectedResponseCode int
		expectedResponse map[string]interface{}
	}{
		{
			name:      "OK",
			body: []byte(`[
				{
					"id": "8xf0y6ziyjabvozdd253nd",
					"text": "Shoutout to all the speakers I know for whom English is not a first language, but can STILL explain a concept well. It's hard enough to give a good talk in your mother tongue!",
					"author": "sarah_edo",
					"timestamp": 1518122597860,
					"likes": ["tylermcginnis"],
					"replies": ["fap8sdxppna8oabnxljzcv", "3km0v4hf1ps92ajf4z2ytg"],
					"replyingTo": null
				}
			]`),
			buildStubs: func(tweetsservice *tweetsservice.MockService) {
				arg := []*model.Tweet{
					{ 
						Id: "8xf0y6ziyjabvozdd253nd",
						Text: "Shoutout to all the speakers I know for whom English is not a first language, but can STILL explain a concept well. It's hard enough to give a good talk in your mother tongue!",
						Author: "sarah_edo",
						Timestamp: model.ChirperAppUnixTime(time.UnixMilli(1518122597860)),
						Likes: []string{"tylermcginnis"},
						Replies: []string{"fap8sdxppna8oabnxljzcv", "3km0v4hf1ps92ajf4z2ytg"},
					},
				}
				tweetsservice.EXPECT().
				BulkSaveTweet(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nil)
			},
			
			expectedResponseCode: http.StatusOK,
			expectedResponse: map[string]interface {}{},
		},
		{
			name:      "empty items",
			body: []byte(`[]`),
			buildStubs: func(tweetsservice *tweetsservice.MockService) {
				
				arg := []*model.Tweet{}
				tweetsservice.EXPECT().
				BulkSaveTweet(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(errors.New("cannot perform action on an empty list"))
			},
			expectedResponseCode:  http.StatusMultiStatus,
			expectedResponse: map[string]interface {}{"message":"cannot perform action on an empty list"},
		},
		{
			name:      "EOF",
			buildStubs: func(tweetsservice *tweetsservice.MockService) {
				tweetsservice.EXPECT().
				BulkSaveTweet(gomock.Any(), gomock.Any()).
					Times(0)
			},
			expectedResponseCode:  http.StatusBadRequest,
			expectedResponse: map[string]interface {}{"message":"EOF"},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tweetsServiceMock := tweetsservice.NewMockService(ctrl)

			tc.buildStubs(tweetsServiceMock)

			server := httptest.NewServer(MigrateTweetsHandler(tweetsServiceMock)) //spin up a test sever that runs our handler
			defer server.Close()

			r, _ := http.NewRequest("POST", server.URL, bytes.NewBuffer(tc.body))
			r.Header.Add("Content-Type", "application/json")

			client := &http.Client{}
			res, _ := client.Do(r)

			checkResponseCode(t, tc.expectedResponseCode, res.StatusCode)

			var resBody map[string]interface{}
			body, _ := io.ReadAll(res.Body)
			_ = json.Unmarshal(body, &resBody);
			require.Equal(t, tc.expectedResponse, resBody)
		})
	}
}

func checkResponseCode(t *testing.T, expected, actual int) {
	require.Equal(t, expected, actual)
}