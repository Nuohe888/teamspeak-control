package tool

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"regexp"
	"strings"
	"time"
)

type Ts struct {
	//服务器信息
	Host     string
	Port     int
	Username string
	Password string
	//TS需要用到参数
	Uuid             string
	DefaultVoicePort string
	QueryPort        string
	FiletransferPort string
	ImageName        string
}

func (i *Ts) connect() (*ssh.Client, error) {
	// 创建SSH客户端配置
	config := &ssh.ClientConfig{
		User: i.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(i.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// 连接到SSH服务器
	addr := fmt.Sprintf("%s:%d", i.Host, i.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("无法连接到服务器: %v", err)
	}

	return client, nil
}

// CheckServer 检查服务器是否存在docker 并且返回docker版本
func (i *Ts) CheckServer() error {
	client, err := i.connect()
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession() // 创建新的 session
	if err != nil {
		return fmt.Errorf("创建 CheckServer SSH 会话失败: %v", err)
	}
	defer session.Close()

	// 执行命令检查Docker是否安装
	_, err = session.CombinedOutput("docker --version")
	if err != nil {
		return fmt.Errorf("docker未安装或无法运行: %v", err)
	}

	// 检查Docker服务是否运行
	session2, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建 SSH 会话失败: %v", err)
	}
	defer session2.Close()

	statusOutput, err := session2.CombinedOutput("systemctl is-active docker")
	if err != nil || string(statusOutput) != "active\n" {
		return fmt.Errorf("docker服务未运行: %v", err)
	}

	return nil
}

// Run docker run ts
func (i *Ts) Run() error {
	// 连接到SSH服务器
	client, err := i.connect()
	if err != nil {
		return err
	}
	defer client.Close()

	// 创建新的会话
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建 Run SSH 会话失败: %v", err)
	}
	defer session.Close()

	// 使用提供的命令运行TeamSpeak容器
	name := "ts" + i.Uuid
	cmd := fmt.Sprintf("docker run --name %s -d -p %s:9987/udp -p %s:10011 -p %s:30033 -e TS3SERVER_LICENSE=accept %s",
		name, i.DefaultVoicePort, i.QueryPort, i.FiletransferPort, i.ImageName)

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return fmt.Errorf("运行TeamSpeak容器失败: %v, 输出: %s", err, string(output))
	}

	// 检查容器是否成功启动
	session2, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建 SSH 会话失败: %v", err)
	}
	defer session2.Close()

	checkCmd := fmt.Sprintf("docker ps -f name=%s --format '{{.Status}}'", name)
	statusOutput, err := session2.CombinedOutput(checkCmd)
	if err != nil || len(statusOutput) == 0 {
		return fmt.Errorf("TeamSpeak容器未成功启动: %v", err)
	}

	return nil
}

// Del 删除Ts
func (i *Ts) Del() error {
	// 连接到SSH服务器
	client, err := i.connect()
	if err != nil {
		return err
	}
	defer client.Close()

	// 容器名称
	name := "ts" + i.Uuid

	// 停止容器
	session1, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建 SSH 会话失败: %v", err)
	}
	stopCmd := fmt.Sprintf("docker stop %s || true", name)
	_, err = session1.CombinedOutput(stopCmd)
	session1.Close()
	// 忽略停止错误，因为容器可能已经停止

	// 删除容器
	session2, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建 SSH 会话失败: %v", err)
	}
	rmCmd := fmt.Sprintf("docker rm %s || true", name)
	_, err = session2.CombinedOutput(rmCmd)
	session2.Close()
	// 忽略删除错误，因为容器可能已经被删除

	return nil
}

// Restart 重启Ts
func (i *Ts) Restart() error {
	// 连接到SSH服务器
	client, err := i.connect()
	if err != nil {
		return err
	}
	defer client.Close()

	// 容器名称
	name := "ts" + i.Uuid

	// 重启容器
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建 SSH 会话失败: %v", err)
	}
	defer session.Close()

	restartCmd := fmt.Sprintf("docker restart %s", name)
	output, err := session.CombinedOutput(restartCmd)
	if err != nil {
		return fmt.Errorf("重启TeamSpeak容器失败: %v, 输出: %s", err, string(output))
	}

	// 检查容器是否成功重启
	session2, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建 SSH 会话失败: %v", err)
	}
	defer session2.Close()

	checkCmd := fmt.Sprintf("docker ps -f name=%s --format '{{.Status}}'", name)
	statusOutput, err := session2.CombinedOutput(checkCmd)
	if err != nil || len(statusOutput) == 0 {
		return fmt.Errorf("TeamSpeak容器未成功重启: %v", err)
	}

	return nil
}

