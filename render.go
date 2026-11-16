// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can be
// found in the LICENSE file.

package termui

import (
	"image"
	"sync"

	"github.com/gdamore/tcell/v2"
)

type Drawable interface {
	GetRect() image.Rectangle
	SetRect(int, int, int, int)
	Draw(*Buffer)
	sync.Locker
}

func Render(items ...Drawable) {
	for _, item := range items {
		buf := NewBuffer(item.GetRect())
		item.Lock()
		item.Draw(buf)
		item.Unlock()
		for point, cell := range buf.CellMap {
			if point.In(buf.Rectangle) {
				screen.SetContent(point.X, point.Y, cell.Rune, nil, cellStyle(cell.Style))
			}
		}
	}
	screen.Show()
}

func cellColor(c Color) tcell.Color {
	if c == ColorClear {
		return tcell.ColorDefault
	}
	return tcell.ColorValid | tcell.Color(c)
}

func cellStyle(s Style) tcell.Style {
	style := tcell.StyleDefault.Foreground(cellColor(s.Fg)).Background(cellColor(s.Bg))
	if s.Modifier&ModifierBold != 0 {
		style = style.Bold(true)
	}
	if s.Modifier&ModifierUnderline != 0 {
		style = style.Underline(true)
	}
	if s.Modifier&ModifierReverse != 0 {
		style = style.Reverse(true)
	}
	return style
}
