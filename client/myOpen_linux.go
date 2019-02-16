package main

import (
	"os"

	"github.com/skratchdot/open-golang/open"
)

var osID uint = 0

func myOpen(path1 string) error {
	switch osID {
	case 0:
		return open.Start(path1)
	case 1:
		return win32Start(path1)
	}
	return nil
}

func win32Start(path1 string) error {
	st, err := os.Stat(path1)
	if err != nil {
		return err
	}
	if st.IsDir() {
		err = open.StartWith(path1, "explorer.exe")
	} else {
		err = open.Start(path1)
	}
	return err
}
