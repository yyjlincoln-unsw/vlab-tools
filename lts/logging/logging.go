package logging

import (
	"fmt"

	"github.com/fatih/color"
)

func Errorf(format string, v ...interface{}) {
	if len(format) == 0 {
		return
	}
	hasTrailingNewLine := false

	if format[len(format)-1] == '\n' {
		format = format[:len(format)-1]
		hasTrailingNewLine = true
	}

	coloredErrorString := color.RedString(fmt.Sprintf(format, v...))

	if hasTrailingNewLine {
		fmt.Printf("%v\n", coloredErrorString)
	} else {
		fmt.Printf("%v", coloredErrorString)
	}
}

func Successf(format string, v ...interface{}) {
	if len(format) == 0 {
		return
	}
	hasTrailingNewLine := false

	if format[len(format)-1] == '\n' {
		format = format[:len(format)-1]
		hasTrailingNewLine = true
	}

	coloredErrorString := color.GreenString(fmt.Sprintf(format, v...))

	if hasTrailingNewLine {
		fmt.Printf("%v\n", coloredErrorString)
	} else {
		fmt.Printf("%v", coloredErrorString)
	}
}

func Warnf(format string, v ...interface{}) {
	if len(format) == 0 {
		return
	}
	hasTrailingNewLine := false

	if format[len(format)-1] == '\n' {
		format = format[:len(format)-1]
		hasTrailingNewLine = true
	}

	coloredErrorString := color.YellowString(fmt.Sprintf(format, v...))

	if hasTrailingNewLine {
		fmt.Printf("%v\n", coloredErrorString)
	} else {
		fmt.Printf("%v", coloredErrorString)
	}
}

func Infof(format string, v ...interface{}) {
	if len(format) == 0 {
		return
	}
	hasTrailingNewLine := false

	if format[len(format)-1] == '\n' {
		format = format[:len(format)-1]
		hasTrailingNewLine = true
	}

	coloredErrorString := color.BlueString(fmt.Sprintf(format, v...))

	if hasTrailingNewLine {
		fmt.Printf("%v\n", coloredErrorString)
	} else {
		fmt.Printf("%v", coloredErrorString)
	}
}
