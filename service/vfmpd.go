package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sevlyar/go-daemon"
	"gopkg.in/yaml.v2"
)

type ConfigDatabase struct {
	Data struct {
		Dir string `yaml:"dir" default:"/var/lib/vfmp"`
	} `yaml:"data"`
	Server struct {
		Port int `yaml:"port" default:"32768"`
	}
}

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func DefaultConfig() ConfigDatabase {
	d := ConfigDatabase{}
	setDefaults(&d)
	return d
}

func setDefaults(s interface{}) {
	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("default")

		if tag != "" {

			switch field.Kind() {
			case reflect.String:
				field.SetString(tag)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				intValue, err := strconv.ParseInt(tag, 0, field.Type().Bits())
				if err != nil {
					log.Fatalf("Unable to parse default value for field %s: %v", t.Field(i).Name, err)
				}
				field.SetInt(intValue)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				uintValue, err := strconv.ParseUint(tag, 0, field.Type().Bits())
				if err != nil {
					log.Fatalf("Unable to parse default value for field %s: %v", t.Field(i).Name, err)
				}
				field.SetUint(uintValue)
			case reflect.Float32, reflect.Float64:
				floatValue, err := strconv.ParseFloat(tag, field.Type().Bits())
				if err != nil {
					log.Fatalf("Unable to parse default value for field %s: %v", t.Field(i).Name, err)
				}
				field.SetFloat(floatValue)
			case reflect.Bool:
				boolValue, err := strconv.ParseBool(tag)
				if err != nil {
					log.Fatalf("Unable to parse default value for field %s: %v", t.Field(i).Name, err)
				}
				field.SetBool(boolValue)
			case reflect.Struct:
				setDefaults(field.Addr().Interface())
			default:
				log.Fatalf("Unsupported field type for field %s: %v", t.Field(i).Name, field.Kind())
			}
		} else if field.Kind() == reflect.Struct {
			setDefaults(field.Addr().Interface())
		} else {
			log.Fatalf("No default value for field %s", t.Field(i).Name)
		}
	}
}

// To terminate the daemon use:
//
//	kill `cat sample.pid`
func main() {
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

	var cfg ConfigDatabase

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Print("Config file does not exist, creating a new one with default values")

		cfg := DefaultConfig()

		log.Println("Data Dir is not set, using default value:", cfg.Data.Dir)

		data, err := yaml.Marshal(&cfg)
		if err != nil {
			log.Fatal("Unable to marshal config:", err)
		}

		err = os.WriteFile(configFile, data, 0644)
		if err != nil {
			log.Fatal("Unable to write config:", err)
		}
		if err != nil {
			log.Fatal("Unable to write config:", err)
		}
	} else {
		err = updateConfigFile(configFile, &cfg)
		if err != nil {
			log.Fatal("Unable to update config file:", err)
		}

		err = cleanenv.ReadConfig(configFile, &cfg)
		if err != nil {
			log.Fatal("Unable to read config:", err)
		}

		if _, err := os.Stat(cfg.Data.Dir); os.IsNotExist(err) {
			log.Print("vfmp data directory does not exist, creating a new one")
			err = os.Mkdir(cfg.Data.Dir, 0755)
			if err != nil {
				log.Fatal("Unable to create vfmp data directory:", err)
			}
		}
	}

	if _, err := os.Stat(cfg.Data.Dir); os.IsNotExist(err) {
		log.Print("vfmp data directory does not exist, creating a new one")
		err = os.Mkdir(cfg.Data.Dir, 0755)
		if err != nil {
			log.Fatal("Unable to create vfmp data directory:", err)
		}
	}

	// Open the pid file, if it exists and there is a pid, check if there is a process running with that pid. If there is, ask if the user would like to kill it. If yes, run the tryKillDaemon function. If no, exit.
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

		// Check if there is a process running with the pid
		p, err := os.FindProcess(pid)

		// On unix systems run p.Signal(syscall.Signal(0)) to check if the process is running
		err = p.Signal(os.Signal(syscall.Signal(0)))
		if err == nil {
			is_running = true
		}

		if is_running {
			// If the process is running, ask if the user would like to kill it
			if YesNoPrompt("vfmpd is already running, would you like to kill it?", false) {
				tryKillDaemon(cfg.Server.Port)
				return
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

	setupIPCServer(cfg.Server.Port)
}

func tryKillDaemon(port int) {
	serverAddr := fmt.Sprintf("localhost:%d", port)

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		return
	}
	defer conn.Close()

	// Set a timeout for the connection
	conn.SetDeadline(time.Now().Add(10 * time.Second))

	// Send kill message
	message := Message{
		Type: "kill",
		Data: "Daemon restart",
	}
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Failed to marshal message: %v\n", err)
		return
	}

	// Write the json + "\n" to the server
	_, err = conn.Write(append(jsonMessage, '\n'))
	if err != nil {
		fmt.Printf("Failed to send message to server: %v\n", err)
		return
	}

	// Relaunch the current process
	args := os.Args
	executable, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}

	err = syscall.Exec(executable, args, os.Environ())
	if err != nil {
		log.Fatalf("Failed to relaunch process: %v", err)
	}
}

