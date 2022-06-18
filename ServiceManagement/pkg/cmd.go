package service

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func Init(cloud_path string) error {
	// cloud_path: aws/alicloud/hwcloud
	// init前切换至对应云资源目录
	if err := os.Chdir("./tf/" + cloud_path); err != nil {
		fmt.Println("Chdir failed ", err)
		// fmt.Println("command: ", "command")
		return errors.New(err.Error())
	}
	// 开始执行terraform命令行
	cmd := exec.Command("terraform", "init")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\n", outStr)
	if err != nil {
		// 执行错误处理
		fmt.Println("cmd run failed:\n", errStr)
		return errors.New(err.Error())
	}
	return nil
}
func Plan() error {
	cmd := exec.Command("terraform", "plan")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\n", outStr)
	if err != nil {
		// 执行错误处理
		fmt.Println("cmd run failed:\n", errStr)
		return errors.New(err.Error())
	}
	return nil
}
func Apply() error {
	cmd := exec.Command("terraform", "apply", "-input=false", "-auto-approve")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\n", outStr)
	if err != nil {
		// 执行错误处理
		fmt.Println("cmd run failed:\n", errStr)
		return errors.New(err.Error())
	}
	return nil
}
