# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image https://hub.docker.com/_/golang
FROM golang:1.20.2-alpine3.17 AS builder

# we define this values during build phase like: docker-compose build --build-arg GOPRIVATE=github.com/okpalaChidiebere or in the docker-compose.yml file otherwise default value is "github.com/okpalaChidiebere/*"
ARG GOPRIVATE="github.com/okpalaChidiebere/*"
ENV GOPRIVATE $GOPRIVATE

# Install git because Golang tooling uses it
RUN apk --update add git

# Move to working directory /app
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Our build container need access to secure file `.netrc` that contains our git credentials. So using the --mount flag we can provide this file without baking into the docker image; o we dont leak our secret
# then build the application's binary. Mark the build as statically linked.
# for the `--mount` flag to work be sure to enable Docker BuildKit by setting termianl env variable `DOCKER_BUILDKIT=1` see https://docs.docker.com/build/buildkit/#getting-started or having Docker buildx plugin installed
# see: https://docs.docker.com/engine/reference/builder/#run---mounttypesecret
RUN --mount=type=secret,id=gitcredentials,target=/root/.netrc \
  CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

# Utilize multi-state to generate a smaller image size for our go app https://docs.docker.com/build/building/multi-stage/
# Build a smaller image that will only contain the application's binary
FROM alpine:latest

# Add Maintainer Info
LABEL maintainer="Chidiebere Okpala <okpalacollins4@gmail.com>"

# Move to working directory /app
WORKDIR /app

# Copy application's binary
COPY --from=builder /app .

EXPOSE 6000 6001

# Run stage
# Command to run the application when starting the container
CMD ["./main"]
