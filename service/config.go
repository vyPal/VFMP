package main

import (
	"fmt"
	"log"
	"os"
	reflect "reflect"
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v2"
)

type ConfigDatabase struct {
	Data struct {
		Dir     string `yaml:"dir" default:"/var/lib/vfmp"`
		RootDir string `yaml:"root_dir" default:"/home/vypal/Dokumenty/GitHub/VFMP"`
	} `yaml:"data"`
	Server struct {
		Port int `yaml:"port" default:"32768"`
	}
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

func ProcessConfig(configFile string, cfg *ConfigDatabase) {
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
		err = updateConfigFile(configFile, cfg)
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
}
