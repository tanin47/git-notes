//go:build !linux

package main

import (
	"log"
)

func makeService() {
	log.Println("Sorry, we have only linux implementation for now")
}

func makeConfig() {}
