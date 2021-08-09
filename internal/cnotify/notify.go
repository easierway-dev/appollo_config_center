package cnotify

import (
	"fmt"
	"io/ioutil"

	"github.com/CodyGuo/dingtalk"
	"github.com/CodyGuo/dingtalk/pkg/robot"
	"github.com/CodyGuo/glog"
)


func SendText(token, textContent string) {
	glog.SetFlags(glog.LglogFlags)
	webHook := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s",token)
	dt := dingtalk.New(webHook, dingtalk.WithSecret(token))

	// text类型
	atMobiles := robot.SendWithAtMobiles([]string{"15311489030"})
	if err := dt.RobotSendText(textContent, atMobiles); err != nil {
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
	glog.Info("发送消息成功, message: ", reqData)
}


