package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	api "github.com/okpalaChidiebere/chirper-app-api-tweet/api"
	"github.com/okpalaChidiebere/chirper-app-api-tweet/config"
	tweetsservice "github.com/okpalaChidiebere/chirper-app-api-tweet/v0/tweets/business_logic"
	tweetsrepo "github.com/okpalaChidiebere/chirper-app-api-tweet/v0/tweets/data_access"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {

	var (
		cfg aws.Config
		err error
		mConfig    = config.NewConfig()
	)

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()


	if mConfig.IsLocal() {
		/*
			Initialize a session that the SDK will use to load
			credentials from the shared credentials file ~/.aws/credentials
			and region from the shared configuration file ~/.aws/config.
		*/
		cfg, err = awsconfig.LoadDefaultConfig(ctx, 
		awsconfig.WithSharedConfigProfile(mConfig.Aws.Aws_profile))
		if err != nil {
			log.Printf("unable to load local SDK config, %v\n", err)
			os.Exit(3)
		}
	} else {
		/*
			Use EC2 Instance Role to assign credentials to application running on an EC2 instance.
			This removes the need to manage credential files in production.
			Make sure to assign the IAM user the limited correct permissions. In this case, access to our S3 bucket and/or RDS
		*/
		cfg, err = awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(mConfig.Aws.Aws_region))
		if err != nil {
			log.Printf("unable to load local SDK config, %v\n", err)
			os.Exit(3)
		}
	}

	dynamodbClient := dynamodb.NewFromConfig(cfg)

	tweetsRepo := tweetsrepo.NewDynamoDbRepo(dynamodbClient,  mConfig.Dev.TweetTable)

	tweetsService := tweetsservice.New(tweetsRepo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "6060"
	}

	s := api.Servers{
		TweetServer: api.NewTweetServer(tweetsService),
		HealthServer: &api.HealthServer{},
	}
	grpcMux := runtime.NewServeMux(runtime.WithHealthzEndpoint(&api.InProcessHealthClient{ Server: s.HealthServer }))
	httpMux := http.NewServeMux()
	httpMux.Handle("/",  allowCORS(grpcMux))

	creds := insecure.NewCredentials()
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}
	apiServer := s.NewAPIServer(httpMux)
	grpcServer := grpc.NewServer(grpc.Creds(creds))

	apiServer.RegisterAllEndpoint(tweetsService)
	if mConfig.IsLocal() {
		//enable reflection to test services in postman. All you need to do is Add a new grpc tab and enter the url of the server with the right port
		//then you can select the messages
		reflection.Register(grpcServer)
	}
	s.RegisterAllService(grpcServer)
	err = s.RegisterAllServiceHandler(ctx, grpcMux)
	if err != nil {
		os.Exit(5)
	}

	nPort, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	httpLis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", nPort))
	if err != nil {
		log.Fatalf("HTTP server: failed to listen: error %v", err)
		os.Exit(2)
	}
	grpcLis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", nPort+1))
	if err != nil {
		log.Fatalf("gRPC server: failed to listen: error %v", err)
		os.Exit(2)
	}

	httpServer := http.Server{
		Handler: httpMux,
		Addr: httpLis.Addr().String(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	http.Handle("/favicon.ico", http.NotFoundHandler())

	go func() {
		log.Printf("grpc server listening at %v", grpcLis.Addr())
		_ = grpcServer.Serve(grpcLis)
	}()

	go func() {
		log.Printf("http/1.1 server listening at %v", httpLis.Addr())
		httpServer.Serve(httpLis)
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	grpcServer.GracefulStop()
	// Perform application shutdown with a maximum timeout of 10 seconds.
	//we will only keep active http requests for the next 10 seconds before we shutdown the server
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Access-Control-Allow-Credentials", "true")
			headers := []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"}
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
			methods := []string{"get", "patch", "post", "head", "options"}
			w.Header().Set("Access-Control-Allow-Methods", strings.ToUpper(strings.Join(methods, ",")))

			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// preflightHandler adds the necessary headers in order to serve
// CORS from any origin using the methods "GET", "HEAD", "POST", "PATCH"
// We insist, don't do this without consideration in production systems.
func preflightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Max-Age", "1728000")
	w.Header().Add("Content-Type", "text/plain; charset=UTF-8")
    w.Header().Add("Content-Length", "0")
	w.WriteHeader(http.StatusNoContent)
}
