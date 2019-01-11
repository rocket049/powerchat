package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	md "github.com/russross/blackfriday"
)

type myHandler struct {
	root string
}

func newMyHandler(pathRoot string) *myHandler {
	os.MkdirAll(pathRoot, os.ModePerm)
	return &myHandler{root: pathRoot}
}

func startMyHttpServe(pathRoot, addr string) {
	log.Println("www root:", pathRoot)
	handler := newMyHandler(pathRoot)
	http.ListenAndServe(addr, handler)
}

func (s *myHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if strings.ToLower(req.Method) != "get" {
		resp.WriteHeader(404)
		return
	}
	p1 := req.URL.Query().Get("m")
	if "ls" == p1 || req.URL.Path == "/" {
		s.lsPath(resp, req)
	} else if p1 == "md" {
		s.readMd(resp, req)
	} else if p1 == "ht" {
		s.readHtml(resp, req)
	} else if req.URL.Path == "/style.css" {
		s.getStyle(resp, req)
	} else {
		http.ServeFile(resp, req, filepath.Join(s.root, req.URL.Path))
	}
}

type LinkData struct {
	Name   string
	Path   string
	Action string
	MTime  time.Time
}

func (s *myHandler) lsPath(resp http.ResponseWriter, req *http.Request) {
	p1 := filepath.Join(s.root, req.URL.Path)
	info1, err := os.Stat(p1)
	if err != nil {
		resp.WriteHeader(404)
		return
	}
	linkList := []LinkData{}
	if info1.IsDir() {
		dir1, _ := os.Open(p1)
		defer dir1.Close()
		items, _ := dir1.Readdir(0)
		for _, item1 := range items {
			var link1 LinkData
			if item1.IsDir() {
				link1.Name = path.Base(item1.Name()) + "/"
				if strings.HasPrefix(link1.Name, "_") || strings.HasPrefix(link1.Name, ".") {
					continue
				}
				link1.MTime = item1.ModTime()
				link1.Path = path.Join(p1, link1.Name)
				link1.Action = "ls"
				linkList = append(linkList, link1)
			} else {
				link1.Name = path.Base(item1.Name())
				if strings.HasPrefix(link1.Name, "_") || strings.HasPrefix(link1.Name, ".") {
					continue
				}
				link1.MTime = item1.ModTime()
				link1.Path = path.Join(p1, link1.Name)
				ext1 := strings.ToLower(path.Ext(link1.Name))
				switch ext1 {
				case ".md":
					link1.Action = "md"
				case ".html":
					link1.Action = "ht"
				case ".htm":
					link1.Action = "ht"
				default:
					link1.Action = ""
				}
				if link1.Action != "" {
					linkList = append(linkList, link1)
				}
			}
		}
	} else {
		var link1 LinkData
		link1.Name = path.Base(info1.Name())
		link1.MTime = info1.ModTime()
		if strings.HasPrefix(link1.Name, "_") == false {
			link1.Path = path.Join(p1, link1.Name)
			ext1 := strings.ToLower(path.Ext(link1.Name))
			switch ext1 {
			case ".md":
				link1.Action = "md"
			case ".html":
				link1.Action = "ht"
			case ".htm":
				link1.Action = "ht"
			default:
				link1.Action = ""
			}
			if link1.Action != "" {
				linkList = append(linkList, link1)
			}
		}
	}
	//sort
	sort.SliceStable(linkList, func(i, j int) bool {
		return linkList[i].MTime.Unix() > linkList[j].MTime.Unix()
	})

	data := make(map[string]interface{})
	data["title"] = req.URL.Path
	data["links"] = linkList
	tpl1 := template.New("")
	_, err = tpl1.Parse(tmpl)
	if err != nil {
		log.Println("tempate parse error:", err)
		resp.WriteHeader(404)
		return
	}
	resp.WriteHeader(200)
	tpl1.Execute(resp, data)
}

func (s *myHandler) readMd(resp http.ResponseWriter, req *http.Request) {
	p1 := filepath.Join(s.root, req.URL.Path)
	_, err := os.Stat(p1)
	if err != nil {
		resp.WriteHeader(404)
		return
	}
	header1 := strings.Replace(header, "{{.title}}", path.Base(p1), -1)
	p, err := ioutil.ReadFile(p1)
	if err != nil {
		resp.WriteHeader(404)
		return
	}
	r := md.HtmlRenderer(md.HTML_SKIP_HTML, "", "")
	opts := md.EXTENSION_FENCED_CODE | md.EXTENSION_TABLES
	body := md.Markdown(p, r, opts)
	resp.WriteHeader(200)
	resp.Write([]byte(header1))
	resp.Write(body)
	resp.Write([]byte(tail))
}

func (s *myHandler) readHtml(resp http.ResponseWriter, req *http.Request) {
	path1 := filepath.Join(s.root, req.URL.Path)
	_, err := os.Stat(path1)
	if err != nil {
		resp.WriteHeader(404)
		return
	}
	http.ServeFile(resp, req, path1)
}

func (s *myHandler) getStyle(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(200)
	resp.Write(style)
}
