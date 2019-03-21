/* goimport.vala.c generated by valac 0.42.6, the Vala compiler
 * generated from goimport.vala, do not modify */



#include <glib.h>
#include <glib-object.h>
#include <stdlib.h>
#include <string.h>
#include "libclient.h"


#define TYPE_NU_PARAM (nu_param_get_type ())
typedef struct _NUParam NUParam;
#define _g_free0(var) (var = (g_free (var), NULL))

#define TYPE_USER_DATA (user_data_get_type ())
typedef struct _UserData UserData;
typedef struct _Block8Data Block8Data;

#define TYPE_SEARCH_DIALG (search_dialg_get_type ())
#define SEARCH_DIALG(obj) (G_TYPE_CHECK_INSTANCE_CAST ((obj), TYPE_SEARCH_DIALG, SearchDialg))
#define SEARCH_DIALG_CLASS(klass) (G_TYPE_CHECK_CLASS_CAST ((klass), TYPE_SEARCH_DIALG, SearchDialgClass))
#define IS_SEARCH_DIALG(obj) (G_TYPE_CHECK_INSTANCE_TYPE ((obj), TYPE_SEARCH_DIALG))
#define IS_SEARCH_DIALG_CLASS(klass) (G_TYPE_CHECK_CLASS_TYPE ((klass), TYPE_SEARCH_DIALG))
#define SEARCH_DIALG_GET_CLASS(obj) (G_TYPE_INSTANCE_GET_CLASS ((obj), TYPE_SEARCH_DIALG, SearchDialgClass))

typedef struct _SearchDialg SearchDialg;
typedef struct _SearchDialgClass SearchDialgClass;

#define TYPE_USER_MSG (user_msg_get_type ())
typedef struct _UserMsg UserMsg;
typedef struct _Block9Data Block9Data;

#define TYPE_STRANGERS_DIALG (strangers_dialg_get_type ())
#define STRANGERS_DIALG(obj) (G_TYPE_CHECK_INSTANCE_CAST ((obj), TYPE_STRANGERS_DIALG, StrangersDialg))
#define STRANGERS_DIALG_CLASS(klass) (G_TYPE_CHECK_CLASS_CAST ((klass), TYPE_STRANGERS_DIALG, StrangersDialgClass))
#define IS_STRANGERS_DIALG(obj) (G_TYPE_CHECK_INSTANCE_TYPE ((obj), TYPE_STRANGERS_DIALG))
#define IS_STRANGERS_DIALG_CLASS(klass) (G_TYPE_CHECK_CLASS_TYPE ((klass), TYPE_STRANGERS_DIALG))
#define STRANGERS_DIALG_GET_CLASS(obj) (G_TYPE_INSTANCE_GET_CLASS ((obj), TYPE_STRANGERS_DIALG, StrangersDialgClass))

typedef struct _StrangersDialg StrangersDialg;
typedef struct _StrangersDialgClass StrangersDialgClass;
typedef struct _Block10Data Block10Data;

#define TYPE_MY_GRID (my_grid_get_type ())
#define MY_GRID(obj) (G_TYPE_CHECK_INSTANCE_CAST ((obj), TYPE_MY_GRID, MyGrid))
#define MY_GRID_CLASS(klass) (G_TYPE_CHECK_CLASS_CAST ((klass), TYPE_MY_GRID, MyGridClass))
#define IS_MY_GRID(obj) (G_TYPE_CHECK_INSTANCE_TYPE ((obj), TYPE_MY_GRID))
#define IS_MY_GRID_CLASS(klass) (G_TYPE_CHECK_CLASS_TYPE ((klass), TYPE_MY_GRID))
#define MY_GRID_GET_CLASS(obj) (G_TYPE_INSTANCE_GET_CLASS ((obj), TYPE_MY_GRID, MyGridClass))

typedef struct _MyGrid MyGrid;
typedef struct _MyGridClass MyGridClass;
typedef struct _Block11Data Block11Data;

#define TYPE_CHAT_CLIENT (chat_client_get_type ())
#define CHAT_CLIENT(obj) (G_TYPE_CHECK_INSTANCE_CAST ((obj), TYPE_CHAT_CLIENT, ChatClient))
#define CHAT_CLIENT_CLASS(klass) (G_TYPE_CHECK_CLASS_CAST ((klass), TYPE_CHAT_CLIENT, ChatClientClass))
#define IS_CHAT_CLIENT(obj) (G_TYPE_CHECK_INSTANCE_TYPE ((obj), TYPE_CHAT_CLIENT))
#define IS_CHAT_CLIENT_CLASS(klass) (G_TYPE_CHECK_CLASS_TYPE ((klass), TYPE_CHAT_CLIENT))
#define CHAT_CLIENT_GET_CLASS(obj) (G_TYPE_INSTANCE_GET_CLASS ((obj), TYPE_CHAT_CLIENT, ChatClientClass))

typedef struct _ChatClient ChatClient;
typedef struct _ChatClientClass ChatClientClass;
typedef struct _ChatClientPrivate ChatClientPrivate;
enum  {
	CHAT_CLIENT_0_PROPERTY,
	CHAT_CLIENT_NUM_PROPERTIES
};
static GParamSpec* chat_client_properties[CHAT_CLIENT_NUM_PROPERTIES];
typedef struct _Block12Data Block12Data;
#define _g_object_unref0(var) ((var == NULL) ? NULL : (var = (g_object_unref (var), NULL)))
#define _g_thread_unref0(var) ((var == NULL) ? NULL : (var = (g_thread_unref (var), NULL)))
typedef struct _Block13Data Block13Data;

