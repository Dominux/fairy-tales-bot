FROM golang:latest AS build
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./ ./cmd/main.go

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/main /
CMD ["/main"]
