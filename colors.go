package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/acarl005/stripansi"
	chromastyles "github.com/alecthomas/chroma/styles"
	"github.com/jwalton/gchalk"
	markdown "github.com/quackduck/go-term-markdown"
)

func makeFlag(colors []string) func(a string) string {
	flag := make([]*gchalk.Builder, len(colors))
	for i := range colors {
		flag[i] = Chalk.WithHex(colors[i])
	}
	return func(a string) string {
		return applyMulticolor(flag, a)
	}
}

func applyHueRange(start, end float64, a string) string {
	a = stripansi.Strip(a)
	buf := strings.Builder{}
	for i, r := range []rune(a) {
		h := start + (end-start)*float64(i)/float64(len(a))
		buf.WriteString(hue(h).Paint(string(r)))
	}
	return buf.String()
}

func applyStyles(styles []*Style, a string) string {
	//a = stripansi.Strip(a)
	buf := strings.Builder{}
	colorOffset := rand.Intn(len(styles))
	for i, s := range tokenizeAnsi(a) {
		buf.WriteString(styles[(colorOffset+i)%len(styles)].apply(s))
	}
	return buf.String()
}

func applyMulticolor(colors []*gchalk.Builder, a string) string {
	a = stripansi.Strip(a)
	buf := strings.Builder{}
	colorOffset := rand.Intn(len(colors))
	for i, r := range []rune(a) {
		buf.WriteString(colors[(colorOffset+i)%len(colors)].Paint(string(r)))
	}
	return buf.String()
}

// splits runes and includes their color codes
func tokenizeAnsi(a string) []string {
	tokens := make([]string, 0, len(a)/3)
	buf := strings.Builder{}
	buildUntilM := false // m delineates end of ansi color code
	for _, r := range a {
		buf.WriteRune(r)
		if r == '\033' {
			buildUntilM = true
			continue
		}
		if buildUntilM {
			if r == 'm' {
				buildUntilM = false
			}
			continue
		}
		tokens = append(tokens, buf.String())
		buf.Reset()
	}
	if buf.Len() > 0 { // that last m could be needed
		tokens = append(tokens, buf.String())
	}
	return tokens
}

var (
	Chalk   = gchalk.New(gchalk.ForceLevel(gchalk.LevelAnsi256))
	Green   = ansi256(1, 5, 1)
	Red     = ansi256(5, 1, 1)
	Cyan    = ansi256(1, 5, 5)
	Magenta = ansi256(5, 1, 5)
	Yellow  = ansi256(5, 5, 1)
	Orange  = ansi256(5, 3, 0)
	Blue    = ansi256(0, 3, 5)
	White   = ansi256(5, 5, 5)
	Styles  = []*Style{
		{"white", buildStyle(White)},
		{"red", buildStyle(Red)},
		{"coral", buildStyle(ansi256(5, 2, 2))},
		{"green", buildStyle(Green)},
		{"sky", buildStyle(ansi256(3, 5, 5))},
		{"cyan", buildStyle(Cyan)},
		{"magenta", buildStyle(Magenta)},
		{"pink", buildStyle(ansi256(5, 3, 4))},
		{"rose", buildStyle(ansi256(5, 0, 2))},
		{"cranberry", buildStyle(ansi256(3, 0, 1))},
		{"lavender", buildStyle(ansi256(4, 2, 5))},
		{"fire", buildStyle(ansi256(5, 2, 0))},
		{"pastelgreen", buildStyle(ansi256(0, 5, 3))},
		{"olive", buildStyle(ansi256(4, 5, 1))},
		{"yellow", buildStyle(Yellow)},
		{"orange", buildStyle(Orange)},
		{"blue", buildStyle(Blue)}}
	SecretStyles = []*Style{
		{"elitedino", buildStyle(ansi256(5, 0, 0))},
		{"ukraine", buildStyle(Chalk.WithHex("#005bbb").WithBgHex("#ffd500"))},
		{"easter", buildStyle(Chalk.WithRGB(255, 51, 255).WithBgRGB(255, 255, 0))},
		{"baby", buildStyle(Chalk.WithRGB(255, 51, 255).WithBgRGB(102, 102, 255))},
		{"hacker", buildStyle(Chalk.WithRGB(0, 255, 0).WithBgRGB(0, 0, 0))},
		{"l33t", buildStyleNoStrip(Chalk.WithBgBrightBlack())},
		{"whiten", buildStyleNoStrip(Chalk.WithBgWhite())},
		{"trans", makeFlag([]string{"#55CDFC", "#F7A8B8", "#FFFFFF", "#F7A8B8", "#55CDFC"})},
		{"gay", makeFlag([]string{"#FF0018", "#FFA52C", "#FFFF41", "#008018", "#0000F9", "#86007D"})},
		{"lesbian", makeFlag([]string{"#D62E02", "#FD9855", "#FFFFFF", "#D161A2", "#A20160"})},
		{"bi", makeFlag([]string{"#D60270", "#D60270", "#9B4F96", "#0038A8", "#0038A8"})},
		{"ace", makeFlag([]string{"#333333", "#A4A4A4", "#FFFFFF", "#810081"})},
		{"pan", makeFlag([]string{"#FF1B8D", "#FFDA00", "#1BB3FF"})},
		{"enby", makeFlag([]string{"#FFF430", "#FFFFFF", "#9C59D1", "#000000"})},
		{"aro", makeFlag([]string{"#3AA63F", "#A8D47A", "#FFFFFF", "#AAAAAA", "#000000"})},
		{"genderfluid", makeFlag([]string{"#FE75A1", "#FFFFFF", "#BE18D6", "#333333", "#333EBC"})},
		{"agender", makeFlag([]string{"#333333", "#BCC5C6", "#FFFFFF", "#B5F582", "#FFFFFF", "#BCC5C6", "#333333"})},
		{"rainbow", func(a string) string {
			//rainbow := []*gchalk.Builder{}
			//rainbow := []*gchalk.Builder{Red, Orange, Yellow, Green, Cyan, Blue, ansi256(2, 2, 5), Magenta}
			span := 360.0
			length := len(stripansi.Strip(a))
			if length < 8 {
				span = 45 * float64(length) // at least 45 degrees per letter
			}
			start := 360 * rand.Float64()
			return applyHueRange(start, start+span, a)
			//return applyMulticolor(rainbow, a)
		}}}
)

