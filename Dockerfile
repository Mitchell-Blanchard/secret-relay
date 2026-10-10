FROM golang:1.27-alpine3.24 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o main ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /app/main /main

CMD ["/main"]