package cnotify

import (
	"io/ioutil"

	"github.com/CodyGuo/dingtalk"
	"github.com/CodyGuo/dingtalk/pkg/robot"
	gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon
)

func sendText(token sting, oldValue, newValue interface{}) {
	glog.SetFlags(glog.LglogFlags)
	webHook := "https://oapi.dingtalk.com/robot/send?access_token=xxx"
	secret := token
	dt := dingtalk.New(webHook, dingtalk.WithSecret(secret))

	// text类型
	textContent := fmt.Sprintf("[global_config changed]\n old:%v \nnew:%v", oldValue, newValue)
	atMobiles := robot.SendWithAtMobiles([]string{"176xxxxxx07", "178xxxxxx28"})
	if err := dt.RobotSendText(textContent, atMobiles); err != nil {
		ccommon.CLogger.Errorf("send ding failed err: ",err)
	}
	printResult(dt)
}

func printResult(dt *dingtalk.DingTalk) {
	response, err := dt.GetResponse()
	if err != nil {
		ccommon.CLogger.Errorf("Parse dingResp failed err: ",err)
	}
	reqBody, err := response.Request.GetBody()
	if err != nil {
		ccommon.CLogger.Errorf("Parse dingResp failed err: ",err)
	}
	reqData, err := ioutil.ReadAll(reqBody)
	if err != nil {
		ccommon.CLogger.Errorf("Parse dingResp failed err: ",err)
	}
	ccommon.CLogger.Infof("发送消息成功, message: ", reqData)
}

