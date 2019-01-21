package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/oto"
)

func getPcmPath() string {
	exe1, _ := os.Executable()
	dir1 := filepath.Dir(exe1)
	res := filepath.Join(dir1, "pcm", "notice.pcm")
	return res
}

type Noticer struct {
	data   []byte
	player *oto.Player
}

func NewNoticer() (*Noticer, error) {
	data, err := ioutil.ReadFile(getPcmPath())
	if err != nil {
		return nil, err
	}
	player, err := oto.NewPlayer(22050, 1, 2, len(data))
	if err != nil {
		return nil, err
	}
	s := &Noticer{data: data, player: player}
	return s, nil
}

func (s *Noticer) Play() {
	s.player.Write(s.data)
}

func (s *Noticer) Close() {
	s.player.Close()
}
