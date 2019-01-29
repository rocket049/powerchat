using Jsonrpc;
using Gtk;
using Gdk;

public class RpcClient:GLib.Object{
    public Jsonrpc.Client c;
    public int counter{get;set;default=0;}
    public bool connect(string host,uint16 port){
        Resolver resolver = Resolver.get_default ();
		List<InetAddress> addresses = resolver.lookup_by_name (host, null);
		for(uint i=0;i<addresses.length();i++){
			try{
				
				InetAddress address = addresses.nth_data (i);
				SocketClient client = new SocketClient ();
				SocketConnection conn = client.connect(new InetSocketAddress (address, port));

				this.c = new Jsonrpc.Client(conn);
				return true;
			}catch (Error e) {
				stdout.printf ("Error: %s\n", e.message);
			}
		}
		return false;
    }
    public int64 login(string name, string pwd,out UserData u){
        var params = new Variant.parsed("{'Name':<%s>,'Pwd':<%s>}",name,pwd);
        Variant res;
        try{
            var ok = c.call("PClient.Login",params,null,out res);
            if(ok==false){
               return -1;
            }
            //stdout.printf("%"+int64.FORMAT+"\n",res.get_int64());
            var id = res.lookup_value("Id",null).get_int64();
            var sex = res.lookup_value("Sex",null).get_int64();
            var uname = res.lookup_value("Name",null).get_string();
            var desc = res.lookup_value("Desc",null).get_string();
            var age = res.lookup_value("Age",null).get_int64();
            u = {id,(int16)sex,uname,desc,(int16)age,"",""};
            return id;
        }catch (Error e) {
            stdout.printf ("Error: %s\n", e.message);
            return -1;
        }
    }
    public bool ping(){
		var params = new Variant.bytestring("");
		try{
            c.call_async.begin("PClient.Ping",params,null,(s,r)=>{c.call_async.end(r,null);});
            return true;
        }catch (Error e) {
            stdout.printf ("update_pwd Error: %s\n", e.message);
            return false;
        }
	}
    public bool update_pwd(string name,string pwd1,string pwd2){
		var params = new Variant("(sss)",name,pwd1,pwd2);
		try{
            c.call_async.begin("PClient.NewPasswd",params,null,(s,r)=>{c.call_async.end(r,null);});
            return true;
        }catch (Error e) {
            stdout.printf ("update_pwd Error: %s\n", e.message);
            return false;
        }
	}
    public bool add_friend(int64 uid){
		var params = new Variant("(i)",uid);
        //Variant res;
        try{
            c.call_async.begin("PClient.AddFriend",params,null,(s,r)=>{c.call_async.end(r,null);});
            return true;
        }catch (Error e) {
            stdout.printf ("Error: %s\n", e.message);
            return false;
        }
	}
	
	public bool search_person_async(string key,SearchCallback f){
		var params = new Variant("(s)",key);
        try{
            c.call_async.begin("PClient.SearchPersons",params,null,(s,r)=>{
				Variant res;
				c.call_async.end(r,out res);
				Variant val;
				var size1 = res.n_children();
				for(size_t i=0;i<size1;i++){
					val = res.get_child_value(i).get_child_value(0);
					var id = val.lookup_value("Id",null).get_int64();
					var sex = val.lookup_value("Sex",null).get_int64();
					var name = val.lookup_value("Name",null).get_string();
					var desc = val.lookup_value("Desc",null).get_string();
					var age = val.lookup_value("Age",null).get_int64();
					//stdout.printf("%s %s\n",name,desc);
					UserData u1 = {id,(int16)sex,name,desc,(int16)age,"",""};
					f(u1);
				}
			});
            return true;
        }catch (Error e) {
            stdout.printf ("Error: %s\n", e.message);
            return false;
        }
	}
	public bool move_stranger_to_friend(int64 fid){
		var params = new Variant("(i)",fid);
        Variant res;
        try{
            var ok = c.call("PClient.MoveStrangerToFriend",params,null,out res);
            return ok;
        }catch (Error e) {
            stdout.printf ("Error: %s\n", e.message);
            return false;
        }
	}
	
