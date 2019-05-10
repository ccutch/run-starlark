package main

//import "encoding/json"
import "fmt"
import "go.starlark.net/starlark"
import "syscall/js"

var (
	quit          chan bool
	result, input js.Value
)

func main() {
	document := js.Global().Get("document")
	codeboxes := document.Call("getElementsByClassName", "codebox-container").Call("item", 0)

	result = document.Call("createElement", "textarea")
	result.Set("className", "codebox")

	input = document.Call("createElement", "textarea")
	input.Set("className", "codebox")
	input.Call("addEventListener", "input", js.FuncOf(handleChange))

	codeboxes.Call("appendChild", input)
	codeboxes.Call("appendChild", result)

	<-quit
}

func handleChange(this js.Value, args []js.Value) interface{} {
	text := args[0].Get("target").Get("value").String()
	thread := &starlark.Thread{Name: "main"}
	res, err := starlark.ExecFile(thread, "", text, nil)
	if err == nil {
		fmt.Println("res: ", res)
		result.Set("innerHTML", res.String())
	} else {
		result.Set("innerHTML", err.Error())
	}

	return nil
}
