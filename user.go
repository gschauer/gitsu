package main

import (
	"errors"
	"sort"

	"github.com/asaskevich/govalidator"
)

// user is a git user represented by name and email.
type user struct {
	Name  string `json:"Name"`
	Email string `json:"Email"`
}

// String formats the user as "Name <Email>".
func (u user) String() string {
	return u.Name + " <" + u.Email + ">"
}

// update adds the given user.
// If user with the same email exists already, it's name is updated instead.
func (c *config) update(u user) {
	// check if there's already a user with the same email
	idx := sort.Search(len(c.Users), func(i int) bool {
		return c.Users[i].Email == u.Email
	})

	// if the user already exists, update it - otherwise create a new one
	if idx < len(c.Users) {
		c.Users[idx].Name = u.Name
	} else {
		c.Users = append(c.Users, u)
		sort.Slice(c.Users, func(i, j int) bool {
			return c.Users[i].Email < c.Users[j].Email
		})
	}
}

// removeUser returns a new slice with the element at index idx removed.
// Note that this function does not do any bound checks.
func (c *config) removeUser(idx int) {
	copy(c.Users[idx:], c.Users[idx+1:])
	c.Users = c.Users[:len(c.Users)-1]
}

// validateEmail checks if the string is an email.
func validateEmail(email string) error {
	if !govalidator.IsEmail(email) {
		return errors.New("invalid email address")
	}
	return nil
}
