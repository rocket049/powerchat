using Gtk;

static SearchDialg search1;
public delegate void SearchCallback(UserMsg user1);

public class SearchDialg:GLib.Object{
	public Gtk.Dialog dlg1;
	public Gtk.Entry key1;
	public Gtk.ListStore store1;
	public GLib.List<UserMsg?> persons;
	public SearchDialg(){
		this.dlg1 = new Gtk.Dialog.with_buttons(_("Find Persons"),app,Gtk.DialogFlags.MODAL);
        this.dlg1.set_size_request(350,400);
		var grid = new Gtk.Grid();
		grid.attach(new Gtk.Label(_("(part of)Name：")),0,0);
		this.key1 = new Gtk.Entry();
        key1.hexpand = true;
		grid.attach(this.key1,1,0);
		var b1 = new Gtk.Button.with_label(_("Search"));
		grid.attach(b1,2,0);
		b1.clicked.connect( this.search );
		
		this.store1 = new Gtk.ListStore (5, typeof(int64), typeof (string), typeof (string), typeof(int16), typeof(string));
		
		var view = new Gtk.TreeView.with_model(this.store1);
		//Gtk.CellRendererText cell = new Gtk.CellRendererText ();
		view.insert_column_with_attributes (0, "ID", new Gtk.CellRendererText(), "text",0);
		view.insert_column_with_attributes (1, _("Name"), new Gtk.CellRendererText(), "text",1);
		view.insert_column_with_attributes (2, _("Sex"), new Gtk.CellRendererText(), "text",2);
		view.insert_column_with_attributes (3, _("Age"), new Gtk.CellRendererText(), "text",3);
		view.insert_column_with_attributes (4, _("Description"), new Gtk.CellRendererText(), "text",4);
		view.headers_visible = true;
		view.show_all();
		
		var scroll1 = new Gtk.ScrolledWindow(null,null);
		scroll1.add(view);
		view.expand = true;
		scroll1.expand = true;
		grid.attach(scroll1,0,1,3,1);
		
		var content = this.dlg1.get_content_area () as Gtk.Box;
		content.pack_start(grid);
		this.dlg1.add_button(_("Close"),2);
		this.dlg1.response.connect((rid)=>{
			this.dlg1.close();
		});
		
		view.row_activated.connect( (tree,path,col)=>{
			//stdout.printf("%s\n",path.to_string());
			Gtk.TreeIter iter;
			var model = tree.get_model();
			model.get_iter(out iter,path);
			var idv = Value(typeof (int64));
			model.get_value(iter,0,out idv);
			int64 id = idv.get_int64();
			client.add_friend(id);
			foreach( UserMsg u in this.persons ){
				if(u.id==id){
					//stdout.printf("add: %s %s\n",u.name,u.desc);
					grid1.add_friend(u);
					//client.tell(id);
				}
			}
		});
	}
	public void add_row(UserMsg u1){
		Gtk.TreeIter iter;
		this.store1.append (out iter);
		string sex=_("Unknown");
		if(u1.sex==1){
			sex = _("Man");
		}else if(u1.sex==2){
			sex = _("Woman");
		}
		this.store1.set (iter, 0, u1.id, 1, u1.name,2,sex,3,u1.age,4,u1.desc);
		//stdout.printf("%s %s\n",u1.name,u1.desc);
	}
	public void seach_callback(UserMsg u){
		this.add_row(u);
		this.persons.append(u);
	}
	public void search(){
		this.store1.clear();
		this.persons = new GLib.List<UserMsg?>();
		this.store1.set_sort_column_id(0,Gtk.SortType.ASCENDING);
		
		client.search_person_async(this.key1.text);
	}
	public void show(){
		this.dlg1.show_all();
	}
}
