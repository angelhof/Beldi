package main

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/beldilib"
)

func ClearAll() {
	beldilib.DeleteLambdaTables("append")
	beldilib.DeleteLambdaTables("nop")
	beldilib.DeleteTable("bappend")
	beldilib.DeleteTable("bnop")
	beldilib.DeleteLambdaTables("tappend")
	beldilib.DeleteLambdaTables("tnop")
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
		"nop", "nop-log", "nop-collector",
		"append", "tappend-log", "tappend-collector",
		"tnop", "tnop-log", "tnop-collector",
		"append", "bnop",
	})
	beldilib.CreateLambdaTables("append")
	beldilib.CreateLambdaTables("nop")

	beldilib.CreateBaselineTable("bappend")
	beldilib.CreateBaselineTable("bnop")

	beldilib.CreateTxnTables("tappend")
	beldilib.CreateTxnTables("tnop")

	// TODO: Modify these to write a list
	time.Sleep(60 * time.Second)
	beldilib.WriteNRows("append", "K", 20)

	beldilib.LibWrite("bappend", aws.JSONValue{"K": "K"}, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(1),
	})

	beldilib.LibWrite("tappend", aws.JSONValue{"K": "K"}, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(1),
	})
}
