package utils

import (
	"github.com/fatih/color"
	"log"
)

var (
	Info = Cyan
	Warn = Yellow
	Fata = Red
	Pur  = Purple
	Tea  = Teal
	OK   = HiGreen
)

func Cyan(format string, a ...interface{}) {
	color.Cyan(format, a...)
}

func Yellow(format string, a ...interface{}) {
	color.Yellow(format, a...)
}

func YellowBg(format string, a ...interface{}) {
	yl := color.New(color.FgYellow, color.BgBlack).PrintfFunc()
	yl(format, a...)
}

func Red(format string, a ...interface{}) {
	color.Red(format, a...)
	log.Printf(format, a...)
}

func Purple(format string, a ...interface{}) {
	color.HiCyan(format, a...)
}

func Teal(format string, a ...interface{}) {
	color.HiBlue(format, a...)
}

func HiGreen(format string, a ...interface{}) {
	color.HiGreen(format, a...)
}
