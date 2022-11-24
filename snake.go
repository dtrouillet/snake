package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Direction int
type GameState int

const (
	UP    = 0
	RIGHT = 1
	DOWN  = 2
	LEFT  = 3
)

const (
	MENU    = 0
	START   = 1
	PLAYING = 2
	ENDING  = 3
)

const PLAYER_INIT_POSITION_X = 160
const PLAYER_INIT_POSITION_Y = 160

type Game struct {
	direction Direction
	state     GameState
	snake     []SnakeSlice
	size      int
	snakeSize int
	speed     float64
	test      bool
	test2     int
	food      Food
}

type SnakeSlice struct {
	positionX float64
	positionY float64
}

type Food struct {
	positionX float64
	positionY float64
}

func (g *Game) Update() error {
	if g.state == ENDING {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = START
		}
		return nil
	}

	if g.state == START {
		snakeSlice := SnakeSlice{
			positionX: PLAYER_INIT_POSITION_X,
			positionY: PLAYER_INIT_POSITION_Y,
		}
		g.snake = make([]SnakeSlice, 1)
		g.speed = 10
		g.size = 10
		g.snakeSize = 1
		g.snake[0] = snakeSlice
		g.direction = RIGHT
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = PLAYING
		}
		g.food.positionX = -1
		g.food.positionY = -1
		g.test = true
		g.test2 = 0
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.direction = UP
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.direction = RIGHT
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.direction = DOWN
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.direction = LEFT
	}

	g.test2++
	if g.test2%10 != 0 {
		return nil
	}

	snakeNew := make([]SnakeSlice, 1)
	snakeSliceNew := SnakeSlice{
		positionX: g.snake[0].positionX,
		positionY: g.snake[0].positionY,
	}
	switch g.direction {
	case UP:
		snakeSliceNew.positionY -= g.speed
	case RIGHT:
		snakeSliceNew.positionX += g.speed
	case DOWN:
		snakeSliceNew.positionY += g.speed
	case LEFT:
		snakeSliceNew.positionX -= g.speed
	}

	snakeNew[0] = snakeSliceNew
	snakeNew = append(snakeNew, g.snake[:g.snakeSize-1]...)
	g.snake = snakeNew

	if (g.food.positionX == -1 || g.food.positionY == -1) && 0 == rand.Intn(1) {
		rand.Seed(time.Now().UnixNano())
		x := rand.Intn(320)
		y := rand.Intn(320)
		g.food.positionX = float64(x - x%10)
		g.food.positionY = float64(y - y%10)
	}

	if g.snake[0].positionX == g.food.positionX && g.snake[0].positionY == g.food.positionY {
		//if false {
		g.food.positionX = -1
		g.food.positionY = -1
		posX := g.snake[len(g.snake)-1].positionX
		posY := g.snake[len(g.snake)-1].positionY
		switch g.direction {
		case UP:
			posY += 10.0
		case RIGHT:
			posX -= 10.0
		case DOWN:
			posY -= 10.0
		case LEFT:
			posX += 10.0
		}
		snakeSlice := SnakeSlice{
			positionX: posX,
			positionY: posY,
		}

		g.snake = append(g.snake, snakeSlice)
		g.snakeSize++
	}

	//IL FAUT CREER UN NVX TABLEAU
	//AJOUTER UN NVX ITEM AVEC LA NOUVELLE POSITION
	//AJOUTER TOUS LES ITEMS DE L'ANCIEN TABLEAU SANS LE DERNIER ITEM
	if g.snake[0].positionX < 0 || g.snake[0].positionX > 320-10 || g.snake[0].positionY < 0 || g.snake[0].positionY > 320-10 {
		g.state = ENDING
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.state == START {
		ebitenutil.DebugPrint(screen, "PRESS START")
		return
	}
	if g.state == ENDING {
		ebitenutil.DebugPrint(screen, "SCORE : "+strconv.Itoa(g.snakeSize))
		return
	}

	img := ebiten.NewImage(g.size, g.size)
	img.Fill(color.RGBA{
		R: 200,
		G: 100,
		B: 10,
		A: 100,
	})

	img2 := ebiten.NewImage(g.size, g.size)
	img2.Fill(color.RGBA{
		R: 100,
		G: 0,
		B: 0,
		A: 100,
	})

	for index, snakeSlice := range g.snake {
		op := &ebiten.DrawImageOptions{}

		op.GeoM.Translate(snakeSlice.positionX, snakeSlice.positionY)

		if index%2 == 0 {
			screen.DrawImage(img, op)
		} else {
			screen.DrawImage(img2, op)
		}
	}

	if g.food.positionX != -1 && g.food.positionY != -1 {
		food := ebiten.NewImage(g.size, g.size)
		food.Fill(color.RGBA{
			R: 0,
			G: 0,
			B: 100,
			A: 100,
		})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(g.food.positionX, g.food.positionY)
		screen.DrawImage(food, op)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f Size: %d / %d", ebiten.ActualTPS(), len(g.snake), g.snakeSize))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 320
}

func main() {
	ebiten.SetWindowSize(320, 320)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{state: START}); err != nil {
		log.Fatal(err)
	}
}
