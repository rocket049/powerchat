using Gtk;
using Gee;
struct IDArray{
	int64[] ids;
}
public class MultiSendUi: GLib.Object{
	Gtk.Dialog frame1;
	Gtk.ListBox friends;
	Gtk.ListBox groups;
	Gtk.Entry entry1;  //edit message
	Gtk.Button send1;  //send message
	Gtk.Button save1;  //save new group
	Gee.HashSet<int64?> ids;   //store id selected
	Gee.HashMap<string,IDArray?> group_map;
	
	public MultiSendUi(){
		frame1 = new Gtk.Dialog();
		frame1.set_modal(true);
		frame1.title=_("MultiSend");
		frame1.set_size_request(480,640);
		friends = new Gtk.ListBox();
		groups = new Gtk.ListBox();
		entry1 = new Gtk.Entry();
		send1 = new Gtk.Button.with_label(_("Send"));
		save1 = new Gtk.Button.with_label("=>");
		save1.set_size_request(24,12);
		save1.expand=false;
		
		ids = new Gee.HashSet<int64?>(null,(a,b)=>{
			return (a==b);
		});
		group_map = new Gee.HashMap<string,IDArray?>();
		
		var grid1 = new Gtk.Grid();
		var box1 = frame1.get_content_area ();
		box1.pack_start(grid1);
		var label1 = new Gtk.Label(_("My Friends"));
		var label2 = new Gtk.Label(_("Groups"));
		grid1.attach(label1,1,0);
		grid1.attach(label2,3,0);
		
		var scrollWin1 = new Gtk.ScrolledWindow(null,null);
		scrollWin1.width_request = 200;
		scrollWin1.expand = true;
		scrollWin1.add(friends);
		grid1.attach(scrollWin1,1,1,1,3);
		
		var lb1 = new Gtk.Label("");
		lb1.expand=true;
		grid1.attach(lb1,2,1);
		grid1.attach(save1,2,2);
		var lb2 = new Gtk.Label("");
		lb2.expand=true;
		grid1.attach(lb2,2,3);
		
		scrollWin1 = new Gtk.ScrolledWindow(null,null);
		scrollWin1.width_request = 200;
		scrollWin1.expand = true;
		scrollWin1.add(groups);
		grid1.attach(scrollWin1,3,1,1,3);
		
		grid1.attach(new Gtk.Label(_("Text:")),0,4);
		grid1.attach(entry1,1,4,3,1);
		grid1.attach(send1,4,4);
		
		grid1.show_all();
		
		set_css();
		set_lists();
		set_send_callback();
		set_save_callback();
		set_group_callback();
	}
	private void set_group_callback(){
		groups.row_selected.connect((r)=>{
			if(group_map.has_key(r.name)==false){
				return;
			}
			var array1 = group_map[r.name];
			for(int i=0;i<array1.ids.length;i++){
				print(@"$(array1.ids[i])\n");
			}
		});
	}
	private void set_save_callback(){
		save1.clicked.connect(()=>{
			var name1 = new Gtk.Entry();
			groups.add(name1);
			name1.text = @"$(ids.size) persons group";
			name1.show();
			name1.activate.connect(()=>{
				if(group_map.has_key(name1.text)){
					return;
				}
				var label1 = new Gtk.Label(name1.text);
				groups.remove(name1);
				name1.destroy();
				groups.add(label1);
				label1.show();
				
				var idsi = new int64[ids.size];
				int i=0;
				foreach( int64 id1 in ids){
					idsi[i] = id1;
					i++;
				}
				group_map[name1.text] = {idsi};
				(label1.parent as Gtk.ListBoxRow).name = name1.text;
			});
		});
	}
	private void set_send_callback(){
		send1.clicked.connect(()=>{
			var idsi = new int64[ids.size];
			int i=0;
			foreach( int64 id1 in ids){
				idsi[i] = id1;
				i++;
			}
			if (entry1.text.length==0){
				return;
			}
			rpc1.multi_send(idsi,entry1.text);
		});
	}
	private void set_css(){
		var cssp = new Gtk.CssProvider();
		var sc = friends.get_style_context ();
		sc.add_provider(cssp,Gtk.STYLE_PROVIDER_PRIORITY_USER);
		sc = groups.get_style_context ();
		sc.add_provider(cssp,Gtk.STYLE_PROVIDER_PRIORITY_USER);
		cssp.load_from_data("""
list{
	background-color:#FFFFFF;
	color:#000000;
}
""");
	}
	public void show(){
		frame1.show_all();
	}
	private void set_lists(){
		var suid = grid1.uid.to_string();
		foreach(string id1 in grid1.frds1.keys){
			if(suid==id1){
				continue;
			}
			add_friend(grid1.frds1[id1]);
		}
	}
	private void add_friend(UserData u1){
		var button1 = new Gtk.CheckButton.with_label(u1.name);
		button1.name = u1.id.to_string();
		friends.add(button1);
		button1.show();
		button1.toggled.connect(()=>{
			if (button1.active){
				ids.add(u1.id);
			}else{
				ids.remove(u1.id);
			}
		});
	}
}
