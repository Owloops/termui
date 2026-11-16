// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can be
// found in the LICENSE file.

package termui

import (
	"github.com/gdamore/tcell/v2"
)

var screen tcell.Screen

// Init initializes the terminal and is required to render anything.
// After initialization, the library must be finalized with `Close`.
func Init() error {
	s, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := s.Init(); err != nil {
		return err
	}
	s.EnableMouse()
	screen = s
	return nil
}

// Close closes the terminal.
func Close() {
	if screen == nil {
		return
	}
	screen.Fini()
	screen = nil
}

func TerminalDimensions() (int, int) {
	screen.Sync()
	return screen.Size()
}

func Clear() {
	style := tcell.StyleDefault.Background(cellColor(Theme.Default.Bg))
	screen.SetStyle(style)
	screen.Clear()
	screen.Show()
}
