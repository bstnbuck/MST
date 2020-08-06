package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

// gets the array of files which should be archived by searchFileByFileSize() and destFolder
func archive(sourceFiles []string, filenames []string, destFolder string, randPrefix string) error {
	if whichOS == "windows" {
		destFolder = filepath.FromSlash(destFolder) //only for windows
		message(1, "[archive]", "OS not supported")
		return nil
	}

	fileCount := 0
	for i, sourceFile := range sourceFiles {

		// Because moving files with os.Rename to other hard-drives is not allowed, files must be created by new
		// os.Rename only works on the same hard-drive
		/*err := os.Rename(sourceFile, destFolder+randPrefix+strconv.Itoa(fileCount)+filenames[i])
		if err != nil {
			return err
		}
		*/

		//open the sourceFile
		inputFile, err := os.Open(sourceFile)
		if err != nil {
			return fmt.Errorf("couldn't open source file: %s", err)
		}

		//create the new file
		outputFile, err := os.Create(destFolder + randPrefix + strconv.Itoa(fileCount) + filenames[i])
		if err != nil {
			err = inputFile.Close()
			return fmt.Errorf("couldn't open dest file: %s", err)
		}
		defer func() {
			err = outputFile.Close()
		}()

		//copy the content from source to destination
		_, err = io.Copy(outputFile, inputFile)
		if err != nil {
			return fmt.Errorf("writing to output file failed: %s", err)
		}
		err = inputFile.Close()
		if err != nil {
			return fmt.Errorf("closing inputFile failed: %s", err)
		}
		// The copy was successful, so now delete the original file
		err = os.Remove(sourceFile)
		if err != nil {
			return fmt.Errorf("failed removing original file: %s", err)
		}

		err = sysLogging(sourceFile, destFolder+randPrefix+strconv.Itoa(fileCount)+filenames[i])
		if err != nil {
			return err
		}

		// only linux using
		//new output name is, something random + file-id + file-name, this prevent matching file-names
		err = os.Symlink(destFolder+randPrefix+strconv.Itoa(fileCount)+filenames[i], sourceFile)
		if err != nil {
			return err
		}

		message(0, "[archive]", sourceFile, destFolder+randPrefix+strconv.Itoa(fileCount)+filenames[i])
		fileCount++
	}
	message(0, "[archive]")
	return nil
}
