# syntax=docker/dockerfile:1
FROM node:lts-alpine as frontend

WORKDIR /app
COPY web/package*.json ./
RUN npm install
COPY web/ .
RUN npm run build

FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
COPY --from=frontend /app/dist ./web/dist

RUN go build -o /wikileet

EXPOSE 8080

CMD [ "/wikileet" ]
