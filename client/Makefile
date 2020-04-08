SRC=ui.vala goimport.vala adduser.vala search.vala strangers.vala popup.vala multisend.vala client.vapi
VALACFLAGS=-X -D'GETTEXT_PACKAGE="powerchat"' -X -Wl,--rpath=. -X -L. -X -lclient --pkg=appindicator3-0.1 --pkg gtk+-3.0 --pkg gee-0.8 --pkg json-glib-1.0

powerchat:$(SRC) libclient.a
	valac -o powerchat $(VALACFLAGS) -X -lasound $(SRC)
powerchat.exe:$(SRC) libclient.dll
	if [ -f libclient.a ] ; then rm libclient.a; fi
	valac -o powerchat.exe -X -mwindows -D WINDOWS $(VALACFLAGS) $(SRC)
libclient.a:./libclient/*.go
	sh ./libclient/linux.sh
libclient.dll:./libclient/*.go
	sh ./libclient/dll.sh
libclient.h:./libclient/*.go
	go build -buildmode=c-archive ./libclient
install:
	cp powerchat ../bin/
clean:
	rm *.c
