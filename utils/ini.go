package utils

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	"gopkg.in/ini.v1"
)

// Default Config
type Config struct {
	// Application
	Version string `default:"1.0.0"`
	DevMode bool   `default:"false"`
	// Server
	IP   string `default:"127.0.0.1"`
	Port int    `default:"3000"`
	// WebView
	Title  string `default:"STLViewer"`
	Width  int    `default:"800"`
	Height int    `default:"600"`
	URL    string `default:"http://127.0.0.1:3000"`
}

var Env Config

// InItialize
func Init() {
	filename := fmt.Sprintf("%v/server.ini", DirPath)

	// 'server.ini' file exists check
	// 1. true  - load server.ini
	// 2. false - make server.ini
	if fileExists := Exists(filename); fileExists {
		// load 'server.ini'
		cfg, err := ini.Load(filename)
		if err != nil {
			logger.Printf("Load 'server.ini':%v \n", err)
			log.Fatalf("Load 'server.ini':%v \n", err)
		}

		version := cfg.Section("Application").Key("Version").String()

		devMode, err := cfg.Section("Application").Key("DevMode").Bool()
		if err != nil {
			logger.Printf("cfg.Section key'DevMode':%v \n", err)
			log.Fatalf("cfg.Section key'DevMode':%v \n", err)
		}

		ip := cfg.Section("Server").Key("IP").String()

		port, err := cfg.Section("Server").Key("Port").Int()
		if err != nil {
			logger.Printf("cfg.Section key'Port':%v \n", err)
			log.Fatalf("cfg.Section key'Port':%v \n", err)
		}

		title := cfg.Section("WebView").Key("Title").String()

		width, err := cfg.Section("WebView").Key("Width").Int()
		if err != nil {
			logger.Printf("cfg.Section key'Width':%v \n", err)
			log.Fatalf("cfg.Section key'Width':%v \n", err)
		}

		height, err := cfg.Section("WebView").Key("Height").Int()
		if err != nil {
			logger.Printf("cfg.Section key'Height':%v \n", err)
			log.Fatalf("cfg.Section key'Height':%v \n", err)
		}

		url := cfg.Section("WebView").Key("URL").String()

		Env = Config{
			Version: version,
			DevMode: devMode,
			IP:      ip,
			Port:    port,
			Title:   title,
			Width:   width,
			Height:  height,
			URL:     url,
		}
	} else {
		// make 'server.ini'
		_, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			logger.Printf("Make 'server.ini':%v \n", err)
			log.Fatalf("Make 'server.ini':%v \n", err)
		}

		// set defalut config
		config := &Config{}
		setError := Set(config, "default")
		if setError != nil {
			logger.Printf("Set Default Config:%v \n", setError)
			log.Fatalf("Set Default Config:%v \n", setError)
		}
		Env = *config

		// load 'server.ini'
		cfg, err := ini.Load(filename)

		cfg.Section("Application").Key("Version").SetValue(fmt.Sprint(config.Version))
		cfg.Section("Application").Key("DevMode").SetValue(fmt.Sprint(config.DevMode))
		cfg.Section("Server").Key("IP").SetValue(fmt.Sprint(config.IP))
		cfg.Section("Server").Key("Port").SetValue(fmt.Sprint(config.Port))
		cfg.Section("WebView").Key("Title").SetValue(fmt.Sprint(config.Title))
		cfg.Section("WebView").Key("Width").SetValue(fmt.Sprint(config.Width))
		cfg.Section("WebView").Key("Height").SetValue(fmt.Sprint(config.Height))
		cfg.Section("WebView").Key("URL").SetValue(fmt.Sprint(config.URL))

		// save 'server.ini'
		cfg.SaveTo(filename)
	}
}

func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func Set(ptr interface{}, tag string) error {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return fmt.Errorf("Not a pointer")
	}

	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		if defaultVal := t.Field(i).Tag.Get(tag); defaultVal != "-" {
			if err := setField(v.Field(i), defaultVal); err != nil {
				return err
			}

		}
	}
	return nil
}

func setField(field reflect.Value, defaultVal string) error {
	if !field.CanSet() {
		return fmt.Errorf("Can't set value\n")
	}

	switch field.Kind() {
	case reflect.Int:
		if val, err := strconv.ParseInt(defaultVal, 10, 64); err == nil {
			field.Set(reflect.ValueOf(int(val)).Convert(field.Type()))
		}
	case reflect.String:
		field.Set(reflect.ValueOf(defaultVal).Convert(field.Type()))
	}

	return nil
}
