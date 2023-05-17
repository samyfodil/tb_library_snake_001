//go:build wasi

package main

//export _ready
func ready()

func main() {
	ready()
}
