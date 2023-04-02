package api

import (
	"context"
	"net/http"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	http_handlers "github.com/okpalaChidiebere/chirper-app-api-tweet/api/http_handlers"
	tweetsservice "github.com/okpalaChidiebere/chirper-app-api-tweet/v0/tweets/business_logic"
	pb "github.com/okpalaChidiebere/chirper-app-gen-protos/tweet/v1"
	health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

type Servers struct {
	TweetServer  pb.TweetServiceServer
	health_v1.HealthServer 
}

type APIServer struct {
	httpMux *http.ServeMux
}

func (a Servers) NewAPIServer(httpMux *http.ServeMux) *APIServer{
	server := &APIServer{ httpMux: httpMux, }
	return server
}

func (server *APIServer) RegisterAllEndpoint(tweetsService tweetsservice.Service) error {
	server.httpMux.HandleFunc("/migrate-tweet", http_handlers.MigrateTweetsHandler(tweetsService))
	return nil
}

// Add endpoints to grpc
func (a Servers) RegisterAllService (s *grpc.Server){
	pb.RegisterTweetServiceServer(s, a.TweetServer)
	health_v1.RegisterHealthServer(s, a.HealthServer)
}

//Add endpoints to runtime.ServeMux for http
func (a Servers) RegisterAllServiceHandler (ctx context.Context, mux *runtime.ServeMux) error {
	err := pb.RegisterTweetServiceHandlerServer(ctx, mux, a.TweetServer)
	if err != nil {
		return err
	}
	return nil
}
