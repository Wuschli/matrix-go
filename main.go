package main

import (
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/nsf/termbox-go"
)

var width, height int
var framerate = time.Duration(1000000/60) * time.Microsecond

var runes = []rune{
	'｡', '｢', '｣', '､', '･', 'ｦ', 'ｧ', 'ｨ', 'ｩ', 'ｪ', 'ｫ', 'ｬ', 'ｭ', 'ｮ', 'ｯ',
	'ｰ', 'ｱ', 'ｲ', 'ｳ', 'ｴ', 'ｵ', 'ｶ', 'ｷ', 'ｸ', 'ｹ', 'ｺ', 'ｻ', 'ｼ', 'ｽ', 'ｾ', 'ｿ',
	'ﾀ', 'ﾁ', 'ﾂ', 'ﾃ', 'ﾄ', 'ﾅ', 'ﾆ', 'ﾇ', 'ﾈ', 'ﾉ', 'ﾊ', 'ﾋ', 'ﾌ', 'ﾍ', 'ﾎ', 'ﾏ',
	'ﾐ', 'ﾑ', 'ﾒ', 'ﾓ', 'ﾔ', 'ﾕ', 'ﾖ', 'ﾗ', 'ﾘ', 'ﾙ', 'ﾚ', 'ﾛ', 'ﾜ', 'ﾝ', 'ﾞ', 'ﾟ',
}

var lenrunes = len(runes)

func rain(i int) {
	y := rand.Intn(height)
	x := i
	var r [16]rune
	runeOffset := 0
	sleeptime := time.Duration(50+rand.Intn(75)) * time.Millisecond

	for i := range r {
		r[i] = runes[rand.Intn(lenrunes)]
	}
	lenr := len(r)

	for {
		fg := termbox.Attribute(0x29)
		for j := range r {
			p := j * j * j * j
			if rand.Intn((lenr-1)*(lenr-1)*(lenr-1)*(lenr-1)) < p {
				r[(j+runeOffset)%lenr] = runes[rand.Intn(lenrunes)]
				fg = 0x31
			}
			if j == lenr-1 {
				fg = termbox.ColorWhite
			}
			termbox.SetCell(x, (y+j)%height, r[(j+runeOffset)%lenr], fg, termbox.ColorBlack)

		}
		termbox.SetCell(x, (y-1)%height, ' ', termbox.ColorBlack, termbox.ColorBlack)
		runeOffset++
		y++
		<-time.After(sleeptime)

	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	width, height = termbox.Size()
	termbox.SetOutputMode(termbox.Output256)

	for i := 0; i < width; i++ {
		go rain(i)
	}

	events := make(chan termbox.Event)

	go func() {
		for {
			event := termbox.PollEvent()
			events <- event
		}
	}()

	sigChan := make(chan os.Signal)
	done := make(chan bool)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	go func() {
		for {
			select {
			case <-sigChan:
				done <- true
				return
			case event := <-events:
				switch event.Type {
				case termbox.EventKey:
					switch event.Key {
					case termbox.KeyCtrlC:
						done <- true
						return
					}
				}
			}
		}
	}()

	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	go func() {
		for {
			<-time.After(framerate)
			err := termbox.Flush()
			if err != nil {
				panic(err)
			}
		}
	}()

	<-done
	defer termbox.Close()
}