struct _NUParam {
	gchar* Name;
	gint Sex;
	gint Birth;
	gchar* Desc;
	gchar* Pwd;
};

struct _UserData {
	gint64 Id;
	gchar* Name;
	gint Sex;
	gint Age;
	gchar* Desc;
};

struct _Block8Data {
	int _ref_count_;
	gint64 id;
	gchar* name;
	gint sex;
	gint age;
	gchar* desc;
	gchar* msg;
	gchar* msg_time;
};

struct _UserMsg {
	gint64 id;
	gint sex;
	gchar* name;
	gchar* desc;
	gint age;
	gchar* msg_offline;
	gchar* timestamp_offline;
};

struct _Block9Data {
	int _ref_count_;
	gint64 id;
	gchar* name;
	gint sex;
	gint age;
	gchar* desc;
	gchar* msg;
	gchar* msg_time;
};

struct _Block10Data {
	int _ref_count_;
	gint64 id;
	gchar* name;
	gint sex;
	gint age;
	gchar* desc;
	gchar* msg;
	gchar* msg_time;
};

struct _Block11Data {
	int _ref_count_;
	gint8 typ;
	gint64 from;
	gchar* msg;
};

struct _ChatClient {
	GObject parent_instance;
	ChatClientPrivate * priv;
};

struct _ChatClientClass {
	GObjectClass parent_class;
};

struct _Block12Data {
	int _ref_count_;
	ChatClient* self;
	gchar* key;
};

struct _Block13Data {
	int _ref_count_;
	ChatClient* self;
	gint64 uid;
	gchar* msg;
};


extern SearchDialg* search1;
extern StrangersDialg* strangers1;
extern MyGrid* grid1;
static gpointer chat_client_parent_class = NULL;

GType nu_param_get_type (void) G_GNUC_CONST;
NUParam* nu_param_dup (const NUParam* self);
void nu_param_free (NUParam* self);
void nu_param_copy (const NUParam* self,
                    NUParam* dest);
void nu_param_destroy (NUParam* self);
GType user_data_get_type (void) G_GNUC_CONST;
UserData* user_data_dup (const UserData* self);
void user_data_free (UserData* self);
void user_data_copy (const UserData* self,
                     UserData* dest);
void user_data_destroy (UserData* self);
void search_add (gint64 id,
                 const gchar* name,
                 gint sex,
                 gint age,
                 const gchar* desc,
                 const gchar* msg,
                 const gchar* msg_time);
static Block8Data* block8_data_ref (Block8Data* _data8_);
static void block8_data_unref (void * _userdata_);
static gboolean __lambda28_ (Block8Data* _data8_);
GType search_dialg_get_type (void) G_GNUC_CONST;
GType user_msg_get_type (void) G_GNUC_CONST;
UserMsg* user_msg_dup (const UserMsg* self);
void user_msg_free (UserMsg* self);
void user_msg_copy (const UserMsg* self,
                    UserMsg* dest);
void user_msg_destroy (UserMsg* self);
void search_dialg_seach_callback (SearchDialg* self,
                                  UserMsg* u);
static gboolean ___lambda28__gsource_func (gpointer self);
void stranger_add (gint64 id,
                   const gchar* name,
                   gint sex,
                   gint age,
                   const gchar* desc,
                   const gchar* msg,
                   const gchar* msg_time);
static Block9Data* block9_data_ref (Block9Data* _data9_);
static void block9_data_unref (void * _userdata_);
static gboolean __lambda13_ (Block9Data* _data9_);
GType strangers_dialg_get_type (void) G_GNUC_CONST;
void strangers_dialg_prepend_row (StrangersDialg* self,
                                  UserMsg* u1);
static gboolean ___lambda13__gsource_func (gpointer self);
void friend_add (gint64 id,
                 const gchar* name,
                 gint sex,
                 gint age,
                 const gchar* desc,
                 const gchar* msg,
                 const gchar* msg_time);
static Block10Data* block10_data_ref (Block10Data* _data10_);
static void block10_data_unref (void * _userdata_);
static gboolean __lambda32_ (Block10Data* _data10_);
GType my_grid_get_type (void) G_GNUC_CONST;
void my_grid_add_friend (MyGrid* self,
                         UserMsg* user1,
                         gboolean tell);
static gboolean ___lambda32__gsource_func (gpointer self);
void client_notify (gint8 typ,
                    gint64 from,
                    gint64 to,
                    const gchar* msg);
static Block11Data* block11_data_ref (Block11Data* _data11_);
static void block11_data_unref (void * _userdata_);
static gboolean __lambda9_ (Block11Data* _data11_);
void my_grid_rpc_callback (MyGrid* self,
                           gint8 typ,
                           gint64 from,
                           const gchar* msg);
static gboolean ___lambda9__gsource_func (gpointer self);
GType chat_client_get_type (void) G_GNUC_CONST;
UserData* chat_client_login (ChatClient* self,
                             const gchar* name,
                             const gchar* pwd);
void chat_client_ping (ChatClient* self);
gint chat_client_update_pwd (ChatClient* self,
                             const gchar* name,
                             const gchar* pwd1,
                             const gchar* pwd2);
void chat_client_add_friend (ChatClient* self,
                             gint64 uid);
