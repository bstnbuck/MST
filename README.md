[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/bstnbuck/Simple-Go-Blockchain/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/bstnbuck/MST)](https://goreportcard.com/report/github.com/bstnbuck/MST)
[![Build Status](https://travis-ci.org/bstnbuck/MST.svg?branch=master)](https://travis-ci.org/bstnbuck/MST)
# MST (Move & Symlink Tool)

MST searches for large or old files and moves them to another drive (e.g. external hard disk) to avoid filling up your own

>**Not all functions are implemented yet! Right now MST is in an alpha stage.
  Attention! It is strongly discouraged to run with sudo rights	or run with system / program files**

**Please report unchecked bugs via Issue!**

## Requirements
* Go compiler (v1.14.4)

## Installation
`git clone https://github.com/bstnbuck/MST.git`

### To build the program run following commands:
#### Linux:    
  * cd mst
  * go build .
  * ./mst + flags

#### Windows:  
  * move into program directory
  * go build .
  * Using CMD: **mst.exe + flags** (only analyze is supported during now)

## Usage:
#### How it works:
- put a external HDD into your server / computer and mount it
- MST moves bigger or older files to new HDD and makes a symbolic link between them
- this is helpful, for files of download-sites.
***
##### Arguments: 	**-m -size -days -dest -src -help -h -log -a -depth**
- -m -> Select running mode 0 = file size (default); 1 = days; 2 = dir size (not implemented yet)
- -size -> Archive files by file size (default 20MB)
- -days -> Archive files by days (last modified) (default 60 days)
- -dest (required) -> Specify destination path (like: "/var/www/newPath/")
- -src (required) -> Specify source path (like: "/var/www/newPath/")
- -help -h -> see this help
- -log	-> turn logging on (default = false)
- -save	-> save system log file with actual date(default = false)
- -a -> analyze all files that could be archived (default = false)
- -reset	-> reset all changes of last run, optional with -log
- -remove -> remove all changes of last run, optional with -log
- -filename	-> choose other systemLog file for -reset or -remove
***
##### Status-Codes
0=Success; 1=Failure; 2=Info; 3=Modified; 4=User Interrupt; 9=Not implemented
***

### Examples
- -m 1 -days 360 -src "/test/testdrive/" -log -a (MST analyze within elapsed 360 days and search in directory and log all commands)
- (-m 0) -size 2 -src "/test/testDrive/" -dest "/test/testPaste/" -log (MST with filesize = 2 MB and logging)
- -h (or -help) (prints help)
- -reset -filename "systemLog2020-6-27-12-11-19.log" -log


### Information
- MST only runs with Linux OS

### The following is still being implemented
* Maybe Windows support.
