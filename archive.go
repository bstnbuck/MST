package main

import (
	"os"
	"path/filepath"
	"strconv"
)

// gets the array of files which should be archived by searchFileByFileSize() and destFolder
func archive(sourceFiles []string, filenames []string, destFolder string)error{
	if whichOS == "windows"{
		destFolder = filepath.FromSlash(destFolder)		//only for windows
		message(1, "[archive]", "OS not supported")
		return nil
	}

	fileCount := 0

	for i, sourceFile := range sourceFiles{
		err := os.Rename(sourceFile, destFolder+strconv.Itoa(fileCount)+filenames[i])
		if err != nil{
			return err
		}

		// only linux using /
		err = os.Symlink(destFolder+strconv.Itoa(fileCount)+filenames[i], sourceFile)
		if err != nil{
			return err
		}

		message(0, "[archive]", sourceFile, destFolder+strconv.Itoa(fileCount)+filenames[i])
		fileCount++
	}
	message(0, "[archive]")
	return nil
}