void chat_client_search_person_async (ChatClient* self,
                                      const gchar* key);
static Block12Data* block12_data_ref (Block12Data* _data12_);
static void block12_data_unref (void * _userdata_);
static gint __lambda27_ (Block12Data* _data12_);
static gpointer ___lambda27__gthread_func (gpointer self);
void chat_client_move_stranger_to_friend (ChatClient* self,
                                          gint64 fid);
void chat_client_get_stranger_msgs_async (ChatClient* self);
static gint __lambda14_ (ChatClient* self);
static gpointer ___lambda14__gthread_func (gpointer self);
void chat_client_get_friends_async (ChatClient* self);
static gint __lambda31_ (ChatClient* self);
static gpointer ___lambda31__gthread_func (gpointer self);
void chat_client_ChatTo (ChatClient* self,
                         gint64 to,
                         const gchar* msg);
void chat_client_tell (ChatClient* self,
                       gint64 to);
void chat_client_multi_send (ChatClient* self,
                             gint64* to,
                             int to_length1,
                             const gchar* msg);
void chat_client_tell_all (ChatClient* self,
                           gint64* to,
                           int to_length1);
void chat_client_set_http_id (ChatClient* self,
                              gint64 uid);
void chat_client_set_tcp_id (ChatClient* self,
                             gint64 uid);
void chat_client_send_file (ChatClient* self,
                            gint64 to,
                            const gchar* pathname);
gint chat_client_add_user (ChatClient* self,
                           const gchar* name,
                           const gchar* pwd,
                           gint sex,
                           gint birthyear,
                           const gchar* desc);
void chat_client_quit (ChatClient* self);
gint chat_client_get_proxy (ChatClient* self);
gchar* chat_client_get_host (ChatClient* self);
gint chat_client_set_proxy (ChatClient* self,
                            gint port);
void chat_client_remove_friend (ChatClient* self,
                                gint64 fid);
void chat_client_open_path (ChatClient* self,
                            const gchar* path1);
void chat_client_update_desc (ChatClient* self,
                              const gchar* desc);
void chat_client_delete_me (ChatClient* self,
                            const gchar* name,
                            const gchar* pwd);
gint chat_client_user_status (ChatClient* self,
                              gint64 uid);
void chat_client_offline_msg_with_id (ChatClient* self,
                                      gint64 uid,
                                      const gchar* msg);
static Block13Data* block13_data_ref (Block13Data* _data13_);
static void block13_data_unref (void * _userdata_);
static gint __lambda12_ (Block13Data* _data13_);
static gpointer ___lambda12__gthread_func (gpointer self);
gchar* chat_client_get_pg_path (ChatClient* self);
ChatClient* chat_client_new (void);
ChatClient* chat_client_construct (GType object_type);


void
nu_param_copy (const NUParam* self,
               NUParam* dest)
{
	const gchar* _tmp0_;
	gchar* _tmp1_;
	gint _tmp2_;
	gint _tmp3_;
	const gchar* _tmp4_;
	gchar* _tmp5_;
	const gchar* _tmp6_;
	gchar* _tmp7_;
	_tmp0_ = (*self).Name;
	_tmp1_ = g_strdup (_tmp0_);
	_g_free0 ((*dest).Name);
	(*dest).Name = _tmp1_;
	_tmp2_ = (*self).Sex;
	(*dest).Sex = _tmp2_;
	_tmp3_ = (*self).Birth;
	(*dest).Birth = _tmp3_;
	_tmp4_ = (*self).Desc;
	_tmp5_ = g_strdup (_tmp4_);
	_g_free0 ((*dest).Desc);
	(*dest).Desc = _tmp5_;
	_tmp6_ = (*self).Pwd;
	_tmp7_ = g_strdup (_tmp6_);
	_g_free0 ((*dest).Pwd);
	(*dest).Pwd = _tmp7_;
}


void
nu_param_destroy (NUParam* self)
{
	_g_free0 ((*self).Name);
	_g_free0 ((*self).Desc);
	_g_free0 ((*self).Pwd);
}


NUParam*
nu_param_dup (const NUParam* self)
{
	NUParam* dup;
	dup = g_new0 (NUParam, 1);
	nu_param_copy (self, dup);
	return dup;
}


void
nu_param_free (NUParam* self)
{
	nu_param_destroy (self);
	g_free (self);
}


GType
nu_param_get_type (void)
{
	static volatile gsize nu_param_type_id__volatile = 0;
	if (g_once_init_enter (&nu_param_type_id__volatile)) {
		GType nu_param_type_id;
		nu_param_type_id = g_boxed_type_register_static ("NUParam", (GBoxedCopyFunc) nu_param_dup, (GBoxedFreeFunc) nu_param_free);
		g_once_init_leave (&nu_param_type_id__volatile, nu_param_type_id);
	}
	return nu_param_type_id__volatile;
}


void
user_data_copy (const UserData* self,
                UserData* dest)
{
	gint64 _tmp0_;
	const gchar* _tmp1_;
	gchar* _tmp2_;
	gint _tmp3_;
	gint _tmp4_;
	const gchar* _tmp5_;
	gchar* _tmp6_;
	_tmp0_ = (*self).Id;
	(*dest).Id = _tmp0_;
	_tmp1_ = (*self).Name;
	_tmp2_ = g_strdup (_tmp1_);
	_g_free0 ((*dest).Name);
	(*dest).Name = _tmp2_;
	_tmp3_ = (*self).Sex;
	(*dest).Sex = _tmp3_;
	_tmp4_ = (*self).Age;
	(*dest).Age = _tmp4_;
	_tmp5_ = (*self).Desc;
	_tmp6_ = g_strdup (_tmp5_);
	_g_free0 ((*dest).Desc);
	(*dest).Desc = _tmp6_;
}


