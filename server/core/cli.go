package core

import (
	"bufio"
	"log"
	"os"
	"strings"
)

var (
	unsafeLog = log.New(os.Stdout, "DONT ", log.LstdFlags)
)

func CliControl() {
	in := bufio.NewReader(os.Stdin)
	for {
		line, _, _ := in.ReadLine()
		if strings.ContainsRune(string(line), 'r') {
			unsafeLog.Println("Initiating unsafe reload of the constant world.")
			world.LoadConstantWorld()
		}
		if strings.ContainsRune(string(line), 'q') {
			os.Exit(0)
		}
	}
}
