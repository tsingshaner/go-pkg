package main

import "github.com/tsingshaner/go-pkg/color"

func main() {
	println(
		color.Red("red"),
		color.Green("green"),
		color.Yellow("yellow"),
		color.Blue("blue"),
		color.Magenta("magenta"),
		color.Cyan("cyan"),
		color.White("white"),
		color.Black("black"),
		color.Gray("gray"),
	)

	println(
		color.RedBg("red"),
		color.GreenBg("green"),
		color.YellowBg("yellow"),
		color.BlueBg("blue"),
		color.MagentaBg("magenta"),
		color.CyanBg("cyan"),
		color.WhiteBg("white"),
		color.BlackBg("black"),
		color.GrayBg("gray"),
	)

	println(
		color.Bold("bold"),
		color.Underline("underline"),
		color.Reverse("reverse"),
		color.Blink("blink"),
		color.Dim("dim"),
		color.Hidden("hidden"),
		color.Italic("italic"),
		color.Strikethrough("strikethrough"),
	)

	println(
		color.Italic(
			color.Red(
				color.Bold(color.Yellow("http://")) + color.Cyan("qingshaner.com"),
			)),
	)

	println(color.Black(color.CyanBg(" INFO ")))
}
