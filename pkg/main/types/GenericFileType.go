package types

import "os"

type GenericFileType struct {
	AbstractFile
	FileName string
	File     os.File
}

func (p GenericFileType) FileType() string {
	return "*"
}

func (p GenericFileType) SetFileName(fileName string) {
	p.FileName = fileName
}

func (p GenericFileType) GetFileName() string {
	return p.FileName
}
