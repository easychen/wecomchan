set -ex

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main && upx -9 main

zip main.zip main