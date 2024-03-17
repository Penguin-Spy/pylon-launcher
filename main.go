package main

import (
	"fmt"

	g "github.com/AllenDang/giu"
	lua "github.com/yuin/gopher-lua"
)

func onClickMe() {
	fmt.Println("Hello world!")
}

func onImSoCute() {
	fmt.Println("Im sooooooo cute!!")
}

func loop() {
	g.SingleWindow().Layout(
		g.Label("Hello world from giu"),
		g.Row(
			g.Button("Click Me").OnClick(onClickMe),
			g.Button("I'm so cute").OnClick(onImSoCute),
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

	wnd := g.NewMasterWindow("Hello world", 400, 200, g.MasterWindowFlagsNotResizable)
	wnd.Run(loop)
}
