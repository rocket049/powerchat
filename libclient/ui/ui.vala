using Gtk;
using Gdk;
using Pango;
using Gee;
using Json;

static AppWin app;
static MyGrid grid1;
static LoginDialog login1;
static AddUserDialog adduser1;
static MultiSendUi msend_ui;
static ChatClient client;
static int RELEASE=22;
static int LATESTVER=0;

public struct UserMsg{
	public int64 id;
	public int sex;
	public string name;
	public string desc;
	public int age;
	public string msg_offline;
	public string timestamp_offline;
}

public class MyGrid: GLib.Object{
	Gtk.ListBox friends;
	Gtk.ListBox msgs = null;
	Gtk.Entry entry1;
	public Gtk.Entry port1;
	int64 to;
	bool running = true;
	//public MyBrowser browser;

	public Gtk.Grid mygrid;
	public Gee.HashMap<string,UserMsg?> frds1;
	Gee.HashMap<string,weak Gtk.Grid?> frd_boxes;
	public Gtk.CssProvider provider1;
	public Gtk.CssProvider mark1;
	public int mark_num=0;
	public Gtk.CssProvider button1;
	public Gtk.CssProvider link_css1;
	public Gtk.Button strangers_btn;
	public Gtk.Button user_btn;
	public Gtk.Button msend_btn;

	public string man_icon = GLib.Path.build_path(GLib.Path.DIR_SEPARATOR_S,prog_path,"..","share","icons","powerchat","man.png");
	public string woman_icon = GLib.Path.build_path(GLib.Path.DIR_SEPARATOR_S,prog_path,"..","share","icons","powerchat","woman.png");
	public int64 uid;
	public int16 usex;
	public string uname;
	public string udesc;
	public int16 uage;
	
