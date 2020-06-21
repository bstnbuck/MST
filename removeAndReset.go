package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// reset last program execution
func runReset() error {
	if whichOS == "linux" {
		var accept string
		fmt.Printf("This will reset all made changes from last run. Run cannot be cancelled!\nconfirm (y/n)\n\n")
		_, err := fmt.Scanln(&accept)
		if accept == "y" && err == nil {
			sysLogFile, err := os.OpenFile(systemLog, os.O_RDWR, 0755)
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
				err = os.Rename(dest, source)
				if err != nil {
					return err
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
func runRemove() error {

	if whichOS == "linux" {
		var accept string
		fmt.Printf("Remove will now start, cannot be cancelled and undone!\nconfirm (y/n)\n\n")
		_, err := fmt.Scanln(&accept)
		if accept == "y" && err == nil {
			sysLogFile, err := os.OpenFile(systemLog, os.O_RDWR, 0755)
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
