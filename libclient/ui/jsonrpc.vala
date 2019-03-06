using Jsonrpc;
using Gtk;
using Gdk;

public struct NUParam {
	public string Name;
	public int Sex;
	public int Birth;
	public string Desc;
	public string Pwd;
}

public struct UserData {
	public int64 Id;
	public string Name;
	public int Sex;
	public int Age;
	public string Desc;
}

public void search_add(int64 id,string name,int sex,int age,string desc,string msg, string msg_time){
	search1.seach_callback({id,sex,name,desc,age,msg,msg_time});
}
public void stranger_add(int64 id,string name,int sex,int age,string desc,string msg, string msg_time){
	strangers1.add_row({id,sex,name,desc,age,msg,msg_time});
}
public void friend_add(int64 id,string name,int sex,int age,string desc,string msg, string msg_time){
	grid1.add_friend({id,sex,name,desc,age,msg,msg_time});
}
public void client_notify(int8 typ,int64 from,int64 to,string msg){
	grid1.rpc_callback(typ,from,msg);
}

public class ChatClient:GLib.Object{
	public UserData u;
    public UserData? login(string name, string pwd){
		var ret = Client_Login(name,pwd);
		if(ret !=null){
			Client_SetNotifyFn((void*)client_notify);
			u=ret;
			print("LOGIN:");
			print(u.Name);
			return u;
		} else
			return null;
    }
    public void ping(){
		Client_Ping();
	}
    public int update_pwd(string name,string pwd1,string pwd2){
		return Client_NewPasswd(name,pwd1,pwd2);
	}
    public void add_friend(int64 uid){
		Client_AddFriend(uid);
	}
	//f(long long id,string name,int sex,int age,string desc,string msg, string msg_time)
	public void search_person_async(string key){
		Client_SearchPersons(key,(void*)search_add);
	}
	public void move_stranger_to_friend(int64 fid){
		Client_MoveStrangerToFriend(fid);
	}
	
	public void get_stranger_msgs_async(){
		Client_GetStrangerMsgs((void*)stranger_add);
	}

    public void get_friends_async(){
		Client_GetFriends((void*)friend_add);
    }
    public void ChatTo(int64 to , string msg){
		Client_ChatTo(to,msg);
	}
	public void tell(int64 to ){
		Client_Tell(to);
	}
	//send text message
	public void multi_send( int64[] to,string msg ){
		Client_MultiSend(msg,to,to.length);
	}
	
	public void tell_all(int64[] to ){
		Client_TellAll(to,to.length);
	}
	
	public void set_http_id(int64 uid){
		Client_SetHttpId(uid);
	}
	public void set_tcp_id(int64 uid){
		Client_SetHttpId(uid);
	}
	public void send_file(int64 to , string pathname){
		Client_SendFile(to,pathname);
	}
	public int add_user(string name,string pwd,int sex,int birthyear,string desc){
		return Client_NewUser({name,sex,birthyear,desc,pwd});
	}
	public void quit(){
		Client_Quit();
	}
	public int get_proxy(){
		return Client_GetProxyPort();
	}
	public string get_host(){
		return Client_GetHost();
	}
	public void set_proxy(int port){
		Client_ProxyPort(port);
	}
	public void remove_friend(int64 fid){
		if(fid==0){
			return;
		}
		Client_RemoveFriend(fid);
	}
	public void open_path(string path1){
		Client_OpenPath(path1);
	}
	public void update_desc(string desc){
		Client_UpdateDesc(desc);
	}
	public void delete_me(string name,string pwd){
		Client_DeleteMe(name,pwd);
	}
    public int user_status(int64 uid){
		return Client_UserStatus(uid);
    }
	
	public void offline_msg_with_id(int64 uid,string msg){
		if(uid==0){
			return;
		}
		var u = Client_QueryID(uid);
		var tm1 = new GLib.DateTime.now_local();
		UserMsg u1 = {u.Id,u.Sex,u.Name,u.Desc,u.Age,msg,tm1.format("%Y-%m-%d %H:%M:%S")};
		//show u1
		strangers1.prepend_row(u1);
		msg_notify(_("Strangers"));
	}
}

