package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"
)

type VersionPID struct {
	Return struct {
		PID int `json:"pid"`
	} `json:"return"`
}

type StatusResponse struct {
	Return struct {
		OutData string `json:"out-data"`
	} `json:"return"`
}

func main() {

	vmName := "debian-12-421ap1"

	payload := `{"execute": "guest-exec", "arguments": {"path": "bash", "arg": ["-c", "dpkg -l"], "capture-output": true}}`

	cmd := exec.Command("sudo", "virsh", "qemu-agent-command", vmName, payload)

	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	var data VersionPID

	err = json.Unmarshal(output, &data)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	basepayload := fmt.Sprintf(`{"execute": "guest-exec-status", "arguments": {"pid": %d}}`, data.Return.PID)

	statusCMD := exec.Command("sudo", "virsh", "qemu-agent-command", vmName, basepayload)

	statusOutput, err := statusCMD.Output()
	if err != nil {
		log.Fatal(err)
	}

	var statusData StatusResponse

	err = json.Unmarshal(statusOutput, &statusData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(data.Return.PID)
	fmt.Println(string(output))
	fmt.Println(string(statusOutput))

	decodbytes, err := base64.StdEncoding.DecodeString(statusData.Return.OutData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(decodbytes))
}
