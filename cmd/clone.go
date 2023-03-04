/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type (
	File struct {
		File  string `json:"file"`
		Sha   string `json:"sha"`
		Defer bool   `json:"defer,omitempty"`
	}

	Fingerprint struct {
		Files   []File `json:"files"`
		Sha     string `json:"sha"`
		Version string `json:"version"`
	}
)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "克隆一个完整的assets到本地",
	Long:  `克隆一个完整的assets到本地`,
	Run:   cloneFunc,
}

func init() {
	rootCmd.AddCommand(cloneCmd)

	cloneCmd.Flags().StringP("domain", "u", "game-assets.clashroyaleapp.com", "资源cdn域名")
	cloneCmd.Flags().StringP("output", "o", "./", "资源下载路径")
	cloneCmd.Flags().StringP("extension", "e", "all", "仅下载指定扩展名的文件")
	cloneCmd.Flags().StringP("fingerprint", "f", "acf932573295414ef92479e9240aecb0854a70a7", "fingerprint文件中的hash值")

	viper.BindPFlag("domain", cloneCmd.Flags().Lookup("domain"))
	viper.BindPFlag("output", cloneCmd.Flags().Lookup("output"))
	viper.BindPFlag("hash", cloneCmd.Flags().Lookup("fingerprint"))
	viper.BindPFlag("extension", cloneCmd.Flags().Lookup("extension"))
}

func cloneFunc(cmd *cobra.Command, args []string) {

	// TODO:  <04-03-23, 校验配置文件> //
	domain, output, hash, _ := viper.GetString("domain"),
		viper.GetString("output"),
		viper.GetString("hash"),
		viper.GetString("extension")

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

	//设置输出
	client.SetOutputDirectory(output)

	//写配置过去下次用
	v := viper.New()
	v.Set("hash", hash)
	v.Set("domain", domain)
	v.WriteConfigAs(output + "/.config.yaml")

	//下载fingerprint
	log.Info("保存fingerprint文件--")
	_, err3 := client.R().
		SetOutput("fingerprint.json").
		Get(fingerprintUrl)
	if err3 != nil {
		fmt.Printf("err:%v", &err3)
	}

	//下载files
	log.Info("开始下载---")
	for i, v := range fingerprint.Files {
		log.Infof("[%d/%d]下载: %s", i+1, len(fingerprint.Files), v.File)
		_, err4 := client.R().
			SetOutput(v.File).
			Get(fmt.Sprintf("%s/%s", endpoint, v.File))
		if err4 != nil {
			fmt.Printf("err:%v", &err4)
		}
	}
	log.Info("下载完成！")

}
