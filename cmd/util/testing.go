// Package util exports some utility functions, that have no dependency to
// other packages of this module, and that may be used in other
// packages.
package util

// PanicIf panics if there is an error. Use for testing purposes only.
func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
