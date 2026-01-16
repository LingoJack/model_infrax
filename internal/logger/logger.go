package logger

import (
	"fmt"
	"log"

	"github.com/fatih/color"
)

func Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func Info(args ...interface{}) {
	log.Println(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func Error(args ...interface{}) {
	log.Println(args...)
}

func JustPrintf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// ColorPrintf 输出彩色日志
// 使用示例: logger.ColorPrintf(logger.ColorRed, "错误信息: %s", err)
// 或使用预定义的颜色: logger.ColorPrintf(logger.ColorGreen, "成功: %s", msg)
func ColorPrintf(colorAttr color.Attribute, format string, args ...interface{}) {
	color.Set(colorAttr)
	defer color.Unset()
	fmt.Printf(format, args...)
}

const (
	ColorRed     = color.FgRed
	ColorGreen   = color.FgGreen
	ColorYellow  = color.FgYellow
	ColorBlue    = color.FgBlue
	ColorMagenta = color.FgMagenta
	ColorCyan    = color.FgCyan
	ColorWhite   = color.FgWhite

	ColorHiRed     = color.FgHiRed
	ColorHiGreen   = color.FgHiGreen
	ColorHiYellow  = color.FgHiYellow
	ColorHiBlue    = color.FgHiBlue
	ColorHiMagenta = color.FgHiMagenta
	ColorHiCyan    = color.FgHiCyan
	ColorHiWhite   = color.FgHiWhite
)
