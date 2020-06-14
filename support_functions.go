package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

//#####################################################
// gets the array of files which should be archived by searchFileByFileSize(), source and destination folder
// returns true if successfully implemented symlinks, error if occoured
func proveSymLink(sourceFiles[] string, filenames []string, destFolder string)(bool, error){
	if whichOS == "windows"{
		destFolder = filepath.FromSlash(destFolder)		//only for windows
		message(1,"[proveSymLink]", "OS not supported")
		return false, nil
	}

	//due to files moved all into one dir, files could have same name, this eliminates the problem
	fileCount := 0

	for i, sourceFile := range sourceFiles{
		fi, err := os.Lstat(sourceFile)
		if err != nil{
			message(1,"[proveSymLink]")
			return false, err
		}
		if (os.ModeSymlink & fi.Mode()) != 0 {
			symlinkDestination, err := os.Readlink(sourceFile)
			if err != nil{
				message(1,"[proveSymLink]")
				return false, err
			}

			if  symlinkDestination == destFolder+strconv.Itoa(fileCount)+filenames[i]{
				message(0, "[proveSymLink]", sourceFile, destFolder+strconv.Itoa(fileCount)+filenames[i])
			}
			fileCount++
		}
	}
	message(0,"[proveSymLink]")
	return true, nil
}
//#####################################################


//#####################################################
func printHelp(){
	fmt.Println("#**************************************#\n"+
		"MST (Move & Symlink Tool) \n"+
		"2020 Bastian Buck (HS-AlbSig) \n"+
		"!No liability is assumed!\n"+
		"#**************************************#\n"+
		"How it works:\n"+
		"-MST only runs with linux os\n"+
		"-put a external hdd into your server / computer and mount it\n"+
		"-MST moves bigger or older files to new hdd and makes symlink between them\n"+
		"-this is helpful, for files of download-sites.\n"+
		"-Attention! It is strongly discouraged to run with sudo rights\n"+
		"	or run with system / program files\n"+
		"#**************************************#\n"+
		"Arguments: 	-m -size -days -dest -src -help -h -log -a -depth\n"+
		"-m 			Select running mode 0 = file size (default); 1 = days; 2 = dir size\n"+
		"-size			Archive files by file size (default 20MB)\n"+
		"-days (beta)		Archive files by days (last modified) (default 60 days)\n"+
		"-dest (required)	Specify destination path (like: \"/var/www/newPath/\")\n"+
		"-src  (required)	Specify source path (like: \"/var/www/newPath/\")\n"+
		"-help -h		see this help\n"+
		"-log			turn logging on (default = false)\n"+
		"-a			analyze all files that could be archived (default = false)\n"+
		"-depth			only in combination with -m 2, depth to search dir's (-1 = all) (default = 3)\n"+
		"#**************************************#\n"+
		"Status-Codes \n" +
		"0=Success; 1=Failure; 3=Modified; 4=User Interrupt; 9=Not implemented\n"+
		"#**************************************#")
}
//#####################################################



//#####################################################
// send special output messages, distinguishes between different modes (see help) and arguments
func message(mode int, messageToOutput ...string){
	if logVar {
		logfile, err := os.OpenFile(logname, os.O_APPEND|os.O_WRONLY, 0755)
		if err != nil {
			fmt.Printf("[MESSAGE ERROR]: %e",err)
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
					sendmsg = "Status-Code " + strconv.Itoa(mode) + ": " + mesg+"\n"
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
			fmt.Printf("[MESSAGE ERROR]: %e",err)
		}
	}else{
		switch mode{
		case 0:
			for i, mesg := range messageToOutput{
				if i==0{
					fmt.Print("Success! "+mesg+" ")
				}else{
					fmt.Print(mesg+" ")
				}
			}
			fmt.Println("")
		case 1:
			for i, mesg := range messageToOutput{
				if i==0{
					fmt.Print("FAILED! "+mesg+" ")
				}else{
					fmt.Print(mesg+" ")
				}
			}
			fmt.Println("")

		default:
			for i, mesg := range messageToOutput{
				if i==0{
					fmt.Printf("Status-Code %d: %s ",mode,mesg)
				}else{
					fmt.Print(mesg+" ")
				}
			}
			fmt.Println("")
		}
	}
}
//#####################################################


//#####################################################
// writes messages into log file if desired
func writeLogFile(file io.Writer, output string) error{

	writer := bufio.NewWriter(file)

	_, err := fmt.Fprintf(writer, "%v", output)
	err = writer.Flush()
	if err != nil{
		return err
	}
	return nil
}
//#####################################################