func init() {
	markdown.CurrentTheme = chromastyles.ParaisoDark
}

type Style struct {
	name  string
	apply func(string) string
}

func buildStyle(c *gchalk.Builder) func(string) string {
	return func(s string) string { return c.Paint(stripansi.Strip(s)) }
}
func buildStyleNoStrip(c *gchalk.Builder) func(string) string {
	return func(s string) string { return c.Paint(s) }
}

// h from 0 to 360
// https://www.desmos.com/calculator/wb91fw4nyj
func hue(h float64) *gchalk.Builder {
	pi := math.Pi
	h = math.Mod(h, 360) / 360.0
	r := math.Round(255.0 * (0.5 + 0.5*math.Sin(2*pi*h+pi/2)))
	g := math.Round(255.0 * (0.5 + 0.5*math.Sin(2*pi*h+pi/2+2*pi/3)))
	b := math.Round(255.0 * (0.5 + 0.5*math.Sin(2*pi*h+pi/2+4*pi/3)))
	//r, g, b, err := colorconv.HSVToRGB(math.Mod(h, 360), s, v)
	//if err != nil {
	//	return Chalk.WithRGB(0, 0, 0)
	//}
	return Chalk.WithRGB(uint8(r), uint8(g), uint8(b))
}

// with r, g and b values from 0 to 5
func ansi256(r, g, b uint8) *gchalk.Builder {
	return Chalk.WithRGB(255/5*r, 255/5*g, 255/5*b)
	//return Chalk.WithRGB(uint8(math.Round(255*float64(r)/5)), uint8(math.Round(255*float64(g)/5)), uint8(math.Round(255*float64(b)/5)))
}

func bgAnsi256(r, g, b uint8) *gchalk.Builder {
	return Chalk.WithBgRGB(255/5*r, 255/5*g, 255/5*b)
}

// Applies color from name
func (u *User) changeColor(colorName string) error {
	if strings.Contains(colorName, "bg") {
		if names := strings.Fields(colorName); len(names) > 1 { // do we need to separate bg and fg colors?
			fgColors := make([]string, 0, len(names)-1)
			bgColors := make([]string, 0, len(names))
			for _, name := range names {
				if strings.HasPrefix(name, "bg") {
					bgColors = append(bgColors, name)
				} else {
					fgColors = append(fgColors, name)
				}
			}
			if len(fgColors) != 0 { // if no fg colors, carry on normally
				err := u.changeColor(strings.Join(fgColors, " "))
				if err != nil {
					return err
				}
				return u.changeColor(strings.Join(bgColors, " "))
			}
		}
	}
	style, err := getStyle(colorName)
	if err != nil {
		return err
	}

	//changedBg := false
	if strings.HasPrefix(colorName, "bg") {
		//changedBg = true
		u.ColorBG = style.name // update bg color
	} else {
		u.Color = style.name // update fg color
	}

	u.Name, _ = applyColorToData(u.Name, u.Color, u.ColorBG)
	//styleFG := &Style{}
	//styleBG := &Style{}
	//if changedBg {
	//	styleFG, _ = getStyle(u.Color) // already checked for errors
	//	styleBG = style
	//} else {
	//	styleBG, err = getStyle(u.ColorBG) // already checked for errors
	//	styleFG = style
	//}
	//u.Name = styleBG.apply(styleFG.apply(u.Name))
	u.term.SetPrompt(u.Name + ": ")
	return nil
}

