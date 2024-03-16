package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/creack/pty"
)

const MAX_BUFFER_SIZE = 16

func main() {
	a := app.New()
	w := a.NewWindow("luiserm")

	ui := widget.NewTextGrid()

	c := exec.Command("/bin/bash")
	p, err := pty.Start(c)
	if err != nil {
		fyne.LogError("Error starting pty", err)
		os.Exit(1)
	}

	defer c.Process.Kill()

	onTypedKey := func(e *fyne.KeyEvent) {
		if e.Name == fyne.KeyEnter || e.Name == fyne.KeyReturn {
			_, _ = p.Write([]byte("\r"))
		}
	}

	onTypedRune := func(r rune) {
		_, _ = p.WriteString(string(r))
	}

	w.Canvas().SetOnTypedKey(onTypedKey)
	w.Canvas().SetOnTypedRune(onTypedRune)

	buffer := [][]rune{}
	reader := bufio.NewReader(p)

	// This go routine reads from the pty
	go func() {
		line := []rune{}
		buffer = append(buffer, line)
		for {
			r, _, err := reader.ReadRune()

			if err != nil {
				if err == io.EOF {
					return
				}
				os.Exit(1)
			}

			line = append(line, r)
			buffer[len(buffer)-1] = line

			if r == '\n' {
				if len(buffer) > MAX_BUFFER_SIZE {
					buffer = buffer[1:]
				}

				line = []rune{}
				buffer = append(buffer, line)
			}
 		}
	}()

		go func() {
			for {
				time.Sleep(100 * time.Millisecond)
				ui.SetText("")
				lines := ""
				for _, line := range buffer {
					lines += string(line)
				}
				ui.SetText(lines)
			}
		}()

	w.SetContent(
		container.NewVBox(
			container.NewGridWrap(fyne.NewSize(420, 200)),
			ui,
		),
	)

	w.ShowAndRun()
}