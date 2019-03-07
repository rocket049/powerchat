using Gtk;
using Gdk;

static MyFriendMenu popup1;

public class MyFriendMenu : Gtk.Menu{
	private string friend_id;
	public void set_id(string id){
		friend_id = id;
	}
	public MyFriendMenu(){
		var item1 = new Gtk.MenuItem.with_label (_("Test Online"));
		item1.activate.connect(()=>{
			//stdout.printf("menu: %s\n",this.friend_id);
			grid1.user_online(this.friend_id.to_int64());
			rpc1.tell(this.friend_id.to_int64());
		});
		this.append(item1);
		var item2 = new Gtk.MenuItem.with_label (_("Delete"));
		item2.activate.connect(()=>{
			//stdout.printf("menu: %s\n",this.friend_id);
			grid1.remove_friend(this.friend_id);
		});
		this.append(item2);
		this.show_all();
	}
}
