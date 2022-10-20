package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"game-of-life/game"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var (
		width  = flag.Int("width", 80, "Grid width")
		height = flag.Int("height", 30, "Grid height")
		seed   = flag.Int64("seed", 0, "Game seed")
		play   = flag.Bool("play", false, "Auto play next frame")
		fps    = flag.Int("fps", 24, "Frames per second")
		ratio  = flag.Int("ratio", 20, "Initial alive ratio (one in N)")
	)
	flag.Parse()

	p := tea.NewProgram(&game.Game{
		Width:  *width,
		Height: *height,
		Seed:   *seed,
		Ratio:  *ratio,
	})

	if *play {
		go func() {
			for range time.Tick(time.Second / time.Duration(*fps)) {
				p.Send(tea.KeyMsg{Type: tea.KeySpace})
			}
		}()
	}

	if err := p.Start(); err != nil {
		return fmt.Errorf("err: %s", err)
	}

	return nil
}
