package main

import (
	"github.com/hajimehoshi/oto"
)

type Noticer struct {
	data   []byte
	player *oto.Player
}

func NewNoticer() (*Noticer, error) {
	player, err := oto.NewPlayer(22050, 1, 2, len(noticeData))
	if err != nil {
		return nil, err
	}
	s := &Noticer{data: noticeData, player: player}
	return s, nil
}

func (s *Noticer) Play() {
	s.player.Write(s.data)
}

func (s *Noticer) Close() {
	s.player.Close()
}