	public string host;
	Gtk.CssProvider cssp;
	Gtk.ScrolledWindow msg_win;
	public MyGrid(){
		//this.browser = new MyBrowser();
		this.frds1 = new Gee.HashMap<string,UserMsg?>();
		this.frd_boxes = new Gee.HashMap<string,weak Gtk.Grid?>();
		this.mygrid = new Gtk.Grid();
		this.mygrid.set_column_spacing(5);
		this.cssp = new Gtk.CssProvider();
		var sc = this.mygrid.get_style_context ();
		sc.add_provider(this.cssp,Gtk.STYLE_PROVIDER_PRIORITY_USER);
		this.cssp.load_from_data("""grid{
	padding:5px 5px 5px 5px;
	background-color:#BABABA;
}
list{
	background-color:#FFFFFF;
	color:#000000;
}
""");

		this.provider1 = new Gtk.CssProvider();
		this.provider1.load_from_data("""grid{color:#FF0000;}
""");

		this.mark1 = new Gtk.CssProvider();
		this.mark1.load_from_data("""grid{background:#F59433;}
""");

		this.button1 = new Gtk.CssProvider();
		this.button1.load_from_data("""button{color:#FF0000;}
""");
		this.link_css1 = new Gtk.CssProvider();
		try{
			this.link_css1.load_from_data("label>link{color:#0000FF;}\nlabel>selection{background: #A8141B; color: white;}\n");
		} catch (Error e) {
            print ("CSS Error: %s\n", e.message);
        }

		var scrollWin1 = new Gtk.ScrolledWindow(null,null);
		scrollWin1.width_request = 240;
		scrollWin1.expand = true;
		this.mygrid.attach(scrollWin1,0,0,2,3);
		this.friends = new Gtk.ListBox();
		scrollWin1.add(this.friends);
		var t1 = new Gtk.Label(_("My Friends"));
		this.friends.add(t1);
		var r0 = (t1.parent as Gtk.ListBoxRow);
		r0.set_selectable(false);
		r0.name = "0";
		this.friends.border_width = 3;
		var sc1 = this.friends.get_style_context ();
		sc1.add_provider(this.cssp,Gtk.STYLE_PROVIDER_PRIORITY_USER);

		this.friends.set_sort_func((row1,row2)=>{
			if(row1.name=="0"){
				return -1;
			}
			var rsc1 = this.frd_boxes[row1.name].get_style_context();
			var rsc2 = this.frd_boxes[row2.name].get_style_context();
			if(rsc1.has_class("off")){
				if(rsc2.has_class("off")){
					if(row1.name.to_int64() > row2.name.to_int64()){
						return 1;
					}else{
						return -1;
					}
				}else{
					//print("row1 is off row2 on\\n");
					return 1;
				}
			}else{
				if(rsc2.has_class("off")){
					//print("row1 is on row2 off\n");
					return -1;
				}else if(row1.name.to_int64() > row2.name.to_int64()){
					return 1;
				}else{
					return -1;
				}
			}
		});
		var bottom_grid = new Gtk.Grid();
		mygrid.attach(bottom_grid,0,3,2,1);
		var b1 = new Gtk.Button.with_label(_("Find Persons"));
		bottom_grid.attach(b1,0,0);

		strangers_btn = new Gtk.Button.with_label(_("Strangers"));
		bottom_grid.attach(strangers_btn,1,0);
		
		msend_btn = new Gtk.Button.with_label(_("MultiSend"));
		bottom_grid.attach(msend_btn,2,0);

		user_btn = new Gtk.Button.with_label(_("Current User"));
		this.mygrid.attach(user_btn,2,0,2,1);
		user_btn.hexpand = true;

		var b4 = new Gtk.Label(_("Proxy Port"));
		this.mygrid.attach(b4,4,0,1,1);

		port1 = new Gtk.Entry();
		port1.set_text(proxy_port.to_string());
		port1.tooltip_text = _("Click `Modify` button to edit.\nSet 0 to restore default value.");
		port1.max_length = 5;
		port1.width_request = 50;
		port1.editable=false;
		this.mygrid.attach(port1,5,0,1,1);

		var b6 = new Gtk.Button.with_label(_("Modify"));
		this.mygrid.attach(b6,6,0,1,1);

		this.msg_win = new Gtk.ScrolledWindow(null,null);
		this.msg_win.height_request = 450;
		this.msg_win.expand = true;
		this.mygrid.attach(this.msg_win,2,1,5,1);

		var grid1 = new Gtk.Grid();
		this.mygrid.attach(grid1,2,2,5,2);
		//文件拖放区
		Gtk.EventBox dropbox = new Gtk.EventBox();
		dropbox.set_size_request(240,40);
		grid1.attach(dropbox,0,0,4,1);
		var droplabel = new Gtk.Label(_("Send File/Image Here: Drag a file Or Press Ctrl-V to paste a file"));
		droplabel.wrap = true;
        droplabel.wrap_mode = Pango.WrapMode.CHAR;
		droplabel.selectable = true;
		dropbox.add(droplabel);
		Gtk.drag_dest_set (dropbox, Gtk.DestDefaults.ALL, null, Gdk.DragAction.COPY);
		Gtk.drag_dest_add_uri_targets(dropbox);
		dropbox.drag_data_received.connect((context, x,y,data, info, time)=>{
			var uris = data.get_uris();
			send_uri1(uris);
		});
		dropbox.key_press_event.connect((e)=>{
			if(e.keyval==Gdk.Key.v && e.state==Gdk.ModifierType.CONTROL_MASK){
				var clipboard1 = Gtk.Clipboard.@get(Gdk.Atom.NONE);
				clipboard1.request_uris((b, uris)=>{
					send_uri1(uris);
				});
			}
			return true;
		});

		this.entry1 = new Gtk.Entry();
		grid1.attach(this.entry1,0,1,3,1);
		this.entry1.hexpand = true;

		var b7 = new Gtk.Button.with_label(_("Send"));
		grid1.attach(b7,3,1,1,1);

		this.mygrid.show.connect(()=>{
			//var mutex1 = new GLib.Mutex();
			stdout.printf("grid show\n");
			
			strangers1 = new StrangersDialg();
			uint16 port2 = (uint16)client.get_proxy();
			if(port2>0){
				proxy_port = port2;
				this.port1.text = proxy_port.to_string();
			}else{
				Gtk.main_quit();
			}
			this.host = client.get_host();
			app.title = _("Everyone Publish!")+@"($(this.mark_num))"+" - "+@"$(this.uname)@$(this.host)";
			app.update_tooltip();
		});

		b1.clicked.connect(()=>{
			search1 = new SearchDialg();
			search1.show();
			return;
		});

		b6.clicked.connect (() => {
			// 修改代理端口
			if(port1.editable==false){
				port1.editable = true;
				b6.set_label(_("Save"));
			}else{
				port1.editable=false;
				b6.set_label(_("Modify"));
				client.set_proxy( (int)port1.text.to_int64() );
			}
		});
        strangers_btn.clicked.connect (() => {
			strangers1.show();
		});
		msend_btn.clicked.connect(()=>{
			msend_ui = new MultiSendUi();
			msend_ui.show();
		});
        user_btn.clicked.connect (() => {
			// Emitted when the button has been activated:
			var dlg_user = new Gtk.MessageDialog(app, Gtk.DialogFlags.MODAL, Gtk.MessageType.INFO, Gtk.ButtonsType.OK,null);
			dlg_user.text = this.uname+_(" Details");
			var sex=_("Man");
			if (this.usex==2)
				sex=_("Woman");
			var blog_dir = GLib.Path.build_path(GLib.Path.DIR_SEPARATOR_S,Environment.get_home_dir(),"ChatShare");
            dlg_user.secondary_text = @"ID:$(this.uid)\n"+_("Age:")+@"$(this.uage)\n"+_("Sex:")+@"$(sex)\n"+_("Description:")+@"$(this.udesc)\n"
					+ _("Blog Directory:")+blog_dir;
            dlg_user.show();
            dlg_user.response.connect((rid)=>{
				dlg_user.destroy();
			});
		});
		this.entry1.activate.connect( ()=>{
			this.send_msg();
		} );
			
        b7.clicked.connect (() => {
			// 发送信息
			this.send_msg();
		});

		this.friends.row_selected.connect((r)=>{
			var id = r.name.to_int64();
			if (id==0)
				return;
			this.to = id;
			var u = this.frds1[id.to_string()];
			//stdout.printf(@"selected $(id) $(u.name) $(u.sex)\n");
			if (this.msgs!=null){
				this.msg_win.remove(this.msgs);
			}
			this.msgs = this.boxes[id.to_string()];
			this.msg_win.add(this.msgs);
			Gtk.Grid grid = this.frd_boxes[id.to_string()];
			var sc3 = grid.get_style_context();
			sc3.remove_provider(this.mark1);
			if ( sc3.has_class("mark") ){
				sc3.remove_class("mark");
				this.mark_num--;
				if (this.mark_num==0)
					app.clear_notify();
				app.title = _("Everyone Publish!")+@"($(this.mark_num))"+" - "+this.host;
				app.show_all();
			}else{
				this.msg_win.show_all();
			}
		});
	}
	public void send_uri1(string[] uris){
		if(uris.length!=1){
			this.add_text(_("too many files"));
			return;
		}
		var fname = GLib.Filename.from_uri(uris[0]);
		if(FileUtils.test(fname, FileTest.IS_REGULAR)==false){
			this.add_text(_("this is not a file")+@" : $(fname)");
			return;
		}
		client.send_file(this.to, fname);
		string text1 = @"<a href='$(uris[0])'>$(GLib.Path.get_basename(fname))</a>";
		this.add_left_name_icon(this.uname,this.usex);
		this.add_text(text1,true,true);
	}
	public void send_msg(){
		// 发送信息
		if(this.to==0)
			return;
		client.ChatTo(this.to,this.entry1.text);
		//var u = this.frds1[this.uid.to_string()];
		this.add_left_name_icon(this.uname,this.usex);
		this.add_text(this.entry1.text);
		this.entry1.text = "";
		GLib.Idle.add(()=>{
			var adj1 = this.msgs.get_adjustment();
			adj1.value = adj1.upper;
			return false;
		});
	}
    public void show_sended_msg_to(int64 to,string msg){
        this.add_left_name_icon_to(to,this.uname,this.usex);
        this.add_text_to(to,msg);
    }
	public void update_pwd(){
		var dlg_pwd = new Gtk.Dialog.with_buttons(_("Update Password"),app,Gtk.DialogFlags.MODAL);
		var grid = new Gtk.Grid();
		grid.attach(new Gtk.Label(_("Input your new password")),0,0,2,1);
		
		grid.attach(new Gtk.Label(_("Old Password:")),0,1);
		var pwd1 = new Gtk.Entry();
		pwd1.set_visibility(false);
		grid.attach(pwd1,1,1);
		
		grid.attach(new Gtk.Label(_("New Password:")),0,2);
		var pwd2 = new Gtk.Entry();
		pwd2.set_visibility(false);
		grid.attach(pwd2,1,2);
		
		grid.attach(new Gtk.Label(_("Confirm Password:")),0,3);
		var pwd3 = new Gtk.Entry();
		pwd3.set_visibility(false);
		grid.attach(pwd3,1,3);
		
		var content = dlg_pwd.get_content_area () as Gtk.Box;
		content.pack_start(grid);
		
		dlg_pwd.add_button(_("Update"),2);
		dlg_pwd.add_button(_("Cancel"),3);
		
		dlg_pwd.response.connect((rid)=>{
			if(rid==3){
				dlg_pwd.destroy();
			}else if(rid==2){
				if (pwd2.text != pwd3.text){
					dlg_pwd.title = _("Confirm Fail!");
					return;
				}
				client.update_pwd(this.uname,pwd1.text,pwd2.text);
				dlg_pwd.destroy();
			}
		});
		dlg_pwd.show_all();
	}
	Gee.HashMap<string,Gtk.ListBox?> boxes = new Gee.HashMap<string,Gtk.ListBox?>();
	//Gtk.ListBox hides = new Gtk.ListBox();
	public void add_listbox_id(int64 uid){
		var box = new Gtk.ListBox();
		this.boxes[uid.to_string()] = box;
		box.selection_mode = Gtk.SelectionMode.NONE;
		box.expand = true;
		box.border_width = 3;
		//this.msgs.modify_bg(Gtk.StateType.NORMAL,color1);
		//this.msg_win.add(box);
		var u1 = this.frds1[uid.to_string()];
		var t2 = new Gtk.Label(_("Chat To: ")+u1.name);
		box.add(t2);
		(t2.parent as Gtk.ListBoxRow).set_selectable(false);

		if (this.msgs!=null){
			this.msg_win.remove(this.msgs);
		}
		this.msgs = box;
		this.msg_win.add(this.msgs);
		this.to = uid;
		if (u1.timestamp_offline.length > 10){
			//insert offline message
			add_right_name_icon(u1.name,(int16)u1.sex);
			add_text(_("Rewrite Offline Message:")+@"[$(u1.timestamp_offline)]\n$(u1.msg_offline)");
			this.msg_mark(u1.id.to_string());
		}
		this.msg_win.show_all();

		var sc2 = box.get_style_context ();
		sc2.add_provider(this.cssp,Gtk.STYLE_PROVIDER_PRIORITY_USER);
	}

