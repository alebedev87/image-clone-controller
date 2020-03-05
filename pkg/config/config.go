package config

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

func init() {
	// registering the program's flags
	pflag.StringVar(&GlobalConfig.Registry, "backup-registry", "", "Backup image registry.")
	pflag.StringVar(&GlobalConfig.Organization, "registry-org", "", "Backup image registry's organization.")
	pflag.StringVar(&GlobalConfig.Username, "registry-username", "", "Username to access the backup image registry.")
	pflag.StringVar(&GlobalConfig.Password, "registry-password", "", "Password to access the backup image registry.")
	pflag.StringSliceVar(&GlobalConfig.AdditionalNamespaceBlacklist, "additional-namespace-blacklist", []string{}, "List of namespace(s) which should NOT be watched.")
	pflag.IntVar(&GlobalConfig.ImageCopyTimeoutSeconds, "img-copy-timeout", defaultImageCopyTimeout, "Timeout for the copy of a single image to the backup registry (in seconds).")
}

const (
	usernameVar             = "IMG_CTR_REGISTRY_USERNAME"
	passwordVar             = "IMG_CTR_REGISTRY_PASSWORD"
	defaultImageCopyTimeout = 60 * 60
)

// GlobalConfig is all program's config
var GlobalConfig *Config = &Config{
	MandatoryNamespaceBlacklist: []string{"kube-system"},
}

// Config stores the configuration to the whole program
type Config struct {
	Registry                     string
	Organization                 string
	Username                     string
	Password                     string
	ImageCopyTimeoutSeconds      int
	MandatoryNamespaceBlacklist  []string
	AdditionalNamespaceBlacklist []string
}

// Validate validates the important fields of the configuration
func (c *Config) Validate() error {
	if len(strings.TrimSpace(c.Registry)) == 0 {
		return errors.New("no backup image registry provided")
	}

	if len(strings.TrimSpace(c.Organization)) == 0 {
		return errors.New("no organization for backup image registry provided")
	}

	if len(strings.TrimSpace(c.Username)) == 0 {
		c.Username = os.Getenv(usernameVar)
		if len(strings.TrimSpace(c.Username)) == 0 {
			return errors.New("no userame to access backup image registry provided")
		}
	}

	if len(strings.TrimSpace(c.Password)) == 0 {
		c.Password = os.Getenv(passwordVar)
		if len(strings.TrimSpace(c.Password)) == 0 {
			return errors.New("no pasword to access backup image registry provided")
		}
	}

	return nil
}

// NamespaceBlacklist returns a set of all blacklisted namespaces
func (c *Config) NamespaceBlacklist() map[string]bool {
	set := map[string]bool{}
	for _, n := range c.MandatoryNamespaceBlacklist {
		set[n] = true
	}
	for _, n := range c.AdditionalNamespaceBlacklist {
		set[n] = true
	}
	return set
}
