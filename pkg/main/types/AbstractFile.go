package types

type AbstractFile interface {
	FileType() string
	GetFileName() string
	SetFileName(filename string)
}
