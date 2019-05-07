package main

//import "encoding/json"
import "fmt"
import "go.starlark.net/starlark"
import "syscall/js"

var (
	quit                                  chan bool
	window, document, body, result, input js.Value
)

func main() {
	quit = make(chan bool)

	window = js.Global()
	document = window.Get("document")
	body = document.Get("body")

	result = document.Call("createElement", "div")
	result.Set("id", "result")

	input = document.Call("createElement", "textarea")
	input.Set("value", "")
	input.Call("addEventListener", "input", js.FuncOf(handleChange))

	body.Call("appendChild", input)
	body.Call("appendChild", result)

	<-quit
}

func handleChange(this js.Value, args []js.Value) interface{} {
	text := args[0].Get("target").Get("value").String()
	thread := &starlark.Thread{Name: "main"}
	res, err := starlark.ExecFile(thread, "", text, nil)
	if err == nil {
		fmt.Println("res: ", res)
		result.Set("innerHTML", res.String())
		// b, err := json.Marshal(&res)
		//fmt.Println("encoding:", b, err)

		//result.Set("innerHTML", string(b))
	} else {
		fmt.Println("err: ", err)
	}

	return nil
}