// Stop 停止Ts
func (i *Ts) Stop() error {
	// 连接到SSH服务器
	client, err := i.connect()
	if err != nil {
		return err
	}
	defer client.Close()

	// 容器名称
	name := "ts" + i.Uuid

	// 停止容器
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建 SSH 会话失败: %v", err)
	}
	defer session.Close()

	stopCmd := fmt.Sprintf("docker stop %s", name)
	output, err := session.CombinedOutput(stopCmd)
	if err != nil {
		return fmt.Errorf("停止TeamSpeak容器失败: %v, 输出: %s", err, string(output))
	}

	// 检查容器是否已停止
	session2, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建 SSH 会话失败: %v", err)
	}
	defer session2.Close()

	checkCmd := fmt.Sprintf("docker ps -f name=%s -f status=running --format '{{.Status}}'", name)
	statusOutput, err := session2.CombinedOutput(checkCmd)
	if err == nil && len(statusOutput) > 0 {
		return fmt.Errorf("TeamSpeak容器未成功停止")
	}

	return nil
}

// Status Ts状态
func (i *Ts) Status() string {
	// 连接到SSH服务器
	client, err := i.connect()
	if err != nil {
		return "未知"
	}
	defer client.Close()

	// 容器名称
	name := "ts" + i.Uuid

	// 检查容器状态
	session, err := client.NewSession()
	if err != nil {
		return "未知"
	}
	defer session.Close()

	// 使用docker ps命令获取容器状态
	statusCmd := fmt.Sprintf("docker ps -a -f name=%s --format '{{.Status}}'", name)
	statusOutput, err := session.CombinedOutput(statusCmd)
	if err != nil || len(statusOutput) == 0 {
		return "未知"
	}

	// 解析状态输出
	status := string(statusOutput)

	// 根据Docker状态返回更友好的状态描述
	if status == "" {
		return "未创建"
	} else if strings.Contains(status, "Up") {
		return "运行中"
	} else if strings.Contains(status, "Exited") {
		return "已停止"
	} else if strings.Contains(status, "Created") {
		return "已创建"
	} else if strings.Contains(status, "Restarting") {
		return "重启中"
	} else {
		// 返回原始状态，去除末尾的换行符
		return strings.TrimSpace(status)
	}
}

// Log Ts日志
func (i *Ts) Log() (string, error) {
	// 连接到SSH服务器
	client, err := i.connect()
	if err != nil {
		return "", fmt.Errorf("连接服务器失败: %v", err)
	}
	defer client.Close()

	// 容器名称
	name := "ts" + i.Uuid

	// 创建新的会话
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("创建 SSH 会话失败: %v", err)
	}
	defer session.Close()

	// 获取容器日志
	logCmd := fmt.Sprintf("docker logs --tail 100 %s 2>&1", name)
	logOutput, err := session.CombinedOutput(logCmd)
	if err != nil {
		return "", fmt.Errorf("获取TeamSpeak容器日志失败: %v", err)
	}

	// 如果容器不存在或没有日志，尝试获取服务器日志文件
	if len(logOutput) == 0 {
		session2, err := client.NewSession()
		if err != nil {
			return "", fmt.Errorf("创建 SSH 会话失败: %v", err)
		}
		defer session2.Close()

		// 尝试读取日志文件
		fileLogCmd := fmt.Sprintf("cat /opt/teamspeak-control/order/%s/ts3server.log 2>/dev/null || echo '没有找到日志文件'", i.Uuid)
		fileLogOutput, err := session2.CombinedOutput(fileLogCmd)
		if err != nil {
			return "", fmt.Errorf("获取TeamSpeak日志文件失败: %v", err)
		}

		return string(fileLogOutput), nil
	}

	return string(logOutput), nil
}

// ParseLog 解析Ts日志并提取重要信息
func (i *Ts) ParseLog() (map[string]string, error) {
	// 获取日志内容
	logContent, err := i.Log()
	if err != nil {
		return nil, fmt.Errorf("获取日志失败: %v", err)
	}

	// 初始化结果map
	result := make(map[string]string)

	// 使用正则表达式提取信息
	// 提取token
	tokenRegex := regexp.MustCompile(`token=([A-Za-z0-9+/=]+)`)
	tokenMatches := tokenRegex.FindStringSubmatch(logContent)
	if len(tokenMatches) > 1 {
		result["token"] = tokenMatches[1]
	}

	// 提取loginname
	loginnameRegex := regexp.MustCompile(`loginname=\s*"([^"]+)"`)
	loginnameMatches := loginnameRegex.FindStringSubmatch(logContent)
	if len(loginnameMatches) > 1 {
		result["loginname"] = loginnameMatches[1]
	}

	// 提取password
	passwordRegex := regexp.MustCompile(`password=\s*"([^"]+)"`)
	passwordMatches := passwordRegex.FindStringSubmatch(logContent)
	if len(passwordMatches) > 1 {
		result["password"] = passwordMatches[1]
	}

	// 提取apikey
	apikeyRegex := regexp.MustCompile(`apikey=\s*"([^"]+)"`)
	apikeyMatches := apikeyRegex.FindStringSubmatch(logContent)
	if len(apikeyMatches) > 1 {
		result["apikey"] = apikeyMatches[1]
	}

	// 检查是否找到任何信息
	if len(result) == 0 {
		return nil, fmt.Errorf("未在日志中找到任何凭据信息")
	}

	return result, nil
}