func applyColorToData(data string, color string, colorBG string) (string, error) {
	styleFG, err := getStyle(color)
	if err != nil {
		return "", err
	}
	styleBG, err := getStyle(colorBG)
	if err != nil {
		return "", err
	}
	return styleBG.apply(styleFG.apply(data)), nil // fg clears the bg color
}

// If the input is a named style, returns it. Otherwise, returns nil.
func getNamedColor(name string) *Style {
	for i := range Styles {
		if Styles[i].name == name {
			return Styles[i]
		}
	}
	for i := range SecretStyles {
		if SecretStyles[i].name == name {
			return SecretStyles[i]
		}
	}
	return nil
}

func getCustomColor(name string) (*Style, error) {
	if strings.HasPrefix(name, "#") {
		return &Style{name, buildStyle(Chalk.WithHex(name))}, nil
	}
	if strings.HasPrefix(name, "bg-#") {
		return &Style{name, buildStyleNoStrip(Chalk.WithBgHex(strings.TrimPrefix(name, "bg-")))}, nil
	}
	if len(name) == 3 || len(name) == 6 {
		rgbCode := name
		if strings.HasPrefix(name, "bg-") {
			rgbCode = strings.TrimPrefix(rgbCode, "bg-")
		}
		a, err := strconv.Atoi(rgbCode)
		if err == nil {
			r := (a / 100) % 10
			g := (a / 10) % 10
			b := a % 10
			if r > 5 || g > 5 || b > 5 || r < 0 || g < 0 || b < 0 {
				return nil, errors.New("custom colors have values from 0 to 5 smh")
			}
			if strings.HasPrefix(name, "bg-") {
				return &Style{name, buildStyleNoStrip(bgAnsi256(uint8(r), uint8(g), uint8(b)))}, nil
			}
			return &Style{name, buildStyle(ansi256(uint8(r), uint8(g), uint8(b)))}, nil
		}
		return nil, err
	}
	return nil, nil
}

// Turns name into a style (defaults to nil)
func getStyle(name string) (*Style, error) {
	name = strings.TrimSpace(name)
	if names := strings.Fields(name); len(names) > 1 {
		styleSlice := make([]*Style, len(names))
		newName := ""
		for i := range names {
			style, err := getStyle(names[i])
			if err != nil {
				return nil, err
			}
			styleSlice[i] = style
			newName += style.name + " "
		}

		return &Style{newName[:len(newName)-1], func(a string) string {
			return applyStyles(styleSlice, a)
		}}, nil
	}
	switch name {
	case "random":
		r, g, b := rand.Intn(6), rand.Intn(6), rand.Intn(6)
		return &Style{fmt.Sprintf("%03d", r*100+g*10+b), buildStyle(ansi256(uint8(r), uint8(g), uint8(b)))}, nil
	case "bg-random":
		r, g, b := rand.Intn(6), rand.Intn(6), rand.Intn(6)
		return &Style{fmt.Sprintf("%03d", r*100+g*10+b), buildStyleNoStrip(bgAnsi256(uint8(r), uint8(g), uint8(b)))}, nil
	case "bg-off":
		return &Style{"bg-off", func(a string) string { return a }}, nil // no background
	}
	namedColor := getNamedColor(name)
	if namedColor != nil {
		return namedColor, nil
	}
	if strings.HasPrefix(name, "#") {
		return &Style{name, buildStyle(Chalk.WithHex(name))}, nil
	}
	customColor, err := getCustomColor(name)
	if err != nil {
		return nil, err
	}
	if customColor != nil {
		return customColor, nil
	}
	//s, err := Chalk.WithStyle(strings.Split(name, "-")...)
	//if err == nil {
	//	return &style{name, buildStyle(s)}, nil
	//}

	return nil, errors.New("Which color? Choose from random, " + strings.Join(func() []string {
		colors := make([]string, 0, len(Styles))
		for i := range Styles {
			colors = append(colors, Styles[i].name)
		}
		return colors
	}(), ", ") + "  \nMake your own colors using hex (#A0FFFF, etc) or RGB values from 0 to 5 (for example, `color 530`, a pretty nice orange). Set bg color like this: color bg-530; remove bg color with color bg-off.\nThere's also a few secret colors :)")
}
