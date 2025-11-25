package main

import (
	"flag"
	"fmt"
)

func main() {
	task := flag.String("task", "", "Task Description")
	priority := flag.Int("priority", 1, "Priority (1-5).")
	flag.Parse()

	if task == nil || *task == "" {
		fmt.Println("Error. Task is empty!")
		return
	}

	fmt.Printf("Added task: %s (Priority: %d)\n", *task, *priority)
}