package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func copyFile(src, dst string, mode os.FileMode) error {
	fp1, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer fp1.Close()
	fp2, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fp2.Close()
	_, err = io.Copy(fp1, fp2)
	return err
}

func makeLauncher() {
	appimage := os.Getenv("APPIMAGE")
	appdir := os.Getenv("APPDIR")
	if len(appimage) == 0 || len(appdir) == 0 {
		return
	}
	home, _ := os.UserHomeDir()
	dst := filepath.Join(home, ".local", "share", "applications", "powerchat.desktop")
	iconSrc := filepath.Join(appdir, "usr/share/icons/powerchat/tank.png")
	iconDir := filepath.Join(home, ".local", "share", "icons", "powerchat")
	os.MkdirAll(iconDir, os.ModePerm)
	iconDst := filepath.Join(iconDir, "powerchat.png")
	copyFile(iconSrc, iconDst, 0644)

	data := struct {
		Name string
		Icon string
	}{appimage, iconDst}

	tpl := `[Desktop Entry]
Name=PowerChat
Comment=powerchat
Exec="{{.Name}}" %U
Icon={{.Icon}}
Terminal=false
Type=Application
StartupNotify=true
Categories=Network;GTK;
	
`
	t := template.New("")
	t.Parse(tpl)
	fp, err := os.Create(dst)
	if err != nil {
		log.Println(err)
		return
	}
	defer fp.Close()
	t.Execute(fp, data)
}