	public delegate void callback_add_row(UserData u1);
	public bool get_stranger_msgs_async( callback_add_row add_row ){
		var params = new Variant.bytestring("");
        try{
			c.call_async.begin("PClient.GetStrangerMsgs",params,null,(s,r)=>{
				Variant val,res;
				c.call_async.end(r,out res);
				var size1 = res.n_children();
				for(size_t i=0;i<size1;i++){
					val = res.get_child_value(i).get_child_value(0);
					var id = val.lookup_value("Id",null).get_int64();
					var sex = val.lookup_value("Sex",null).get_int64();
					var name = val.lookup_value("Name",null).get_string();
					var desc = val.lookup_value("Desc",null).get_string();
					var age = val.lookup_value("Age",null).get_int64();
					var obj_offline = val.lookup_value("MsgOffline",null);
					var msg_offline = obj_offline.lookup_value("Msg",null).get_string();
					var timestamp_offline = obj_offline.lookup_value("Timestamp",null).get_string();
					//stdout.printf("%s %s\n",name,desc);
					UserData u1 = {id,(int16)sex,name,desc,(int16)age,msg_offline,timestamp_offline};
					add_row(u1);
				}
				if(size1>0){
					var sc1 = grid1.strangers_btn.get_style_context();
					sc1.add_provider(grid1.button1,Gtk.STYLE_PROVIDER_PRIORITY_USER);
					sc1.add_class("off");
				}
			});
            return true;
		}catch (Error e) {
            stdout.printf ("Error: %s\n", e.message);
            return false;
        }
	}

    public bool get_friends_async(){
        var params = new Variant.bytestring("");
        
        try{
            c.call_async.begin("PClient.GetFriends",params,null,(s,r)=>{
				Variant res;
				c.call_async.end(r,out res);
				Variant val;
				var size1 = res.n_children();
				for(size_t i=0;i<size1;i++){
					val = res.get_child_value(i).get_child_value(0);
					var id = val.lookup_value("Id",null).get_int64();
					var sex = val.lookup_value("Sex",null).get_int64();
					var name = val.lookup_value("Name",null).get_string();
					var desc = val.lookup_value("Desc",null).get_string();
					var age = val.lookup_value("Age",null).get_int64();
					var obj_offline = val.lookup_value("MsgOffline",null);
					var msg_offline = obj_offline.lookup_value("Msg",null).get_string();
					var timestamp_offline = obj_offline.lookup_value("Timestamp",null).get_string();
					//stdout.printf("%s %s\n",name,desc);
					UserData u1 = {id,(int16)sex,name,desc,(int16)age,msg_offline,timestamp_offline};
					grid1.add_friend(u1);
				}
			});
            
            return true;
        }catch (Error e) {
            stdout.printf ("Error: %s\n", e.message);
            return false;
        }
    }
    public bool ChatTo(int64 to , string msg){
		//stdout.printf("%s\n",msg);
		var v1 = new Variant.int64(to);
		var v2 = new Variant.string(msg);
		var params = new Variant.parsed("{'To':%v,'Msg':%v}",v1,v2 );
        //Variant res;
        try{
            c.call_async.begin("PClient.ChatTo",params,null,(s,r)=>{c.call_async.end(r,null);});
            return true;
        }catch (Error e) {
            stdout.printf ("Error: %s\n", e.message);
            return false;
        }
	}
	public bool tell(int64 to ){
		var params = new Variant("(i)",to);
        //Variant res;
        try{
            c.call_async.begin("PClient.Tell",params,null,(s,r)=>{c.call_async.end(r,null);});
            return true;
        }catch (Error e) {
            stdout.printf ("Error: %s\n", e.message);
            return false;
        }
	}
	
