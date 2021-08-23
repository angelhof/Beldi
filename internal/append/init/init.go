package main

import (
	"os"
	"time"

	"github.com/eniac/Beldi/pkg/beldilib"
)

func ClearAll() {
	beldilib.DeleteLambdaTables("append")
	beldilib.DeleteTable("append-local")
}

func main() {
	if len(os.Args) >= 2 {
		option := os.Args[1]
		if option == "clean" {
			ClearAll()
			return
		}
	}
	ClearAll()
	beldilib.WaitUntilAllDeleted([]string{
		"append", "append-log", "append-collector",
		"append-local", // For transactions
	})
	beldilib.CreateLambdaTables("append")
	beldilib.CreateMainTable("append-local")
	beldilib.WaitUntilActive("append-local")

	// TODO: Modify these to write a list
	time.Sleep(60 * time.Second)

	beldilib.Populate("append", "K", "", false)
}
