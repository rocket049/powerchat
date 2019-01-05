using Gtk;
using Gdk;

static MyFriendMenu popup1;

public class MyFriendMenu : Gtk.Menu{
	private string friend_id;
	public void set_id(string id){
		friend_id = id;
	}
	public MyFriendMenu(){
		var item1 = new Gtk.MenuItem.with_label (_("Delete"));
		item1.activate.connect(()=>{
			stdout.printf("menu: %s\n",this.friend_id);
			grid1.remove_friend(this.friend_id);
		});
		this.attach(item1,0,1,0,1);
		this.show_all();
	}
}
