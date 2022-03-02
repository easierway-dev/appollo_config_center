package cworker

import (
	"context"
	"fmt"
	"github.com/shima-park/agollo"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/capi"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
	"strings"
)

//work run
func (cw *CWorker) Run1(ctx context.Context) {
	errorCh := cw.AgolloClient.Start()
	watchCh := cw.AgolloClient.Watch()
	var deletedKeys []string
	updateContent := ""
	var updatedKeys []string
	willUpdateConsul := true
	var modifierList []string
	go func(cw *CWorker) {
		for {
			select {
			case <-ctx.Done():
				ccommon.CLogger.Info(ccommon.DefaultDingType, cw.WkInfo.Cluster, "watch quit...")
				return
			case err := <-errorCh:
				if ccommon.AppConfiger.AppConfigMap != nil {
					if _, ok := ccommon.AppConfiger.AppConfigMap[ccommon.DefaultPollDingType]; ok {
						ccommon.ChklogRate = ccommon.AppConfiger.AppConfigMap[ccommon.DefaultPollDingType].ChklogRate
					}
				}
				if ccommon.ChklogRamdom < ccommon.ChklogRate {
					ccommon.CLogger.Info(ccommon.DefaultPollDingType, "Error:", err)
				}
			case update := <-watchCh:
				// 全局配置
				ccommon.Configer = ccommon.InitAppCfgMap(ccommon.AppConfiger, cw.WkInfo.AppID, update.Namespace)
				fmt.Println("ccommon.Configer=", ccommon.Configer)

				// 是否更新consul
				if ccommon.Configer.EnUpdateConsul != 1 {
					ccommon.CLogger.Warn(cw.WkInfo.AppID, "is not permit to update consul")
					ccommon.CLogger.Info(ccommon.DefaultDingType, "Apollo cluster(", cw.WkInfo.Cluster, ") namespace(", update.Namespace, ") \nold_value:(", update.OldValue, ") \nnew_value:(", update.NewValue, ") \n error:(", update.Error, ")\n")
				} else {
					url := fmt.Sprintf("http://%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s", ccommon.AgolloConfiger.PortalURL, "DEV", cw.WkInfo.AppID, cw.WkInfo.Cluster, update.Namespace)
					nsInfo, _ := capi.GetNamespaceInfo(url, ccommon.Configer.AccessToken)
					fmt.Println("cw.WkInfo.AppID=", cw.WkInfo.AppID)
					url1 := fmt.Sprintf("http://%s/openapi/v1/apps/%s/envclusters", ccommon.AgolloConfiger.PortalURL, cw.WkInfo.AppID)
					ecinfo, _ := capi.GetEnvClustersInfo(url1, ccommon.Configer.AccessToken)
					fmt.Println("url1=", url1)
					fmt.Println("ecinfo=", ecinfo)
					url2 := fmt.Sprintf("http://%s/openapi/v1/apps", ccommon.AgolloConfiger.PortalURL)
					appInfo, _ := capi.GetAppInfo(url2, ccommon.Configer.AccessToken)
					fmt.Println("url2=", url2)
					fmt.Println("appInfo=", appInfo)
					// 除dsp之外的业务线
					isSuccess := isContainsExceptDsp(cw, update, "write", updateContent, updatedKeys, deletedKeys, modifierList, willUpdateConsul, nsInfo)
					if !isSuccess {
						// 获取更新的key更新consul
						updatedKeys, modifierList = getUpdatedKey(cw, update, nsInfo, "write")
						// 删除操作并更新consul
						if ccommon.Configer.EnDelConsul == 1 {
							deletedKeys, willUpdateConsul = getDeleteKey(cw, update, "del")
						}
					}
					//只有abtest显示更新内容的详情，其他都只提示变更的key
					if find := strings.Contains(cw.WkInfo.AppID, ccommon.ABTestAppid); !find && len(updatedKeys) > 0 {
						updateContent = strings.Join(updatedKeys, "\n")
					}
					//记录删除的key
					if len(deletedKeys) > 0 {
						updateContent = fmt.Sprintf("%s\n\ndelelte_key=%s", updateContent, strings.Join(deletedKeys, "#"))
					}
					ccommon.CLogger.Info(ccommon.DefaultDingType, "Apollo cluster(", cw.WkInfo.Cluster, ") namespace(", update.Namespace, ") \nold_value:(", update.OldValue, ") \nnew_value:(", update.NewValue, ") \n error:(", update.Error, ")\n")
					//isLogUpdateContentToDingDing(cw,willUpdateConsul,update,updateContent,modifierList)
				}
			}
		}
	}(cw)
}
func isLogUpdateContentToDingDing(cw *CWorker, willUpdateConsul bool, update *agollo.ApolloResponse, updateContent string, modifierList []string) {
	if willUpdateConsul {
		if updateContent == "" {
			updateContent = fmt.Sprintf("nothing to update !!!\nisSupportDelete=%s", string(ccommon.Configer.EnDelConsul), " (1: support)")
		}
		if len(modifierList) > 0 {
			ccommon.CLogger.Warn(RemoveDuplicatesSlice(modifierList), cw.WkInfo.AppID, "#", cw.WkInfo.Cluster, "#", update.Namespace, ": \nupdatecontent:\n", updateContent)
		} else {
			ccommon.CLogger.Warn(cw.WkInfo.AppID, "#", cw.WkInfo.Cluster, "#", update.Namespace, ": \nupdatecontent:\n", updateContent)
		}
	} else {
		ccommon.CLogger.Warn(cw.WkInfo.AppID, "#", cw.WkInfo.Cluster, "#", update.Namespace, ": !!! invalid config will not update consul !!!")
	}
}
func isContainsExceptDsp(cw *CWorker, update *agollo.ApolloResponse, consulMode string, updateContent string, updatedKeys, deletedKeys, modifierList []string, willUpdateConsul bool, nsInfo *capi.NamespaceInfo) bool {
	if strings.Contains(cw.WkInfo.AppID, ccommon.ABTestAppid) || strings.Contains(cw.WkInfo.AppID, ccommon.BidForceAppid) {
		updateConsulValue := ""
		path := ""
		updateContent, updateConsulValue, path, updatedKeys, modifierList, willUpdateConsul = MergeUpdate(cw.WkInfo.AppID, cw.WkInfo.Cluster, update.NewValue, update.OldValue, nsInfo)
		if path != "" && willUpdateConsul {
			UpdateConsul(cw.WkInfo.AppID, update.Namespace, cw.WkInfo.Cluster, path, updateConsulValue, consulMode)
		}
		//delete keys
		for k, _ := range update.OldValue {
			if _, ok := update.NewValue[k]; !ok {
				deletedKeys = append(deletedKeys, k)
			}
		}
		return true
	}
	return false
}
func getUpdatedKey(cw *CWorker, update *agollo.ApolloResponse, nsInfo *capi.NamespaceInfo, consulMode string) (updatedKeys []string, modifierList []string) {
	//新增、更新
	for path, value := range update.NewValue {
		if oValue, ok := update.OldValue[path]; ok {
			//未发生变化的key，跳过不更新
			if oValue.(string) == value.(string) {
				continue
			}
		}
		modifier := GetModifyInfo(nsInfo, path)
		//updatecontent = fmt.Sprintf("%s\nkey=%s\nold=%s\dunnew=%s\nmodifier=%s\n", updatecontent, path, ovalue, value, modifier)
		updatedKeys = append(updatedKeys, fmt.Sprintf("update_key=%s__changedby=%s", path, modifier))
		if modifier != "" {
			modifierList = append(modifierList, modifier)
		}
		UpdateConsul(cw.WkInfo.AppID, update.Namespace, cw.WkInfo.Cluster, path, value.(string), consulMode)
	}
	return
}
func getDeleteKey(cw *CWorker, update *agollo.ApolloResponse, consulMode string) (deletedKeys []string,willUpdateConsul bool) {
	//删除
	if len(update.NewValue) == 0 {
		return
	}
	for path, value := range update.OldValue {
		if _, ok := update.NewValue[path]; !ok {
			deletedKeys = append(deletedKeys, path)
			UpdateConsul(cw.WkInfo.AppID, update.Namespace, cw.WkInfo.Cluster, path, value.(string), consulMode)
		}
	}
	willUpdateConsul = true
	return
}
