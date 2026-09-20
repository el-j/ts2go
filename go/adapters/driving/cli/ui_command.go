package cli

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/el-j/ts2go/adapters/driving/web"
)

// UICommand starts the web-based API server and optionally opens a browser
func (app *Application) UICommand(args []string) error {
	port := 8080
	openBrowser := false

	// Parse command line options
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--port", "-p":
			if i+1 < len(args) {
				p, err := strconv.Atoi(args[i+1])
				if err != nil {
					return fmt.Errorf("invalid port number: %s", args[i+1])
				}
				port = p
				i++
			}
		case "--open", "-o":
			openBrowser = true
		}
	}

	addr := fmt.Sprintf("localhost:%d", port)
	server, err := web.NewServer(addr)
	if err != nil {
		return fmt.Errorf("failed to create web server: %w", err)
	}

	if openBrowser {
		go func() {
			url := fmt.Sprintf("http://localhost:%d", port)
			fmt.Printf("Opening browser at %s\n", url)
			_ = openBrowserURL(url)
		}()
	}

	return server.Start()
}

// openBrowserURL opens the specified URL in the default browser based on OS
func openBrowserURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start()
}
