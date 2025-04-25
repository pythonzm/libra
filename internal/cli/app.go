package cli

import (
	"libra/internal/commands"
	"libra/internal/config"

	"slices"

	"github.com/urfave/cli/v2"
)

const (
	appName    = "libra"
	appUsage   = "A AI helper for CLI"
	appVersion = "1.3"
)

var GlobalFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    config.ConfigFlagConfigFile,
		Value:   config.DefaultConfigFile,
		Aliases: []string{"c"},
		Usage:   "指定配置文件路径",
	},
	&cli.StringFlag{
		Name:    config.ConfigFlagGitHubToken,
		Value:   config.DefaultGitHubToken,
		Aliases: []string{"g"},
		Usage:   "指定github token",
	},
	&cli.StringFlag{
		Name:    config.ConfigFlagApiUrl,
		Value:   config.DefaultAPIUrl,
		Aliases: []string{"u"},
		Usage:   "指定api地址，一般情况下默认即可",
	},
	&cli.StringFlag{
		Name:    config.ConfigFlagReqTimeout,
		Value:   config.DefaultReqTimeout.String(),
		Aliases: []string{"t"},
		Usage:   "设置请求超时时间",
	},
	&cli.StringFlag{
		Name:    config.ConfigFlagIdleTimeout,
		Value:   config.DefaultIdleTimeout.String(),
		Aliases: []string{"T", "idle"},
		Usage:   "设置空闲超时时间",
	},
	&cli.StringFlag{
		Name:    config.ConfigFlagModel,
		Value:   config.DefaultModel,
		Aliases: []string{"m"},
		Usage:   "指定模型，更多可用模型：https://github.com/marketplace?type=models",
	},
	&cli.StringFlag{
		Name:    config.ConfigFlagChatPath,
		Value:   config.DefaultChatPath,
		Aliases: []string{"C"},
		Usage:   "指定对话内容存储路径",
	},
	&cli.StringFlag{
		Name:    config.ConfigFlagPrompt,
		Value:   "",
		Aliases: []string{"p"},
		Usage:   "指定自定义prompt",
	},
}

func SetupApp() *cli.App {
	app := &cli.App{
		Name:                 appName,
		Usage:                appUsage,
		Version:              appVersion,
		Flags:                GlobalFlags,
		EnableBashCompletion: true,
		Action:               commands.SeekAction,

		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "生成配置文件",
				Flags:   slices.Clone(GlobalFlags),
				Action:  commands.InitAction,
			},
			{
				Name:    "exec",
				Aliases: []string{"e"},
				Flags: append([]cli.Flag{
					&cli.StringFlag{
						Name:    config.ConfigFlagScriptPath,
						Value:   config.DefaultScriptPath,
						Aliases: []string{"S"},
						Usage:   "运行脚本存储路径",
					},
				}, GlobalFlags...),
				Usage:  "将AI回答的答案存储至脚本文件用于执行",
				Action: commands.ExecAction,
			},
		},

		// Action: func(c *cli.Context) error {
		// 	fmt.Printf("%s: 无效的指令. 使用 --help/-h/h 获取帮助.\n", appName)
		// 	return nil
		// },
	}
	return app
}