	public bool set_http_id(int64 uid){
		var params = new GLib.Variant("(i)",uid);
		//GLib.Variant res;
		try{
            c.call_async.begin("PClient.HttpConnect",params,null,(s,r)=>{
				c.call_async.end(r,null);
				Gtk.show_uri(null,@"http://localhost:$(server_port)/",Gdk.CURRENT_TIME);
			});
            return true;
        }catch (Error e) {
            stdout.printf ("Error: %s\n", e.message);
            return false;
        }
	}
	public bool set_tcp_id(int64 uid){
		var params = new GLib.Variant("(i)",uid);
		//GLib.Variant res;
		try{
            c.call_async.begin("PClient.HttpConnect",params,null,(s,r)=>{
				c.call_async.end(r,null);
			});
            return true;
        }catch (Error e) {
            stdout.printf ("Error: %s\n", e.message);
            return false;
        }
	}
	public bool send_file(int64 to , string pathname){
		var v1 = new Variant.int64(to);
		var v2 = new Variant.string(pathname);
		var params = new Variant.parsed("{'To':%v,'PathName':%v}",v1,v2 );
		//GLib.Variant res;
		try{
            c.call_async.begin("PClient.SendFile",params,null,(s,r)=>{
				try{
					c.call_async.end(r,null);
				}catch(Error e){
					//stdout.printf ("sendfile return Error: %s\n", e.message);
					grid1.add_text(_("[TransferBlocked]"),true);
				}
			});
            return true;
        }catch (Error e) {
            stdout.printf ("sendfile Error: %s\n", e.message);
            return false;
        }
	}
	public bool add_user(string name,string pwd,int sex,int birthyear,string desc){
		var params = new GLib.Variant.parsed("{'Name':<%s>,'Sex':<%i>,'Birth':<%i>,'Desc':<%s>,'Pwd':<%s>}",name,sex,birthyear,desc,pwd);
		GLib.Variant res;
		try{
            var ok = c.call("PClient.NewUser",params,null,out res);
            if(ok==false)
				return false;
			if(res.get_int64()==1)
				return true;
			else
				return false;
        }catch (Error e) {
            stdout.printf ("add user Error: %s\n", e.message);
            return false;
        }
	}
	public bool quit(){
		var params = new Variant.bytestring("");
        //Variant res;
        try{
            c.call_async.begin("PClient.Quit",params,null,(s,r)=>{c.call_async.end(r,null);});
            return true;
        }catch (Error e) {
            stdout.printf ("Error: %s\n", e.message);
            return false;
        }
	}
	public bool get_proxy(out uint16 port){
		var params = new Variant.bytestring("");
        Variant res;
        try{
            var ok = c.call("PClient.GetProxyPort",params,null,out res);
            if(ok == false){
				return false;
			}
            port = (uint16)res.get_int64();
            return true;
        }catch (Error e) {
            stdout.printf ("Error: %s\n", e.message);
            return false;
        }
	}
	public bool get_host(out string host){
		var params = new Variant.bytestring("");
        Variant res;
        try{
            var ok = c.call("PClient.GetHost",params,null,out res);
            if(ok == false){
				return false;
			}
            host = res.get_string();
            return true;
        }catch (Error e) {
            stdout.printf ("Error: %s\n", e.message);
            return false;
        }
	}
	public bool set_proxy(uint16 port){
		var params = new Variant("(i)",port);
        
        try{
            c.call_async.begin("PClient.ProxyPort",params,null,(s,r)=>{
				Variant res;
				c.call_async.end(r,out res);
				proxy_port = (uint16)res.get_int64();
				grid1.port1.text = proxy_port.to_string();
			});
            return true;
        }catch (Error e) {
            stdout.printf ("Error: %s\n", e.message);
            return false;
        }
	}
	public bool remove_friend(int64 fid){
		if(fid==0){
			return false;
		}
		var params = new Variant("(i)",fid);
		try{
            c.call_async.begin("PClient.RemoveFriend",params,null,(s,r)=>{c.call_async.end(r,null);});
            return true;
        }catch (Error e) {
            stdout.printf ("Error: %s\n", e.message);
            return false;
        }
	}
	public bool open_path(string path1){
		var params = new Variant("(s)",path1);
		try{
            c.call_async.begin("PClient.OpenPath",params,null,(s,r)=>{c.call_async.end(r,null);});
            return true;
        }catch (Error e) {
            stdout.printf ("Error: %s\n", e.message);
            return false;
        }
	}
	
	public delegate void QueryCallback(UserData u);
	public bool offline_msg_with_id(int64 uid,string msg, QueryCallback f){
		if(uid==0){
			return false;
		}
		var params = new Variant("(i)",uid);
		try{
            c.call_async.begin("PClient.QueryID",params,null,(s,r)=>{
				GLib.Variant res;
				c.call_async.end(r,out res);
				var id = res.lookup_value("Id",null).get_int64();
				var name = res.lookup_value("Name",null).get_string();
				var age = (int16) res.lookup_value("Age",null).get_int64();
				var sex = (int16) res.lookup_value("Sex",null).get_int64();
				var desc = res.lookup_value("Desc",null).get_string();
				//grid1.add_text(@"$(id) : $(name)");
				var tm1 = new GLib.DateTime.now_local();
				UserData u = {id,sex,name,desc,age,msg,tm1.format("%Y-%m-%d %H:%M:%S")};
				f(u);
			});
            return true;
        }catch (Error e) {
            stdout.printf ("Error: %s\n", e.message);
            return false;
        }
	}
}

