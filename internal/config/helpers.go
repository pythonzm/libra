package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

const (
	ConfigKeyConfigFile  = "config"
	ConfigKeyGitHubToken = "github_token"
	ConfigKeyApiUrl      = "apiurl"
	ConfigKeyModel       = "model"
	ConfigKeyScriptPath  = "script_path"
	ConfigKeyChatPath    = "chat_path"
	ConfigKeyReqTimeout  = "req_timeout"
	ConfigKeyIdleTimeout = "idle_timeout"
	ConfigKeySeekPrompt  = "seek_prompt"
	ConfigKeyExecPrompt  = "exec_prompt"

	ConfigFlagConfigFile  = "config"
	ConfigFlagGitHubToken = "github-token"
	ConfigFlagApiUrl      = "apiurl"
	ConfigFlagModel       = "model"
	ConfigFlagScriptPath  = "script-path"
	ConfigFlagChatPath    = "chat-path"
	ConfigFlagReqTimeout  = "req-timeout"
	ConfigFlagPrompt      = "prompt"
	ConfigFlagIdleTimeout = "idle-timeout"

	DefaultGitHubToken = "YOUR_GITHUB_TOKEN"
	DefaultAPIUrl      = "https://models.github.ai/inference/chat/completions"
	DefaultModel       = "openai/gpt-4o-mini"
	DefaultScriptPath  = "~/libra_data/scripts"
	DefaultChatPath    = "~/libra_data/chats"
	DefaultReqTimeout  = 20 * time.Second
	DefaultIdleTimeout = 5 * time.Minute

	DefaultSeekPrompt = `
请扮演命令行助手。用户会提出希望实现的功能，你的任务是只返回相关的命令。

要求如下：
- 根据操作系统信息输出对应命令。
- 回复时不要使用markdown语法（即不加代码块符号、标题、强调等格式标记）。
- 只返回命令，不要添加任何解释或额外的文字。
- 如果用户的请求不明确，询问用户以获取更多信息。
- 如果用户的请求不合法，告诉用户并提供正确的命令格式。
- 如果用户的请求不完整，询问用户以获取更多信息。
- 如果用户的请求不符合常规用法，告诉用户并提供正确的命令格式。

以下是一些示例：
用户：我想要在当前目录下创建一个名为test.txt的文件，并写入内容"Hello World"
助手：echo "Hello World" > test.txt
	`
	DefaultExecPrompt = `
请扮演命令行助手。用户会提出希望实现的功能，你的任务是只返回相关的命令。并以脚本的形式返回。

要求如下：
- 根据操作系统信息输出对应脚本格式：
  - Linux：bash脚本（以 “#!/bin/bash” 开头）。
  - Windows：bat脚本（以 “@echo off” 开头）。
  - MacOS：zsh脚本（以 “#!/bin/zsh” 开头）。
- 回复时不要使用markdown语法（即不加代码块符号、标题、强调等格式标记）。
- 只返回命令，不要添加任何解释或额外的文字。
- 如果用户的请求不明确，询问用户以获取更多信息。
- 如果用户的请求不合法，告诉用户并提供正确的命令格式。
- 如果用户的请求不完整，询问用户以获取更多信息。
- 如果用户的请求不符合常规用法，告诉用户并提供正确的命令格式。

以下是一些示例：
Linux用户：我想要在当前目录下创建一个名为test.txt的文件，并写入内容"Hello World"
助手：
#!/bin/bash
echo "Hello World" > test.txt
Windows用户：我想要在当前目录下创建一个名为test.txt的文件，并写入内容"Hello World"
助手：
@echo off
echo "Hello World" > test.txt
MacOS用户：我想要在当前目录下创建一个名为test.txt的文件，并写入内容"Hello World"
助手：
#!/bin/zsh
echo "Hello World" > test.txt
	`
)

var DefaultConfigFile string = filepath.Join("~/.libra.conf")

func GetFinalConfigValue(
	c *cli.Context,
	flagName string,
	configMap map[string]any,
	configKey string,
	defaultValue any,
) any {
	// 1. 环境变量优先
	envKey := fmt.Sprintf("LIBRA_%s", strings.ReplaceAll(strings.ToUpper(flagName), "-", "_"))
	if envVal, ok := os.LookupEnv(envKey); ok && envVal != "" {
		switch defaultValue.(type) {
		case time.Duration:
			if d, err := time.ParseDuration(envVal); err == nil {
				return d
			}
			return defaultValue
		default:
			return envVal
		}
	}

	// 2. CLI flag 覆盖
	if c.IsSet(flagName) {
		switch defaultValue.(type) {
		case time.Duration:
			raw := c.String(flagName)
			if d, err := parseDurationRaw(raw); err != nil {
				fmt.Printf("警告: 无法解析 '%s' 的值 '%s': %v\n", flagName, raw, err)
				return defaultValue
			} else {
				return d
			}
		default:
			return c.String(flagName)
		}
	}

	// 3. 配置文件中的值
	if val, ok := configMap[configKey]; ok {
		switch def := defaultValue.(type) {
		case time.Duration:
			switch v := val.(type) {
			case time.Duration:
				return v
			case string:
				if d, err := time.ParseDuration(v); err == nil {
					return d
				}
				return def
			default:
				return def
			}
		default:
			return val
		}
	}

	// 4. 回退到默认值
	return defaultValue
}
