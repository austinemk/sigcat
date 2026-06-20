// Package app is the package entry

package app

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	tea "charm.land/bubbletea/v2"
	"github.com/austinemk/sigcat/pkg/dashboard"
	"github.com/austinemk/sigcat/pkg/helpers"
	"github.com/austinemk/sigcat/pkg/popup"
)

func Execute() {
	runFlag := flag.Bool("run-daemon", false, "Start the sigcat engine context")
	stopFlag := flag.Bool("stop-daemon", false, "Stop the active background engine")
	uiMode := flag.String("ui", "", "Launch UI ('dashboard' or 'popup')")
	taskID := flag.String("task-id", "", "Target task reference for the popup renderer")
	flag.Parse()

	if *stopFlag {
		cmd := exec.Command("pkill", "-f", "sigcat --run-daemon")
		if err := cmd.Run(); err != nil {
			fmt.Println("❌ No running sigcat daemon found.")
			return
		}
		fmt.Println("🛑 Sigcat daemon stopped successfully.")
		return
	}

	if *runFlag {
		if os.Getenv("SIGCAT_BACKGROUND") != "true" {
			executable, _ := os.Executable()
			cmd := exec.Command(executable, "--run-daemon")
			cmd.Env = append(os.Environ(), "SIGCAT_BACKGROUND=true")

			logFile, _ := os.OpenFile(os.Getenv("HOME")+"/.config/sigcat/daemon.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
			cmd.Stdout = logFile
			cmd.Stderr = logFile

			if err := cmd.Start(); err != nil {
				fmt.Printf("❌ Failed to split engine thread: %v\n", err)
				return
			}
			fmt.Println("🚀 Sigcat tracking platform engaged in background!")
			return
		}
		helpers.RunDaemon()
		return
	}

	// Route user interface environments
	if *uiMode == "dashboard" || (*uiMode == "" && len(os.Args) == 1) {
		p := tea.NewProgram(dashboard.InitialDashboardModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Dashboard crash: %v\n", err)
		}
		return
	}

	if *uiMode == "popup" || *uiMode == "break" { // Backward compatible fallback support included
		p := tea.NewProgram(popup.InitialPopupModel(*taskID))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Popup display engine failure: %v\n", err)
		}
		return
	}

	fmt.Println("Usage:\n  ./sigcat                  (Launches Config Dashboard)\n  ./sigcat --run-daemon     (Starts background engine)\n  ./sigcat --stop-daemon    (Stops background engine)")
}
