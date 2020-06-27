/*****************
*
* MST (Move & Symlink Tool)
* Bastian Buck
* HS-AlbSig
* 2020
*
******************/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

var logVar, saveLog bool //log variable is used in message-function, dir for archive whole dirs
var logname = "mstLog.log"

var systemLog = "systemLog.log"
var whichOS = runtime.GOOS

func main() {
	fmt.Println("\nWelcome to MST (Move & Symlink Tool)")
	fmt.Println("")

	var m int
	var size, days int64
	var dest, src, sysLogFileName string
	var analyze bool

	help := flag.Bool("help", false, "help")
	h := flag.Bool("h", false, "help")
	flag.IntVar(&m, "m", 0, "mode 0 = biggest file; 1 = days")
	flag.Int64Var(&size, "size", 0, "file size in megabyte")
	flag.Int64Var(&days, "days", 0, "date in days")
	flag.StringVar(&dest, "dest", "", "destination (/which/folder/)")
	flag.StringVar(&src, "src", "", "source (/which/folder/)")
	flag.BoolVar(&logVar, "log", false, "logging (bool)")
	flag.BoolVar(&saveLog, "save", false, "logging with current date (bool)")
	flag.BoolVar(&analyze, "a", false, "analyze all files or dir's and make output")
	flag.StringVar(&sysLogFileName, "filename", "systemLog.log", "other systemLog name for reset and remove")

	//new flags to reset last execution and remove all moved files and symlinks
	reset := flag.Bool("reset", false, "reset last move & symlink execution")
	remove := flag.Bool("remove", false, "remove all files from last move & symlink execution")

	flag.Parse()

	if logVar {
		//if log is set, create logfile
		file, err := os.Create(logname)
		if err != nil {
			log.Fatal(err)
		}
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}

	switch {
	case *help || *h:
		printHelp()
		return
	case *reset:
		err := runReset(sysLogFileName)
		if err != nil {
			log.Fatal(err)
		}
		return
	case *remove:
		err := runRemove(sysLogFileName)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	sysLogFile, err := os.Create(systemLog)
	if err != nil {
		log.Fatal(err)
	}
	err = sysLogFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	//rename logFile to filename + actual date
	if saveLog {
		defer func() {
			y, m, d := time.Now().Date()
			h, min, s := time.Now().Clock()
			saveLogName := fmt.Sprintf("systemLog%d-%d-%d-%d-%d-%d.log", y, m, d, h, min, s)
			err := os.Rename(systemLog, saveLogName)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	if whichOS == "linux" || whichOS == "windows" {

		//#####################################################
		// if analyze is set and arguments are passed, analyze which files can be archived, after that exit
		if analyze && m == 0 && size > 0 && src[len(src)-1:] == "/" && days == 0 {
			fmt.Println("Analyze the situation, this will take few minutes...")
			// function exists, therefore only print array
			size *= 1000000 //megabyte to byte
			_, _, err := searchFileByFileSize(size, src)
			if err != nil {
				message(1, "Error while analyzing situation", "Exit")
				fmt.Println(err)
				return
			}
			message(0, "Files successfully analyzed\n")
			return
		} else if analyze && m == 1 && src[len(src)-1:] == "/" && days > 0 {
			fmt.Println("Analyze the situation, this will take few minutes...")
			// function exists, therefore only print array
			size *= 1000000 //megabyte to byte
			_, _, err := searchFileByDays(days, src)
			if err != nil {
				message(1, "Error while analyzing situation", "Exit")
				fmt.Println(err)
				return
			}
			message(0, "Files successfully analyzed\n")
			return
		} else if analyze {
			printHelp()
			fmt.Println("[ERROR] Bad arguments!")
			message(1, "")
			return
		}
		//#####################################################

		//#####################################################
		//check correct src and dest arguments
		if src == "" || dest == "" {
			printHelp()
			fmt.Println("[ERROR] Missing arguments!")
			message(1, "")
			return
		} else if src == dest {
			printHelp()
			fmt.Println("[ERROR] Source can not be same as Destination!")
			message(1, "")
			return
		}
		//#####################################################

		//#####################################################
		var accept string
		//check arguments and partition into modes
		if m == 0 && size > 0 && dest[len(dest)-1:] == "/" && src[len(src)-1:] == "/" && dest[0] == '/' && src[0] == '/' && analyze == false {
			if size == 0 {
				message(3, "[Modified] default parameter size=20")
				size = 20
			}

			fmt.Printf("Program will now analyze the situation with arguments: mode=%d, size=%d, src=%s, dest=%s log=%t\nconfirm (y/n)\n\n", m, size, src, dest, logVar)
			_, err := fmt.Scanln(&accept)
			if accept == "y" && err == nil {
				size *= 1000000 //megabyte to byte
				runByFileSize(size, dest, src)
			} else if accept == "n" && err == nil {
				message(4, "User Interrupt")
			} else {
				fmt.Println(err)
				return
			}
		} else if m == 1 && days > 0 && dest[len(dest)-1:] == "/" && src[len(src)-1:] == "/" && dest[0] == '/' && src[0] == '/' && analyze == false {
			if days == 0 {
				message(3, "[Modified] default parameter days=60")
				days = 60
			}

			fmt.Printf("Program will now analyze the situation with arguments: mode=%d, days=%d, src=%s, dest=%s log=%t\nconfirm (y/n)\n\n", m, days, src, dest, logVar)
			_, err := fmt.Scanln(&accept)
			if accept == "y" && err == nil {
				runByDays(days, dest, src)
			} else if accept == "n" && err == nil {
				message(4, "User Interrupt")
			} else {
				fmt.Println(err)
				return
			}
			//if arguments check failed, print to user
		} else {
			printHelp()
			fmt.Println("[ERROR] Bad arguments!")
			message(1, "")
			return
		}
		//#####################################################
	} else {
		//if OS is not supported
		fmt.Printf("OS %s not supported!", whichOS)
		return
	}
}
