package types

import "github.com/tealeg/xlsx/v3"

type XlsxFileType struct {
	AbstractFile
	File     *xlsx.File
	FileName string
}

func (p XlsxFileType) FileType() string {
	return "xlsx"
}

func (p XlsxFileType) SetFileName(fileName string) {
	p.FileName = fileName
}

func (p XlsxFileType) GetFileName() string {
	return p.FileName
}
