FROM golang:1.20 AS build
workdir /go/src

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN go build -v -buildvcs=false -a -installsuffix cgo -o scraper .

FROM busybox AS runtime
WORKDIR /go/app

COPY --from=build /go/src/scraper .
ENTRYPOINT ["./scraper"]
