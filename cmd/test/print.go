package test

import "fmt"

func stringError(msg string) string {
	return stringWithSymbol("ERROR ❌", msg)
}

func stringOK(msg string) string {
	return stringWithSymbol("OK ✅", msg)
}

func stringDetail(msg string) string {
	return stringWithSymbol("", msg)
}

func stringWithSymbol(symbol, msg string) string {
	return fmt.Sprintf(
		"%7s %s\n",
		symbol,
		msg,
	)
}
