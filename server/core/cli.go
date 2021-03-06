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
		if strings.ContainsRune(string(line), 's') {
			unsafeLog.Println("Initiating unsafe save of the live world.")
			if err := world.Save(); err != nil {
				panic(err)
			}
			unsafeLog.Println("Done saving.")
		}
		if strings.ContainsRune(string(line), 'r') {
			unsafeLog.Println("Initiating unsafe reload of the constant world.")
			if err := world.ReloadMaps(); err != nil {
				panic(err)
			}
			unsafeLog.Println("Done reloading.")
		}
		if strings.ContainsRune(string(line), 'q') {
			os.Exit(0)
		}
	}
}
