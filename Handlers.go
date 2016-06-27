package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	todos := Todos{
		Todo{Name: "Write presentation"},
		Todo{Name: "Host meetup"},
	}

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoId)
}

func ListContainers(w http.ResponseWriter, r *http.Request) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
	if err != nil {
		panic(err)
	}

	options := types.ContainerListOptions{All: false}
	containers, err := cli.ContainerList(context.Background(), options)
	if err != nil {
		panic(err)
	}

	if err := json.NewEncoder(w).Encode(containers); err != nil {
		panic(err)
	}
}

func LogContainer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
	if err != nil {
		panic(err)
	}

	options := types.ContainerLogsOptions{ShowStdout: true, Follow: true, Timestamps: true, Details: true}
	reader, err := cli.ContainerLogs(ctx, "efe57b174696627e4f2d291abd1b683099b621a377929d19098fe44b3367bb37", options)
	if err != nil {
		log.Fatal(err)
	}

	// convert reader in string
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	s := buf.String()

	if err := json.NewEncoder(w).Encode(s); err != nil {
		panic(err)
	}

}
