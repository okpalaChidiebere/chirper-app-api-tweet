package api_http_handlers

import (
	"encoding/json"
	"net/http"

	tweetsservice "github.com/okpalaChidiebere/chirper-app-api-tweet/v0/tweets/business_logic"
	model "github.com/okpalaChidiebere/chirper-app-api-tweet/v0/tweets/model"
)

func MigrateTweetsHandler(tweetsService tweetsservice.Service) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		items := make([]*model.Tweet, 0)

		err := json.NewDecoder(r.Body).Decode(&items)
		if err != nil {
			JSONError(w, map[string]interface{}{
				"message": err.Error(),
			},  http.StatusBadRequest)
			return
		}

		err = tweetsService.BulkSaveTweet(ctx, items)
		if err != nil {
			JSONError(w, map[string]interface{}{
				"message": err.Error(),
			}, http.StatusMultiStatus)
			return
			//https://aws.github.io/aws-sdk-go-v2/docs/handling-errors/
			//https://www.mscharhag.com/api-design/bulk-and-batch-operations#:~:text=Which%20HTTP%20status%20code%20is,simply%20return%20HTTP%20200%20OK.
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response, _ := json.Marshal(map[string]string{})
		w.Write(response)
	}
}

//The default, http.Error func returns a plain ext, we had to create our own custom error to return a JSON
//https://stackoverflow.com/questions/59763852/can-you-return-json-in-golang-http-error
func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}