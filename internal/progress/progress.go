package progress

import (
	"fmt"
	"os"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/term"
)

type Progress struct {
	Bar *progressbar.ProgressBar
}

func NewProgress(max int64, desc string) *Progress {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width < 100 {
		// If the size of the terminal cannot be obtained or is smaller than 100, set the width roughly to 10.
		width = 10
	} else {
		// adjust the width roughly...
		width = width / 4
	}

	return &Progress{
		Bar: progressbar.NewOptions64(
			max,
			progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionSetWidth(width),
			progressbar.OptionSetDescription(fmt.Sprintf("🐟 %s", desc)),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[blue]=[reset]",
				SaucerHead:    "[blue]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
		),
	}
}

func (p *Progress) Increment() {
	_ = p.Bar.Add(1)
}
