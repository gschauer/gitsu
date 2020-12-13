package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/manifoldco/promptui"
)

func usage() {
	format := `Usage:
  gitsu [flags]

Flags:
  --global              Set user as global.

Author:
  matsuyoshi30 <sfbgwm30@gmail.com>
`
	fmt.Fprintln(os.Stderr, format)
}

var (
	isGlobal = flag.Bool("global", false, "Set user as global")

	// these are set in build step
	version = "unversioned"
	commit  = "?"
	date    = "?"
)

const (
	sel = "Select git user"
	add = "Add new git user"
	del = "Delete git user"
)

func main() {
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	prompt := promptui.Select{
		Label: "Select action",
		Items: []string{sel, add, del},
	}

	_, res, err := prompt.Run()
	if err != nil {
		log.Fatalln("Failed to select action:", err)
	}
	run(res)
}

func run(res string) {
	c, err := readConfig()
	if err != nil && (res != add || !errors.Is(err, errNoUser)) {
		log.Fatalln("Cannot read config:", err)
	}

	switch res {
	case sel:
		if err := setUser(c); err != nil {
			log.Fatalln("Failed to select user:", err)
		}
	case add:
		u, err := addUser()
		if err != nil {
			log.Fatalln("Failed to add user:", err)
		}
		c.update(u)
		if err = writeConfig(c); err != nil {
			log.Fatalln("Failed to add user:", err)
		}
	case del:
		if err := deleteUser(c); err != nil {
			log.Fatalln("Failed to delete user:", err)
		}
	default:
		log.Fatalln("Unexpected action type")
	}
}

func setUser(c *config) error {
	idx, err := selUser(c)
	if err != nil {
		return err
	}

	option := "--local"
	if *isGlobal {
		option = "--global"
	}

	nameCmd := exec.Command("git", "config", option, "user.name", c.Users[idx].Name)
	if err := nameCmd.Run(); err != nil {
		return err
	}
	emailCmd := exec.Command("git", "config", option, "user.email", c.Users[idx].Email)
	if err := emailCmd.Run(); err != nil {
		return err
	}
	return nil
}

func addUser() (user, error) {
	var u user
	prompt := promptui.Prompt{
		Label: "Input git user name",
	}
	name, err := prompt.Run()
	if err != nil {
		return u, err
	}

	prompt = promptui.Prompt{
		Label:    "Input git email address",
		Validate: validateEmail,
	}
	email, err := prompt.Run()
	if err != nil {
		return u, err
	}
	return user{name, email}, nil
}

func deleteUser(c *config) error {
	idx, err := selUser(c)
	if err != nil {
		return err
	}
	c.removeUser(idx)
	return writeConfig(c)
}

func selUser(c *config) (int, error) {
	if len(c.Users) == 0 {
		return 0, errNoUser
	}

	prompt := promptui.Select{
		Label: "Select git user",
		Items: c.Users,
	}
	idx, _, err := prompt.Run()
	return idx, err
}