	public void release_resource(){
		this.running = false;
		//this.conn.close();
	}
	string pressed = "";
	public void add_friend(UserMsg user1,bool tell=true){
		if(this.frds1.has_key(user1.id.to_string()))
			return;
		else
			this.frds1[user1.id.to_string()] = user1;
		if(tell){
			client.tell(user1.id);
		}
		string iconp;
		if (user1.sex==1)
			iconp = this.man_icon;
		else
			iconp = this.woman_icon;
		var pix1 = new Gdk.Pixbuf.from_file(iconp);
        var grid2 = new Gtk.Grid();
        var img2 = new Gtk.Image();
        img2.set_from_pixbuf(pix1);
        grid2.attach(img2,0,0);

        var l2 = new Gtk.Label(user1.name);
		l2.xalign = (float)0;
		l2.hexpand = true;
        grid2.attach(l2,1,0);

        var b2 = new Gtk.Button.with_label("WEB");
        b2.set_tooltip_text(@"http://localhost:$(server_port)");
        grid2.attach(b2,2,0);
        grid2.set_column_spacing(5);
        
        var b3 = new Gtk.Button.with_label("TCP");
        b3.set_tooltip_text(_("TCP Tunnel, PORT:")+@"$(server_port)");
        grid2.attach(b3,3,0);

        this.frd_boxes[@"$(user1.id)"] = grid2;
        this.friends.add(grid2);
        //var row2 = new Gtk.ListBoxRow();
        //row2.add(grid2);
        var row2 = grid2.get_parent() as Gtk.ListBoxRow;
        row2.name = @"$(user1.id)";
        //this.friends.add(row2);
        //grid2.parent.name = @"$(user1.id)";
        grid2.show_all();

        img2.tooltip_text = @"$(user1.age)岁\n$(user1.desc)";

        b2.clicked.connect(()=>{
			//stdout.printf(@"open %$(uint64.FORMAT)\n",user1.id);
			client.set_http_id(user1.id);
		});
		
		b3.clicked.connect(()=>{
			//stdout.printf(@"open %$(uint64.FORMAT)\n",user1.id);
			client.set_tcp_id(user1.id);
		});
		this.friends.button_release_event.connect((e)=>{
			if(e.button!=3)
				return false;

			//stdout.printf("button:%u %f\n",e.button,e.y);
			Gtk.ListBoxRow r = this.friends.get_row_at_y((int)e.y);
			this.friends.select_row(r);
			popup1.set_id( r.name );
			popup1.popup_at_pointer(e);
			return true;
		});

		this.add_listbox_id(user1.id);
	}
	public void remove_friend(string fid){
		var grid = this.frd_boxes[fid];
		this.frd_boxes.unset(fid);
		this.frds1.unset(fid);
		client.remove_friend(fid.to_int64());
		//hide row
		Gtk.ListBoxRow r = grid.get_parent() as Gtk.ListBoxRow;
		r.set_selectable(false);
		r.name="";
		//r.hide();
		this.friends.remove( r );
	}
	public void add_right_name_icon(string name,int16 sex){
		string iconp;
		if (sex==1)
			iconp = this.man_icon;
		else
			iconp = this.woman_icon;
        var pix1 = new Gdk.Pixbuf.from_file(iconp);
        var grid2 = new Gtk.Grid();
        var img2 = new Gtk.Image();
        img2.set_from_pixbuf(pix1);
        grid2.attach(img2,1,0);
        grid2.halign = Gtk.Align.END;
		var l2 = new Gtk.Label(name);
		l2.xalign = (float)1;
        grid2.attach(l2,0,0);
        grid2.set_column_spacing(5);
		this.msgs.add(grid2);
		grid2.show_all();
    }
    public void add_left_name_icon(string name,int16 sex){
		string iconp;
		if (sex==1)
			iconp = this.man_icon;
		else
			iconp = this.woman_icon;
        var grid1 = new Gtk.Grid();
        var pix1 = new Gdk.Pixbuf.from_file(iconp);
        var img1 = new Gtk.Image();
        img1.set_from_pixbuf(pix1);
        grid1.attach(img1,0,0);
		var l1 = new Gtk.Label(name);
		l1.xalign = (float)0;
        grid1.attach(l1,1,0);
        grid1.set_column_spacing(5);
		this.msgs.add(grid1);
		grid1.show_all();
    }
    public void add_left_name_icon_to(int64 to , string name,int16 sex){
        var box1 = boxes[to.to_string()];
        if(box1==null){
            return;
        }
		string iconp;
		if (sex==1)
			iconp = this.man_icon;
		else
			iconp = this.woman_icon;
        var grid = new Gtk.Grid();
        var pix1 = new Gdk.Pixbuf.from_file(iconp);
        var img1 = new Gtk.Image();
        img1.set_from_pixbuf(pix1);
        grid.attach(img1,0,0);
		var l1 = new Gtk.Label(name);
		l1.xalign = (float)0;
        grid.attach(l1,1,0);
        grid.set_column_spacing(5);

        box1.add(grid);

		grid.show_all();
    }
    public void add_text_to(int64 to,string text,bool center=false ,bool markup=false){
        var box1 = boxes[to.to_string()];
        if(box1==null){
            return;
        }
        var lb = new Gtk.Label("");
        lb.set_selectable(true);
        var sc1 = lb.get_style_context();
		sc1.add_provider(this.link_css1,Gtk.STYLE_PROVIDER_PRIORITY_USER);
        if(markup){
			lb.set_markup(text);	
		} else
			lb.set_label(text);
		lb.wrap = true;
        lb.wrap_mode = Pango.WrapMode.CHAR;
        if(!center){
            lb.xalign = (float)0;
        }
        lb.width_request = 360;
        lb.max_width_chars = 15;
        var grid=new Gtk.Grid();
        var lb1 = new Gtk.Label("");
        lb1.width_request = 5;
        grid.attach(lb1,0,0);
        grid.attach(lb,1,0);
        var lb2 = new Gtk.Label("");
        lb2.width_request = 5;
        grid.attach(lb2,2,0);
        grid.halign = Gtk.Align.CENTER;
        
		box1.add(grid);
        
		grid.show_all();
    }
    public void add_image(string pathname){
        var p1 = new Gdk.Pixbuf.from_file(pathname);
        var image = new Gtk.Image();
        if(p1.width>300){
            var xs = (double)300/(double)p1.width;
            var h2 = (int)(p1.height*xs);
            var p2 = new Gdk.Pixbuf(Gdk.Colorspace.RGB,true,8,300,h2);
            p1.scale(p2, 0, 0, 300, h2, 0.0, 0.0, xs, xs,Gdk.InterpType.NEAREST);
            image.set_from_pixbuf(p2);
        }else{
            image.set_from_pixbuf(p1);
        }
		this.msgs.add(image);
		image.show();
    }
    public void add_text(string text,bool center=false ,bool markup=false){
        var lb = new Gtk.Label("");
        lb.set_selectable(true);
        var sc1 = lb.get_style_context();
		sc1.add_provider(this.link_css1,Gtk.STYLE_PROVIDER_PRIORITY_USER);
        if(markup){
			lb.set_markup(text);	
		} else
			lb.set_label(text);
		lb.wrap = true;
        lb.wrap_mode = Pango.WrapMode.CHAR;
        if(!center){
            lb.xalign = (float)0;
        }
        lb.width_request = 360;
        lb.max_width_chars = 15;
        var grid=new Gtk.Grid();
        var lb1 = new Gtk.Label("");
        lb1.width_request = 5;
        grid.attach(lb1,0,0);
        grid.attach(lb,1,0);
        var lb2 = new Gtk.Label("");
        lb2.width_request = 5;
        grid.attach(lb2,2,0);
        grid.halign = Gtk.Align.CENTER;
		this.msgs.add(grid);
		grid.show_all();
    }
	private void add_operate_buttons(string pathname){
		var dir1 = GLib.Path.get_dirname(pathname);
		var grid = new Gtk.Grid();
		var lb1 = new Gtk.Label(" ");
		var lb2 = new Gtk.Label(" ");
		lb1.expand=true;
		lb2.expand=true;
		var bt_open = new Gtk.Button.with_label(_("OpenFile"));
		bt_open.tooltip_text = pathname;
		var bt_dir = new Gtk.Button.with_label(_("OpenDir"));
		bt_dir.tooltip_text = dir1;
		var bt_del = new Gtk.Button.with_label(_("RemoveFile"));
		grid.attach(lb1,0,0);
		grid.attach(bt_open,1,0);
		grid.attach(bt_dir,2,0);
		grid.attach(bt_del,3,0);
		grid.attach(lb2,4,0);
		
		this.msgs.add( grid );
		grid.show_all();
		
		bt_open.clicked.connect(()=>{
			client.open_path(pathname);
		});
		bt_dir.clicked.connect(()=>{
			client.open_path(dir1);
		});
		bt_del.clicked.connect(()=>{
			GLib.FileUtils.remove(pathname);
		});
	}
	//callback in rpc msg
	public void rpc_callback(int8 typ,int64 from,string msg){
		//Msg　开头可以带着类型标记 JSON/TEXT
		//print(@"ID: $(from) ,Msg: $(msg)\n");
		//from==0  "Offline id"
		if(from==0){
			if(msg.length<=8){
				return;
			}
			if(msg[0:8]=="Offline "){
				string off_id = msg[8:msg.length];
				//print("ID:%s : %s\n",off_id,msg);
				var grid = this.frd_boxes[off_id];
				var sc = grid.get_style_context();
				sc.add_provider(this.provider1,Gtk.STYLE_PROVIDER_PRIORITY_USER);
				if (sc.has_class("off")==false){
					sc.add_class("off");
				}
				this.friends.invalidate_sort ();
				//show msg
				var tmp = this.msgs;
				this.msgs = this.boxes[off_id];
				this.add_text(_("[Offline]"),false);
				this.msgs = tmp;
				GLib.Idle.add(()=>{
					var adj1 = this.msgs.get_adjustment();
					adj1.value = adj1.upper;
					return false;
				});
			}else if(msg[0:8]=="Version:"){
				LATESTVER = msg[8:msg.length].to_int();
				if (LATESTVER>RELEASE){
					version_notify();
				}
			}else if(msg[0:8]=="DELETE 1"){
                Gtk.main_quit();
            }else if(msg[0:8]=="DELETE 0"){
                this.add_text(_("[Operate Fail]"));
            }
			//print("Cmd:%i From:%"+int64.FORMAT+" Msg:%s\n",typ,from,msg);
			return;
		}
		//from>0
		string typ1 = msg[0:4];
		string msg1 = msg[4:msg.length];
		string fname="";
		int16 fsex=2;
		var u = this.frds1[from.to_string()];
		var display = this.boxes[from.to_string()];
		var bak_msgs = this.msgs;
		if( u!=null ){
			//print(@"has_key:$(u.name)\n");
			fname = u.name;
			fsex = (int16)u.sex;
			this.msgs = display;
		}else if (typ1 == "TEXT"){
			fname = @"ID:$(from)";
			client.offline_msg_with_id(from,msg1);
			this.msgs = bak_msgs;
			return;
		}else if(typ1=="JSON"){
			fname = @"ID:$(from)";
		}

		switch(typ1){
		case "TEXT":
			this.add_right_name_icon(fname,fsex);
			this.add_text(msg1);
			msg_mark(from.to_string());
			msg_notify(fname);
			break;
		case "JSON":
			var p2 = new Json.Parser();
			if(p2.load_from_data(msg1)==false){
				break;
			}
			var node2 = p2.get_root();
			if (node2==null){
				break;
			}
			var obj2 = node2.get_object();
			string name1 = obj2.get_string_member("Name");
			string mime1 = obj2.get_string_member("Mime");
			this.add_right_name_icon(fname,fsex);
			if(mime1[0:5]=="image"){
				this.add_image(name1);
			}
			this.add_text(GLib.Path.get_basename(name1),true);
			add_operate_buttons(name1);
			msg_mark(from.to_string());
			this.msgs.show_all();
			msg_notify(fname);
			break;
		case "LOGI":
			if( u==null )
				break;
			user_online(from);
			break;
		}
		GLib.Idle.add(()=>{
			var adj1 = this.msgs.get_adjustment();
			adj1.value = adj1.upper;
			return false;
		});
		this.msgs = bak_msgs;
	}
    public void user_online(int64 uid1){
        Gtk.Grid grid = this.frd_boxes[uid1.to_string()];
        var sc3 = grid.get_style_context();
        sc3.remove_provider(this.provider1);
        sc3.remove_class("off");
        this.friends.invalidate_sort ();
    }
	public void msg_mark(string uid){
		Gtk.ListBoxRow r = this.frd_boxes[uid].get_parent() as Gtk.ListBoxRow;
		if(app.counter==1){
			app.tray_notify();
		}
		if(r.is_selected())
			return;
		Gtk.Grid grid = this.frd_boxes[uid];
		var sc3 = grid.get_style_context();
		sc3.add_provider(this.mark1,Gtk.STYLE_PROVIDER_PRIORITY_USER);
		if ( sc3.has_class("mark")==false ){
			sc3.add_class("mark");
			this.mark_num++;
			if(this.mark_num==1)
				app.tray_notify();
			
			app.title = _("Everyone Publish!")+@"($(this.mark_num))"+" - "+this.host;
			grid.show_all();
		}else{
			grid.show_all();
		}
		//print(@"mark: $(uid)\n");
	}
}
Gtk.Application application1;
public void msg_notify(string uname){
	//var pname = GLib.Environment.get_prgname();
    var x = GLib.Environ.get_variable(GLib.Environ.@get(),"DISPLAY");
	if( x != null ){
        var app = application1;
	    app.hold();
		var notify1 = new Notification(_("New message"));
		notify1.set_body(_("From: ")+uname);
		notify1.set_default_action("app.show-win");
		app.send_notification(null,notify1);
        app.release();
	}
}
public void version_notify(){
	var pname = GLib.Environment.get_prgname();
	if( pname=="ui" ){
        var app = application1;
	    app.hold();
		var notify1 = new Notification(_("New Version Released!"));
		notify1.set_body(_("Click here or click menu item 'Help->Upgrade' to get new version."));
		notify1.set_default_action("app.down-page");
		app.send_notification(null,notify1);
        app.release();
	}
}
public class AppWin:Gtk.ApplicationWindow{
	Gtk.StatusIcon tray1;
	Gdk.Pixbuf icon1;
	Gdk.Pixbuf icon2;
	Gtk.VBox box1;
	public int counter=0;
	public AppWin(){
		// Sets the title of the Window:
		var dt1 = new DateTime.now_local();
		
		application1 = new Gtk.Application(@"app.powerchat.id$(dt1.to_unix())",GLib.ApplicationFlags.FLAGS_NONE);
		application1.register();
		application1.add_window(this as Gtk.Window);
		this.title = _("Everyone Publish!");

		// Center window at startup:
		this.window_position = Gtk.WindowPosition.CENTER;

		// Sets the default size of a window:
		this.set_default_size(640,480);
		// Whether the titlebar should be hidden during maximization.
		this.hide_titlebar_when_maximized = true;

        this.set_resizable(false);
        var icon_path = GLib.Path.build_path(GLib.Path.DIR_SEPARATOR_S,prog_path,"..","share","icons","powerchat","tank.png");
        this.set_icon_from_file(icon_path);
        this.icon1 = new Gdk.Pixbuf.from_file(icon_path);
        this.icon2 = new Gdk.Pixbuf.from_file(GLib.Path.build_path(GLib.Path.DIR_SEPARATOR_S,prog_path,"..","share","icons","powerchat","msg.png"));
        this.tray1 = new Gtk.StatusIcon.from_pixbuf(this.icon1);
		this.tray1.set_visible(false);
		this.tray1.activate.connect(()=>{
			if (counter == 0){
				this.hide();
				counter = 1;
			} else{
				this.show();
				counter = 0;
                this.set_keep_above(true);
			}
		});
		this.set_focus_child.connect((w)=>{
			this.set_keep_above(false);
		});
		this.show.connect(()=>{
			this.tray1.set_visible(true);
			if(grid1.mark_num==0){
				this.clear_notify();
			}
		});
		// Method called on pressing [X]
		this.set_destroy_with_parent(false);
		this.delete_event.connect((e)=>{ 
			counter = 1;
			return this.hide_on_delete ();
		});
		this.box1 = new Gtk.VBox(false,0);
		this.add(this.box1);
		this.setup_menubar();
	}
	public void update_tooltip(){
		this.tray1.set_tooltip_text(_("Everyone Publish!")+" - "+@"$(grid1.uname)@$(grid1.host)"+"\n"+_("(Click to Hide/Show)"));
	}
	public void tray_notify(){
		this.tray1.set_from_pixbuf(this.icon2);
	}
	public void clear_notify(){
		this.tray1.set_from_pixbuf(this.icon1);
	}
	public void append(Gtk.Widget w){
		this.box1.pack_start(w);
	}
	public void setup_menubar(){
        var menu1 = new GLib.Menu();
        var menubar =new GLib.Menu();
        
        var item1 = new GLib.MenuItem(_("PreviewWeb"),"app.preview-web");
        menu1.append_item(item1);
        
        item1 = new GLib.MenuItem(_("OpenBlogDir"),"app.blog-dir");
        menu1.append_item(item1);
        
        item1 = new GLib.MenuItem(_("OpenDownloadDir"),"app.down-dir");
        menu1.append_item(item1);
        
        item1 = new GLib.MenuItem(_("ModifyDesc"),"app.modify-desc");
        menu1.append_item(item1);
        
        item1 = new GLib.MenuItem(_("UpdatePasswd"),"app.update-pwd");
        menu1.append_item(item1);
        
        item1 = new GLib.MenuItem(_("Quit"),"app.quit");
        menu1.append_item(item1);
        
        menubar.append_submenu(_("Operate"),menu1);
        
        menu1 = new GLib.Menu();
        
        item1 = new GLib.MenuItem(_("Homepage"),"app.homepage");
        menu1.append_item(item1);
        
        item1 = new GLib.MenuItem(_("Upgrade"),"app.down-page");
        menu1.append_item(item1);
        
        item1 = new GLib.MenuItem(_("About"),"app.about");
        menu1.append_item(item1);
        
        item1 = new GLib.MenuItem(_("Pay"),"app.pay");
        menu1.append_item(item1);
        
        item1 = new GLib.MenuItem(_("DeleteMe"),"app.delete-me");
        menu1.append_item(item1);
        
        menubar.append_submenu(_("Help"),menu1);
        
        application1.set_menubar(menubar as GLib.MenuModel);
        
        add_actions();
    }
    private void add_actions () {
		SimpleAction act1 = new SimpleAction ("show-win", null);
		act1.activate.connect (() => {
			application1.hold ();
			this.show();
            //窗口顶置
            this.set_keep_above(true);
			counter=0;
			application1.release ();
		});
        act1.set_enabled(true);
		application1.add_action (act1);
		
		SimpleAction act2 = new SimpleAction ("about", null);
		act2.activate.connect (() => {
			application1.hold ();
			var dlg_about = new Gtk.MessageDialog(this, Gtk.DialogFlags.MODAL, Gtk.MessageType.INFO, Gtk.ButtonsType.OK,null);
			dlg_about.text = _("Copy Right:");
            dlg_about.secondary_text = "Fu Huizhong <fuhuizn@163.com>";
            dlg_about.show();
            dlg_about.response.connect((rid)=>{
				dlg_about.destroy();
			});
			application1.release ();
		});
        act2.set_enabled(true);
		application1.add_action (act2);
		
		SimpleAction act3 = new SimpleAction ("homepage", null);
		act3.activate.connect (() => {
			application1.hold ();
			//Gtk.show_uri(null,"https://gitee.com/rocket049/powerchat",Gdk.CURRENT_TIME);
			client.open_path("https://gitee.com/rocket049/powerchat");
			application1.release ();
		});
        act3.set_enabled(true);
		application1.add_action (act3);
		
		SimpleAction act4 = new SimpleAction ("pay", null);
		act4.activate.connect (() => {
			application1.hold ();
			//Gtk.show_uri(null,"https://gitee.com/rocket049/powerchat/wikis/powerchat?sort_id=1325779",Gdk.CURRENT_TIME);
			client.open_path("https://gitee.com/rocket049/powerchat/wikis/powerchat?sort_id=1325779");
			application1.release ();
		});
        act4.set_enabled(true);
		application1.add_action (act4);
		
		SimpleAction act5 = new SimpleAction ("preview-web", null);
		act5.activate.connect (() => {
			application1.hold ();
			//Gtk.show_uri(null,@"http://localhost:$(proxy_port)/",Gdk.CURRENT_TIME);
			client.open_path(@"http://localhost:$(proxy_port)/");
			application1.release ();
		});
        act5.set_enabled(true);
		application1.add_action (act5);
		
		SimpleAction act6 = new SimpleAction ("blog-dir", null);
		act6.activate.connect (() => {
			application1.hold ();
			var home1 = GLib.Environment.get_home_dir();
			var blog1 = GLib.Path.build_path(GLib.Path.DIR_SEPARATOR_S,home1,"ChatShare");
			client.open_path(blog1);
			application1.release ();
		});
        act6.set_enabled(true);
		application1.add_action (act6);
		
		SimpleAction act7 = new SimpleAction ("quit", null);
		act7.activate.connect (() => {
			application1.hold ();
			// Print "Bye!" to our console:
			print ("Bye!\n");
			grid1.release_resource();
			// Terminate the mainloop: (main returns 0)
			Gtk.main_quit ();
			application1.release ();
		});
        act7.set_enabled(true);
		application1.add_action (act7);
		
		SimpleAction act8 = new SimpleAction ("delete-me", null);
		act8.activate.connect (() => {
			application1.hold ();
			var dlg = new Gtk.Dialog ();
			dlg.add_button(_("Yes"),1);
			dlg.add_button(_("No"),2);
			var text1 = new Gtk.Label(_("  You will Delete this user!  \n  Are you sure?"));
            var text2 = new Gtk.Label(_("Password："));
            var pwd1 = new Gtk.Entry();
            var grid = new Gtk.Grid();
            grid.attach(text1,0,0,2,1);
            grid.attach(text2,0,1);
            grid.attach(pwd1,1,1);
			dlg.get_content_area().pack_start(grid);
			grid.show_all();
			int r = dlg.run();
			if(r==1){
				client.delete_me(grid1.uname,pwd1.text);
			}
            dlg.destroy();
			application1.release ();
		});
        act8.set_enabled(true);
		application1.add_action (act8);
		
		SimpleAction act9 = new SimpleAction ("modify-desc", null);
		act9.activate.connect (() => {
			application1.hold ();
			//client.update_desc(desc);
			var dlg = new Gtk.Dialog ();
			dlg.add_button(_("Yes"),1);
			dlg.add_button(_("No"),2);
			var text1 = new Gtk.Label(_("New description:"));
			var desc1 = new Gtk.Entry();
			desc1.width_chars = 40;
			var g1 = new Gtk.Grid();
			g1.attach(text1,0,0);
			g1.attach(desc1,0,1);
			dlg.get_content_area().pack_start(g1);
			g1.show_all();
			int r = dlg.run();
			if(r==1){
				client.update_desc(desc1.text);
			}
			dlg.destroy();
			application1.release ();
		});
        act9.set_enabled(true);
		application1.add_action (act9);
		
		SimpleAction act10 = new SimpleAction ("down-page", null);
		act10.activate.connect (() => {
			application1.hold ();
			//client.open_path("https://github.com/rocket049/powerchat/releases");
			Gtk.Dialog dlg1 = new Gtk.Dialog.with_buttons(_("Upgrade"),app,Gtk.DialogFlags.MODAL);
			dlg1.set_size_request(300,200);
			var area = dlg1.get_content_area() as Gtk.Box;
			Gtk.Label label1 = new Gtk.Label("");
			label1.expand = true;
			string ln1 = _("This is the latest version!");
			string text1 = ln1;
			if(LATESTVER>RELEASE){
				ln1 = _("New Version Released!");
				string github = "https://github.com/rocket049/powerchat/releases";
				string gitee = "https://gitee.com/rocket049/powerchat/releases";
				text1 = @"\n  $(ln1)\n\n  <a href='$(github)'>$(github)</a>  \n\n  <a href='$(gitee)'>$(gitee)</a>  \n\n";
			}
			print(text1);
			label1.set_markup(text1);
			Gtk.Grid g1 = new Gtk.Grid();
			g1.attach(label1,0,0);
			g1.show_all();
			area.pack_start(g1);
			area.show_all();
			dlg1.show();
			application1.release ();
		});
        act10.set_enabled(true);
		application1.add_action (act10);
        
        SimpleAction act11 = new SimpleAction ("update-pwd", null);
		act11.activate.connect (() => {
			application1.hold ();
            grid1.update_pwd();
            application1.release ();
        });
        act11.set_enabled(true);
		application1.add_action (act11);
        
        SimpleAction act12 = new SimpleAction ("down-dir", null);
		act12.activate.connect (() => {
			application1.hold ();
			var home1 = GLib.Environment.get_home_dir();
			var blog1 = GLib.Path.build_path(GLib.Path.DIR_SEPARATOR_S,home1,".powerchat", "RecvFiles");
			client.open_path(blog1);
			application1.release ();
		});
        act12.set_enabled(true);
		application1.add_action (act12);
	}
}

