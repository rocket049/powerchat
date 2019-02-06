export CGO_ENABLED=1
export GOOS=windows
export GOARCH=386
go build -ldflags "-s -H windowsgui"
