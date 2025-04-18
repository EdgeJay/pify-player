package utils

import (
	"log"
	"os"
	"os/exec"
)

func TriggerSystemShutdown() error {
	if !ShellCommandsAllowed() {
		log.Println("Shell commands are not allowed.")
		return nil
	}

	cmd := exec.Command("sh", "-c", "sleep 10 && sudo /host_sbin/shutdown -h now")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Println("Failed to shutdown:", err)
	}

	log.Println("Shutdown command issued successfully.")

	return err
}

func TriggerSystemRestart() error {
	if !ShellCommandsAllowed() {
		log.Println("Shell commands are not allowed.")
		return nil
	}

	cmd := exec.Command("sh", "-c", "sleep 10 && sudo /host_sbin/shutdown -r now")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Println("Failed to restart:", err)
	}

	log.Println("Restart command issued successfully.")

	return err
}
