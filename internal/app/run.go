package app

import (
	"fmt"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/tbuen/sms-viewer/internal/gui"
)

const name = "SMS Viewer"

var version = ""
var runIdle = true

func Run() int {
	application, err := gtk.ApplicationNew("com.github.tbuen.sms-viewer", glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		fmt.Println("could not create application:", err)
		os.Exit(1)
	}

	application.Connect("activate", func() { onActivate(application) })

	return application.Run(os.Args)
}

func onActivate(application *gtk.Application) {
	gui.ShowMain(application, name, version)
}
