package game

import (
	"log"
	"math/rand"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	pxTrue  = '█'
	pxFalse = ' '
)

type Game struct {
	Width, Height int
	Seed          int64
	Ratio         int

	grid []bool
	mu   sync.RWMutex
}

func (g *Game) Init() tea.Cmd {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.grid = make([]bool, g.Width*g.Height)

	rand := rand.New(rand.NewSource(g.Seed))
	for i := 0; i < len(g.grid); i++ {
		g.grid[i] = rand.Intn(g.Ratio) == 0
	}

	return nil
}

func (g *Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		log.Printf("Unexpected message type: %T (%#v)", msg, msg)
		return g, nil
	}

	switch key.String() {
	default:
		log.Printf("Seed: %d, Width: %d, Height: %d, Ratio: %d", g.Seed, g.Width, g.Height, g.Ratio)
		return g, tea.Quit
	case " ":
		g.nextFrame()
	case "n":
		g.Seed++
		g.Init()
	}

	return g, nil
}

func (g Game) View() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var str strings.Builder
	/*
		for i := 0; i < g.Width; i++ {
			str.WriteString(strconv.Itoa(i % 10))
		}
		str.WriteRune('\n')
	*/

	for i := 0; i < len(g.grid); i++ {
		r := pxFalse
		if g.grid[i] {
			r = pxTrue
		}
		str.WriteRune(r)
		if i%g.Width == g.Width-1 {
			str.WriteRune('\n')
		}
	}
	return str.String()
}

func (g *Game) nextFrame() {
	g.mu.Lock()
	defer g.mu.Unlock()

	next := make([]bool, len(g.grid))

	for i := 0; i < len(g.grid); i++ {
		alive := g.grid[i]
		neighbors := g.sumNeighbors(i)

		next[i] = (alive && neighbors == 2 || neighbors == 3) ||
			(!alive && neighbors == 3)
	}

	g.grid = next
}

func (g *Game) sumNeighbors(i int) int {
	var sum int

	// north
	if i >= g.Width {
		// north west
		if i%g.Width != 0 &&
			g.grid[i-g.Width-1] {
			sum++
		}
		// north
		if g.grid[i-g.Width] {
			sum++
		}
		// north east
		if i%g.Width != g.Width-1 &&
			g.grid[i-g.Width+1] {
			sum++
		}
	}

	// west
	if i%g.Width != 0 &&
		g.grid[i-1] {
		sum++
	}

	// east
	if i%g.Width != g.Width-1 &&
		g.grid[i+1] {
		sum++
	}

	// south
	if i/g.Width < g.Height-1 {
		// south west
		if i%g.Width != 0 &&
			g.grid[i+g.Width-1] {
			sum++
		}
		// south
		if g.grid[i+g.Width] {
			sum++
		}
		// south east
		if i%g.Width != g.Width-1 &&
			g.grid[i+g.Width+1] {
			sum++
		}
	}

	return sum
}
