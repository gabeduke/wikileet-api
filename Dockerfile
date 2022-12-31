# syntax=docker/dockerfile:1
FROM node:lts-alpine as frontend

WORKDIR /app
COPY ./frontend/wikileet-ui/package*.json ./
RUN npm install
COPY ./frontend/wikileet-ui/ .
RUN npm run build

FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
COPY --from=frontend /app/dist ./frontend/wikileet-ui/dist

RUN go build -o /wikileet

EXPOSE 8080

CMD [ "/wikileet" ]
