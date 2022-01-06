package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func main() {

	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Required argument is port number")
	}

	port, portErr := strconv.Atoi(args[0])

	if portErr != nil {
		fmt.Printf("Failed to start server. '%s' is not a valid port\n", args[0])
		return
	}

	fmt.Printf("Started server on port %d\n", port)

	err := http.ListenAndServe(
		fmt.Sprintf(":%d", port),
		http.FileServer(http.Dir("../../web")),
	)
	if err != nil {
		fmt.Println("Failed to start server", err)
		return
	}

}
