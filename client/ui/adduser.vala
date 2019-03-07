using Gtk;

public class AddUserDialog: GLib.Object{
	public Gtk.Dialog dlg1;
	public Gtk.Entry name;
	public Gtk.Entry pwd;
	public int birthyear;
	public Gtk.Entry desc;
	public int sex = 1;
	public AddUserDialog(){
		this.dlg1 = new Gtk.Dialog.with_buttons(_("Register"),app,Gtk.DialogFlags.MODAL);
		var grid = new Gtk.Grid();
		grid.attach(new Gtk.Label(_("Input your information")),0,0,2,1);
		grid.attach(new Gtk.Label(_("Name:")),0,1);
		this.name = new Gtk.Entry();
		grid.attach(this.name,1,1);
		
		grid.attach(new Gtk.Label(_("Password:")),0,2);
		this.pwd = new Gtk.Entry();
		this.pwd.set_visibility(false);
		grid.attach(this.pwd,1,2);
		
		grid.attach(new Gtk.Label(_("Confirm Password:")),0,3);
		var cfpwd = new Gtk.Entry();
		cfpwd.set_visibility(false);
		grid.attach(cfpwd,1,3);
		
		var date1 = new GLib.DateTime.now_local();
		int year1 = date1.get_year();
		Gtk.ComboBoxText yearSelect = new Gtk.ComboBoxText ();
		for(int i=5;i<100;i++){
			yearSelect.append_text(@"$(year1-i)");
		}
		yearSelect.active=10;
		this.birthyear = year1-15;
		grid.attach(new Gtk.Label(_("Born Year:")),0,4);
		grid.attach(yearSelect,1,4);
		yearSelect.changed.connect(()=>{
			string title = yearSelect.get_active_text ();
			this.birthyear = (int) title.to_int64();
		});
		
		Gtk.ComboBoxText sexSelect = new Gtk.ComboBoxText ();
		sexSelect.append_text (_("Man"));
		sexSelect.append_text (_("Woman"));
		sexSelect.active = 0;
		grid.attach(new Gtk.Label(_("Sex:")),0,5);
		grid.attach(sexSelect,1,5);
		sexSelect.changed.connect(()=>{
			this.sex = sexSelect.active + 1;
		});
		
		grid.attach(new Gtk.Label(_("Description:")),0,6);
		this.desc = new Gtk.Entry();
		grid.attach(this.desc,1,6);
		
		var content = this.dlg1.get_content_area () as Gtk.Box;
		content.pack_start(grid);
		//content.show_all();
		this.dlg1.add_button(_("Register"),2);
		this.dlg1.add_button(_("Cancel"),3);
		
		this.dlg1.response.connect((rid)=>{
			if(rid==2){
				//stdout.printf("%s %d %d\n",this.name.text,this.sex,this.birthyear);
				if(cfpwd.text != this.pwd.text){
					this.dlg1.title = _("Confirm Fail!");
					return;
				}
					
				if( client.add_user(this.name.text,this.pwd.text,this.sex,this.birthyear,this.desc.text)==0 ){
					this.dlg1.title=_("Register Fail!");
					return;
				}
				var u = client.login(this.name.text,this.pwd.text);
				if (u!=null){
					stdout.printf("login ok\n");
					grid1.uid = u.Id;
					grid1.uname = u.Name;
					grid1.usex = (int16)u.Sex;
					grid1.uage = (int16)u.Age;
					grid1.udesc = u.Desc;
					grid1.user_btn.label = _("About: ")+u.Name;
					save_name(this.name.text);
				}else{
					this.dlg1.title = _("Name/Password Error!");
					stdout.printf("login fail\n");
					return;
				}
				client.get_friends_async();
				this.dlg1.destroy();
				app.show_all();
			}else{
				Gtk.main_quit();
			}
		});
	}
	public void show(){
		this.dlg1.show_all();
	}
}
