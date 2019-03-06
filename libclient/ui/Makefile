SRC=ui.vala jsonrpc.vala adduser.vala search.vala strangers.vala popup.vala multisend.vala client.vapi
VALACFLAGS=-X -D'GETTEXT_PACKAGE="powerchat"' -X -Wl,-rpath=. -X -L. -X -lclient --pkg gtk+-3.0 --pkg jsonrpc-glib-1.0 --pkg gee-0.8 --pkg json-glib-1.0
powerchat:$(SRC) libclient.h libclient.so
	valac -o powerchat  $(VALACFLAGS) $(SRC)
powerchat.exe:$(SRC) libclient.h libclient.so
	valac -o powerchat.exe -X -mwindows $(VALACFLAGS) $(SRC)
libclient.h libclient.so:../*.go
	go build -o libclient.so -buildmode=c-shared ..
clean:
	rm *.c
