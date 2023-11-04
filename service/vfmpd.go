package main

import (
	"flag"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/sevlyar/go-daemon"
)

func main() {
	startCommand := flag.NewFlagSet("start", flag.ExitOnError)
	stopCommand := flag.NewFlagSet("stop", flag.ExitOnError)
	restartCommand := flag.NewFlagSet("restart", flag.ExitOnError)
	statusCommand := flag.NewFlagSet("status", flag.ExitOnError)

	cfg := ConfigDatabase{}
	loadConfig(&cfg)

	if len(os.Args) < 2 {
		startDaemon(&cfg)
	} else {
		switch os.Args[1] {
		case "start":
			startCommand.Parse(os.Args[2:])
			startDaemon(&cfg)
		case "stop":
			stopCommand.Parse(os.Args[2:])
			tryKillDaemon(cfg.Server.Port)
		case "restart":
			restartCommand.Parse(os.Args[2:])
			tryKillDaemon(cfg.Server.Port)
			startDaemon(&cfg)
		case "status":
			statusCommand.Parse(os.Args[2:])
			if tryConnect(strconv.Itoa(cfg.Server.Port)) {
				log.Print("vfmpd is running")
			} else {
				log.Print("vfmpd is not running")
			}
		default:
			log.Fatal("Unknown command")
		}
	}
}

func loadConfig(cfg *ConfigDatabase) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Print("Error getting user config directory:", err)
		return
	}
	log.Print("User config directory:", configDir)

	vfmpDir := filepath.Join(configDir, "vfmp")
	if _, err := os.Stat(vfmpDir); os.IsNotExist(err) {
		log.Print("vfmp config directory does not exist, creating a new one")
		err = os.Mkdir(vfmpDir, 0755)
		if err != nil {
			log.Fatal("Unable to create vfmp config directory:", err)
		}
	}

	configFile := filepath.Join(vfmpDir, "config.yaml")
	log.Print("Config file:", configFile)

	ProcessConfig(configFile, cfg)
}

func tryConnect(port string) bool {
	conn, err := net.DialTimeout("tcp", "localhost:"+port, time.Second*3)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func startDaemon(cfg *ConfigDatabase) {
	if _, err := os.Stat(cfg.Data.Dir); os.IsNotExist(err) {
		log.Print("vfmp data directory does not exist, creating a new one")
		err = os.Mkdir(cfg.Data.Dir, 0755)
		if err != nil {
			log.Fatal("Unable to create vfmp data directory:", err)
		}
	}

	pidFile := filepath.Join(cfg.Data.Dir, "vfmpd.pid")
	if _, err := os.Stat(pidFile); err == nil {
		is_running := false

		data, err := os.ReadFile(pidFile)
		if err != nil {
			log.Fatal("Unable to read pid file:", err)
		}

		pid, err := strconv.Atoi(string(data))
		if err != nil {
			log.Fatal("Unable to parse pid:", err)
		}

		p, err := os.FindProcess(pid)

		// On unix systems run p.Signal(syscall.Signal(0)) to check if the process is running
		err = p.Signal(os.Signal(syscall.Signal(0)))
		if err == nil {
			is_running = true
		}

		if is_running {
			if YesNoPrompt("vfmpd is already running, would you like to kill it?", false) {
				tryKillDaemon(cfg.Server.Port)
			} else {
				return
			}
		}
	}

	cntxt := &daemon.Context{
		PidFileName: filepath.Join(cfg.Data.Dir, "vfmpd.pid"),
		PidFilePerm: 0644,
		LogFileName: filepath.Join(cfg.Data.Dir, "vfmpd.log"),
		LogFilePerm: 0640,
		WorkDir:     cfg.Data.Dir,
		Umask:       027,
		Args:        []string{"[vfmpd/daemon]"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Print("- - - - - - - - - - - - - - -")
	log.Print("daemon started")

	setupIPCServer(cfg)
}
