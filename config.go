package main

import (
	"crypto/sha256"
	"errors"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

var globalConfig *config

// const value
const (
	DefaultConfigFile = "config.yaml"

	DefaultZone         = "zh-hangzhou"
	DefaultDomain       = "dysmsapi.aliyuncs.com"
	DefaultTemplateCode = "SMS_000000000"
	DefaultSignName     = "test"
	DefaultBaseAddr     = "127.0.0.1:8899"
)

type config struct {
	AliyunAccessKey    string
	AliyunAccessSecret string

	SmsZone   string
	SmsDomain string

	SmsSignName     string
	SmsTemplateCode string

	MongoDB struct {
		URI      string
		User     string
		Password string
	}
	BaseAddr string
}

func convertTo32(origin string) string {
	if origin == "" {
		origin = "default_string"
	}
	data := sha256.Sum256([]byte(origin))
	return string(data[:])
}

func (c *config) validate() error {
	if c.MongoDB.URI == "" {
		return errors.New("MongoDB.URI不能为空")
	}
	if len(c.AliyunAccessKey) != 16 {
		log.Println(c.AliyunAccessKey, len(c.AliyunAccessKey))
		return errors.New("AliyunAccessKey长度必须是16位")
	}
	if len(c.AliyunAccessSecret) != 30 {
		return errors.New("AliyunAccessSecret长度必须是30位")
	}
	if c.SmsZone == "" {
		c.SmsZone = DefaultZone
	}
	if c.SmsDomain == "" {
		c.SmsDomain = DefaultDomain
	}
	if c.SmsSignName == "" {
		c.SmsSignName = DefaultSignName
	}
	if c.SmsTemplateCode == "" {
		c.SmsTemplateCode = DefaultTemplateCode
	}
	if c.BaseAddr == "" {
		c.BaseAddr = DefaultBaseAddr
	}
	return nil
}

func readConfig(fileName string) *config {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err)
		return nil
	}
	c := &config{}
	err = yaml.Unmarshal(content, c)
	if err != nil {
		log.Println("Parse config file", fileName, "fail:", err)
		return nil
	}
	if err := c.validate(); err != nil {
		log.Println("Parse config file", fileName, "fail:", err)
		return nil
	}
	return c
}
