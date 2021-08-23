package main

import (
	"os"
	"time"

	"github.com/eniac/Beldi/pkg/beldilib"
)

func ClearAll() {
	beldilib.DeleteLambdaTables("append")
	beldilib.DeleteLambdaTables("tappend")
	beldilib.DeleteTable("tappend-local")
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
		"tappend", "tappend-log", "tappend-collector",
		"tappend-local", // For transactions
	})
	beldilib.CreateLambdaTables("append")
	beldilib.CreateLambdaTables("tappend")
	beldilib.CreateMainTable("tappend-local")
	beldilib.WaitUntilActive("tappend-local")

	// TODO: Modify these to write a list
	time.Sleep(60 * time.Second)

	beldilib.Populate("append", "K", "", false)
	beldilib.Populate("tappend", "K", "", false)
}
