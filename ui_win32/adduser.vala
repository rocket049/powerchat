using Gtk;

public class AddUserDialog: GLib.Object{
	public Gtk.Dialog dlg1;
	public Gtk.Entry name;
	public Gtk.Entry pwd;
	public int birthyear;
	public Gtk.Entry desc;
	public int sex = 1;
	public AddUserDialog(){
		this.dlg1 = new Gtk.Dialog.with_buttons("注册",app,Gtk.DialogFlags.MODAL);
		var grid = new Gtk.Grid();
		grid.attach(new Gtk.Label("输入注册信息"),0,0,2,1);
		grid.attach(new Gtk.Label("名　　称："),0,1);
		this.name = new Gtk.Entry();
		grid.attach(this.name,1,1);
		
		grid.attach(new Gtk.Label("密　　码："),0,2);
		this.pwd = new Gtk.Entry();
		this.pwd.set_visibility(false);
		grid.attach(this.pwd,1,2);
		
		var date1 = new GLib.DateTime.now_local();
		int year1 = date1.get_year();
		Gtk.ComboBoxText yearSelect = new Gtk.ComboBoxText ();
		for(int i=5;i<100;i++){
			yearSelect.append_text(@"$(year1-i)");
		}
		yearSelect.active=10;
		this.birthyear = year1-15;
		grid.attach(new Gtk.Label("出生年份："),0,3);
		grid.attach(yearSelect,1,3);
		yearSelect.changed.connect(()=>{
			string title = yearSelect.get_active_text ();
			this.birthyear = (int) title.to_int64();
		});
		
		Gtk.ComboBoxText sexSelect = new Gtk.ComboBoxText ();
		sexSelect.append_text ("男性");
		sexSelect.append_text ("女性");
		sexSelect.active = 0;
		grid.attach(new Gtk.Label("性别："),0,4);
		grid.attach(sexSelect,1,4);
		sexSelect.changed.connect(()=>{
			this.sex = sexSelect.active + 1;
		});
		
		grid.attach(new Gtk.Label("自我描述："),0,5);
		this.desc = new Gtk.Entry();
		grid.attach(this.desc,1,5);
		
		var content = this.dlg1.get_content_area () as Gtk.Box;
		content.pack_start(grid);
		//content.show_all();
		this.dlg1.add_button("注册",2);
		this.dlg1.add_button("取消",3);
		
		this.dlg1.response.connect((rid)=>{
			if(rid==2){
				//stdout.printf("%s %d %d\n",this.name.text,this.sex,this.birthyear);
				if( rpc1.add_user(this.name.text,this.pwd.text,this.sex,this.birthyear,this.desc.text)==false ){
					this.dlg1.title="注册发生错误！";
					return;
				}
				UserData u;
				var res = rpc1.login(this.name.text,this.pwd.text,out u);
				if (res>0){
					stdout.printf("login ok\n");
					grid1.uid = res;
					grid1.uname = u.name;
					grid1.usex = u.sex;
					grid1.uage = u.age;
					grid1.udesc = u.desc;
					grid1.user_btn.label = u.name;
				}else{
					this.dlg1.title = "用户／密码错误。";
					stdout.printf("login fail\n");
					return;
				}
				if(rpc1.get_friends_async()==false){
					print("RPC error");
					Gtk.main_quit();
				}
				this.dlg1.hide();
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
