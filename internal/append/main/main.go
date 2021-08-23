package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"
)

var TXN = "DISABLE"

func SerializeSlice(slice []string) string {
	return strings.Join(slice, ",")
}

func DeserializeSlice(ser_slice string) []string {
	if ser_slice == "" {
		return []string{}
	}

	return strings.Split(ser_slice, ",")
}

func Handler(env *beldilib.Env) interface{} {
	// var slice []string = make([]string, 0)
	// len := 100
	// for i := 0; i < len; i++ {
	// 	slice = append(slice, "hi")
	// }

	// a := SerializeSlice(slice)
	// fmt.Printf("Serialized slice %s\n", a)

	start := time.Now()
	if TXN == "ENABLE" {
		// Start a transaction to read the list, append to it, and then write it
		beldilib.BeginTxn(env)

		start_read := time.Now()
		ok, item := beldilib.TPLRead(env, "append", "K", []string{"V"})
		if !ok {
			return false
		}
		fmt.Printf("DURATION TPLRead %s\n", time.Since(start_read))

		start_append := time.Now()
		var old_ser_slice string
		beldilib.CHECK(mapstructure.Decode(item["V"], &old_ser_slice))
		var slice []string = DeserializeSlice(old_ser_slice)

		fmt.Printf("Retrieved slice %s\n", SerializeSlice(slice))

		slice = append(slice, "hi")

		ser_slice := SerializeSlice(slice)
		fmt.Printf("New slice %s\n", ser_slice)
		fmt.Printf("DURATION Append %s\n", time.Since(start_append))

		start_write := time.Now()
		ok = beldilib.TPLWrite(env, "append", "K",
			aws.JSONValue{"V": ser_slice})
		if !ok {
			return false
		}
		fmt.Printf("DURATION TPLWrite %s\n", time.Since(start_write))

		beldilib.CommitTxn(env)
	} else {
		// Start a transaction to read the list, append to it, and then write it

		start_read := time.Now()
		item := beldilib.Read(env, "append", "K")
		fmt.Printf("DURATION TPLRead %s\n", time.Since(start_read))

		start_append := time.Now()
		var old_ser_slice string
		beldilib.CHECK(mapstructure.Decode(item, &old_ser_slice))
		var slice []string = DeserializeSlice(old_ser_slice)

		fmt.Printf("Retrieved slice %s\n", SerializeSlice(slice))

		slice = append(slice, "hi")

		ser_slice := SerializeSlice(slice)
		fmt.Printf("New slice %s\n", ser_slice)
		fmt.Printf("DURATION Append %s\n", time.Since(start_append))

		start_write := time.Now()
		beldilib.Write(env, "append", "K",
			map[expression.NameBuilder]expression.OperandBuilder{
				expression.Name("V"): expression.Value(ser_slice),
			})
		fmt.Printf("DURATION TPLWrite %s\n", time.Since(start_write))
	}

	fmt.Printf("DURATION Txn %s\n", time.Since(start))
	return 0
}

func main() {
	lambda.Start(beldilib.Wrapper(Handler))
}
