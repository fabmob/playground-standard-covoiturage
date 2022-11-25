package util

// PanicIf panics if there is an error. Use for testing purposes only.
func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
