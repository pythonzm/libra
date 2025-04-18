package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

type Config struct {
	GitHubToken string
	APIUrl      string
	Model       string
	ScriptPath  string
	ChatPath    string
	ReqTimeout  time.Duration
	IdleTimeout time.Duration
	SeekPrompt  string
	ExecPrompt  string
}

func LoadConfig(filePath string) (*Config, error) {
	cfg := new(Config) // Start with an empty config

	resolvedPath, err := ResolvePath(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法解析配置文件路径 '%s': %v", filePath, err)
	}

	if _, err := os.Stat(resolvedPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件 '%s' 不存在，将使用环境变量/参数/默认值", filePath)
	}

	iniCfg, err := ini.LoadSources(ini.LoadOptions{
		Insensitive:                true, // key 不区分大小写
		AllowShadows:               true,
		AllowPythonMultilineValues: true,
		IgnoreContinuation:         true,
		IgnoreInlineComment:        false, // 忽略行内注释
	}, resolvedPath) // Use resolvedPath

	if err != nil {
		// Return the error if loading/parsing failed for an existing file
		return nil, fmt.Errorf("无法加载或解析配置文件 '%s': %w", resolvedPath, err)
	}

	defaultSection := iniCfg.Section("") // Get the default section

	// Map keys, using Key().String() which returns "" if key missing
	cfg.GitHubToken = defaultSection.Key(ConfigKeyGitHubToken).String()
	cfg.APIUrl = defaultSection.Key(ConfigKeyApiUrl).String()
	cfg.Model = defaultSection.Key(ConfigKeyModel).String()
	cfg.ScriptPath = defaultSection.Key(ConfigKeyScriptPath).String()
	cfg.ChatPath = defaultSection.Key(ConfigKeyChatPath).String()

	reqTimeout := defaultSection.Key(ConfigKeyReqTimeout).String()
	if reqTimeout != "" {
		if d, err := parseDurationRaw(reqTimeout); err == nil {
			cfg.ReqTimeout = d
		} else {
			return nil, fmt.Errorf("无法解析请求超时 '%s': %w", reqTimeout, err)
		}
	} else {
		cfg.ReqTimeout = DefaultReqTimeout
	}
	idleTimeout := defaultSection.Key(ConfigKeyIdleTimeout).String()
	if idleTimeout != "" {
		if d, err := parseDurationRaw(idleTimeout); err == nil {
			cfg.IdleTimeout = d
		} else {
			return nil, fmt.Errorf("无法解析空闲超时 '%s': %w", idleTimeout, err)
		}
	} else {
		cfg.IdleTimeout = DefaultIdleTimeout
	}

	// Handle multi-line prompts
	cfg.SeekPrompt = defaultSection.Key(ConfigKeySeekPrompt).String()
	cfg.ExecPrompt = defaultSection.Key(ConfigKeyExecPrompt).String()

	return cfg, nil
}

func parseDurationRaw(raw string) (time.Duration, error) {
	if sec, err := strconv.Atoi(raw); err == nil {
		return time.Duration(sec) * time.Second, nil
	}
	return time.ParseDuration(raw)
}

func ConfigMap(cfg *Config) map[string]any {
	if cfg == nil {
		return map[string]any{
			ConfigKeyGitHubToken: DefaultGitHubToken,
			ConfigKeyApiUrl:      DefaultAPIUrl,
			ConfigKeyModel:       DefaultModel,
			ConfigKeyScriptPath:  DefaultScriptPath,
			ConfigKeyChatPath:    DefaultChatPath,
			ConfigKeyReqTimeout:  DefaultReqTimeout,
			ConfigKeyIdleTimeout: DefaultIdleTimeout,
			ConfigKeySeekPrompt:  DefaultSeekPrompt,
			ConfigKeyExecPrompt:  DefaultExecPrompt,
		}
	}
	return map[string]any{
		ConfigKeyGitHubToken: cfg.GitHubToken,
		ConfigKeyApiUrl:      cfg.APIUrl,
		ConfigKeyModel:       cfg.Model,
		ConfigKeyScriptPath:  cfg.ScriptPath,
		ConfigKeyChatPath:    cfg.ChatPath,
		ConfigKeyReqTimeout:  cfg.ReqTimeout,
		ConfigKeyIdleTimeout: cfg.IdleTimeout,
		ConfigKeySeekPrompt:  cfg.SeekPrompt,
		ConfigKeyExecPrompt:  cfg.ExecPrompt,
	}
}

func ResolvePath(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("无法获取用户主目录: %w", err)
	}
	if path == "~" {
		return homeDir, nil
	}
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:]), nil
	}
	return "", fmt.Errorf("不支持的波浪线扩展格式: %s", path)
}
