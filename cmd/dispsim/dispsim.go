package main

import (
	"flamingsteve/cmd"
	"flamingsteve/pkg/muthur"
	"github.com/aarzilli/nucular"
	nstyle "github.com/aarzilli/nucular/style"
)

func main() {
	cmd.SetupLoggers()

	muthur.Connect("virtual-display")
	defer muthur.Close()

	gui := NewGui(DefaultResolutionX, DefaultResolutionY)

	gui.MainWindow = nucular.NewMasterWindowSize(nucular.WindowClosable|nucular.WindowNoScrollbar, "dispsim", gui.WindowSize(), gui.render)
	gui.MainWindow.SetStyle(nstyle.FromTheme(nstyle.DarkTheme, 1.0))
	gui.MainWindow.Main()
}

