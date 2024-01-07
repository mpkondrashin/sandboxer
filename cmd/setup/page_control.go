package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type PageControl struct {
	statusLabel   *widget.Label
	bypassButton  *widget.Button
	inspectButton *widget.Button
}

var _ Page = &PageOptions{}

func (p *PageControl) Name() string {
	return "Control"
}

func (p *PageControl) GetStatus(model *Model) {
}

func (p *PageControl) Content(win fyne.Window, model *Model) fyne.CanvasObject {
	return container.NewVBox(
		container.NewHBox(widget.NewLabel("Current State:"), p.statusLabel),
		p.bypassButton,
		p.inspectButton,
	)
}

func (p *PageControl) AquireData(model *Model) error {
	return nil
}

type progress struct {
	chars string
	size  int
}

const (
	ProgressTypePulse = iota
	ProgressTypeRot1
	ProgressTypebdqp
	ProgressTypeFlyBy
	ProgressTypeRotVWide
	ProgressTypeRotVNarrow
	ProgressTypeFlap
	ProgressTypeRot
	ProgressTypeRotFly
	ProgressTypeRotFlyWide
	ProgressTypeFlipAndBack
	ProgressTypeFlip
	ProgressTypeFlipShort
	ProgressTypeDrop
	ProgressTypeRotWide
	ProgressTypeX
)

var progressTypes = []progress{
	{
		chars: `.oOo`,
		size:  1,
	},
	{
		chars: `/-\|`,
		size:  1,
	},
	{
		chars: ` booo dooo oboo odoo oobo oodo ooob oood oooq ooop ooqo oopo oqoo opoo qooo pooo`,
		size:  5,
	},
	{ //        123451234512345123451234512345123451234512345123451234512345
		chars: `     .     o     O     o     .         .   o   O   o   .    `,
		size:  5,
	},
	{
		chars: `<      ^      >  v  `,
		size:  5,
	},
	{
		chars: `<   ^   > v `,
		size:  3,
	},
	{
		chars: `->|<-<|>-`,
		size:  1,
	},
	{ //        123123123123123123123
		chars: `-- \   |   / --`,
		size:  3,
	},
	{ //        12345123451234512345123451234512345123451234512345
		chars: `--    \     |     /    --   \   |   /   `,
		size:  5,
	},
	{ //        12345123451234512345123451234512345123451234512345
		chars: `__    \     |     /    __  __  __  `,
		size:  5,
	},
	{ //        12345123451234512345123451234512345123451234512345
		chars: `__    \     |     /    __   /   |   \   `,
		size:  5,
	},
	{ //        12345123451234512345123451234512345123451234512345
		chars: `__    \     |     / `,
		size:  5,
	},
	{
		chars: "':,",
		size:  1,
	},
	{
		chars: ` |  / --- \ `,
		size:  3,
	},
	{
		chars: `.-+x+-. `,
		size:  1,
	},
}

func InfiniteProgressFunc(callback func(s string)) func() {
	stop := make(chan struct{})
	pick := ProgressTypeRotFly
	chars := progressTypes[pick].chars
	size := progressTypes[pick].size
	go func() {
		i := 0
		sleepTime := 200 * time.Millisecond
		for {
			select {
			case <-stop:
				return
			default:
				callback(chars[i*size : i*size+size])
				i++
				if i*size == len(chars) {
					i = 0
				}
				time.Sleep(sleepTime)
			}
		}
	}()
	return func() {
		stop <- struct{}{}
	}
}