public class LoginDialog :GLib.Object{
	public Gtk.Entry name;
	public Gtk.Entry passwd;
	public Gtk.Dialog dlg1;
	public LoginDialog(){
		this.dlg1 = new Gtk.Dialog.with_buttons(_("login"),app,Gtk.DialogFlags.MODAL);
		var grid = new Gtk.Grid();
		grid.attach(new Gtk.Label(_("Input name and password:")),0,0,2,1);
		grid.attach(new Gtk.Label(_("Login Name：")),0,1,1,1);
		grid.attach(new Gtk.Label(_("Password：")),0,2,1,1);
		this.name = new Gtk.Entry();
		grid.attach(this.name,1,1,1,1);
		this.passwd = new Gtk.Entry();
		this.passwd.set_visibility(false);
		grid.attach(this.passwd,1,2,1,1);
		var content = this.dlg1.get_content_area () as Gtk.Box;
		content.pack_start(grid);
		content.show_all();
		this.dlg1.add_button(_("Login"),2);
		this.dlg1.add_button(_("Register"),4);
		this.dlg1.add_button(_("Cancel"),3);
        this.load_name();
        this.passwd.activate.connect(()=>{
			this.login();
		});
		this.dlg1.response.connect((rid)=>{
			if (rid==2){
				//stdout.printf("next %d\n%s\n%s\n",rid,this.name.text,this.passwd.text);
				this.login();
			}else if(rid==4){
				this.dlg1.hide();
				adduser1 = new AddUserDialog();
				adduser1.show();
			}else{
				Gtk.main_quit();
			}
		});
	}
	public void login(){
		var u = client.login(this.name.text,this.passwd.text);
		if (u!=null){
			stdout.printf("login ok\n");
			grid1.uid = u.Id;
			grid1.uname = u.Name;
			grid1.usex = (int16)u.Sex;
			grid1.uage = (int16)u.Age;
			grid1.udesc = u.Desc;
			grid1.user_btn.label = _("About: ")+u.Name;
		}else{
			this.dlg1.title = _("Name/Password Error!");
			stdout.printf("login fail\n");
			return;
		}
		GLib.Idle.add(()=>{
			client.get_friends_async();
			return false;
		});
		save_name(this.name.text);
		app.show_all();
		this.dlg1.hide();
	}
	public int run(){
		return this.dlg1.run();
	}
	public void hide(){
		this.dlg1.hide();
	}
    
