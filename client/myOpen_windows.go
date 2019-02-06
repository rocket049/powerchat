package main

import (
	"os"

	"github.com/skratchdot/open-golang/open"
)

func myOpen(path1 string) error {
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