void
user_data_destroy (UserData* self)
{
	_g_free0 ((*self).Name);
	_g_free0 ((*self).Desc);
}


UserData*
user_data_dup (const UserData* self)
{
	UserData* dup;
	dup = g_new0 (UserData, 1);
	user_data_copy (self, dup);
	return dup;
}


void
user_data_free (UserData* self)
{
	user_data_destroy (self);
	g_free (self);
}


GType
user_data_get_type (void)
{
	static volatile gsize user_data_type_id__volatile = 0;
	if (g_once_init_enter (&user_data_type_id__volatile)) {
		GType user_data_type_id;
		user_data_type_id = g_boxed_type_register_static ("UserData", (GBoxedCopyFunc) user_data_dup, (GBoxedFreeFunc) user_data_free);
		g_once_init_leave (&user_data_type_id__volatile, user_data_type_id);
	}
	return user_data_type_id__volatile;
}


static Block8Data*
block8_data_ref (Block8Data* _data8_)
{
	g_atomic_int_inc (&_data8_->_ref_count_);
	return _data8_;
}


static void
block8_data_unref (void * _userdata_)
{
	Block8Data* _data8_;
	_data8_ = (Block8Data*) _userdata_;
	if (g_atomic_int_dec_and_test (&_data8_->_ref_count_)) {
		_g_free0 (_data8_->name);
		_g_free0 (_data8_->desc);
		_g_free0 (_data8_->msg);
		_g_free0 (_data8_->msg_time);
		g_slice_free (Block8Data, _data8_);
	}
}


static gboolean
__lambda28_ (Block8Data* _data8_)
{
	gboolean result = FALSE;
	SearchDialg* _tmp0_;
	UserMsg _tmp1_ = {0};
	_tmp0_ = search1;
	_tmp1_.id = _data8_->id;
	_tmp1_.sex = _data8_->sex;
	_g_free0 (_tmp1_.name);
	_tmp1_.name = _data8_->name;
	_g_free0 (_tmp1_.desc);
	_tmp1_.desc = _data8_->desc;
	_tmp1_.age = _data8_->age;
	_g_free0 (_tmp1_.msg_offline);
	_tmp1_.msg_offline = _data8_->msg;
	_g_free0 (_tmp1_.timestamp_offline);
	_tmp1_.timestamp_offline = _data8_->msg_time;
	search_dialg_seach_callback (_tmp0_, &_tmp1_);
	result = FALSE;
	return result;
}


static gboolean
___lambda28__gsource_func (gpointer self)
{
	gboolean result;
	result = __lambda28_ (self);
	return result;
}


void
search_add (gint64 id,
            const gchar* name,
            gint sex,
            gint age,
            const gchar* desc,
            const gchar* msg,
            const gchar* msg_time)
{
	Block8Data* _data8_;
	gchar* _tmp0_;
	gchar* _tmp1_;
	gchar* _tmp2_;
	gchar* _tmp3_;
	g_return_if_fail (name != NULL);
	g_return_if_fail (desc != NULL);
	g_return_if_fail (msg != NULL);
	g_return_if_fail (msg_time != NULL);
	_data8_ = g_slice_new0 (Block8Data);
	_data8_->_ref_count_ = 1;
	_data8_->id = id;
	_tmp0_ = g_strdup (name);
	_g_free0 (_data8_->name);
	_data8_->name = _tmp0_;
	_data8_->sex = sex;
	_data8_->age = age;
	_tmp1_ = g_strdup (desc);
	_g_free0 (_data8_->desc);
	_data8_->desc = _tmp1_;
	_tmp2_ = g_strdup (msg);
	_g_free0 (_data8_->msg);
	_data8_->msg = _tmp2_;
	_tmp3_ = g_strdup (msg_time);
	_g_free0 (_data8_->msg_time);
	_data8_->msg_time = _tmp3_;
	g_idle_add_full (G_PRIORITY_DEFAULT_IDLE, ___lambda28__gsource_func, block8_data_ref (_data8_), block8_data_unref);
	block8_data_unref (_data8_);
	_data8_ = NULL;
}


static Block9Data*
block9_data_ref (Block9Data* _data9_)
{
	g_atomic_int_inc (&_data9_->_ref_count_);
	return _data9_;
}


static void
block9_data_unref (void * _userdata_)
{
	Block9Data* _data9_;
	_data9_ = (Block9Data*) _userdata_;
	if (g_atomic_int_dec_and_test (&_data9_->_ref_count_)) {
		_g_free0 (_data9_->name);
		_g_free0 (_data9_->desc);
		_g_free0 (_data9_->msg);
		_g_free0 (_data9_->msg_time);
		g_slice_free (Block9Data, _data9_);
	}
}


