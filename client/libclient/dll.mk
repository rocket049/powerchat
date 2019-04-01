libclient.dll:*.go
	go build -o libclient.dll -buildmode=c-shared .
