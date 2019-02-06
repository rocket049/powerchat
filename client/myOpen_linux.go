package main

import "github.com/skratchdot/open-golang/open"

func myOpen(path1 string) error {
	return open.Start(path1)
}
