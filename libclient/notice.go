package main

import (
	"sync"

	"github.com/hajimehoshi/oto"
)

type Noticer struct {
	data   []byte
	player *oto.Player
	lock1  *sync.Mutex
}

func NewNoticer() (*Noticer, error) {
	player, err := oto.NewPlayer(22050, 1, 2, len(noticeData))
	if err != nil {
		return nil, err
	}
	s := &Noticer{data: noticeData, player: player, lock1: new(sync.Mutex)}
	return s, nil
}

func (s *Noticer) Play() {
	s.lock1.Lock()
	s.player.Write(s.data)
	s.lock1.Unlock()
}

func (s *Noticer) Close() {
	s.lock1.Lock()
	s.player.Close()
	s.lock1.Unlock()
}
