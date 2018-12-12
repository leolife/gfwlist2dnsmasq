package cmd

import (
	"os"
	"path/filepath"
)

func generate(list map[string]bool) error {

	outputFilePath := "gfwlist.txt"
	if outputFlag != "" {
		outputFilePath = outputFlag
	}
	outputFilePath, _ = filepath.Abs(outputFilePath)

	f, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	for i := range list {
		f.WriteString("server=/." + i)
		if portFlag != "" {
			f.WriteString("#" + portFlag)
		}
		f.WriteString("\n")
		if ipsetFlag != "" {
			f.WriteString("ipset=/." + i + "/" + ipsetFlag + "\n")
		}
	}

	return nil
}
