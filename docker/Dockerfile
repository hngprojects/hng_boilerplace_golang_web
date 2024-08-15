FROM golang:1.20.1-alpine3.17 as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the source code into the container and build the app
COPY . .
RUN go build -v -o /dist/golang_app

# Deployment stage
FROM alpine:3.17
WORKDIR /usr/src/app
COPY --from=build /usr/src/app ./
COPY --from=build /dist/golang_app /usr/local/bin/golang_app

# Start the application
CMD golang_app