func setupIPCServer(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("Unable to listen: ", err)
	}
	defer listener.Close()

	log.Printf("IPC server listening on port %d", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print("Error accepting connection: ", err)
			continue
		}

		go handleConnection(conn)
	}
}

func updateConfigFile(configFile string, cfg *ConfigDatabase) error {
	// Read existing config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("unable to read config file: %w", err)
	}

	// Unmarshal config file into ConfigDatabase struct
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return fmt.Errorf("unable to unmarshal config file: %w", err)
	}

	// Check if each field exists in the config file
	defaultCfg := DefaultConfig()
	updated := false

	// Write a dynamic function to go through the struct recursively
	var checkField func(interface{}, interface{}, string)
	checkField = func(cfg interface{}, defaultCfg interface{}, path string) {
		v := reflect.ValueOf(cfg).Elem()
		t := v.Type()

		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			tag := t.Field(i).Tag.Get("default")

			if tag != "" {
				switch field.Kind() {
				case reflect.Struct:
					checkField(field.Addr().Interface(), reflect.ValueOf(defaultCfg).Elem().Field(i).Addr().Interface(), path+t.Field(i).Name+".")
				case reflect.String:
					if field.String() == "" {
						log.Printf("%s%s is not set, using default value: %s", path, t.Field(i).Name, tag)
						field.SetString(tag)
						updated = true
					}
				case reflect.Int:
					if field.Int() == 0 {
						intValue, err := strconv.ParseInt(tag, 0, field.Type().Bits())
						if err != nil {
							log.Fatalf("Unable to parse default value for field %s: %v", t.Field(i).Name, err)
						}
						log.Printf("%s%s is not set, using default value: %d", path, t.Field(i).Name, intValue)
						field.SetInt(intValue)
						updated = true
					}
				case reflect.Bool:
					if !field.Bool() {
						boolValue, err := strconv.ParseBool(tag)
						if err != nil {
							log.Fatalf("Unable to parse default value for field %s: %v", t.Field(i).Name, err)
						}
						log.Printf("%s%s is not set, using default value: %t", path, t.Field(i).Name, boolValue)
						field.SetBool(boolValue)
						updated = true
					}
				default:
					log.Fatalf("Unsupported field type for field %s: %v", t.Field(i).Name, field.Kind())
				}
			} else if field.Kind() == reflect.Struct {
				checkField(field.Addr().Interface(), reflect.ValueOf(defaultCfg).Elem().Field(i).Addr().Interface(), path+t.Field(i).Name+".")
			}
		}
	}

	checkField(cfg, &defaultCfg, "")

	// Save updated config file
	if updated {
		data, err = yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("unable to marshal config: %w", err)
		}

		err = os.WriteFile(configFile, data, 0644)
		if err != nil {
			return fmt.Errorf("unable to write config file: %w", err)
		}
	}

	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Print("New connection established")

	for {
		// Read incoming message
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Print("Error reading message: ", err)
			return
		}

		// Check message type
		var data map[string]interface{}
		err = json.Unmarshal([]byte(msg), &data)
		if err != nil {
			log.Print("Error unmarshaling message: ", err)
			return
		}

		if data["type"] == "kill" {
			log.Print("Received kill message")

			if data["data"] != nil {
				log.Print("Kill reason: ", data["data"])
			}

			os.Exit(0)
			return
		}

		if data["type"] == "ping" {
			log.Print("Received ping message")

			// Send pong message
			pong := map[string]interface{}{
				"type": "pong",
			}

			pongData, err := json.Marshal(pong)
			if err != nil {
				log.Print("Error marshaling pong message: ", err)
				return
			}

			_, err = conn.Write(pongData)
			if err != nil {
				log.Print("Error writing pong message: ", err)
				return
			}
		}
	}
}

func YesNoPrompt(label string, def bool) bool {
	choices := "Y/n"
	if !def {
		choices = "y/N"
	}

	r := bufio.NewReader(os.Stdin)
	var s string

	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return def
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}
