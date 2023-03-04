/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "更新当前资源包内资源",
	Long:  `更新当前资源包内资源`,
	Run:   pullFunc,
}

func init() {
	rootCmd.AddCommand(pullCmd)

}

func pullFunc(cmd *cobra.Command, args []string) {
	log.Info("开始更新资源!")

	viper.SetConfigName(".config") // name of config file (without extension)
	viper.SetConfigType("yaml")    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")       // optionally look for config in the working directory
	err := viper.ReadInConfig()    // Find and read the config file
	if err != nil {                // Handle errors reading the config file
		log.Fatalf("fatal error config file: %w", err)
	}

	domain, output, hash := viper.GetString("domain"), "./", viper.GetString("hash")

	//构造各种地址
	endpoint := fmt.Sprintf("https://%s/%s", domain, hash)
	output = output + "/" + hash
	fingerprintUrl := fmt.Sprintf("%s/fingerprint.json", endpoint)

	//获取fingerprint
	log.Info("获取fingerprint--")
	var fingerprint Fingerprint
	client := resty.New()
	_, err2 := client.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&fingerprint).
		Get(fingerprintUrl)
	if err2 != nil {
		log.Error(err2)
		log.Fatal("获取失败")
	}
	log.Info("成功！")

	//校验文件
	log.Info("校验文件Hash--")
	files := []File{}
	for _, v := range fingerprint.Files {
		if !checkShaFile(v.Sha, "./"+v.File) {
			files = append(files, v)
		}
	}
	log.Info("校验完成！")
	log.Infof("共有%d个文件需要下载/更新!", len(files))
	if len(files) <= 0 {
		log.Info("更新任务结束！")
		return
	}

	//下载files
	log.Info("开始下载---")
	for i, v := range files {
		log.Infof("[%d/%d]下载: %s", i+1, len(files), v.File)
		_, err4 := client.R().
			SetOutput(v.File).
			Get(fmt.Sprintf("%s/%s", endpoint, v.File))
		if err4 != nil {
			log.Errorf("err:%v", &err4)
		}
	}
	log.Info("下载完成！")
}

func checkShaFile(hash string, filePath string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		//文件不存在，认为sha不一样
		return false
	}
	defer f.Close()
	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		//文件不存在，认为sha不一样
		return false
	}
	return hash == fmt.Sprintf("%x", h.Sum(nil))
}
