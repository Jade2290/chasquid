// Fuzz testing for package normalize.

//go:build gofuzz
// +build gofuzz

package normalize

func Fuzz(data []byte) int {
	s := string(data)
	User(s)
	Domain(s)
	Addr(s)
	DomainToUnicode(s)

	return 0
}
