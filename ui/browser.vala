using Gtk;
using WebKit;

public class MyBrowser:GLib.Object{
    public Gtk.Window window;
    public WebKit.WebView view;
    public Gtk.Entry addr;
    bool created = false;
    public void create(){
        this.window = new Gtk.Window();
        this.window.title = "Web Browser";
        this.window.set_default_size(800,600);
        this.window.destroy.connect(()=>{
			this.view.try_close();
			this.window.close();
            this.created = false;
        });
        this.created = true;
        
        var grid  = new Gtk.Grid();
        
        var bt_back = new Gtk.Button.with_label("<-");
        var bt_forward = new Gtk.Button.with_label("->");
        grid.attach(bt_back,0,0);
        grid.attach(bt_forward,1,0);
        
        this.addr = new Gtk.Entry();
        this.addr.hexpand = true;
        this.addr.editable = false;
        grid.attach(this.addr,2,0);
        
        this.view = new WebKit.WebView();
        var scrollWin = new Gtk.ScrolledWindow(null,null);
        scrollWin.expand = true;
        scrollWin.add(this.view);
        grid.attach(scrollWin,0,1,3,1);
        
        this.window.add(grid);
        
        WebKit.Settings settings = new WebKit.Settings();
        settings.set_user_agent_with_application_details ("Chrome","63.0.3239.132");
        this.view.settings = settings;

        this.view.load_changed.connect((event)=>{
            switch(event){
                case WebKit.LoadEvent.FINISHED:
                    this.window.title = this.view.get_title();
                    this.addr.text = this.view.get_uri();
                    break;
                case WebKit.LoadEvent.STARTED:
                    this.window.title = "loading "+this.view.get_uri();
                    break;
                case WebKit.LoadEvent.REDIRECTED:
                    this.window.title = "redirected "+this.view.get_uri();
                    this.addr.text = this.view.get_uri();
                    break;
            }
        });
        this.view.create.connect(()=>{
            return this.view;
        });
        this.view.ready_to_show.connect(()=>{
            this.view.show_all();
        });
        
        bt_back.clicked.connect(()=>{
			this.view.go_back();
		});
		bt_forward.clicked.connect(()=>{
			this.view.go_forward();
		});
    }
    public void open(string uri){
        if(this.created==false)
            this.create();
        this.window.show_all();
        this.view.load_uri(uri);
    }
}
