using Gtk;

static StrangersDialg strangers1;

public class StrangersDialg:GLib.Object{
	public Gtk.Dialog dlg1;
	public Gtk.Entry key1;
	public Gtk.ListStore store1;
	public GLib.List<UserData?> persons;
	public StrangersDialg(){
		this.store1 = new Gtk.ListStore (8, typeof (string),typeof(string),typeof(string), typeof(int64), typeof (string), typeof(int16), typeof(string), typeof(string));
		this.get_msgs();
	}
	public void create_dialg(){
		this.dlg1 = new Gtk.Dialog.with_buttons(_("Information of Strangers"),app,Gtk.DialogFlags.MODAL);
		var view = new Gtk.TreeView.with_model(this.store1);
		//Gtk.CellRendererText cell = new Gtk.CellRendererText ();
		view.insert_column_with_attributes (0, _("Name"), new Gtk.CellRendererText (), "text",0,"background",7);
		view.insert_column_with_attributes (1, _("Said"), new Gtk.CellRendererText (), "text",1,"background",7);
		view.insert_column_with_attributes (2, _("Time"), new Gtk.CellRendererText (), "text",2,"background",7);
		view.insert_column_with_attributes (3, "ID", new Gtk.CellRendererText (), "text",3,"background",7);
		view.insert_column_with_attributes (4, _("Sex"), new Gtk.CellRendererText (), "text",4,"background",7);
		view.insert_column_with_attributes (5, _("Age"), new Gtk.CellRendererText (), "text",5,"background",7);
		view.insert_column_with_attributes (6, _("Description"), new Gtk.CellRendererText (), "text",6,"background",7);
		view.headers_visible = true;
		view.show_all();
		
		var scroll1 = new Gtk.ScrolledWindow(null,null);
		scroll1.add(view);
		scroll1.set_size_request(480,480);
		view.expand = true;
		scroll1.expand = true;
		
		var content = this.dlg1.get_content_area () as Gtk.Box;
		content.pack_start(scroll1);
		
		this.dlg1.add_button(_("Close"),2);
		this.dlg1.response.connect((rid)=>{
			this.dlg1.close();
			(this.store1 as Gtk.TreeModel).foreach((m,p,iter)=>{
				this.store1.set (iter,7,"#FFFFFF");
				return false;
			});
		});
		
		view.row_activated.connect( (tree,path,col)=>{
			//stdout.printf("%s\n",path.to_string());
			Gtk.TreeIter iter;
			var model = tree.get_model();
			model.get_iter(out iter,path);
			var idv = Value(typeof (int64));
			model.get_value(iter,3,out idv);
			int64 id = idv.get_int64();
			rpc1.move_stranger_to_friend(id);
			foreach( UserData u in this.persons ){
				if(u.id==id){
					//stdout.printf("add: %s %s\n",u.name,u.desc);
					grid1.add_friend(u);
				}
			}
		});
	}
	public void add_row(UserData u1){
		this.persons.append(u1);
		Gtk.TreeIter iter;
		this.store1.append (out iter);
		string sex=_("Unknown");
		if(u1.sex==1){
			sex = _("Man");
		}else if(u1.sex==2){
			sex = _("Woman");
		}
		this.store1.set (iter, 0, u1.name, 1, u1.msg_offline,2,u1.timestamp_offline,3,u1.id,4,sex,5,u1.age,6,u1.desc,7,"#FFFFFF");
		//stdout.printf("%s %s\n",u1.name,u1.desc);
	}
	public void prepend_row(UserData u1){
		this.persons.prepend(u1);
		Gtk.TreeIter iter;
		this.store1.prepend (out iter);
		string sex=_("Unknown");
		if(u1.sex==1){
			sex = _("Man");
		}else if(u1.sex==2){
			sex = _("Woman");
		}
		this.store1.set (iter, 0, u1.name, 1, u1.msg_offline,2,u1.timestamp_offline,3,u1.id,4,sex,5,u1.age,6,u1.desc,7,"#F75656");
		var sc1 = grid1.strangers_btn.get_style_context();
		sc1.add_provider(grid1.button1,Gtk.STYLE_PROVIDER_PRIORITY_USER);
		sc1.add_class("off");
	}
	private void get_msgs(){
		this.persons = new GLib.List<UserData?>();
		rpc1.get_stranger_msgs_async(this.add_row);
		grid1.strangers_btn.show_all();
	}
	public void show(){
		this.create_dialg();
		var sc1 = grid1.strangers_btn.get_style_context();
		sc1.remove_provider(grid1.button1);
		sc1.remove_class("off");
		this.dlg1.show_all();
	}
}
