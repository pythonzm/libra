package commands

import (
	"libra/internal/config"

	urfaveCli "github.com/urfave/cli/v2"
)

func InitAction(c *urfaveCli.Context) error {
	err := config.CreateDefaultConfigFile(c)
	return err
}
