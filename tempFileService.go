package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

// TempFileService 安全存储临时文件工具类
type TempFileService struct {
	tempDir   string
	filePaths map[string]bool
	mutex     sync.Mutex
}

// NewTempFileManager 创建临时文件管理器
func NewTempFileService(dirname string) (*TempFileService, error) {
	tempDir, err := os.MkdirTemp("", dirname)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	manager := &TempFileService{
		tempDir:   tempDir,
		filePaths: make(map[string]bool),
	}
	manager.registerCleanupHook()

	return manager, nil
}

// WriteTempFile 将字节数据写入临时文件
func (t *TempFileService) WriteTempFile(data []byte, name string) (string, error) {
	filePath := filepath.Join(t.tempDir, name)

	t.mutex.Lock()
	defer t.mutex.Unlock()

	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		t.Cleanup() // 写入失败时立即清理
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}

	t.filePaths[filePath] = true
	return filePath, nil
}

// Cleanup 安全清理临时资源
func (t *TempFileService) Cleanup() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if _, err := os.Stat(t.tempDir); !os.IsNotExist(err) {
		os.RemoveAll(t.tempDir)
		t.filePaths = make(map[string]bool)
	}
}

// registerCleanupHook 注册进程退出/异常时的自动清理
func (t *TempFileService) registerCleanupHook() {
	// 正常退出时清理
	defer t.Cleanup()

	// 处理中断信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Received interrupt signal, cleaning temp files")
		t.Cleanup()
		os.Exit(0)
	}()
}
