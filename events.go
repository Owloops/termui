// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can be
// found in the LICENSE file.

package termui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

/*
List of events:
	mouse events:
		<MouseLeft> <MouseRight> <MouseMiddle>
		<MouseMotion> <MouseWheelUp> <MouseWheelDown>
	keyboard events:
		any uppercase or lowercase letter like j or J
		<C-d> etc
		<M-d> etc
		<Up> <Down> <Left> <Right>
		<Insert> <Delete> <Home> <End> <PageUp> <PageDown>
		<Backspace> <Tab> <Enter> <Escape> <Space>
	terminal events:
        <Resize>
*/

type EventType uint

const (
	KeyboardEvent EventType = iota
	MouseEvent
	ResizeEvent
)

type Event struct {
	Type    EventType
	ID      string
	Payload interface{}
}

// Mouse payload.
type Mouse struct {
	Drag bool
	X    int
	Y    int
}

// Resize payload.
type Resize struct {
	Width  int
	Height int
}

// PollEvents gets events from the terminal, converts them, then sends them to each of its channels.
func PollEvents() <-chan Event {
	ch := make(chan Event)
	go func() {
		for {
			ch <- convertTcellEvent(screen.PollEvent())
		}
	}()
	return ch
}

var keyboardMap = map[tcell.Key]string{
	tcell.KeyF1:        "<F1>",
	tcell.KeyF2:        "<F2>",
	tcell.KeyF3:        "<F3>",
	tcell.KeyF4:        "<F4>",
	tcell.KeyF5:        "<F5>",
	tcell.KeyF6:        "<F6>",
	tcell.KeyF7:        "<F7>",
	tcell.KeyF8:        "<F8>",
	tcell.KeyF9:        "<F9>",
	tcell.KeyF10:       "<F10>",
	tcell.KeyF11:       "<F11>",
	tcell.KeyF12:       "<F12>",
	tcell.KeyInsert:    "<Insert>",
	tcell.KeyDelete:    "<Delete>",
	tcell.KeyHome:      "<Home>",
	tcell.KeyEnd:       "<End>",
	tcell.KeyPgUp:      "<PageUp>",
	tcell.KeyPgDn:      "<PageDown>",
	tcell.KeyUp:        "<Up>",
	tcell.KeyDown:      "<Down>",
	tcell.KeyLeft:      "<Left>",
	tcell.KeyRight:     "<Right>",
	tcell.KeyBackspace: "<Backspace>", // tcell.KeyCtrlH
	tcell.KeyTab:       "<Tab>",       // tcell.KeyCtrlI
	tcell.KeyEnter:     "<Enter>",     // tcell.KeyCtrlM
	tcell.KeyEsc:       "<Escape>",

	tcell.KeyCtrlSpace:      "<C-<Space>>",
	tcell.KeyCtrlA:          "<C-a>",
	tcell.KeyCtrlB:          "<C-b>",
	tcell.KeyCtrlC:          "<C-c>",
	tcell.KeyCtrlD:          "<C-d>",
	tcell.KeyCtrlE:          "<C-e>",
	tcell.KeyCtrlF:          "<C-f>",
	tcell.KeyCtrlG:          "<C-g>",
	tcell.KeyCtrlJ:          "<C-j>",
	tcell.KeyCtrlK:          "<C-k>",
	tcell.KeyCtrlL:          "<C-l>",
	tcell.KeyCtrlN:          "<C-n>",
	tcell.KeyCtrlO:          "<C-o>",
	tcell.KeyCtrlP:          "<C-p>",
	tcell.KeyCtrlQ:          "<C-q>",
	tcell.KeyCtrlR:          "<C-r>",
	tcell.KeyCtrlS:          "<C-s>",
	tcell.KeyCtrlT:          "<C-t>",
	tcell.KeyCtrlU:          "<C-u>",
	tcell.KeyCtrlV:          "<C-v>",
	tcell.KeyCtrlW:          "<C-w>",
	tcell.KeyCtrlX:          "<C-x>",
	tcell.KeyCtrlY:          "<C-y>",
	tcell.KeyCtrlZ:          "<C-z>",
	tcell.KeyCtrlBackslash:  "<C-4>",
	tcell.KeyCtrlRightSq:    "<C-5>",
	tcell.KeyCtrlCarat:      "<C-6>",
	tcell.KeyCtrlUnderscore: "<C-7>",
}

// convertTcellKeyboardEvent converts a tcell keyboard event to a more friendly string format.
// Combines modifiers into the string instead of having them as additional fields in an event.
func convertTcellKeyboardEvent(e *tcell.EventKey) Event {
	var ID string
	if e.Key() == tcell.KeyRune {
		if e.Rune() == ' ' {
			ID = "<Space>"
		} else {
			ID = string(e.Rune())
		}
	} else {
		converted, ok := keyboardMap[e.Key()]
		if !ok {
			converted = ""
		}
		ID = converted
	}
	if e.Modifiers()&tcell.ModAlt != 0 {
		ID = fmt.Sprintf("<M-%s>", ID)
	}
	return Event{
		Type: KeyboardEvent,
		ID:   ID,
	}
}

func convertTcellMouseEvent(e *tcell.EventMouse) Event {
	btns := e.Buttons()
	ID := "Unknown_Mouse_Button"
	drag := false
	switch {
	case btns&tcell.ButtonPrimary != 0:
		ID = "<MouseLeft>"
	case btns&tcell.ButtonSecondary != 0:
		ID = "<MouseRight>"
	case btns&tcell.ButtonMiddle != 0:
		ID = "<MouseMiddle>"
	case btns&tcell.WheelUp != 0:
		ID = "<MouseWheelUp>"
	case btns&tcell.WheelDown != 0:
		ID = "<MouseWheelDown>"
	case btns&tcell.WheelLeft != 0:
		ID = "<MouseWheelLeft>"
	case btns&tcell.WheelRight != 0:
		ID = "<MouseWheelRight>"
	default:
		ID = "<MouseMotion>"
		drag = true
	}
	x, y := e.Position()
	return Event{
		Type: MouseEvent,
		ID:   ID,
		Payload: Mouse{
			X:    x,
			Y:    y,
			Drag: drag,
		},
	}
}

// convertTcellEvent turns a tcell event into a termui event.
func convertTcellEvent(e tcell.Event) Event {
	switch e := e.(type) {
	case *tcell.EventError:
		panic(e.Error())
	case *tcell.EventKey:
		return convertTcellKeyboardEvent(e)
	case *tcell.EventMouse:
		return convertTcellMouseEvent(e)
	case *tcell.EventResize:
		width, height := e.Size()
		return Event{
			Type: ResizeEvent,
			ID:   "<Resize>",
			Payload: Resize{
				Width:  width,
				Height: height,
			},
		}
	}
	return Event{}
}
