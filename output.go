package main

import "fmt"

const Reset = "\033[0m"
const Red = "\033[31m"
const Green = "\033[32m"
const Yellow = "\033[33m"
const Blue = "\033[34m"
const Gray = "\033[37m"
const White = "\033[97m"
const Bold = "\033[1m"

func errorOut(msg string) {
	fmt.Printf("%s\n", Red+msg+Reset)
}

func warningOut(msg string) {
	fmt.Printf("%s\n", Yellow+msg+Reset)
}

func stringOut(key string, value string) {
	fmt.Printf("%-32s: %s\n", Bold+Green+key+Reset+White, value+Reset)
}

func arrayOut(key string, arr []string) {
	if len(arr) > 0 {
		fmt.Printf("%-32s: %s\n", Bold+Green+key+Reset+White, "AVAILABLE"+Reset)
	} else {
		fmt.Printf("%-32s: %s\n", Bold+Green+key+Reset+White, "UNKNOWN"+Reset)
	}
}

func mapOut(key string, kv map[string]string) {
	value, err := kv[key]
	if err {
		fmt.Printf("%-32s: %s\n", Bold+Green+key+Reset+White, value+Reset)
	} else {
		fmt.Printf("%-32s: %s\n", Bold+Green+key+Reset+White, "UNKNOWN"+Reset)
	}
}