static gboolean
__lambda13_ (Block9Data* _data9_)
{
	gboolean result = FALSE;
	StrangersDialg* _tmp0_;
	UserMsg _tmp1_ = {0};
	_tmp0_ = strangers1;
	_tmp1_.id = _data9_->id;
	_tmp1_.sex = _data9_->sex;
	_g_free0 (_tmp1_.name);
	_tmp1_.name = _data9_->name;
	_g_free0 (_tmp1_.desc);
	_tmp1_.desc = _data9_->desc;
	_tmp1_.age = _data9_->age;
	_g_free0 (_tmp1_.msg_offline);
	_tmp1_.msg_offline = _data9_->msg;
	_g_free0 (_tmp1_.timestamp_offline);
	_tmp1_.timestamp_offline = _data9_->msg_time;
	strangers_dialg_prepend_row (_tmp0_, &_tmp1_);
	result = FALSE;
	return result;
}


static gboolean
___lambda13__gsource_func (gpointer self)
{
	gboolean result;
	result = __lambda13_ (self);
	return result;
}


void
stranger_add (gint64 id,
              const gchar* name,
              gint sex,
              gint age,
              const gchar* desc,
              const gchar* msg,
              const gchar* msg_time)
{
	Block9Data* _data9_;
	gchar* _tmp0_;
	gchar* _tmp1_;
	gchar* _tmp2_;
	gchar* _tmp3_;
	g_return_if_fail (name != NULL);
	g_return_if_fail (desc != NULL);
	g_return_if_fail (msg != NULL);
	g_return_if_fail (msg_time != NULL);
	_data9_ = g_slice_new0 (Block9Data);
	_data9_->_ref_count_ = 1;
	_data9_->id = id;
	_tmp0_ = g_strdup (name);
	_g_free0 (_data9_->name);
	_data9_->name = _tmp0_;
	_data9_->sex = sex;
	_data9_->age = age;
	_tmp1_ = g_strdup (desc);
	_g_free0 (_data9_->desc);
	_data9_->desc = _tmp1_;
	_tmp2_ = g_strdup (msg);
	_g_free0 (_data9_->msg);
	_data9_->msg = _tmp2_;
	_tmp3_ = g_strdup (msg_time);
	_g_free0 (_data9_->msg_time);
	_data9_->msg_time = _tmp3_;
	g_idle_add_full (G_PRIORITY_DEFAULT_IDLE, ___lambda13__gsource_func, block9_data_ref (_data9_), block9_data_unref);
	block9_data_unref (_data9_);
	_data9_ = NULL;
}


static Block10Data*
block10_data_ref (Block10Data* _data10_)
{
	g_atomic_int_inc (&_data10_->_ref_count_);
	return _data10_;
}


static void
block10_data_unref (void * _userdata_)
{
	Block10Data* _data10_;
	_data10_ = (Block10Data*) _userdata_;
	if (g_atomic_int_dec_and_test (&_data10_->_ref_count_)) {
		_g_free0 (_data10_->name);
		_g_free0 (_data10_->desc);
		_g_free0 (_data10_->msg);
		_g_free0 (_data10_->msg_time);
		g_slice_free (Block10Data, _data10_);
	}
}


static gboolean
__lambda32_ (Block10Data* _data10_)
{
	gboolean result = FALSE;
	MyGrid* _tmp0_;
	UserMsg _tmp1_ = {0};
	_tmp0_ = grid1;
	_tmp1_.id = _data10_->id;
	_tmp1_.sex = _data10_->sex;
	_g_free0 (_tmp1_.name);
	_tmp1_.name = _data10_->name;
	_g_free0 (_tmp1_.desc);
	_tmp1_.desc = _data10_->desc;
	_tmp1_.age = _data10_->age;
	_g_free0 (_tmp1_.msg_offline);
	_tmp1_.msg_offline = _data10_->msg;
	_g_free0 (_tmp1_.timestamp_offline);
	_tmp1_.timestamp_offline = _data10_->msg_time;
	my_grid_add_friend (_tmp0_, &_tmp1_, TRUE);
	result = FALSE;
	return result;
}


static gboolean
___lambda32__gsource_func (gpointer self)
{
	gboolean result;
	result = __lambda32_ (self);
	return result;
}


void
friend_add (gint64 id,
            const gchar* name,
            gint sex,
            gint age,
            const gchar* desc,
            const gchar* msg,
            const gchar* msg_time)
{
	Block10Data* _data10_;
	gchar* _tmp0_;
	gchar* _tmp1_;
	gchar* _tmp2_;
	gchar* _tmp3_;
	g_return_if_fail (name != NULL);
	g_return_if_fail (desc != NULL);
	g_return_if_fail (msg != NULL);
	g_return_if_fail (msg_time != NULL);
	_data10_ = g_slice_new0 (Block10Data);
	_data10_->_ref_count_ = 1;
	_data10_->id = id;
	_tmp0_ = g_strdup (name);
	_g_free0 (_data10_->name);
	_data10_->name = _tmp0_;
	_data10_->sex = sex;
	_data10_->age = age;
	_tmp1_ = g_strdup (desc);
	_g_free0 (_data10_->desc);
	_data10_->desc = _tmp1_;
	_tmp2_ = g_strdup (msg);
	_g_free0 (_data10_->msg);
	_data10_->msg = _tmp2_;
	_tmp3_ = g_strdup (msg_time);
	_g_free0 (_data10_->msg_time);
	_data10_->msg_time = _tmp3_;
	g_idle_add_full (G_PRIORITY_DEFAULT_IDLE, ___lambda32__gsource_func, block10_data_ref (_data10_), block10_data_unref);
	block10_data_unref (_data10_);
	_data10_ = NULL;
}


static Block11Data*
block11_data_ref (Block11Data* _data11_)
{
	g_atomic_int_inc (&_data11_->_ref_count_);
	return _data11_;
}


