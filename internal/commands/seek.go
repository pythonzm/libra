package commands

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"libra/internal/config"
	"libra/internal/utils"

	urfaveCli "github.com/urfave/cli/v2"
)

func SeekAction(c *urfaveCli.Context) error {
	// --- 获取配置和参数 ---
	configMap, _ := LoadAndMergeConfig(c)

	conf, err := PrepareRequestConf(c, configMap)
	if err != nil {
		log.Fatal(err)
	}
	conf.SeekPrompt = config.GetFinalConfigValue(c, config.ConfigFlagPrompt, configMap, config.ConfigKeySeekPrompt, config.DefaultSeekPrompt).(string)

	chatDatePath, err := CreateDateFolder(conf.ChatPath)
	if err != nil {
		log.Fatalf("创建chat folder失败: %v", err)
	}

	messages := []Message{NewSystemMessage(conf.SeekPrompt)}
	inputCh, errCh := InputWatcher()
	fmt.Println("命令行助手已开启. 输入 'q', 'quit' 或者 'exit' 退出.")

	var chatFile *os.File
	defer func() {
		if chatFile != nil {
			chatFile.Close()
		}
	}()

	var chatID string

	for {
		fmt.Print(">>> ")
		timer := time.NewTimer(conf.IdleTimeout)
		select {
		case line, ok := <-inputCh:
			timer.Stop()
			if !ok {
				// 输入流关闭
				return nil
			}

			userInput := strings.TrimSpace(line)
			if userInput == "" {
				continue
			}
			if isExit(userInput) {
				return nil
			}

			// 第一次输入时附加操作系统信息
			if len(messages) == 1 {
				osInfo := utils.GetOSPlatformInfo()
				userInput = userInput + fmt.Sprintf("\n当前操作系统信息是: %s", osInfo)
			}

			messages = append(messages, NewUserMessage(userInput))
			userMessageEntry := NewMessageEntry(time.Now().Format("2006-01-02 15:04:05"), NewUserMessage(userInput))

			// 准备一个临时文件接收 Assistant 返回的片段
			tempFile, err := os.CreateTemp("", "assistant_*")
			if err != nil {
				log.Fatalf("创建临时文件失败: %v", err)
			}
			defer os.Remove(tempFile.Name())

			// 定义写入回调，把响应内容同时打印到终端和写入临时文件
			writeContent := func(content string) {
				fmt.Print(content)
				tempFile.WriteString(content)
			}

			// 调用公共的发送请求函数
			err = SendStreamRequest(c, conf, messages, &chatID, writeContent)
			if err != nil {
				log.Fatalf("请求处理出错: %v", err)
			}

			var chatFileName string
			if chatID != "" {
				chatFileName = fmt.Sprintf("%s_%s.jsonl", chatID, time.Now().Format("150405"))
			} else {
				chatFileName = fmt.Sprintf("%s.jsonl", time.Now().Format("150405"))
			}

			// 记录对话为 JSONL 日志
			if chatFile == nil {
				chatFile, err = os.OpenFile(path.Join(chatDatePath, chatFileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Fatalf("打开 chat file 失败: %v", err)
				}
			}

			// 读取 Assistant 的响应内容并追加到消息中
			tempFile.Seek(0, 0)
			assistantContent, _ := io.ReadAll(tempFile)
			assistantMessage := NewAssistantMessage(string(assistantContent))
			messages = append(messages, assistantMessage)

			// 将对话记录写入日志文件
			writeJSONLine(chatFile, userMessageEntry)
			writeJSONLine(chatFile, NewMessageEntry(time.Now().Format("2006-01-02 15:04:05"), assistantMessage))
		case err := <-errCh:
			if err != nil {
				log.Printf("读取输入错误: %v", err)
			}
			return nil
		case <-timer.C:
			fmt.Printf("\n[%s] 超过 %v 无输入，自动断开。\n",
				time.Now().Format("15:04:05"), conf.IdleTimeout)
			return nil
		}
	}
}
