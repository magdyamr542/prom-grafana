package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/magdyamr542/prom-grafana/prom"
)

func fib(input int) int {
	if input < 0 {
		panic(fmt.Sprintf("input %q < 0 not allowed\n", input))
	}
	if input == 0 {
		return 0
	}

	if input == 1 {
		return 1
	}

	return fib(input-1) + fib(input-2)
}

func main() {

	// logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	// api
	http.HandleFunc("/fib", func(w http.ResponseWriter, r *http.Request) {
		theInput := -1
		defer func() {
			prom.OnNewFibRequest(theInput)
		}()

		if r.Method != "GET" {
			http.Error(w, "only GET requests are supported", http.StatusBadRequest)
			return
		}

		input := r.URL.Query().Get("input")
		if input == "" {
			http.Error(w, "missing 'input' query parameters", http.StatusBadRequest)
			return
		}

		theInput, err := strconv.Atoi(input)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not convert %q to a number", input), http.StatusBadRequest)
			return
		}

		if theInput < 0 {
			http.Error(w, fmt.Sprintf("only accepts >= 0 inputs. got %q", theInput), http.StatusBadRequest)
			return
		}

		logger.Printf("the number is %q\n", theInput)

		result := fib(theInput)
		fmt.Fprintf(w, "%d", result)
	})

	// prometheus
	promHandler := prom.Init()
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promHandler.ServeHTTP(w, r)
	})

	// spin up the server
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "1234"
		logger.Printf("No PORT supplied. Using port %s\n", PORT)
	}
	portStr := fmt.Sprintf(":%s", PORT)
	logger.Printf("Server is up on localhost:%s\n", PORT)
	if err := http.ListenAndServe(portStr, nil); err != nil {
		log.Fatal("error ListenAndServe", err)
	}
}
