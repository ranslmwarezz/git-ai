package git

import "os/exec"

type GitClient interface {
	DiffCached() (string, error)
}

type Client struct {
}

func (c *Client) DiffCached() (string, error) {

	cmd := exec.Command("git", "diff", "--cached")

	output, err := cmd.Output()
	if err != nil {
		return "", nil
	}
	return string(output), nil
}