static void
block11_data_unref (void * _userdata_)
{
	Block11Data* _data11_;
	_data11_ = (Block11Data*) _userdata_;
	if (g_atomic_int_dec_and_test (&_data11_->_ref_count_)) {
		_g_free0 (_data11_->msg);
		g_slice_free (Block11Data, _data11_);
	}
}


static gboolean
__lambda9_ (Block11Data* _data11_)
{
	gboolean result = FALSE;
	MyGrid* _tmp0_;
	_tmp0_ = grid1;
	my_grid_rpc_callback (_tmp0_, _data11_->typ, _data11_->from, _data11_->msg);
	result = FALSE;
	return result;
}


static gboolean
___lambda9__gsource_func (gpointer self)
{
	gboolean result;
	result = __lambda9_ (self);
	return result;
}


void
client_notify (gint8 typ,
               gint64 from,
               gint64 to,
               const gchar* msg)
{
	Block11Data* _data11_;
	gchar* _tmp0_;
	g_return_if_fail (msg != NULL);
	_data11_ = g_slice_new0 (Block11Data);
	_data11_->_ref_count_ = 1;
	_data11_->typ = typ;
	_data11_->from = from;
	_tmp0_ = g_strdup (msg);
	_g_free0 (_data11_->msg);
	_data11_->msg = _tmp0_;
	g_idle_add_full (G_PRIORITY_DEFAULT_IDLE, ___lambda9__gsource_func, block11_data_ref (_data11_), block11_data_unref);
	block11_data_unref (_data11_);
	_data11_ = NULL;
}


static gpointer
_user_data_dup0 (gpointer self)
{
	return self ? user_data_dup (self) : NULL;
}


UserData*
chat_client_login (ChatClient* self,
                   const gchar* name,
                   const gchar* pwd)
{
	UserData* result = NULL;
	UserData u = {0};
	gchar* _tmp0_;
	gchar* _tmp1_;
	UserData _tmp2_ = {0};
	gint ret = 0;
	gint _tmp3_;
	gint _tmp4_;
	g_return_val_if_fail (self != NULL, NULL);
	g_return_val_if_fail (name != NULL, NULL);
	g_return_val_if_fail (pwd != NULL, NULL);
	_tmp0_ = g_strdup ("");
	_tmp1_ = g_strdup ("");
	_tmp2_.Id = (gint64) 0;
	_g_free0 (_tmp2_.Name);
	_tmp2_.Name = _tmp0_;
	_tmp2_.Sex = 0;
	_tmp2_.Age = 0;
	_g_free0 (_tmp2_.Desc);
	_tmp2_.Desc = _tmp1_;
	u = _tmp2_;
	_tmp3_ = Client_Login (name, pwd, &u);
	ret = _tmp3_;
	_tmp4_ = ret;
	if (_tmp4_ == 1) {
		UserData _tmp5_;
		UserData* _tmp6_;
		UserData* _tmp7_;
		Client_SetNotifyFn ((void*) client_notify);
		_tmp5_ = u;
		_tmp6_ = _user_data_dup0 (&_tmp5_);
		_tmp7_ = _tmp6_;
		user_data_destroy (&_tmp5_);
		result = _tmp7_;
		return result;
	} else {
		result = NULL;
		user_data_destroy (&u);
		return result;
	}
	user_data_destroy (&u);
}


void
chat_client_ping (ChatClient* self)
{
	g_return_if_fail (self != NULL);
	Client_Ping ();
}


gint
chat_client_update_pwd (ChatClient* self,
                        const gchar* name,
                        const gchar* pwd1,
                        const gchar* pwd2)
{
	gint result = 0;
	g_return_val_if_fail (self != NULL, 0);
	g_return_val_if_fail (name != NULL, 0);
	g_return_val_if_fail (pwd1 != NULL, 0);
	g_return_val_if_fail (pwd2 != NULL, 0);
	result = Client_NewPasswd (name, pwd1, pwd2);
	return result;
}


void
chat_client_add_friend (ChatClient* self,
                        gint64 uid)
{
	g_return_if_fail (self != NULL);
	Client_AddFriend (uid);
}


static Block12Data*
block12_data_ref (Block12Data* _data12_)
{
	g_atomic_int_inc (&_data12_->_ref_count_);
	return _data12_;
}


static void
block12_data_unref (void * _userdata_)
{
	Block12Data* _data12_;
	_data12_ = (Block12Data*) _userdata_;
	if (g_atomic_int_dec_and_test (&_data12_->_ref_count_)) {
		ChatClient* self;
		self = _data12_->self;
		_g_free0 (_data12_->key);
		_g_object_unref0 (self);
		g_slice_free (Block12Data, _data12_);
	}
}


static gint
__lambda27_ (Block12Data* _data12_)
{
	ChatClient* self;
	gint result = 0;
	self = _data12_->self;
	Client_SearchPersons (_data12_->key, (void*) search_add);
	result = 0;
	return result;
}


static gpointer
___lambda27__gthread_func (gpointer self)
{
	gpointer result;
	result = (gpointer) ((gintptr) __lambda27_ (self));
	block12_data_unref (self);
	return result;
}


