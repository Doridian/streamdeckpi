package homeassistant

import (
	"image"
	"image/color"
	"strconv"

	"github.com/Doridian/streamdeckpi/agent"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
)

func coerceNumber(val interface{}) (float64, bool) {
	valNum, ok := val.(float64)
	if ok {
		return valNum, true
	}

	valInt, ok := val.(int)
	if ok {
		return float64(valInt), true
	}

	valStr, ok := val.(string)
	if !ok {
		return 0, false
	}

	valNum, err := strconv.ParseFloat(valStr, 64)
	if err == nil {
		return 0, false
	}

	return valNum, true
}

func mustLoadFont(embedFile string) *truetype.Font {
	data, err := agent.FS.ReadFile(embedFile)
	if err != nil {
		panic(err)
	}
	font, err := freetype.ParseFont(data)
	if err != nil {
		panic(err)
	}
	return font
}

var trueTypeFont = mustLoadFont("embed/font.ttf")

func drawCenteredText(img *image.RGBA, col color.RGBA, x, y int, label string) {
	fontSize := 48.0

	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	ttCtx := freetype.NewContext()
	ttCtx.SetFont(trueTypeFont)
	ttCtx.SetFontSize(fontSize)

	labelHeight := ttCtx.PointToFixed(fontSize)

	labelSize, err := ttCtx.DrawString(label, fixed.Point26_6{})
	if err != nil {
		return
	}

	pointFont := fixed.Point26_6{X: point.X - (labelSize.X / 2), Y: point.Y + (labelHeight / 2)}

	ttCtx.SetClip(img.Rect)
	ttCtx.SetDst(img)
	ttCtx.SetSrc(image.NewUniform(col))

	ttCtx.DrawString(label, pointFont)
}
