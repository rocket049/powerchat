SRC=ui.vala jsonrpc.vala browser.vala adduser.vala search.vala strangers.vala popup.vala
ui:$(SRC)
	valac -o ui --pkg gtk+-3.0 --pkg jsonrpc-glib-1.0 --pkg gee-0.8 --pkg json-glib-1.0 --pkg webkit2gtk-4.0 $(SRC)