void
chat_client_search_person_async (ChatClient* self,
                                 const gchar* key)
{
	Block12Data* _data12_;
	gchar* _tmp0_;
	GThread* thread = NULL;
	GThread* _tmp1_;
	g_return_if_fail (self != NULL);
	g_return_if_fail (key != NULL);
	_data12_ = g_slice_new0 (Block12Data);
	_data12_->_ref_count_ = 1;
	_data12_->self = g_object_ref (self);
	_tmp0_ = g_strdup (key);
	_g_free0 (_data12_->key);
	_data12_->key = _tmp0_;
	_tmp1_ = g_thread_new ("search_person", ___lambda27__gthread_func, block12_data_ref (_data12_));
	thread = _tmp1_;
	_g_thread_unref0 (thread);
	block12_data_unref (_data12_);
	_data12_ = NULL;
}


void
chat_client_move_stranger_to_friend (ChatClient* self,
                                     gint64 fid)
{
	g_return_if_fail (self != NULL);
	Client_MoveStrangerToFriend (fid);
}


static gint
__lambda14_ (ChatClient* self)
{
	gint result = 0;
	Client_GetStrangerMsgs ((void*) stranger_add);
	result = 0;
	return result;
}


static gpointer
___lambda14__gthread_func (gpointer self)
{
	gpointer result;
	result = (gpointer) ((gintptr) __lambda14_ ((ChatClient*) self));
	g_object_unref (self);
	return result;
}


void
chat_client_get_stranger_msgs_async (ChatClient* self)
{
	GThread* thread = NULL;
	GThread* _tmp0_;
	g_return_if_fail (self != NULL);
	_tmp0_ = g_thread_new ("get_strangers", ___lambda14__gthread_func, g_object_ref (self));
	thread = _tmp0_;
	_g_thread_unref0 (thread);
}


static gint
__lambda31_ (ChatClient* self)
{
	gint result = 0;
	Client_GetFriends ((void*) friend_add);
	result = 0;
	return result;
}


static gpointer
___lambda31__gthread_func (gpointer self)
{
	gpointer result;
	result = (gpointer) ((gintptr) __lambda31_ ((ChatClient*) self));
	g_object_unref (self);
	return result;
}


void
chat_client_get_friends_async (ChatClient* self)
{
	GThread* thread = NULL;
	GThread* _tmp0_;
	g_return_if_fail (self != NULL);
	_tmp0_ = g_thread_new ("get_friends", ___lambda31__gthread_func, g_object_ref (self));
	thread = _tmp0_;
	_g_thread_unref0 (thread);
}


void
chat_client_ChatTo (ChatClient* self,
                    gint64 to,
                    const gchar* msg)
{
	g_return_if_fail (self != NULL);
	g_return_if_fail (msg != NULL);
	Client_ChatTo (to, msg);
}


void
chat_client_tell (ChatClient* self,
                  gint64 to)
{
	g_return_if_fail (self != NULL);
	Client_Tell (to);
}


void
chat_client_multi_send (ChatClient* self,
                        gint64* to,
                        int to_length1,
                        const gchar* msg)
{
	g_return_if_fail (self != NULL);
	g_return_if_fail (msg != NULL);
	Client_MultiSend (msg, to, to_length1);
}


void
chat_client_tell_all (ChatClient* self,
                      gint64* to,
                      int to_length1)
{
	g_return_if_fail (self != NULL);
	Client_TellAll (to, to_length1);
}


void
chat_client_set_http_id (ChatClient* self,
                         gint64 uid)
{
	g_return_if_fail (self != NULL);
	Client_SetHttpId (uid);
}


void
chat_client_set_tcp_id (ChatClient* self,
                        gint64 uid)
{
	g_return_if_fail (self != NULL);
	Client_SetHttpId (uid);
}


void
chat_client_send_file (ChatClient* self,
                       gint64 to,
                       const gchar* pathname)
{
	g_return_if_fail (self != NULL);
	g_return_if_fail (pathname != NULL);
	Client_SendFile (to, pathname);
}


gint
chat_client_add_user (ChatClient* self,
                      const gchar* name,
                      const gchar* pwd,
                      gint sex,
                      gint birthyear,
                      const gchar* desc)
{
	gint result = 0;
	NUParam _tmp0_ = {0};
	g_return_val_if_fail (self != NULL, 0);
	g_return_val_if_fail (name != NULL, 0);
	g_return_val_if_fail (pwd != NULL, 0);
	g_return_val_if_fail (desc != NULL, 0);
	_g_free0 (_tmp0_.Name);
	_tmp0_.Name = name;
	_tmp0_.Sex = sex;
	_tmp0_.Birth = birthyear;
	_g_free0 (_tmp0_.Desc);
	_tmp0_.Desc = desc;
	_g_free0 (_tmp0_.Pwd);
	_tmp0_.Pwd = pwd;
	result = Client_NewUser (&_tmp0_);
	return result;
}


void
chat_client_quit (ChatClient* self)
{
	g_return_if_fail (self != NULL);
	Client_Quit ();
}


gint
chat_client_get_proxy (ChatClient* self)
{
	gint result = 0;
	g_return_val_if_fail (self != NULL, 0);
	result = Client_GetProxyPort ();
	return result;
}


gchar*
chat_client_get_host (ChatClient* self)
{
	gchar* result = NULL;
	gchar* p = NULL;
	gchar* _tmp0_ = NULL;
	gchar* _tmp1_;
	g_return_val_if_fail (self != NULL, NULL);
	Client_GetHost (&_tmp0_);
	_g_free0 (p);
	p = _tmp0_;
	_tmp1_ = g_strdup (p);
	result = _tmp1_;
	_g_free0 (p);
	return result;
}


