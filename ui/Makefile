SRC=ui.vala jsonrpc.vala adduser.vala search.vala strangers.vala popup.vala
ui:$(SRC)
	valac -o ui -X -D'GETTEXT_PACKAGE="powerchat"' --pkg gtk+-3.0 --pkg jsonrpc-glib-1.0 --pkg gee-0.8 --pkg json-glib-1.0 $(SRC)
ui.exe:$(SRC)
	valac -o ui -X -D'GETTEXT_PACKAGE="powerchat"' -X -mwindows  --pkg gtk+-3.0 --pkg jsonrpc-glib-1.0 --pkg gee-0.8 --pkg json-glib-1.0 $(SRC)
