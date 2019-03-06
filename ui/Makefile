SRC=ui.vala jsonrpc.vala adduser.vala search.vala strangers.vala popup.vala multisend.vala
VALACFLAGS=-X -D'GETTEXT_PACKAGE="powerchat"' --pkg gtk+-3.0 --pkg jsonrpc-glib-1.0 --pkg gee-0.8 --pkg json-glib-1.0
ui:$(SRC)
	valac -o ui  $(VALACFLAGS) $(SRC)
ui.exe:$(SRC)
	valac -o ui.exe -X -mwindows $(VALACFLAGS) $(SRC)
clean:
	rm *.c
