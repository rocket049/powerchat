SRC=ui.vala jsonrpc.vala adduser.vala search.vala strangers.vala popup.vala multisend.vala client.vapi
VALACFLAGS=-X -D'GETTEXT_PACKAGE="powerchat"' -X -L. -X -lclient --pkg gtk+-3.0 --pkg jsonrpc-glib-1.0 --pkg gee-0.8 --pkg json-glib-1.0
powerchat:$(SRC) libclient.h libclient.a
	valac -o powerchat $(VALACFLAGS) -X -lasound $(SRC)
powerchat.exe:$(SRC) libclient.h libclient.dll
	valac -o powerchat.exe -X -mwindows $(VALACFLAGS) $(SRC)
libclient.a:../*.go
	go build -o libclient.a -buildmode=c-archive ..
libclient.dll:../*.go
	go build -o libclient.dll -buildmode=c-shared ..
clean:
	rm *.c
