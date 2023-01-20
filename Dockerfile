FROM golang:1.19.3-buster AS build

WORKDIR /app

COPY . .

RUN go mod download && go mod verify

COPY main.go ./

RUN GOOS=linux go build main.go

FROM gcr.io/distroless/base-debian11

COPY --from=build /main /main

ENTRYPOINT ["/main"]