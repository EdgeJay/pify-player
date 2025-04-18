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

	cmd := exec.Command("sleep 3 && /app/scripts/shutdown.sh")
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

	cmd := exec.Command("sleep 3 && /app/scripts/restart.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Println("Failed to restart:", err)
	}

	log.Println("Restart command issued successfully.")

	return err
}
