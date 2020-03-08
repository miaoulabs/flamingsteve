package main

import (
	"flamingsteve/cmd"
	"flamingsteve/pkg/muthur"
	"github.com/aarzilli/nucular"
	nstyle "github.com/aarzilli/nucular/style"
	"github.com/spf13/pflag"
)

var args struct {
	name   string
	model  string
	width  uint
	height uint
}

func init() {
	pflag.StringVarP(&args.name, "name", "n", "simulator", "display id to use for identification")
	pflag.StringVarP(&args.model, "model", "m", "virtual", "display id to use for identification")
	pflag.UintVarP(&args.width, "width", "w", DefaultResolutionX, "width")
	pflag.UintVarP(&args.height, "height", "h", DefaultResolutionY, "height")
}

func main() {
	pflag.Parse()
	cmd.SetupLoggers()

	muthur.Connect("virtual-display")
	defer muthur.Close()

	gui := NewGui()

	gui.MainWindow = nucular.NewMasterWindowSize(nucular.WindowClosable|nucular.WindowNoScrollbar, "dispsim", gui.WindowSize(), gui.render)
	gui.MainWindow.SetStyle(nstyle.FromTheme(nstyle.DarkTheme, 1.0))
	gui.MainWindow.Main()
}
