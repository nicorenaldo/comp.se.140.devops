package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type Message struct {
	IP       string   `json:"ip"`
	PS       []string `json:"ps"`
	DF       []string `json:"df"`
	LastBoot string   `json:"last_boot"`
}

func getIP() string {
	out, err := exec.Command("hostname", "-I").Output()
	if err != nil {
		log.Println("Error getting IP: ", err)
		return ""
	}

	ip := strings.TrimSpace(string(out))
	return ip
}

func getPS() []string {
	out, err := exec.Command("ps", "aux").Output()
	if err != nil {
		log.Println("Error getting PS: ", err)
		return []string{}
	}

	lines := strings.Split(string(out), "\n")
	return lines
}

func getDF() []string {
	out, err := exec.Command("df", "-h").Output()
	if err != nil {
		log.Println("Error getting DF: ", err)
		return []string{}
	}

	lines := strings.Split(string(out), "\n")
	return lines
}

func getLastBoot() string {
	out, err := exec.Command("last", "reboot", "|", "tail", "-1").Output()
	if err != nil {
		log.Println("Error getting last boot: ", err)
		return ""
	}

	outStr := strings.TrimSpace(string(out))
	log.Println("Last boot: ", outStr)
	outStr = strings.Join(
		strings.Fields(outStr)[2:],
		" ",
	)
	log.Println("Last boot: ", outStr)

	return outStr
}

func handle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/json")

	message := Message{
		IP:       getIP(),
		PS:       getPS(),
		DF:       getDF(),
		LastBoot: getLastBoot(),
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonMessage)
	if err != nil {
		log.Println("Error writing response: ", err)
		return
	}
}

func main() {
	http.HandleFunc("/", handle)

	log.Println("Starting server on port 8199")
	err := http.ListenAndServe(":8199", nil)
	if err != nil {
		log.Println("Error starting server: ", err)
		return
	}
}
