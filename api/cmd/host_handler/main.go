package main

import (
	"log"
	"net/http"
	"os/exec"
)

// shutdownHandler handles POST /shutdown
func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	go func() {
		cmd := exec.Command("sudo", "shutdown", "-h", "now")
		if err := cmd.Run(); err != nil {
			log.Printf("Failed to shutdown: %v", err)
		}
	}()

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Shutdown triggered."))
}

// rebootHandler handles POST /reboot
func rebootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	go func() {
		cmd := exec.Command("sudo", "shutdown", "-r", "now")
		if err := cmd.Run(); err != nil {
			log.Printf("Failed to reboot: %v", err)
		}
	}()

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Reboot triggered."))
}

func main() {
	http.HandleFunc("/shutdown", shutdownHandler)
	http.HandleFunc("/reboot", rebootHandler)

	port := "8081"
	log.Printf("Host control server listening on :%s", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
