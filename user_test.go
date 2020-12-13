package main

import (
	"testing"

	. "github.com/stretchr/testify/require"
)

func TestUser_String(t *testing.T) {
	Equal(t, "name <email>", user{"name", "email"}.String())
}

func TestConfig_removeUser(t *testing.T) {
	users := []user{
		{Name: "a", Email: "a@localhost"},
		{Name: "b", Email: "b@localhost"},
		{Name: "c", Email: "c@localhost"},
	}

	tests := []struct {
		name string
		idx  int
		want []user
	}{
		{"removeFirst", 0, []user{users[1], users[2]}},
		{"removeMid", 1, []user{users[0], users[2]}},
		{"removeLast", 2, []user{users[0], users[1]}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := config{make([]user, len(users))}
			copy(c.Users, users)
			c.removeUser(tt.idx)
			Equal(t, tt.want, c.Users)
		})
	}
}

func TestConfig_update_sorted(t *testing.T) {
	c := &config{}
	c.update(user{"b", "b@b.b"})
	c.update(user{"a", "a@a.a"})
	Len(t, c.Users, 2)
	Equal(t, user{"a", "a@a.a"}, c.Users[0])
	Equal(t, user{"b", "b@b.b"}, c.Users[1])
}

func TestConfig_update_updateExisting(t *testing.T) {
	c := &config{}
	c.update(user{"b", "a@a.a"})
	c.update(user{"a", "a@a.a"})
	Len(t, c.Users, 1)
	Equal(t, user{"a", "a@a.a"}, c.Users[0])
}

func Test_validateEmail(t *testing.T) {
	NoError(t, validateEmail("a@a.a"))
	Error(t, validateEmail(""))
}
