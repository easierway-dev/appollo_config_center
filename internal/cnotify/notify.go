package cnotify

import (
	"fmt"
	"io/ioutil"

	"github.com/CodyGuo/dingtalk"
	"github.com/CodyGuo/dingtalk/pkg/robot"
	"github.com/CodyGuo/glog"
)

const (
	DingLimit    = 10000
)

func SendText(tokens []string, textContent string, dingusers []string) {
	textByteList := []byte(textContent)
	dingContent := textContent
	if len(textByteList) > DingLimit {
		dingContent = string(textByteList[:DingLimit])
	}
	for _, token := range tokens {
		SendTextUnit(token, dingContent, dingusers)
	}
}

func SendTextUnit(token, dingContent string, dingusers []string) {
	glog.SetFlags(glog.LglogFlags)
	webHook := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s",token)
	dt := dingtalk.New(webHook, dingtalk.WithSecret(token))

	// text类型
	atMobiles := robot.SendWithAtMobiles(dingusers)

	if err := dt.RobotSendText(dingContent, atMobiles); err != nil {
		glog.Fatal("send ding failed err: ",err)
	}
	printResult(dt)
}

func printResult(dt *dingtalk.DingTalk) {
	response, err := dt.GetResponse()
	if err != nil {
		glog.Fatal("Parse dingResp failed err: ",err)
	}
	reqBody, err := response.Request.GetBody()
	if err != nil {
		glog.Fatal("Parse dingResp failed err: ",err)
	}
	reqData, err := ioutil.ReadAll(reqBody)
	if err != nil {
		glog.Fatal("Parse dingResp failed err: ",err)
	}
	glog.Info("发送消息成功, message: ", string(reqData))
}


