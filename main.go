package main

import (
	"fmt"
	"image"

	g "github.com/AllenDang/giu"
	lua "github.com/yuin/gopher-lua"
)

func onClickMe() {
	fmt.Println("Hello world!")
}

func onImSoCute() {
	fmt.Println("Im sooooooo cute!!")
}

type FullImageWidget struct {
	id string
}

func FullImage(id string) *FullImageWidget {
	return &FullImageWidget{
		id: id,
	}
}

var rgba *image.RGBA

func (*FullImageWidget) Build() {
	width, height := g.GetAvailableRegion()
	fmt.Println("fiw: ", width, height)
	height = float32(rgba.Rect.Dy()) * (width / float32(rgba.Rect.Dx()))
	g.ImageWithFile("hero.jpg").Size(width, height).Build()
}

var wnd *g.MasterWindow

func loop() {
	fmt.Println("test")
	fmt.Println(wnd.GetSize())
	width, height := wnd.GetSize()
	padX, padY := g.GetWindowPadding()
	width2, height2 := float32(width)-padX*2, float32(height)-padY*2
	fmt.Println(width2, height2)

	height2 = 1240 * (width2 / 3840)

	fmt.Println(width2, height2)

	image := g.ImageWithFile("hero.jpg").Size(width2, height2)

	g.SingleWindow().Layout(
		g.Label("Hello world from giu"),
		g.Row(
			g.Column(
				g.Button("Click Me").OnClick(onClickMe),
				g.Button("I'm so cute").OnClick(onImSoCute),
			),
			g.Column(
				g.Label("I'm a centered label"),
				g.Button("I'm a centered button"),
				image,
				FullImage("a"),
			),
		),
	)
}

func main() {
	fmt.Println("hello world")

	L := lua.NewState()
	defer L.Close()
	launcher_open_lua_lib(L)
	if err := L.DoFile("testplugin.lua"); err != nil {
		panic(err)
	}

	wnd = g.NewMasterWindow("Hello world", 400, 200, 0)
	rgba, _ = g.LoadImage("hero.jpg")
	wnd.Run(loop)
}
