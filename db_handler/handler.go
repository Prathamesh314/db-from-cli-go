package db_handler

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)


type DBHandler struct{
	Handler any
}

func NewDBHandler(handler any) *DBHandler {
	return &DBHandler{Handler: handler}
}


func PrettyPrint(data []map[string]interface{}) {
	if len(data) == 0 {
		fmt.Println("No data to display")
		return
	}

	maxKeyLength := 0
	for _, row := range data {
		for key := range row {
			if len(key) > maxKeyLength {
				maxKeyLength = len(key)
			}
		}
	}

	for i, row := range data {
		fmt.Printf("Row %d:\n", i+1)
		printRow(row, maxKeyLength, 1)
		fmt.Println(strings.Repeat("-", 50))
	}
}


func printRow(row map[string]interface{}, maxKeyLength, indentLevel int) {
	indent := strings.Repeat("  ", indentLevel)

	for key, value := range row {
		switch v := value.(type) {
		case map[string]interface{}:
			
			fmt.Printf("%s%-*s :\n", indent, maxKeyLength, key)
			printRow(v, maxKeyLength, indentLevel+1) 
		default:
			
			fmt.Printf("%s%-*s : %v\n", indent, maxKeyLength, key, value)
		}
	}
}

func ClearTerminal() {
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println("Unsupported platform")
	}
}

