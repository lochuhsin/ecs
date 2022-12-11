package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	_ "net/http/pprof"
	"time"
)

var (
	emptyImage    = ebiten.NewImage(3, 3)
	emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	emptyImage.Fill(color.White)
}

const (
	screenWidth  = 1280 * 1.5
	screenHeight = 960 * 1.5
)

type Circle struct {
	isReducing bool
	radius     float32
	accelerate float32
	x          float32
	y          float32
	vx         float32
	vy         float32
}

type Game struct {
	circles []Circle
}

func (g *Game) Update() error {
	Movement(g)
	RadiusGrowth(g)
	IsOverlap(g)
	RadiusReduce(g)
	return nil
}

func drawCircle(screen *ebiten.Image, circles []Circle) {
	for i := 0; i < len(circles); i++ {

		x, y, radius := circles[i].x, circles[i].y, circles[i].radius

		var path vector.Path
		path.Arc(x, y, radius, 0, 360, vector.Clockwise)
		op := &ebiten.DrawTrianglesOptions{
			FillRule: ebiten.FillAll,
		}
		vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
		for i := range vs {
			vs[i].SrcX = 1
			vs[i].SrcY = 1
			//vs[i].ColorR = 173
			//vs[i].ColorG = 111
			//vs[i].ColorB = 47
		}
		screen.DrawTriangles(vs, is, emptySubImage, op)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	drawCircle(screen, g.circles)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\n", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func InitCircle(game *Game, count int, radius float32) {

	circles := make([]Circle, count)
	for i := 0; i < count; i++ {
		circles[i] = CreateNewCircle(radius)
	}
	game.circles = circles
}

func Direction() float32 {
	if rand.Float32() < 0.5 {
		return -1
	}
	return 1
}

func ExceedBoundary(val float32, boundary float32) bool {

	if val >= boundary || val <= 0 {
		return true
	}
	return false
}

func Movement(game *Game) {

	for i := 0; i < len(game.circles); i++ {
		// this should move to entity manager (do query like thing)
		if game.circles[i].isReducing {
			continue
		}

		game.circles[i].vx += Direction() * game.circles[i].accelerate
		game.circles[i].vy += Direction() * game.circles[i].accelerate

		game.circles[i].x += game.circles[i].vx
		game.circles[i].y += game.circles[i].vy

		if ExceedBoundary(game.circles[i].x, screenWidth) {
			game.circles[i].vx *= -1
		}

		if ExceedBoundary(game.circles[i].y, screenHeight) {
			game.circles[i].vy *= -1
		}
	}
}

func RadiusGrowth(game *Game) {
	for i := 0; i < len(game.circles); i++ {
		if game.circles[i].isReducing {
			continue
		}
		game.circles[i].radius += rand.Float32() * 0.05
	}
}

func IsOverlap(game *Game) {
	for i := 0; i < len(game.circles); i++ {
		for j := i + 1; j < len(game.circles); j++ {

			radiusSum := game.circles[i].radius + game.circles[j].radius
			centerDistance := math.Sqrt(math.Pow(float64(game.circles[i].x-game.circles[j].x), 2) + math.Pow(float64(game.circles[i].y-game.circles[j].y), 2))

			if centerDistance < float64(radiusSum) {
				game.circles[i].isReducing = true
				game.circles[j].isReducing = true
			}
		}
	}
}

func RadiusReduce(game *Game) {

	for i := 0; i < len(game.circles); i++ {
		// if radius < 0 create a new one at the same array
		// this should move to entity manager (do query like thing)
		if !game.circles[i].isReducing {
			continue
		}

		game.circles[i].radius -= 0.5
		if game.circles[i].radius < 0 {
			game.circles[i] = CreateNewCircle(1)
		}

	}

}

func CreateNewCircle(radius float32) Circle {
	return Circle{
		x:          rand.Float32() * screenWidth,
		y:          rand.Float32() * screenHeight,
		vx:         rand.Float32() * Direction(),
		vy:         rand.Float32() * Direction(),
		radius:     radius,
		isReducing: false,
		accelerate: 0.01,
	}
}

func main() {

	log.Println("Start initializing game world")
	rand.Seed(time.Now().UnixNano())
	g := &Game{}
	InitCircle(g, 1000, 5)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Vector (Ebiten Demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
