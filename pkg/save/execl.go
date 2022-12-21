package save

import (
	"fmt"
	"strconv"

	"github.com/tealeg/xlsx"
	//"os/exec"
	//"path/filepath"
	//"os"
	//"strings"
)

func SaveExcel(data []interface{}, fileName string) error {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	//var fileName string

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}

	for k, v := range data {
		names, boolean := v.([]string)
		if boolean {
			row = sheet.AddRow()
			cell = row.AddCell()
			cell.Value = strconv.Itoa(k + 1)
			for _, v2 := range names {
				cell = row.AddCell()
				cell.Value = v2
			}
		}
	}

	//runFile, _ := exec.LookPath(os.Args[0])
	//runPath, _ := filepath.Abs(runFile)
	//index := strings.LastIndex(runPath, string(os.PathSeparator))
	//runPathDir := runPath[:index]

	//t := time.Now()
	//fileName = runPathDir + string(os.PathSeparator) + strconv.Itoa(t.Nanosecond()) + ".xlsx"
	//fileName = strconv.Itoa(t.Nanosecond()) + ".xlsx"

	err = file.Save(fileName)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	return nil
}
