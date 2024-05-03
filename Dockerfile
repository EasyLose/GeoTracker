FROM node:14 as react-build
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

FROM golang:1.16 as go-build
WORKDIR /go/src/app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=go-build /go/src/app/server .
COPY --from=react-build /app/build ./frontend/build
EXPOSE 8080
CMD ["./server"]
