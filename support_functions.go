package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

//#####################################################
// gets the array of files which should be archived by searchFileByFileSize(), source and destination folder
// returns true if successfully implemented symlinks, error if occoured
func proveSymLink(sourceFiles []string, filenames []string, destFolder string, randPrefix string) (bool, error) {
	if whichOS == "windows" {
		destFolder = filepath.FromSlash(destFolder) //only for windows
		//message(1, "[proveSymLink]", "OS not supported")
		//return false, nil
	}

	//due to files moved all into one dir, files could have same name, this eliminates the problem
	fileCount := 0

	for i, sourceFile := range sourceFiles {
		fi, err := os.Lstat(sourceFile)
		if err != nil {
			message(1, "[proveSymLink]")
			return false, err
		}
		if (os.ModeSymlink & fi.Mode()) != 0 {
			symlinkDestination, err := os.Readlink(sourceFile)
			if err != nil {
				message(1, "[proveSymLink]")
				return false, err
			}

			if symlinkDestination == destFolder+randPrefix+strconv.Itoa(fileCount)+filenames[i] {
				message(0, "[proveSymLink]", sourceFile, destFolder+randPrefix+strconv.Itoa(fileCount)+filenames[i])
			} else {
				message(1, "[proveSymLink]", "Proving Symlink failed!")
				return false, errors.New("[proveSymLink] Failed proving Symlink")
			}
			fileCount++
		}
	}
	message(0, "[proveSymLink]")
	return true, nil
}

//#####################################################

//#####################################################
func printHelp() {
	fmt.Println("#**************************************#\n" +
		"MST (Move & Symlink Tool) \n" +
		"2020 Bastian Buck (HS-AlbSig) \n" +
		"!No liability is assumed!\n" +
		"#**************************************#\n" +
		"How it works:\n" +
		"-MST only runs with linux and windows os\n" +
		"-put a external hdd into your server / computer and mount it\n" +
		"-MST moves bigger or older files to new hdd and makes symlink between them\n" +
		"-this is helpful, for files of download-sites.\n" +
		"-Attention! It is strongly discouraged to run with sudo rights\n" +
		"	or run with system / program files\n" +
		"#**************************************#\n" +
		"Arguments: 	-m -size -days -dest -src -help -h -log -a -depth\n" +
		"-m 			Select running mode 0 = file size (default); 1 = days\n" +
		"-size			Archive files by file size (default 20MB)\n" +
		"-days (beta)		Archive files by days (last modified) (default 60 days)\n" +
		"-dest (required)	Specify destination path (like: \"/var/www/newPath/\")\n" +
		"-src  (required)	Specify source path (like: \"/var/www/newPath/\")\n" +
		"-help -h		see this help\n" +
		"-log			turn logging on (default = false)\n" +
		"-save			save system log file with actual date(default = false)\n" +
		"-a				analyze all files that could be archived (default = false)\n" +
		"-reset			reset all changes of last run, optional with -log\n" +
		"-remove			remove all changes of last run, optional with -log\n" +
		"-filename		choose other systemLog file for -reset or -remove\n" +
		"#**************************************#\n" +
		"Status-Codes \n" +
		"0=Success; 1=Failure; 2=Info; 3=Modified; 4=User Interrupt; 9=Not implemented\n" +
		"#**************************************#")
}

//#####################################################

