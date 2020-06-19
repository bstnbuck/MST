package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func runByDays(days int64, dest string, source string) {

	if _, err := os.Stat(dest); os.IsNotExist(err) {
		message(1, "Destination does not exist!")
		fmt.Println(err)
		return
	}

	//analyze which files should be archived
	filesToArchive, filenames, err := searchFileByDays(days, source)
	if err != nil {
		fmt.Println(err)
		message(1, "")
		return
	}

	randPrefix := GenerateRandomString()

	var accept string
	fmt.Printf("Archiving will now start and cannot be cancelled!\nconfirm (y/n)\n\n")
	_, err = fmt.Scanln(&accept)
	if accept == "y" && err == nil {
		//archive this files
		err = archive(filesToArchive, filenames, dest, randPrefix)
		if err != nil {
			fmt.Println(err)
			message(1, "")
			return
		}
		//proof if symlinks successfully set and files archived
		isArchived, err := proveSymLink(filesToArchive, filenames, dest, randPrefix)
		if isArchived && err == nil {
			fmt.Println("[runByDays] Succesfully archived!")
		} else {
			fmt.Println("[ERROR - runByDays] While proving symlinks")
			fmt.Println(err)
			message(1, "")
			return
		}
		//if all done, make success
		message(0, "")
	} else if accept == "n" && err == nil {
		message(4, "User Interrupt")
	} else {
		log.Fatal(err)
	}
}

//#####################################################
// analyze which files should be archived, gets size of files should be archived and source folder
// returns a string array of filenames and error if occoured
func searchFileByDays(days int64, source string) ([]string, []string, error) {
	fmt.Printf("[searchFileByFileSize] days: %d; source: %s\n", days, source)
	var filesToArchiveWithPath []string
	var fileNames []string
	if whichOS == "windows" {
		source = filepath.FromSlash(source) //only for windows
	}
	var totalFileSize int64
	filesCount := 0

	if _, err := os.Stat(source); os.IsNotExist(err) {
		message(1, "Source does not exist!")
		return nil, nil, err
	}

	daysToTime := time.Now().AddDate(0, 0, -int(days))
	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		fi, err := os.Lstat(path)
		if err != nil {
			message(1, "[searchFileByFileSize]")
			return err
		}
		if (os.ModeSymlink&fi.Mode()) != 0 || fi.Mode().IsDir() { //check if symlink or directory
			return nil
		} else if info.ModTime().Before(daysToTime) {
			filesToArchiveWithPath = append(filesToArchiveWithPath, path)
			fileNames = append(fileNames, info.Name())
			totalFileSize += info.Size() / 1000000
			filesCount++
			message(0, "[searchFileByFileSize]", "File: \""+path+"\" with size: "+strconv.FormatInt(info.Size()/1000000, 10)+" MB")
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	if filesCount == 0 {
		return nil, nil, errors.New("[searchFileByFileSize]: Nothing to move")
	}
	message(0, "[searchFileByDays]", "Successfully analyzed!", "Total File Size: "+strconv.FormatInt(totalFileSize, 10)+" MB")
	return filesToArchiveWithPath, fileNames, nil
}

//#####################################################
