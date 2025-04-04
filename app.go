package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"net/http"
	"os"
)

// App struct
type App struct {
	ctx context.Context
}

// Config struct to match the config.json structure
type Config struct {
	Endpoint string `json:"endpoint"`
}

// Request struct
type Request struct{}

// 全局变量
var Endpoint string
var ContextX context.Context

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 初始化全局变量 Endpoint
	configFile, err := os.ReadFile("./config.json")
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("读取配置文件失败: %v", err))
		return
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("解析配置文件失败: %v", err))
		return
	}

	Endpoint = config.Endpoint
	ContextX = a.ctx
	runtime.LogInfo(a.ctx, "Endpoint 已初始化: "+Endpoint)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

type FileManager struct{}

func (a *FileManager) SetEndpoint(endpoint string) bool {
	// 1. 读取配置文件
	configFile, err := os.ReadFile("./config.json")
	if err != nil {
		runtime.LogError(ContextX, fmt.Sprintf("读取配置文件失败: %v", err))
		return false
	}

	// 2. 解析配置文件
	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		runtime.LogError(ContextX, fmt.Sprintf("解析配置文件失败: %v", err))
		return false
	}

	// 3. 更新内存中的Endpoint和配置结构体
	Endpoint = endpoint        // 将传入的endpoint赋值给全局变量
	config.Endpoint = endpoint // 更新配置结构体中的endpoint

	// 4. 将修改后的配置写回文件
	updatedConfig, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		runtime.LogError(ContextX, fmt.Sprintf("生成新配置失败: %v", err))
		return false
	}

	err = os.WriteFile("./config.json", updatedConfig, 0644)
	if err != nil {
		runtime.LogError(ContextX, fmt.Sprintf("写入配置文件失败: %v", err))
		return false
	}

	runtime.LogDebug(ContextX, "MDF ---> Endpoint: "+Endpoint)
	return true
}

// OpenFile opens a file dialog and returns the file path and base64 encoded content
func (a *FileManager) UploadFile() string {
	options := runtime.OpenDialogOptions{
		DefaultDirectory:           "",
		DefaultFilename:            "",
		Title:                      "上传文件",
		Filters:                    []runtime.FileFilter{{DisplayName: "所有文件，老弟，所有文件", Pattern: "*.*"}},
		ShowHiddenFiles:            true,
		CanCreateDirectories:       true,
		ResolvesAliases:            true,
		TreatPackagesAsDirectories: true,
	}

	// 打开文件对话框，获取文件路径
	filePath, err := runtime.OpenFileDialog(ContextX, options)
	if err != nil {
		return fmt.Sprintf("打开失败: %v", err)
	}

	// 读取文件内容
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Sprintf("读取文件失败: %v", err)
	}

	// 将文件内容编码为 Base64
	base64Encoded := base64.StdEncoding.EncodeToString(fileContent)
	runtime.LogDebug(ContextX, "BASE64 Start => "+base64Encoded)

	// 获取文件名和文件大小
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Sprintf("获取文件信息失败: %v", err)
	}
	fileName := fileInfo.Name()
	fileSize := fileInfo.Size()

	// 组装请求体
	requestBody := map[string]interface{}{
		"name": fileName,
		"size": fileSize,
		"data": base64Encoded,
	}

	// 将请求体转换为 JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Sprintf("JSON 编码失败: %v", err)
	}

	// 发送 POST 请求到指定的端点
	url := Endpoint + "upload"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Sprintf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("读取响应失败: %v", err)
	}

	// 返回文件路径、Base64 编码后的内容和服务器响应
	return fmt.Sprintf("选择了文件: %s\nBase64 编码内容:\n%s\n服务器响应:\n%s", filePath, base64Encoded, string(responseBody))
}

// GetAllFilesInfo fetches all files info from the endpoint specified in config.json
func (a *Request) GetAllFilesInfo() string {
	// 1. Check if Endpoint is set
	if Endpoint == "" {
		runtime.LogError(ContextX, "Endpoint is not set")
		return "{\"error\":\"Endpoint not configured\"}"
	}

	// 2. Create request with error handling
	url := Endpoint + "getAllFilesInfo"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		runtime.LogError(ContextX, fmt.Sprintf("Failed to create request: %v", err))
		return fmt.Sprintf("{\"error\":\"Request creation failed: %v\"}", err)
	}

	// 3. Send request with error handling
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		runtime.LogError(ContextX, fmt.Sprintf("Request failed: %v", err))
		return fmt.Sprintf("{\"error\":\"Request failed: %v\"}", err)
	}
	defer res.Body.Close()

	// 4. Read response with error handling
	body, err := io.ReadAll(res.Body)
	if err != nil {
		runtime.LogError(ContextX, fmt.Sprintf("Failed to read response: %v", err))
		return fmt.Sprintf("{\"error\":\"Failed to read response: %v\"}", err)
	}

	runtime.LogDebug(ContextX, fmt.Sprintf("Response: %s", string(body)))
	return string(body)
}
