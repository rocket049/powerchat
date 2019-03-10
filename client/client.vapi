[CCode (cheader_filename="libclient.h",has_type_id = false)]

public static void Client_SetHttpId(int64 p0);

public static int Client_NewUser(NUParam p0);

public static int Client_NewPasswd(string p0, string p1, string p2);

public static int Client_Login(string p0, string p1,ref UserData u0);

public static void Client_GetFriends(void* p0);

public static int Client_UserStatus(int64 p0);

public static void Client_QueryID(int64 p0,string msg,void* callfn);

public static void Client_MoveStrangerToFriend(int64 p0);

public static void Client_GetStrangerMsgs(void* p0);

public static void Client_SearchPersons(string p0, void* p1);

public static void Client_ChatTo(int64 p0, string p1);

public static void Client_Tell(int64 p0);

public static void Client_TellAll(int64* p0, int p1);

public static void Client_MultiSend(string p0, int64* p1, int p2);

public static void Client_Ping();

public static void Client_HttpConnect(int64 p0);

public static int Client_ProxyPort(int p0);

public static void Client_SetNotifyFn(void* p0);

public static void Client_SendFile(int64 p0, string p1);

public static void Client_AddFriend(int64 p0);

public static void Client_RemoveFriend(int64 p0);

public static int Client_GetProxyPort();

public static void Client_Quit();

public static void Client_GetHost(out string p0);

public static void Client_OpenPath(string p0);

public static void Client_UpdateDesc(string p0);

public static int Client_DeleteMe(string p0, string p1);

public static void Client_GetPgPath(out string p0);
