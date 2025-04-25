# libra

命令行下的AI助理

## 安装

在 [release](https://github.com/pythonzm/libra/releases) 页最新版本，找到适合自己机器的安装包下载即可

### 要求

本工具使用 [GitHub Models](https://docs.github.com/en/github-models) 编写，需要提前准备GitHub Token，至少要拥有 ` models:read` 权限

配置 `GitHub Token` 

```bash
# 可以通过环境变量
$ export LIBRA_GITHUB_TOKEN=xxxxxx

# 也可以先生成配置文件，在配置文件中修改
$ ./libra i
配置已成功写入文件: /home/zm/.libra.conf

# 也可以通过命令行参数指定
$ ./libra s -g xxxxxxxxxx
```

## 使用

[![asciicast](https://asciinema.org/a/716947.svg)](https://asciinema.org/a/716947)

### 示例

#### 查看帮助信息

```bash
$ ./libra h
NAME:
   libra - A AI helper for CLI

USAGE:
   libra [global options] command [command options]

VERSION:
   1.3

COMMANDS:
   init, i  生成配置文件
   exec, e  将AI回答的答案存储至脚本文件用于执行
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value                      指定配置文件路径 (default: "/home/zm/.libra.conf")
   --github-token value, -g value                指定github token (default: "YOUR_GITHUB_TOKEN")
   --apiurl value, -u value                      指定api地址，一般情况下默认即可 (default: "https://models.github.ai/inference/chat/completions")
   --req-timeout value, -t value                 设置请求超时时间 (default: "20s")
   --idle-timeout value, -T value, --idle value  设置空闲超时时间 (default: "5m0s")
   --model value, -m value                       指定模型，更多可用模型：https://github.com/marketplace?type=models (default: "openai/gpt-4o-mini")
   --chat-path value, -C value                   指定对话内容存储路径 (default: "/home/zm/libra_data/chats")
   --prompt value, -p value                      指定自定义prompt
   --help, -h                                    show help
   --version, -v                                 print the version
```

#### 生成配置文件

```bash
$ ./libra i
配置已成功写入文件: /home/zm/.libra.conf
```

#### 咨询问题

```bash
$ ./libra
2025/04/18 10:44:58 common.go:26: 警告: 解析配置文件 '~/.libra.conf' 失败: 配置文件 '~/.libra.conf' 不存在，将使用环境变量/参数/默认值
命令行助手已开启. 输入 'q', 'quit' 或者 'exit' 退出.
>>> 查看当前占用内存最多的5个进程
ps aux --sort=-%mem | head -n 6
>>> 占用CPU最多的5个进程
ps aux --sort=-%cpu | head -n 6
>>> q
$ 
```

#### 执行脚本

```bash
$ ./libra e
2025/04/18 10:45:48 common.go:26: 警告: 解析配置文件 '~/.libra.conf' 失败: 配置文件 '~/.libra.conf' 不存在，将使用环境变量/参数/默认值
命令行助手已开启. 输入 'q', 'quit' 或 'exit' 退出.
>>> 
>>> 查看当前占用内存最多的5个进程
#!/bin/bash
ps aux --sort=-%mem | head -n 6
>>> 是否执行脚本内容【y/n/q】 y
2025/04/18 10:46:08 platform.go:57: 开始运行 Shell (bash) 脚本: /home/zm/libra_data/scripts/2025-04-18/chatcmpl-BNW1uKrAzEtZ9jR04hkqa0Tx2SkyN_104605.sh
2025/04/18 10:46:08 platform.go:70: Shell (bash) 脚本成功执行。
--- 输出 ---
USER         PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
zm        978232  1.3  6.6 76324004 530088 pts/0 Sl+  Apr17  12:53 /home/zm/.vscode-server/bin/4949701c880d4bdb949e3c0e6b400288da7f474b/node --dns-result-order=ipv4first /home/zm/.vscode-server/bin/4949701c880d4bdb949e3c0e6b400288da7f474b/out/bootstrap-fork --type=extensionHost --transformURIs --useHostProxy=true
zm        978496  0.0  2.3 3152468 189852 pts/0  Sl+  Apr17   0:20 /home/zm/go/bin/gopls -mode=stdio
zm           351  0.0  1.4 11822156 114216 pts/0 Sl+  Apr14   3:05 /home/zm/.vscode-server/bin/4949701c880d4bdb949e3c0e6b400288da7f474b/node /home/zm/.vscode-server/bin/4949701c880d4bdb949e3c0e6b400288da7f474b/out/server-main.js --host=127.0.0.1 --port=0 --connection-token=2117758677-4276087677-485001766-151388202 --use-host-proxy --without-browser-env-var --disable-websocket-compression --accept-server-license-terms --telemetry-level=all
zm        978307  0.0  1.0 1025256 84892 pts/0   Sl+  Apr17   0:02 /home/zm/.vscode-server/bin/4949701c880d4bdb949e3c0e6b400288da7f474b/node /home/zm/.vscode-server/extensions/github.vscode-github-actions-0.27.1/dist/server-node.js --node-ipc --clientProcessId=978232
zm           458  0.1  0.9 11620172 78420 pts/0  Sl+  Apr14   5:57 /home/zm/.vscode-server/bin/4949701c880d4bdb949e3c0e6b400288da7f474b/node /home/zm/.vscode-server/bin/4949701c880d4bdb949e3c0e6b400288da7f474b/out/bootstrap-fork --type=ptyHost --logsPath /home/zm/.vscode-server/data/logs/20250414T100659

----------
>>> q
```

#### 其他用途

```bash
$ ./libra -p "你是一名出色的解梦大师，精通各种梦境解读"
2025/04/18 10:50:31 common.go:26: 警告: 解析配置文件 '~/.libra.conf' 失败: 配置文件 '~/.libra.conf' 不存在，将使用环境变量/参数/默认值
命令行助手已开启. 输入 'q', 'quit' 或者 'exit' 退出.
>>> 
>>> 梦到牙齿掉落
梦见牙齿掉落是一种比较常见的梦境，通常被认为与个人的焦虑、烦恼或压力有关。不同文化对这个梦的解释可能存在差异，但以下是一些普遍的解读：

1. **失去和不安**：梦到牙齿掉落可能象征着对失去的恐惧，或者对自我形象的担忧。这种梦境常常出现在生活中面临重大变化或压力时。

2. **沟通问题**：牙齿在梦中也可能代表表达能力，掉落的牙齿可能意味着你在表达自己方面感到无能或困惑。

3. **衰老和时间流逝**：牙齿掉落也可能象征对衰老的焦虑，或对时间流逝的敏感。

4. **生活变化**：这种梦也可能预示着即将发生的一些重要变化，无论是积极的还是消极的。

为了更深入解读，可以结合你最近的生活状态、情绪和心理压力。如果最近经历了一些重大事件或变化，这可能帮助你更好地理解这个梦的含义。
>>> q
$ 
```

## 注意

既然AI接口可以免费使用，那肯定是有限制的，具体哪个模型有多少限制，可查看官方文档：[https://docs.github.com/en/github-models/prototyping-with-ai-models#rate-limits](https://docs.github.com/en/github-models/prototyping-with-ai-models#rate-limits)

对于个人使用个人认为完全够用，一个模型达到限制之后继续换另一个模型就好了

最新的接口地址已更新为：https://models.github.ai/inference/chat/completions，国内访问有时候可能会超时，可以使用原来的接口地址：https://models.inference.ai.azure.com/chat/completions，这个地址国内访问目前比较稳定，但是需要修改模型名称，比如 `openai/gpt-4o-mini` 要换成 `gpt-4o-mini` 

## License
[Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0)