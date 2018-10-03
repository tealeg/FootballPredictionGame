# Build stage
FROM golang:alpine AS build-env
ADD . /go/src/github.com/tealeg/FootballPredictionGame
RUN go version
RUN cd /go/src/github.com/tealeg/FootballPredictionGame && go build -o fpg

# Final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/github.com/tealeg/FootballPredictionGame/fpg /app/
ENTRYPOINT ./fpg
