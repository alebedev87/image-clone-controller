package registry

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"image-clone-controller/pkg/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	defaultSkopeoTransport = "docker://"
)

var registryAliases = map[string][]string{
	"registry-1.docker.io": {
		"docker.io",
		"registry-1.docker.io",
	},
	"quay.io": {
		"quay.io",
	},
}

// Client is a high level abstraction over image registry client
type Client struct {
	registry           string
	organization       string
	username           string
	password           string
	copyTimeoutSeconds int
	transport          string
}

// NewClient returns new registry client
func NewClient(registry, org, username, password string, timeout int) *Client {
	return &Client{
		registry:           registry,
		organization:       org,
		username:           username,
		password:           password,
		copyTimeoutSeconds: timeout,
		transport:          defaultSkopeoTransport,
	}
}

// NewClientFromConfig returns new registry client set from the program's config
func NewClientFromConfig() *Client {
	return &Client{
		registry:           config.GlobalConfig.Registry,
		organization:       config.GlobalConfig.Organization,
		username:           config.GlobalConfig.Username,
		password:           config.GlobalConfig.Password,
		copyTimeoutSeconds: config.GlobalConfig.ImageCopyTimeoutSeconds,
		transport:          defaultSkopeoTransport,
	}
}

// Belongs returns true if given full image name
// belongs to the registry the client is currently using
func (c *Client) Belongs(fullName string) bool {
	fullName = strings.TrimSpace(fullName)
	substr := strings.Split(fullName, "/")
	if len(substr) < 2 {
		// cannot be from the backup repository
		// as it starts with registry and organization
		return false
	}

	for _, r := range registryAliases[c.registry] {
		if substr[0] == r {
			if substr[1] == c.organization {
				return true
			}
		}
	}

	return false
}

// Backup pulls the given image to the backup registry.
// New image full name is returned.
func (c *Client) Backup(fullName string) (string, error) {
	newName := c.newFullName(fullName)
	if err := c.copyImage(fullName, newName); err != nil {
		return "", err
	}
	return newName, nil
}

// newFullName compacts the given image name to a single repository
// and prepends the backup destination, tag remains untouched
func (c *Client) newFullName(fullName string) string {
	return fmt.Sprintf("%s/%s/%s", registryAliases[c.registry][0], c.organization, strings.ReplaceAll(strings.TrimSpace(fullName), "/", "-"))
}

// copyImage mirrors the image from source to destination
func (c *Client) copyImage(src, dst string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.copyTimeoutSeconds)*time.Second)
	defer cancel()
	cmdStr := c.skopeoCopyCmd(src, dst)
	logf.Log.WithName("registry_client").V(1).Info("Command", cmdStr)
	cmdSl := strings.Split(cmdStr, " ")
	return exec.CommandContext(ctx, cmdSl[0], cmdSl[1:]...).Run()
}

// skopeoCopyCmd constructs skopeo copy command
func (c *Client) skopeoCopyCmd(src, dst string) string {
	cred := fmt.Sprintf("%s:%s", c.username, c.password)
	src = fmt.Sprintf("%s%s", c.transport, src)
	dst = fmt.Sprintf("%s%s", c.transport, dst)
	return fmt.Sprintf("skopeo copy --dest-creds %s %s %s", cred, src, dst)
}
