package api_http_handlers

import "net/http"


type Interface interface{
	MigrateTweetsHandler() http.HandlerFunc
}