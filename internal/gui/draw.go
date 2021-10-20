package gui

import (
	"math"
	"regexp"
	"strconv"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"

	"github.com/tbuen/sms-viewer/internal/backend"
)

var colorRegexp = regexp.MustCompile("^([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})$")

func setSourceColor(ctx *cairo.Context, color string) {
	var c [3]float64
	cols := colorRegexp.FindStringSubmatch(color)
	if len(cols) == 4 {
		for i := 0; i < 3; i++ {
			if ii, err := strconv.ParseUint(cols[i+1], 16, 8); err == nil {
				c[i] = float64(ii) / 255
			}
		}
	}
	ctx.SetSourceRGB(c[0], c[1], c[2])
}

func drawRect(ctx *cairo.Context, x, y, width, height, radius float64) {
	degrees := math.Pi / 180.0
	ctx.NewPath()
	ctx.Arc(x+width-radius, y+radius, radius, -90*degrees, 0*degrees)
	ctx.Arc(x+width-radius, y+height-radius, radius, 0*degrees, 90*degrees)
	ctx.Arc(x+radius, y+height-radius, radius, 90*degrees, 180*degrees)
	ctx.Arc(x+radius, y+radius, radius, 180*degrees, 270*degrees)
	ctx.ClosePath()
	ctx.Fill()
}

func onDraw(da *gtk.DrawingArea, ctx *cairo.Context, id string) {
	messages := backend.Messages(id)

	width := float64(da.GetAllocatedWidth())
	offset := 10.0

	layout := pango.CairoCreateLayout(ctx)
	layout.SetWrap(pango.WRAP_WORD)
	layout.SetWidth(int(width * 0.7 * pango.PANGO_SCALE))

	date := ""
	for _, msg := range messages {
		newdate := msg.Time.Format("Mon, 02.01.2006")
		if newdate != date {
			layout.SetFontDescription(pango.FontDescriptionFromString("DejaVu Sans 10"))
			layout.SetText(newdate, -1)
			w, h := layout.GetSize()
			setSourceColor(ctx, "000000")
			drawRect(ctx, (width-(float64(w)/pango.PANGO_SCALE))/2.0-10.0, offset, float64(w)/pango.PANGO_SCALE+20.0, float64(h)/pango.PANGO_SCALE+10.0, 5.0)
			setSourceColor(ctx, "FFFFFF")
			ctx.MoveTo((width-(float64(w)/pango.PANGO_SCALE))/2.0, offset+5.0)
			pango.CairoShowLayout(ctx, layout)
			offset += float64(h)/pango.PANGO_SCALE + 30.0
			date = newdate
		}

		layout.SetFontDescription(pango.FontDescriptionFromString("DejaVu Sans 12"))
		layout.SetText(msg.Text, -1)
		w, h := layout.GetSize()

		if msg.Sent {
			setSourceColor(ctx, "33DD33")
			drawRect(ctx, width-(float64(w)/pango.PANGO_SCALE)-30.0, offset, float64(w)/pango.PANGO_SCALE+20.0, float64(h)/pango.PANGO_SCALE+40.0, 10.0)
		} else {
			setSourceColor(ctx, "A0A0A0")
			drawRect(ctx, 10.0, offset, float64(w)/pango.PANGO_SCALE+20.0, float64(h)/pango.PANGO_SCALE+40.0, 10.0)
		}

		setSourceColor(ctx, "000000")
		if msg.Sent {
			ctx.MoveTo(width-(float64(w)/pango.PANGO_SCALE)-20.0, offset+10.0)
		} else {
			ctx.MoveTo(20.0, offset+10.0)
		}
		pango.CairoShowLayout(ctx, layout)

		layout.SetFontDescription(pango.FontDescriptionFromString("DejaVu Sans 10"))
		layout.SetText(msg.Time.Format("15:04"), -1)
		w2, _ := layout.GetSize()
		setSourceColor(ctx, "444444")
		if msg.Sent {
			ctx.MoveTo(width-(float64(w2)/pango.PANGO_SCALE)-20.0, offset+float64(h)/pango.PANGO_SCALE+15.0)
		} else {
			ctx.MoveTo(float64(w)/pango.PANGO_SCALE+20.0-float64(w2)/pango.PANGO_SCALE, offset+float64(h)/pango.PANGO_SCALE+15.0)
		}
		pango.CairoShowLayout(ctx, layout)

		offset += float64(h)/pango.PANGO_SCALE + 60.0
	}
	pw, _ := da.GetPreferredWidth()
	da.SetSizeRequest(pw, int(offset))
}
