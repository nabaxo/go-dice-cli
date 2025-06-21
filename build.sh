go build -o dist/go-dice -ldflags="-s -w" .
upx --brute dist/go-dice
