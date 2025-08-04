# # Builder
# FROM golang:1.21.5-alpine3.17 as builder

# RUN apk update && apk upgrade && \
#     apk --update add git make bash build-base


# WORKDIR /home
# COPY go.mod .
# COPY go.sum .
# RUN go mod download

# COPY . .


# RUN go build \
#     -trimpath  \
#     -o app/engine ./app
# # COPY /app/locale /locale

# EXPOSE 9090
# CMD /home/app/engine

# Builder
FROM golang:1.22.6-alpine3.20 as builder

RUN apk update && apk upgrade && \
  apk --update add git make bash build-base

WORKDIR /app

COPY . .

RUN go build \
  -trimpath  \
  -o engine \
  ./cmd/all

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
  apk --update --no-cache add tzdata && \
  mkdir -p /app/cmd/all && mkdir -p /asset && mkdir -p /configs 
#   mkdir /home/app/locale && mkdir /home/app/static
# && chmod +x /home/app/.bin/webp/webpc

WORKDIR /app

EXPOSE 9090

COPY --from=builder /app/configs /configs
COPY --from=builder /app/asset /asset
COPY --from=builder /app/engine /app/cmd/all/
COPY --from=builder /app/cmd/app/.env /app/.env


CMD /app/cmd/all/engine