    public void load_name(){
        var loguser = GLib.Path.build_path(GLib.Path.DIR_SEPARATOR_S,Environment.get_home_dir(),".powerchat","manual","loguser.txt");
        try{
            GLib.File fp = GLib.File.new_for_path(loguser);
            var fs = fp.read();
            DataInputStream dis = new DataInputStream (fs as InputStream);
            string ns = dis.read_line();
            if (ns != null ){
                this.name.text = ns;
                this.passwd.grab_focus_without_selecting();
            }
        }catch (Error e) {
            print ("load name Error: %s\n", e.message);
        }
    } 
}
public void save_name(string name1){
	var loguser = GLib.Path.build_path(GLib.Path.DIR_SEPARATOR_S,Environment.get_home_dir(),".powerchat","manual","loguser.txt");
	GLib.File fp = GLib.File.new_for_path(loguser);
	GLib.FileOutputStream fs;
	try{
		fs = fp.create(FileCreateFlags.PRIVATE);
	}catch (Error e1) {
		try{
			fs = fp.replace(null,false,FileCreateFlags.PRIVATE);
		}catch (Error e2){
			print ("write name Error: %s\n", e2.message);
			return;
		}
	}
	DataOutputStream dos = new DataOutputStream (fs as OutputStream);
	dos.put_string( name1 );
	print("write name. \n");
}
public string get_cfg_dir(string name){
	var home1 = GLib.Environment.get_home_dir();
	var res = GLib.Path.build_path(GLib.Path.DIR_SEPARATOR_S,home1,".powerchat", name);
	GLib.DirUtils.create_with_parents(res,0755);
	return res;
}
public static string prog_path;
public void set_my_locale(string path1){
	var dir1 = GLib.Path.get_dirname(path1);
	prog_path = dir1;
	var textpath = GLib.Path.build_path(GLib.Path.DIR_SEPARATOR_S,prog_path,"..","share","locale");
	GLib.Intl.setlocale(GLib.LocaleCategory.ALL,"");
	GLib.Intl.textdomain("powerchat");
	GLib.Intl.bindtextdomain("powerchat",textpath);
	GLib.Intl.bind_textdomain_codeset ("powerchat", "UTF-8");
}
static uint16 server_port=7890;
static uint16 proxy_port;
public static int main(string[] args){
	set_my_locale(args[0]);
	if (!Thread.supported()) {
		stderr.printf("Cannot run without threads.\n");
		return 1;
	}
	if(args.length==2){
		server_port = (uint16)args[1].to_int64();
	}
	//proxy_port = server_port + 2000;
	client = new ChatClient();

	Gtk.init(ref args);
	grid1 = new MyGrid();
	app = new AppWin();
	app.append(grid1.mygrid);

	login1 = new LoginDialog();
	login1.dlg1.show_all();

	popup1 = new MyFriendMenu();
	GLib.Timeout.add_seconds(60,()=>{
		client.ping();
		return true;
	});
	Gtk.main ();
	client.quit();
	return 0;
}
