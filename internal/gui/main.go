package gui

import (
	"path/filepath"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"

	"github.com/tbuen/sms-viewer/internal/backend"
)

var main struct {
	header  *gtk.HeaderBar
	nb      *gtk.Notebook
	threads []backend.Thread
}

func ShowMain(app *gtk.Application, name, version string) {
	w := app.GetActiveWindow()
	if w != nil {
		w.Present()
		return
	}

	// window
	window, _ := gtk.ApplicationWindowNew(app)
	window.SetBorderWidth(18)
	window.SetDefaultSize(800, 600)

	// header bar
	main.header, _ = gtk.HeaderBarNew()
	main.header.SetTitle(name + " " + version)
	main.header.SetShowCloseButton(true)

	open, _ := gtk.ButtonNewFromIconName("document-open-symbolic", gtk.ICON_SIZE_BUTTON)
	open.Connect("clicked", func() { openFile() })
	main.header.PackStart(open)

	window.SetTitlebar(main.header)

	// notebook
	main.nb, _ = gtk.NotebookNew()
	main.nb.SetTabPos(gtk.POS_LEFT)
	window.Add(main.nb)

	// show
	window.ShowAll()
}

func openFile() {
	fc, err := gtk.FileChooserDialogNewWith2Buttons("Open File", nil, gtk.FILE_CHOOSER_ACTION_OPEN, "Open", gtk.RESPONSE_OK, "Cancel", gtk.RESPONSE_CANCEL)
	if err != nil {
		return
	}
	defer fc.Destroy()
	if fc.Run() == gtk.RESPONSE_OK {
		for main.nb.GetNPages() > 0 {
			main.nb.RemovePage(0)
		}
		f := fc.GetFilename()
		if err := backend.OpenDatabase(f); err != nil {
			main.header.SetSubtitle("")
			return
		}
		main.header.SetSubtitle(filepath.Base(f))
		threads := backend.Threads()
		for _, t := range threads {
			id := t.Id
			l, _ := gtk.LabelNew(t.Name)
			s, _ := gtk.ScrolledWindowNew(nil, nil)
			s.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_ALWAYS)
			d, _ := gtk.DrawingAreaNew()
			s.Add(d)
			d.Connect("draw", func(d *gtk.DrawingArea, ctx *cairo.Context) { onDraw(d, ctx, id) })
			main.nb.AppendPage(s, l)
		}
		main.nb.ShowAll()
	}
}
