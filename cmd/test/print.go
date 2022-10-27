package test

import "fmt"

func printError(msg string) {
	printWithSymbol("ERROR ❌", msg)
}

func printOK(msg string) {
	printWithSymbol("OK ✅", msg)
}

func printDetail(msg string) {
	printWithSymbol("", msg)
}

func printWithSymbol(symbol, msg string) {
	fmt.Printf(
		"%7s %s",
		symbol,
		msg,
	)
}
