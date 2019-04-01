cd ./libclient
go build -o libclient.dll -buildmode=c-shared .
cp libclient.dll ..
cp libclient.h ..
cd ..
