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
	outStr = strings.Join(
		strings.Fields(outStr)[2:],
		" ",
	)

	return outStr
}

func getService2Info() Message {
	resp, err := http.Get("http://service2:8199/request")
	if err != nil {
		log.Println("Error getting service2 info: ", err)
		return Message{}
	}
	defer resp.Body.Close()

	service2Info := Message{}
	err = json.NewDecoder(resp.Body).Decode(&service2Info)
	if err != nil {
		log.Println("Error decoding service2 info: ", err)
		return Message{}
	}

	return service2Info
}
