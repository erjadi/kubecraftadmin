FROM golang:latest AS build-go
RUN mkdir /app 
ADD ./src/app/ /app/ 
WORKDIR /app
RUN go get -d
RUN go build -o main . 

FROM alpine
RUN mkdir /app
WORKDIR /app
COPY --from=build-go /app/main /app/
CMD ["/app/main"]
EXPOSE 8000/tcp