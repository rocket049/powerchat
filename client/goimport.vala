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
	Idle.add(()=>{
		search1.seach_callback({id,sex,name,desc,age,msg,msg_time});
		return false;});
}
public void stranger_add(int64 id,string name,int sex,int age,string desc,string msg, string msg_time){
	Idle.add(()=>{
		strangers1.prepend_row({id,sex,name,desc,age,msg,msg_time});
		return false;});
}
public void friend_add(int64 id,string name,int sex,int age,string desc,string msg, string msg_time){
	Idle.add(()=>{
		grid1.add_friend({id,sex,name,desc,age,msg,msg_time});
		return false;});
}
public void client_notify(int8 typ,int64 from,int64 to,string msg){
	Idle.add(()=>{
		grid1.rpc_callback(typ,from,msg);
		return false;});
}

public class ChatClient:GLib.Object{
    public UserData? login(string name, string pwd){
		UserData u = {0,"",0,0,""};
		var ret = Client_Login(name,pwd,ref u);
		if(ret == 1){
			Client_SetNotifyFn((void*)client_notify);
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
		var thread = new Thread<int>("search_person", ()=>{
			Client_SearchPersons(key,(void*)search_add);
			return 0;});
	}
	public void move_stranger_to_friend(int64 fid){
		Client_MoveStrangerToFriend(fid);
	}
	
	public void get_stranger_msgs_async(){
		var thread = new Thread<int>("get_strangers", ()=>{
			Client_GetStrangerMsgs((void*)stranger_add);
			return 0;});
	}

    public void get_friends_async(){
		var thread = new Thread<int>("get_friends", ()=>{
			Client_GetFriends((void*)friend_add);
			return 0;});
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
		var ret = Client_SendFile(to,pathname);
		switch (ret){
		case 1:
			Idle.add(()=>{
				grid1.add_text(_("Fail send file : unknow error"));
				return false;
			});
			
			break;
		case 2:
			Idle.add(()=>{
				grid1.add_text(_("Fail send file : peer block"));
				return false;
			});
			
			break;
		}
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
		string p;
		Client_GetHost(out p);
		return p.dup();
	}
	public int set_proxy(int port){
		return Client_ProxyPort(port);
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
		var thread = new Thread<int>("search_person", ()=>{
			Client_QueryID(uid,msg,(void*)stranger_add);
			return 0;});
	}
	public string get_pg_path(){
		string p;
		Client_GetPgPath(out p);
		return p.dup();
	}
}

