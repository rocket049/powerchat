cd ./libclient
go build -buildmode=c-archive .
cp libclient.a ..
cp libclient.h ..
cd ..