//#####################################################
// send special output messages, distinguishes between different modes (see help) and arguments
func message(mode int, messageToOutput ...string) {
	if logVar {
		logfile, err := os.OpenFile(logname, os.O_APPEND|os.O_WRONLY, 0755)
		if err != nil {
			fmt.Printf("[MESSAGE ERROR]: %e", err)
		}
		var sendmsg string

		switch mode {
		case 0:
			for i, mesg := range messageToOutput {
				switch {
				case len(messageToOutput) == 1:
					sendmsg = "Success! " + mesg + "\n"
					fmt.Print(sendmsg)
					err = writeLogFile(logfile, sendmsg)
					if err != nil {
						fmt.Printf("[MESSAGE ERROR]: %e", err)
					}
				case i == 0:
					sendmsg = "Success! " + mesg + " "
					fmt.Print(sendmsg)
					err = writeLogFile(logfile, sendmsg)
					if err != nil {
						fmt.Printf("[MESSAGE ERROR]: %e", err)
					}
				case i == len(messageToOutput)-1:
					sendmsg = mesg + "\n"
					fmt.Print(sendmsg)
					err = writeLogFile(logfile, sendmsg)
					if err != nil {
						fmt.Printf("[MESSAGE ERROR]: %e", err)
					}
				default:
					sendmsg = mesg + " "
					fmt.Print(sendmsg)
					err = writeLogFile(logfile, sendmsg)
					if err != nil {
						fmt.Printf("[MESSAGE ERROR]: %e", err)
					}
				}
			}
		case 1:
			for i, mesg := range messageToOutput {
				switch {
				case len(messageToOutput) == 1:
					sendmsg = "FAILED! " + mesg + "\n"
					fmt.Print(sendmsg)
					err = writeLogFile(logfile, sendmsg)
					if err != nil {
						fmt.Printf("[MESSAGE ERROR]: %e", err)
					}
				case i == 0:
					sendmsg = "FAILED! " + mesg + " "
					fmt.Print(sendmsg)
					err = writeLogFile(logfile, sendmsg)
					if err != nil {
						fmt.Printf("[MESSAGE ERROR]: %e", err)
					}
				case i == len(messageToOutput)-1:
					sendmsg = mesg + "\n"
					fmt.Print(sendmsg)
					err = writeLogFile(logfile, sendmsg)
					if err != nil {
						fmt.Printf("[MESSAGE ERROR]: %e", err)
					}
				default:
					sendmsg = mesg + " "
					fmt.Print(sendmsg)
					err = writeLogFile(logfile, sendmsg)
					if err != nil {
						fmt.Printf("[MESSAGE ERROR]: %e", err)
					}
				}
			}
		default:
			for i, mesg := range messageToOutput {
				switch {
				case len(messageToOutput) == 1:
					sendmsg = "Status-Code " + strconv.Itoa(mode) + ": " + mesg + "\n"
					fmt.Print(sendmsg)
					err = writeLogFile(logfile, sendmsg)
					if err != nil {
						fmt.Printf("[MESSAGE ERROR]: %e", err)
					}
				case i == 0:
					sendmsg = "Status-Code " + strconv.Itoa(mode) + ": " + mesg
					fmt.Print(sendmsg)
					err = writeLogFile(logfile, sendmsg)
					if err != nil {
						fmt.Printf("[MESSAGE ERROR]: %e", err)
					}
				case i == len(messageToOutput)-1:
					sendmsg = mesg + "\n"
					fmt.Print(sendmsg)
					err = writeLogFile(logfile, sendmsg)
					if err != nil {
						fmt.Printf("[MESSAGE ERROR]: %e", err)
					}
				default:
					sendmsg = mesg + " "
					fmt.Print(sendmsg)
					err = writeLogFile(logfile, sendmsg)
					if err != nil {
						fmt.Printf("[MESSAGE ERROR]: %e", err)
					}
				}
			}
		}
		err = logfile.Close()
		if err != nil {
			fmt.Printf("[MESSAGE ERROR]: %e", err)
		}
	} else {
		switch mode {
		case 0:
			for i, mesg := range messageToOutput {
				if i == 0 {
					fmt.Print("Success! " + mesg + " ")
				} else {
					fmt.Print(mesg + " ")
				}
			}
			fmt.Println("")
		case 1:
			for i, mesg := range messageToOutput {
				if i == 0 {
					fmt.Print("FAILED! " + mesg + " ")
				} else {
					fmt.Print(mesg + " ")
				}
			}
			fmt.Println("")

		default:
			for i, mesg := range messageToOutput {
				if i == 0 {
					fmt.Printf("Status-Code %d: %s ", mode, mesg)
				} else {
					fmt.Print(mesg + " ")
				}
			}
			fmt.Println("")
		}
	}
}

//#####################################################

//logs all changes to file which is easy readable for program
func sysLogging(source string, dest string) error {

	sysLogFile, err := os.OpenFile(systemLog, os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(sysLogFile)

	_, err = fmt.Fprintf(writer, "%s<separator>%s\n", source, dest)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	err = sysLogFile.Close()
	if err != nil {
		return err
	}
	return nil
}

//#####################################################
// writes messages into log file if desired
func writeLogFile(file io.Writer, output string) error {

	writer := bufio.NewWriter(file)

	_, err := fmt.Fprintf(writer, "%v", output)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

//#####################################################

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomString() string {
	b := GenerateRandomBytes(3)
	return base64.URLEncoding.EncodeToString(b) //encode random byte array to base64 encoding
}

// make random bytes and return them
func GenerateRandomBytes(n int) []byte {
	b := make([]byte, n)   //new byte array of length n
	_, err := rand.Read(b) //fill array with random
	if err != nil {        //if error print
		println(err)
	}
	return b //return array
}
