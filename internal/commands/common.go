// internal/commands/common.go
package commands

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"libra/internal/config"

	urfaveCli "github.com/urfave/cli/v2"
)

func LoadAndMergeConfig(c *urfaveCli.Context) (map[string]any, error) {
	configPath := c.String(config.ConfigFlagConfigFile)
	cfgData, err := config.LoadConfig(configPath)
	if err != nil {
		log.Printf("警告: 解析配置文件 '%s' 失败: %v", configPath, err)
	}
	return config.ConfigMap(cfgData), nil
}

func PrepareRequestConf(c *urfaveCli.Context, configMap map[string]any) (*config.Config, error) {
	conf := new(config.Config)
	// 必填参数：githubToken
	githubToken := config.GetFinalConfigValue(c, config.ConfigFlagGitHubToken, configMap, config.ConfigKeyGitHubToken, config.DefaultGitHubToken)
	if githubToken.(string) == "" || githubToken == config.DefaultGitHubToken {
		return nil, fmt.Errorf("GitHub token is required. 使用 --github-token 或设置环境变量 LIBRA_GITHUB_TOKEN 或使用 init 初始化配置文件并设置 token")
	}
	conf.GitHubToken = githubToken.(string)
	conf.APIUrl = config.GetFinalConfigValue(c, config.ConfigFlagApiUrl, configMap, config.ConfigKeyApiUrl, config.DefaultAPIUrl).(string)
	conf.Model = config.GetFinalConfigValue(c, config.ConfigFlagModel, configMap, config.ConfigKeyModel, config.DefaultModel).(string)
	conf.ChatPath = config.GetFinalConfigValue(c, config.ConfigFlagChatPath, configMap, config.ConfigKeyChatPath, config.DefaultChatPath).(string)
	conf.ReqTimeout = config.GetFinalConfigValue(c, config.ConfigFlagReqTimeout, configMap, config.ConfigKeyReqTimeout, config.DefaultReqTimeout).(time.Duration)
	conf.IdleTimeout = config.GetFinalConfigValue(c, config.ConfigFlagIdleTimeout, configMap, config.ConfigKeyIdleTimeout, config.DefaultIdleTimeout).(time.Duration)
	return conf, nil
}

func SendStreamRequest(c *urfaveCli.Context, conf *config.Config, messages []Message, streamIDPtr *string, writeContent func(content string)) error {
	reqBody := RequestBody{
		Messages: messages,
		Stream:   true,
		Model:    conf.Model,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("解析请求体失败: %w", err)
	}

	req, err := http.NewRequest("POST", conf.APIUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+conf.GitHubToken)

	client := &http.Client{
		Timeout: conf.ReqTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("请求失败，响应码： %d\n响应体: %s", resp.StatusCode, string(bodyBytes))
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		// 处理 SSE 数据格式
		if strings.HasPrefix(line, "data:") {
			jsonDataPart := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if jsonDataPart == "[DONE]" {
				break
			}
			var streamResp StreamResponse
			err := json.Unmarshal([]byte(jsonDataPart), &streamResp)
			if err != nil {
				log.Printf("Warning: Could not unmarshal stream chunk '%s': %v", jsonDataPart, err)
				continue
			}
			// 第一次捕获 streamResp.ID
			if len(streamResp.Choices) > 0 && *streamIDPtr == "" {
				*streamIDPtr = streamResp.ID
			}
			if len(streamResp.Choices) > 0 {
				content := streamResp.Choices[0].Delta.Content
				writeContent(content)
				// 检查是否流结束
				if streamResp.Choices[0].FinishReason != nil {
					fmt.Println()
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取响应体失败: %w", err)
	}
	return nil
}

func CreateDateFolder(basePath string) (string, error) {
	absPath, err := config.ResolvePath(basePath)
	if err != nil {
		return "", err
	}
	datePath := path.Join(absPath, time.Now().Format("2006-01-02"))
	if err := os.MkdirAll(datePath, 0755); err != nil {
		return "", err
	}
	return datePath, nil
}

func isExit(input string) bool {
	switch strings.ToLower(input) {
	case "q", "quit", "exit":
		return true
	default:
		return false
	}
}

func writeJSONLine(file *os.File, v interface{}) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Printf("序列化 JSON 时出错: %v", err)
		return
	}
	file.Write(jsonData)
	file.WriteString("\n")
	file.Sync()
}

func InputWatcher() (inputCh <-chan string, errCh <-chan error) {
	scanner := bufio.NewScanner(os.Stdin)
	in := make(chan string)
	er := make(chan error)

	go func() {
		defer close(in)
		defer close(er)

		for {
			if scanner.Scan() {
				in <- scanner.Text()
			} else {
				er <- scanner.Err()
				return
			}
		}
	}()

	return in, er
}
