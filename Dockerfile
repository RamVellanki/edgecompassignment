FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git

WORKDIR /app
COPY . .
RUN go mod download

RUN go build -o gometrics

FROM alpine:latest
WORKDIR /

COPY --from=builder /app/gometrics gometrics
#COPY --from=builder /bin/sh sh

#RUN chmod +x gometrics

# Expose port 2112 to the outside world
#EXPOSE 2112

# Command to run the executable
ENTRYPOINT ["./gometrics"]
#CMD [ "sh", "-c", "ls -l" ]
