package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/lithammer/shortuuid"
)

var TXN = "DISABLE"

func SerializeSlice(slice []string) string {
	return strings.Join(slice, ",")
}

func Handler(env *beldilib.Env) interface{} {
	if TXN == "ENABLE" {
		var slice []string = make([]string, 0)
		len := 10
		for i := 0; i < len; i++ {
			slice = append(slice, "hi")
		}

		a := SerializeSlice(slice)

		start := time.Now()
		beldilib.TWrite(env, "tappend", "K", a)
		fmt.Printf("DURATION DWrite %s\n", time.Since(start))

		start = time.Now()
		beldilib.TCondWrite(env, "tappend", "K", a, true)
		fmt.Printf("DURATION CWriteT %s\n", time.Since(start))

		start = time.Now()
		beldilib.TCondWrite(env, "tappend", "K", a, false)
		fmt.Printf("DURATION CWriteF %s\n", time.Since(start))

		start = time.Now()
		beldilib.TRead(env, "tappend", "K")
		fmt.Printf("DURATION Read %s\n", time.Since(start))
		return 0
	}
	if beldilib.TYPE == "BELDI" {
		a := shortuuid.New()
		start := time.Now()
		beldilib.Write(env, "append", "K", map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V"): expression.Value(a),
		})
		fmt.Printf("DURATION DWrite %s\n", time.Since(start))

		start = time.Now()
		beldilib.CondWrite(env, "append", "K", map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V2"): expression.Value(1),
		}, expression.Name("V").Equal(expression.Value(a)))
		fmt.Printf("DURATION CWriteT %s\n", time.Since(start))

		start = time.Now()
		beldilib.CondWrite(env, "append", "K", map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V2"): expression.Value(a),
		}, expression.Name("V").Equal(expression.Value(2)))
		fmt.Printf("DURATION CWriteF %s\n", time.Since(start))

		start = time.Now()
		beldilib.Read(env, "append", "K")
		fmt.Printf("DURATION Read %s\n", time.Since(start))

		return 0
	}
	if beldilib.TYPE == "BASELINE" {
		a := shortuuid.New()
		start := time.Now()
		beldilib.Write(env, "bappend", "K", map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V"): expression.Value(a),
		})
		fmt.Printf("DURATION DWrite %s\n", time.Since(start))

		start = time.Now()
		beldilib.CondWrite(env, "bappend", "K", map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V2"): expression.Value(1),
		}, expression.Name("V").Equal(expression.Value(a)))
		fmt.Printf("DURATION CWriteT %s\n", time.Since(start))

		start = time.Now()
		beldilib.CondWrite(env, "bappend", "K", map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V2"): expression.Value(a),
		}, expression.Name("V").Equal(expression.Value(2)))
		fmt.Printf("DURATION CWriteF %s\n", time.Since(start))

		start = time.Now()
		beldilib.Read(env, "bappend", "K")
		fmt.Printf("DURATION Read %s\n", time.Since(start))

		return 0
	}
	return 1
}

func main() {
	lambda.Start(beldilib.Wrapper(Handler))
}
