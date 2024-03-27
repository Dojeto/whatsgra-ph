package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/Dojeto/whatsgra-ph/utils"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	utils.LoadEnvVariables()
}

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go utils.ConnectToWP()
	go utils.HandleRequests()

	<-quit

	fmt.Println("Received termination signal. Exiting...")

	os.Exit(0)
}
