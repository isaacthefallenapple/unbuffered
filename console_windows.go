package unbuffered

import (
	"golang.org/x/sys/windows"
)

const CONSOLE_MODE uint32 = windows.ENABLE_WINDOW_INPUT | windows.ENABLE_MOUSE_INPUT | windows.ENABLE_PROCESSED_INPUT |
	windows.ENABLE_VIRTUAL_TERMINAL_INPUT

var (
	Stdin = windows.Stdin
)

func SetUpConsole() (reset func(), err error) {
	var mode uint32 = 503
	_ = windows.GetConsoleMode(windows.Stdin, &mode)
	return ResetConsole(mode), windows.SetConsoleMode(Stdin, CONSOLE_MODE)
}

func ResetConsole(mode uint32) func() {
	return func() {
		_ = windows.SetConsoleMode(windows.Stdin, mode)
	}
}

func ReadChar() (buf uint16, err error) {
	var (
		toread       uint32 = 1
		read         uint32
		inputControl byte
	)
	err = windows.ReadConsole(windows.Stdin, &buf, toread, &read, &inputControl)
	return
}

func Byte(b *byte) (err error) {
	buf, err := ReadChar()
	*b = byte(buf)
	return
}

func Rune(r *rune) (err error) {
	buf, err := ReadChar()
	*r = rune(buf)
	return
}

func RuneStream() (ch chan rune, cancel func()) {
	ch = make(chan rune)
	quit := make(chan struct{})
	cancel = func() {
		close(quit)
	}
	go func() {
		defer close(ch)
		var r rune
		for Rune(&r) == nil {
			select {
			case <-quit:
				return
			default:
				ch <- r
			}
		}
	}()
	return
}

func ReadRune() (r rune, err error) {
	err = Rune(&r)
	return
}

func ReadRunes(delim rune) chan rune {
	ch := make(chan rune)
	go func() {
		defer close(ch)
		var r rune
		for Rune(&r) == nil && r != delim {
			ch <- r
		}
	}()
	return ch
}

func ByteStream() (ch chan byte, cancel func()) {
	ch = make(chan byte)
	quit := make(chan struct{})
	cancel = func() {
		close(quit)
	}
	go func() {
		defer close(ch)
		var b byte
		for Byte(&b) == nil {
			select {
			case <-quit:
				return
			default:
				ch <- b
			}
		}
	}()
	return
}

func ReadByte() (b byte, err error) {
	err = Byte(&b)
	return
}

func ReadBytes(delim byte) chan byte {
	ch := make(chan byte)
	go func() {
		defer close(ch)
		var b byte
		for Byte(&b) == nil && b != delim {
			ch <- b
		}
	}()
	return ch
}

type Stream struct {
	v uint16
	err error
}

func NewStream() *Stream {
	return &Stream{}
}

func (s *Stream) Next() {
	s.v, s.err = ReadChar()
}

func (s *Stream) Rune() rune {
	return rune(s.v)
}

func (s *Stream) Byte() byte {
	return byte(s.v)
}

func (s *Stream) Err() error {
	return s.err
}

