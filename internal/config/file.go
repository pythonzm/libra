package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	urfaveCli "github.com/urfave/cli/v2"
	"gopkg.in/ini.v1"
)

func CreateDefaultConfigFile(c *urfaveCli.Context) error {
	iniFile := ini.Empty()
	section := iniFile.Section("")
	createKey := func(name, defaultValue, comment string) {
		var value string
		flagName := strings.ReplaceAll(name, "_", "-")
		if c.IsSet(flagName) {
			value = c.String(flagName)
		} else {
			value = defaultValue
		}
		key, err := section.NewKey(name, value)
		if err == nil && comment != "" {
			key.Comment = comment
		} else if err != nil {
			fmt.Printf("警告: 无法为 '%s' 创建 key: %v", name, err)
		}
	}

	createKey(ConfigKeyGitHubToken, DefaultGitHubToken, "")
	createKey(ConfigKeyApiUrl, DefaultAPIUrl, "")
	createKey(ConfigKeyModel, DefaultModel, "")
	createKey(ConfigKeyScriptPath, DefaultScriptPath, "存放执行脚本的目录")
	createKey(ConfigKeyChatPath, DefaultChatPath, "存放对话历史的目录")
	createKey(ConfigKeyReqTimeout, DefaultReqTimeout.String(), "请求超时时间")
	createKey(ConfigKeyIdleTimeout, DefaultIdleTimeout.String(), "空闲超时时间")

	createKey(ConfigKeySeekPrompt, DefaultSeekPrompt, "")
	createKey(ConfigKeyExecPrompt, DefaultExecPrompt, "")

	var filePath string
	if c.IsSet(ConfigFlagConfigFile) {
		filePath, _ = ResolvePath(c.String(ConfigFlagConfigFile))
	} else {
		filePath, _ = ResolvePath(DefaultConfigFile)
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("无法创建目录 '%s' 用于配置文件: %w", dir, err)
	}

	if err := iniFile.SaveTo(filePath); err != nil {
		return fmt.Errorf("无法将配置写入文件 '%s': %w", filePath, err)
	}
	fmt.Printf("配置已成功写入文件: %s\n", filePath)

	return nil
}
