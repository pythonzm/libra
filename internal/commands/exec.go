package commands

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"libra/internal/config"
	"libra/internal/utils"

	urfaveCli "github.com/urfave/cli/v2"
)

func ExecAction(c *urfaveCli.Context) error {
	// 加载配置
	configMap, _ := LoadAndMergeConfig(c)
	conf, err := PrepareRequestConf(c, configMap)
	if err != nil {
		log.Fatal(err)
	}

	conf.ScriptPath = config.GetFinalConfigValue(c, config.ConfigFlagScriptPath, configMap, config.ConfigKeyScriptPath, config.DefaultScriptPath).(string)
	conf.ExecPrompt = config.GetFinalConfigValue(c, config.ConfigFlagPrompt, configMap, config.ConfigKeyExecPrompt, config.DefaultExecPrompt).(string)

	// 获取脚本存放目录和聊天记录目录
	chatDatePath, err := CreateDateFolder(conf.ChatPath)
	if err != nil {
		log.Fatalf("创建chat folder失败: %v", err)
	}
	scriptDatePath, err := CreateDateFolder(conf.ScriptPath)
	if err != nil {
		log.Fatalf("创建script folder失败: %v", err)
	}

	messages := []Message{NewSystemMessage(conf.ExecPrompt)}
	inputCh, errCh := InputWatcher()
	fmt.Println("命令行助手已开启. 输入 'q', 'quit' 或 'exit' 退出.")
	var chatFile *os.File
	var lastAssistantFile *os.File
	defer func() {
		if chatFile != nil {
			chatFile.Close()
		}
		if lastAssistantFile != nil {
			lastAssistantFile.Close()
		}
	}()

	var chatID string
	// 标志是否已获得有效的脚本内容（由助手返回）
	var scriptReady bool = false

	for {
		promptText := ">>> "
		if scriptReady {
			promptText = ">>> 是否执行脚本内容【y/n/q】 "
		}

		fmt.Print(promptText)
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
			userInputLower := strings.ToLower(userInput)
			if isExit(userInput) {
				return nil
			}
			// 如果脚本已经准备好，则这里响应用户是否执行脚本的选项
			if scriptReady {
				if userInputLower == "y" && lastAssistantFile != nil {
					// 执行脚本
					utils.RunScript(lastAssistantFile.Name())
					// 执行后可以选择重置 scriptReady 状态，或根据需求保留当前脚本以供多次执行
					scriptReady = false
					continue
				} else if userInputLower == "n" {
					// 用户选择不执行，将 scriptReady 重置为 false，等待下次有效脚本返回后再询问
					scriptReady = false
					// 继续让用户输入新的内容
					continue
				}
				// 如果输入既不属于 y/n/q，则按新一轮对话处理
			}

			// 第一次输入时附加操作系统信息
			if len(messages) == 1 {
				osInfo := utils.GetOSPlatformInfo()
				userInput = userInput + fmt.Sprintf("\n当前操作系统信息是: %s", osInfo)
			}

			messages = append(messages, NewUserMessage(userInput))
			userMessageEntry := NewMessageEntry(time.Now().Format("2006-01-02 15:04:05"), NewUserMessage(userInput))

			// 创建一个临时文件用于接收 Assistant 返回的内容
			tempFile, err := os.CreateTemp("", "assistant_exec_*")
			if err != nil {
				log.Fatalf("创建临时文件失败: %v", err)
			}
			defer os.Remove(tempFile.Name())

			// 定义写入回调，同时写入终端和临时文件
			writeContent := func(content string) {
				fmt.Print(content)
				tempFile.WriteString(content)
			}

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
			// 检查内容是否有效（非空），标记为 scriptReady
			s := strings.TrimSpace(string(assistantContent))
			scriptReady = s != "" && (strings.HasPrefix(s, "#") || strings.HasPrefix(s, "@"))
			messages = append(messages, assistantMessage)

			// 将对话记录写入 chat log 文件
			writeJSONLine(chatFile, userMessageEntry)
			writeJSONLine(chatFile, NewMessageEntry(time.Now().Format("2006-01-02 15:04:05"), assistantMessage))

			if lastAssistantFile == nil && scriptReady {
				var scriptFileName string
				if chatID != "" {
					if runtime.GOOS == "windows" {
						scriptFileName = fmt.Sprintf("%s_%s.bat", chatID, time.Now().Format("150405"))
					} else {
						scriptFileName = fmt.Sprintf("%s_%s.sh", chatID, time.Now().Format("150405"))
					}
				} else {
					if runtime.GOOS == "windows" {
						scriptFileName = fmt.Sprintf("%s.bat", time.Now().Format("150405"))
					} else {
						scriptFileName = fmt.Sprintf("%s.sh", time.Now().Format("150405"))
					}
				}
				lastAssistantFile, err = os.Create(path.Join(scriptDatePath, scriptFileName))
				if err != nil {
					log.Fatalf("创建脚本文件失败: %v", err)
				}
			}

			// 更新脚本文件内容
			if lastAssistantFile != nil {
				lastAssistantFile.Truncate(0)
				lastAssistantFile.Seek(0, 0)
				lastAssistantFile.Write(assistantContent)
				lastAssistantFile.Sync()
			}
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
