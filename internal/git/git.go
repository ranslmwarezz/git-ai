package git

import (
	"os/exec"
)

type GitClient interface {
	DiffCached() (string, error)
	Commit(m string) (error)
}

type CommandFunc func(string, ...string) *exec.Cmd

type Client struct {
	command CommandFunc
}

func NewClient() *Client {
	return &Client{command: exec.Command,
	}
}

func (c *Client) DiffCached() (string, error) {

	cmd := c.command("git", "diff", "--cached")

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (c *Client) Commit(message string) error {

	cmd := c.command("git", "commit", "-m", message)

	return cmd.Run()
}
