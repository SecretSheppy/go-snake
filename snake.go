package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	GridX              = 10
	GridY              = 10
	GridGapX           = 2
	GridGapY           = 0
	AppleTile          = 'A'
	SnakeBodyTile      = 'S'
	SnakeHeadTile      = 'H'
	SnakeVictoryLength = 20
)

type Vector struct {
	X, Y int
}

func (v *Vector) equals(v2 *Vector) bool {
	return v.X == v2.X && v.Y == v2.Y
}

func (v *Vector) String() string {
	return fmt.Sprintf("(%d,%d)", v.X, v.Y)
}

func pad(count int) string {
	return strings.Repeat(" ", count)
}

func drawXAxisLabel() {
	xAxisLabel := pad(1+2*GridGapX) + "|"
	for x := 0; x < GridX; x++ {
		xAxisLabel += pad(GridGapX) + strconv.Itoa(x)
	}
	fmt.Println(xAxisLabel)
}

func drawXAxisDivider() {
	fmt.Println(strings.Repeat("-", GridX+(GridX+4)*GridGapX))
}

func drawYAxisDivider() {
	fmt.Println(pad(1+2*GridGapX) + "|" + pad(GridX+(GridX+4)*GridGapX))
}

func yAxisLabelGenerator() func() string {
	labelNumber := 0
	return func() string {
		labelNumber++
		return pad(GridGapX) + strconv.Itoa(labelNumber-1) + pad(GridGapX) + "|"
	}
}

func populateGrid(apple *Vector, snake []Vector) [][]byte {
	grid := make([][]byte, GridY)

	for y := 0; y < GridY; y++ {
		grid[y] = make([]byte, GridX)
		for x := 0; x < GridX; x++ {
			grid[y][x] = ' '
		}
	}

	grid[apple.Y][apple.X] = AppleTile
	grid[snake[0].Y][snake[0].X] = SnakeHeadTile

	for _, vector := range snake[1:] {
		grid[vector.Y][vector.X] = SnakeBodyTile
	}

	return grid
}

// only implemented for windows
func clearConsole() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println("Failed to clear console")
	}
}

func drawGrid(populatedGrid [][]byte) {
	clearConsole()
	drawXAxisLabel()
	drawXAxisDivider()

	getLabel := yAxisLabelGenerator()

	// Drawing the main grid body
	for y := 0; y < GridY; y++ {
		row := getLabel()
		for x := 0; x < GridX; x++ {
			row += pad(GridGapX) + string(populatedGrid[y][x])
		}
		fmt.Println(row)

		for i := 0; i < GridGapY; i++ {
			drawYAxisDivider()
		}
	}

	drawXAxisDivider()
	drawXAxisLabel()
}

func main() {
	apple := &Vector{0, 0}
	snake := []Vector{{rand.Intn(GridX), rand.Intn(GridY)}}
	drawGrid(populateGrid(apple, snake))

	for {
		head := snake[0]

		if snake[0].equals(apple) {
			snake = append([]Vector{*apple}, snake...)
			apple.X = rand.Intn(GridX)
			apple.Y = rand.Intn(GridY)
		} else {
			switch {
			case snake[0].X < apple.X:
				head.X++
			case snake[0].X > apple.X:
				head.X--
			case snake[0].Y < apple.Y:
				head.Y++
			case snake[0].Y > apple.Y:
				head.Y--
			}
			snake = append([]Vector{head}, snake[:len(snake)-1]...)
		}

		drawGrid(populateGrid(apple, snake))

		if len(snake) == SnakeVictoryLength {
			break
		}

		time.Sleep(300 * time.Millisecond)
	}
}
