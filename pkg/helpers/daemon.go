package helpers

import (
	"log"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
)

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func RunDaemon() {
	log.Println("🐱 sigcat background runtime scheduler listening...")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	terminalApp := FindTerminal()
	executable, _ := os.Executable()

	for range ticker.C {
		tasks, err := LoadTasks()
		if err != nil {
			continue
		}

		changed := false
		now := time.Now()
		activeCount := 0

		// In daemon.go (inside RunDaemon's loop)

		for i, task := range tasks {
			if !task.IsActive {
				continue
			}

			activeCount++

			if now.After(task.NextRun) {
				log.Printf("⏰ Target hit for profile [%s]: %s\n", task.ID, task.Title)

				// 1. Spawn the floating window
				_ = SpawnFloatingWindow(terminalApp, executable, task.ID)

				// 2. Turn off IsActive temporarily so the daemon stops tracking it
				// until the user interacts with the popup window and closes it.
				tasks[i].IsActive = false
				activeCount--

				changed = true
			}
		}

		if changed {
			_ = SaveTasks(tasks)
		}

		// Self-termination safety logic if no automated tasks remain active
		if activeCount == 0 {
			log.Println("💤 No active profiles found running. Giving workspace windows a second to map before exit...")

			time.Sleep(2 * time.Second)

			log.Println("💤 Shutting down daemon context automatically.")
			return // Exiting main loop shuts down the background daemon process safely!
		}
	}
}
