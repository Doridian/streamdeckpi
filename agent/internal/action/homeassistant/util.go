package homeassistant

import (
	"image"
	"image/color"
	"io"
	"strconv"

	"github.com/Doridian/streamdeckpi/agent/internal/controller"
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

func loadFont(ctrl controller.Controller, file string) (*truetype.Font, error) {
	reader, err := ctrl.ResolveFile(file)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return freetype.ParseFont(data)
}

var trueTypeFonts map[string]*truetype.Font

type FontAlign = int

const (
	FontAlignLeft FontAlign = iota
	FontAlignCenter
	FontAlignRight
)

const (
	FontAlignTop FontAlign = iota
	FontAlignMiddle
	FontAlignBottom
)

func drawCenteredText(ctrl controller.Controller, img *image.RGBA, font string, col color.RGBA, x, y int, fontSize float64, align, verticalAlign FontAlign, label string) error {
	if trueTypeFonts == nil {
		trueTypeFonts = make(map[string]*truetype.Font)
	}

	trueTypeFont := trueTypeFonts[font]
	if trueTypeFont == nil {
		var err error
		trueTypeFont, err = loadFont(ctrl, font)
		if err != nil {
			return err
		}
		trueTypeFonts[font] = trueTypeFont
	}

	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	ttCtx := freetype.NewContext()
	ttCtx.SetFont(trueTypeFont)
	ttCtx.SetFontSize(fontSize)

	labelHeight := ttCtx.PointToFixed(fontSize)

	labelSize, err := ttCtx.DrawString(label, fixed.Point26_6{})
	if err != nil {
		return err
	}

	pointFont := fixed.Point26_6{X: 0, Y: 0}

	switch align {
	case FontAlignLeft:
		pointFont.X = point.X
	case FontAlignRight:
		pointFont.X = point.X - labelSize.X
	case FontAlignCenter:
		fallthrough
	default:
		pointFont.X = point.X - (labelSize.X / 2)
	}
	switch verticalAlign {
	case FontAlignTop:
		pointFont.Y = point.Y + labelHeight
	case FontAlignBottom:
		pointFont.Y = point.Y
	case FontAlignMiddle:
		fallthrough
	default:
		pointFont.Y = point.Y + (labelHeight / 2)
	}

	ttCtx.SetClip(img.Rect)
	ttCtx.SetDst(img)
	ttCtx.SetSrc(image.NewUniform(col))

	ttCtx.DrawString(label, pointFont)

	return nil
}
