package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/1mtrue/dukou_check/types"
	"github.com/1mtrue/dukou_check/utils"
)

func main() {
	
	if err := setUpConfig(); err != nil {
		logrus.Fatal("init config error")
	}
	username := viper.GetString("username")

	logrus.Infof("user: %s start check-in", username)
	token, err := login()
	if err != nil {
		log.Fatal("login error", err)
	}
	logrus.Info("login success")
	if err := CheckIn(token); err != nil {
		log.Fatal("check error", err)

	}
}

func setUpConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 添加配置文件的搜索路径
	viper.AddConfigPath(".")

	// 从环境变量中读取配置
	viper.AutomaticEnv()

	// 解析配置文件
	return viper.ReadInConfig()
}

func login() (token string, err error) {
	username := viper.GetString("username")
	password := viper.GetString("password")
	if username == "" || password == "" {
		return "", errors.New("pleaese set user info in config.yaml or set in you env")
	}
	baseUrl := viper.GetString("baseUrl")
	if baseUrl == "" {
		viper.Set("baseUrl", "https://dukou.io")
	}
	loginURL := viper.GetString("baseUrl") + "/api/token"

	loginData := types.LoginRequest{
		Email:  username,
		Passwd: password,
	}
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		logrus.Fatal("pleaese set user info in config.yaml or set in you env")
	}
	client := &http.Client{}

	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err

	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {

		return "", err
	}
	defer resp.Body.Close()

	// 检查登录请求的响应状态码
	if resp.StatusCode != http.StatusOK {

		return "", errors.New("login status not ok ")
	}

	var respData types.LoginResp

	body, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &respData)
	if err != nil {
		return "", err
	}

	if respData.Ret != 1 {
		return "", errors.New("login status not ok ")
	}
	return respData.Token, nil
}

func CheckIn(token string) error {
	checkInUrl := viper.GetString("baseUrl") + "/api/user/checkin"
	request, err := utils.NewLoginedRequest(token, checkInUrl)
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {

		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// 解析 JSON 数据
	var data map[string]interface{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	if data["ret"].(float64) != 1 || data["result"].(string) == "" {
		return errors.New("checkin  failed " + data["result"].(string))
	}
	count, err := utils.ExtractNum(data["result"].(string))
	logrus.Info(data["result"].(string))
	if err != nil {
		return err
	}

	return ConverTraffic(count, token)
}

func ConverTraffic(count int, token string) error {
	//api/user/koukanntraffic?traffic=239
	converUrl := fmt.Sprintf(viper.GetString("baseUrl")+"/api/user/koukanntraffic?traffic=%d", count)

	request, err := utils.NewLoginedRequest(token, converUrl)
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {

		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// 解析 JSON 数据
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	logrus.Info(data)
	return nil
}
