package main

import "fmt"

func warn(format string, v ...interface{}) {
	fmt.Printf("\033[0;33m"+format+"\033[0m", v...)
}
func info(format string, v ...interface{}) {
	fmt.Printf("\033[0;32m"+format+"\033[0m", v...)
}
func fatal(format string, v ...interface{}) {
	fmt.Printf("\033[0;31m"+format+"\033[0m", v...)
}
