package utils

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func GetOSPlatformInfo() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}

func executeScript(interpreter string, scriptPath string, interpreterArgs ...string) (string, error) {
	args := append(interpreterArgs, scriptPath)

	// 创建命令对象
	cmd := exec.Command(interpreter, args...)

	outputBytes, err := cmd.CombinedOutput()
	outputString := string(outputBytes) // 将字节切片转换为字符串

	// CombinedOutput 会在命令成功执行（退出码为0）时返回 err == nil
	// 如果命令以非零状态退出，err 会是 *exec.ExitError 类型
	// 如果命令无法启动（例如找不到解释器），err 会是其他类型
	if err != nil {
		// 返回捕获到的输出和错误信息
		return outputString, fmt.Errorf("执行 %s %s 时出错: %w", interpreter, scriptPath, err)
	}

	// 命令成功执行，返回输出和 nil 错误
	return outputString, nil
}

func RunScript(scriptPath string) error {
	osType := runtime.GOOS

	var interpreter string
	var scriptType string // 用于日志记录
	var interpreterArgs []string

	switch osType {
	case "windows":
		interpreter = "cmd.exe"
		interpreterArgs = []string{"/c"} // cmd.exe 需要 /c 参数
		scriptType = "Bat"
	case "linux":
		interpreter = "bash" // Linux 通常使用 bash
		scriptType = "Shell (bash)"
	case "darwin": // macOS
		interpreter = "zsh"
		scriptType = "Shell (zsh)"
	default:
		return fmt.Errorf("不支持的操作系统: %s", osType) // 返回错误而不是仅打印
	}

	log.Printf("开始运行 %s 脚本: %s", scriptType, scriptPath)

	// 调用核心执行函数
	output, err := executeScript(interpreter, scriptPath, interpreterArgs...)

	// 处理执行结果
	if err != nil {
		// 记录错误信息和脚本输出（即使出错也可能有输出）
		log.Printf("运行 %s 脚本出错: %v\n--- 输出 ---\n%s\n----------", scriptType, err, output)
		return err // 将错误传递给调用者
	}

	// 命令成功执行
	log.Printf("%s 脚本成功执行。\n--- 输出 ---\n%s\n----------", scriptType, output)
	return nil // 表示成功
}
