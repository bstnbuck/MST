package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// reset last program execution
func runReset(sysLogFileName string) error {
	if whichOS == "linux" {
		var accept string
		fmt.Printf("This will reset all made changes from last run. Run cannot be cancelled!\nconfirm (y/n)\n\n")
		_, err := fmt.Scanln(&accept)
		if accept == "y" && err == nil {
			sysLogFile, err := os.OpenFile(sysLogFileName, os.O_RDWR, 0755)
			if err != nil {
				return err
			}
			scanner := bufio.NewScanner(sysLogFile)

			for scanner.Scan() {
				line := scanner.Text()

				s := strings.Split(line, ":")
				source, dest := s[0], s[1]

				err = os.Remove(source)
				if err != nil {
					return err
				}

				// Because moving files with os.Rename to other hard-drives is not allowed, files must be created by new
				// os.Rename only works on the same hard-drive
				/*err = os.Rename(dest, source)
				if err != nil {
					return err
				}*/
				//open the sourceFile
				inputFile, err := os.Open(source)
				if err != nil {
					return fmt.Errorf("couldn't open source file: %s", err)
				}

				//create the new file
				outputFile, err := os.Create(dest)
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
				err = os.Remove(source)
				if err != nil {
					return fmt.Errorf("failed removing original file: %s", err)
				}
				message(0, "[runReset] Successfully reset changes to ", dest)
			}
			err = sysLogFile.Truncate(0)
			if err != nil {
				return err
			}
			_, err = sysLogFile.Seek(0, 0)
			if err != nil {
				return err
			}
			err = sysLogFile.Close()
			if err != nil {
				return err
			}
			message(0, "[runReset]")
		}
	} else {
		fmt.Printf("OS %s not supported!", whichOS)
		return nil
	}
	return nil
}

// remove all made changes
func runRemove(sysLogFileName string) error {

	if whichOS == "linux" {
		var accept string
		fmt.Printf("Remove will now start, cannot be cancelled and undone!\nconfirm (y/n)\n\n")
		_, err := fmt.Scanln(&accept)
		if accept == "y" && err == nil {
			sysLogFile, err := os.OpenFile(sysLogFileName, os.O_RDWR, 0755)
			if err != nil {
				return err
			}

			scanner := bufio.NewScanner(sysLogFile)

			for scanner.Scan() {
				line := scanner.Text()
				s := strings.Split(line, ":")
				source, dest := s[0], s[1]

				err = os.Remove(source)
				if err != nil {
					return err
				}
				err = os.Remove(dest)
				if err != nil {
					return err
				}
				message(0, "[runRemove] Successfully removed", source, dest)
			}
			err = sysLogFile.Truncate(0)
			if err != nil {
				return err
			}
			_, err = sysLogFile.Seek(0, 0)
			if err != nil {
				return err
			}
			err = sysLogFile.Close()
			if err != nil {
				return err
			}
			message(0, "[runRemove]")
		}
	} else {
		fmt.Printf("OS %s not supported! ", whichOS)
		return nil
	}
	return nil
}
