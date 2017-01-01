package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/markdaws/gohome/pkg/gohome"
	"github.com/markdaws/gohome/pkg/log"
	"github.com/markdaws/gohome/pkg/store"
)

// This is injected by the build process and read from the VERSION file
var VERSION string

func main() {

	version := flag.Bool(
		"version",
		false,
		"View the version of the ghadmin too")

	init := flag.Bool(
		"init",
		false,
		"Inits a new goHOME config and system file in the current directory. You must specify the path to the root of the goHOME source code as the only argument to --init Usage: ghadmin --init /path/to/gohome")

	setPassword := flag.Bool(
		"set-password",
		false,
		"Set the password for a user. Creates a user if the login is not found, you must specify the location to the goHOME config file. e.g. ghadmin --config=./myconfig.json --set-password guest password12345")

	configPath := flag.String("config", "", "Specifies the path and file name to the goHOME config file")

	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		return
	}

	if *init {
		if flag.Arg(0) == "" {
			fmt.Println("You must specify the path to the goHOME source code folder\n\n")
			flag.PrintDefaults()
			os.Exit(1)
		}

		webUIPath, err := filepath.Abs(path.Join(flag.Arg(0), "dist"))
		if err != nil {
			fmt.Println("Invalid dist folder path specified")
			flag.PrintDefaults()
			os.Exit(1)
		}

		cfg := initConfig(webUIPath)
		initSystem(cfg)
		return
	}

	if *setPassword {
		if configPath == nil || *configPath == "" {
			fmt.Println("The config option must be specified when setting the password\n\n")
			flag.PrintDefaults()
			os.Exit(1)
		}

		setPass(flag.Arg(0), flag.Arg(1), *configPath)
		return
	}

	fmt.Println("Please specify an option\n\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func initConfig(webUIPath string) *gohome.Config {
	fmt.Println("Init project with web UI path:", webUIPath)

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to determine the current directory:", err)
		os.Exit(1)
	}

	cfgPath := path.Join(currentDir, "config.json")
	if _, err := os.Stat(cfgPath); err == nil {
		fmt.Printf("The file %s already exists, please remove and then re-run init", cfgPath)
		os.Exit(1)
	}

	cfg := gohome.NewDefaultConfig(currentDir, webUIPath)
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Println("Failed to JSON encode config file: ", err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(cfgPath, b, 0644)
	if err != nil {
		fmt.Println("Failed to write file to disk:", cfgPath)
		os.Exit(1)
	}

	fmt.Println("Config file written to: ", cfgPath)
	return cfg
}

func initSystem(cfg *gohome.Config) {
	if _, err := os.Stat(cfg.SystemPath); err == nil {
		fmt.Printf("The file %s already exists, please remove and then re-run init", cfg.SystemPath)
		os.Exit(1)
	}

	sys := gohome.NewSystem("My goHOME system")
	err := store.SaveSystem(cfg.SystemPath, sys)
	if err != nil {
		fmt.Println("Failed to write system file to disk: ", err)
		os.Exit(1)
	}

	fmt.Println("System file written to: ", cfg.SystemPath)
}

func setPass(login, password, configPath string) {
	if login == "" || password == "" {
		fmt.Println("missing values, --set-password <login> <password>")
		os.Exit(1)
	}

	var cfg *gohome.Config
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Println("Error trying to open:", configPath)
		os.Exit(1)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Println("Failed to parse:", err)
		os.Exit(1)
	}

	if cfg.SystemPath == "" {
		fmt.Println("systemPath key/value not found in:", configPath)
		os.Exit(1)
	}

	log.Silent = true
	sys := loadSystem(cfg.SystemPath)
	log.Silent = false

	var user *gohome.User
	for _, u := range sys.Users() {
		if u.Login == login {
			user = u
			break
		}
	}

	addedUser := false
	if user == nil {
		user = &gohome.User{
			ID:    sys.NewID(),
			Login: login,
		}
		err := user.Validate()
		if err != nil {
			fmt.Println("Failed to add user", err)
			os.Exit(1)
		}

		sys.AddUser(user)
		addedUser = true
	}

	err = user.SetPassword(password)
	if err != nil {
		fmt.Println("Failed to set the password:", err)
		os.Exit(1)
	}

	err = store.SaveSystem(cfg.SystemPath, sys)
	if err != nil {
		fmt.Println("Failed to save the user changes to disk:" + err.Error())
		os.Exit(1)
	}

	if addedUser {
		fmt.Println("Successfully added user:", login, " to:", cfg.SystemPath)
	} else {
		fmt.Println("Successfully updated password for user:", login, "to:", cfg.SystemPath)
	}
}

func loadSystem(systemPath string) *gohome.System {
	sys, err := store.LoadSystem(systemPath)
	if err == store.ErrFileNotFound {
		fmt.Println("System file not found at: ", systemPath)
		os.Exit(1)
	} else if err != nil {
		fmt.Println("Failed to load system file: ", err)
		os.Exit(1)
	}

	return sys
}
