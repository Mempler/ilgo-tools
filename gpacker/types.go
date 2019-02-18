package gpacker

type entrytype byte

const (
	TBinary entrytype = 0x01 + iota
	TText
	TImage
	TFont
)