gint
chat_client_set_proxy (ChatClient* self,
                       gint port)
{
	gint result = 0;
	g_return_val_if_fail (self != NULL, 0);
	result = Client_ProxyPort (port);
	return result;
}


void
chat_client_remove_friend (ChatClient* self,
                           gint64 fid)
{
	g_return_if_fail (self != NULL);
	if (fid == ((gint64) 0)) {
		return;
	}
	Client_RemoveFriend (fid);
}


void
chat_client_open_path (ChatClient* self,
                       const gchar* path1)
{
	g_return_if_fail (self != NULL);
	g_return_if_fail (path1 != NULL);
	Client_OpenPath (path1);
}


void
chat_client_update_desc (ChatClient* self,
                         const gchar* desc)
{
	g_return_if_fail (self != NULL);
	g_return_if_fail (desc != NULL);
	Client_UpdateDesc (desc);
}


void
chat_client_delete_me (ChatClient* self,
                       const gchar* name,
                       const gchar* pwd)
{
	g_return_if_fail (self != NULL);
	g_return_if_fail (name != NULL);
	g_return_if_fail (pwd != NULL);
	Client_DeleteMe (name, pwd);
}


gint
chat_client_user_status (ChatClient* self,
                         gint64 uid)
{
	gint result = 0;
	g_return_val_if_fail (self != NULL, 0);
	result = Client_UserStatus (uid);
	return result;
}


static Block13Data*
block13_data_ref (Block13Data* _data13_)
{
	g_atomic_int_inc (&_data13_->_ref_count_);
	return _data13_;
}


static void
block13_data_unref (void * _userdata_)
{
	Block13Data* _data13_;
	_data13_ = (Block13Data*) _userdata_;
	if (g_atomic_int_dec_and_test (&_data13_->_ref_count_)) {
		ChatClient* self;
		self = _data13_->self;
		_g_free0 (_data13_->msg);
		_g_object_unref0 (self);
		g_slice_free (Block13Data, _data13_);
	}
}


static gint
__lambda12_ (Block13Data* _data13_)
{
	ChatClient* self;
	gint result = 0;
	self = _data13_->self;
	Client_QueryID (_data13_->uid, _data13_->msg, (void*) stranger_add);
	result = 0;
	return result;
}


static gpointer
___lambda12__gthread_func (gpointer self)
{
	gpointer result;
	result = (gpointer) ((gintptr) __lambda12_ (self));
	block13_data_unref (self);
	return result;
}


void
chat_client_offline_msg_with_id (ChatClient* self,
                                 gint64 uid,
                                 const gchar* msg)
{
	Block13Data* _data13_;
	gchar* _tmp0_;
	GThread* thread = NULL;
	GThread* _tmp1_;
	g_return_if_fail (self != NULL);
	g_return_if_fail (msg != NULL);
	_data13_ = g_slice_new0 (Block13Data);
	_data13_->_ref_count_ = 1;
	_data13_->self = g_object_ref (self);
	_data13_->uid = uid;
	_tmp0_ = g_strdup (msg);
	_g_free0 (_data13_->msg);
	_data13_->msg = _tmp0_;
	if (_data13_->uid == ((gint64) 0)) {
		block13_data_unref (_data13_);
		_data13_ = NULL;
		return;
	}
	_tmp1_ = g_thread_new ("search_person", ___lambda12__gthread_func, block13_data_ref (_data13_));
	thread = _tmp1_;
	_g_thread_unref0 (thread);
	block13_data_unref (_data13_);
	_data13_ = NULL;
}


gchar*
chat_client_get_pg_path (ChatClient* self)
{
	gchar* result = NULL;
	gchar* p = NULL;
	gchar* _tmp0_ = NULL;
	gchar* _tmp1_;
	g_return_val_if_fail (self != NULL, NULL);
	Client_GetPgPath (&_tmp0_);
	_g_free0 (p);
	p = _tmp0_;
	_tmp1_ = g_strdup (p);
	result = _tmp1_;
	_g_free0 (p);
	return result;
}


ChatClient*
chat_client_construct (GType object_type)
{
	ChatClient * self = NULL;
	self = (ChatClient*) g_object_new (object_type, NULL);
	return self;
}


ChatClient*
chat_client_new (void)
{
	return chat_client_construct (TYPE_CHAT_CLIENT);
}


static void
chat_client_class_init (ChatClientClass * klass)
{
	chat_client_parent_class = g_type_class_peek_parent (klass);
}


static void
chat_client_instance_init (ChatClient * self)
{
}


GType
chat_client_get_type (void)
{
	static volatile gsize chat_client_type_id__volatile = 0;
	if (g_once_init_enter (&chat_client_type_id__volatile)) {
		static const GTypeInfo g_define_type_info = { sizeof (ChatClientClass), (GBaseInitFunc) NULL, (GBaseFinalizeFunc) NULL, (GClassInitFunc) chat_client_class_init, (GClassFinalizeFunc) NULL, NULL, sizeof (ChatClient), 0, (GInstanceInitFunc) chat_client_instance_init, NULL };
		GType chat_client_type_id;
		chat_client_type_id = g_type_register_static (G_TYPE_OBJECT, "ChatClient", &g_define_type_info, 0);
		g_once_init_leave (&chat_client_type_id__volatile, chat_client_type_id);
	}
	return chat_client_type_id__volatile;
}



