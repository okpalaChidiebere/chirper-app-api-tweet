## Chirper-app-api-tweet

This is a microservice responsible for everything tweets for the chirper-app. Its runs a gRPC server and http/1.1 server at two different ports

## Starting server

- Clone this repo and run `go run main.go`.
- You can access the HTTP/1.1 backend endpoint at [https://localhost:6060](https://localhost:6060) or [http://localhost:6060](http://localhost:6060)
- You can access the gRPC backend with postman at `localhost:6061`
- For grPC Reflection, you will need to load the refection in postman from the insecure port (6061) in the 'new > gRPC Request' tab. After you have load the reflection, it does not matter which port us use to test all the services exposed by the reflection. The only gotcha is if you are to you want to use the secure port, you will need to upload your server cert and key and well as your Authority cert to postman from the preference screen of the app. Learn more about reflection [here](https://www.youtube.com/watch?v=yluYiCj71ss). See this [blog](https://learning.postman.com/docs/sending-requests/certificates/) on how to add SSL to postman; For me i uploaded authority cert generated from [Openssl](https://man.openbsd.org/openssl.1#x509) for the 'CA Certificates' section, server cert and server key for the 'Client Certificates' section.

## Server services

- The three main services for the demo of this project for tweets is defined [here](https://github.com/okpalaChidiebere/chirper-app-apis/blob/master/tweet/v1/api.proto)
- Read this [documentation](https://cloud.google.com/endpoints/docs/grpc/transcoding) to see furthermore on how to interpret the api definitions
- if you want to understand the idea of how the services logic work, you can take a look at the `tweets/business_logic/service.go`

## Useful links about gRPC-Gateway

- [gRPC-Gateway](https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/grpc_api_configuration/). This link have very useful information you need to know about Gateway as well as this google [link](https://cloud.google.com/endpoints/docs/grpc/transcoding)
- This [example](https://github.com/philips/grpc-gateway-example/blob/master/cmd/serve.go) combines an existing grpcServer with grpc-Gateway server. Good to know
- [https://www.clarifai.com/blog/muxing-together-grpc-and-http-traffic-with-grpc-gateway](https://www.clarifai.com/blog/muxing-together-grpc-and-http-traffic-with-grpc-gateway)

## Useful information about the CI build

- For the golangci-lint see more info [here](https://golangci-lint.run/usage/configuration/#command-line-options). For some reason when the lint exits with error the travis pipeline continue to run so i had to manually terminate the build myself if it fails. I learned this from this [example](https://medium.com/@manjula.cse/how-to-stop-the-execution-of-travis-pipeline-if-script-exits-with-an-error-f0e5a43206bf)
- I had to add docker buildx to build our image. There is a difference between regular docker build and docker buildx. See documentation [here](https://docs.docker.com/build/#:~:text=docker%20buildx%20build%20command%20provides,caching%2C%20and%20specifying%20target%20platform.) . With buildx you can build your image for specific platform; see [here](https://docs.docker.com/build/building/multi-platform/)
- We used Travis CI for our build which basically spins up a computer for use remotely and build our app. That computer has git in it. So just provided our github credentials to it which our the computer to build our app with the private modules. It was a good learning. Now if you want do that github step in docker you can checkout this [link](https://jwenz723.medium.com/fetching-private-go-modules-during-docker-build-5b76aa690280). Remember there are ways to provide credentials to github for it to use. See them [here](https://docs.travis-ci.com/user/private-dependencies/). I prefer to use API token
- Learn how to make a go private module with docker [here](https://medium.com/the-godev-corner/how-to-create-a-go-private-module-with-docker-b705e4d195c4)
