#!/usr/bin/env python
# -*- coding: utf-8 -*-
# @Time    : 2021/08/16 18:52
# @Author  : Xuexia.li
# @Site    : 
# @File    : apollo_client_api.py
# @Description:
import requests
import json,toml
import sys

from ast import literal_eval

class RequestClient(object):
    def __init__(self, timeout=60, authorization="0bcbd744e2c08203a384a740f5aa9ab13f7cc24c"):
        self._timeout = timeout
        self._authorization = authorization

    def _request_get(self, url, token):
        if token.strip() != "" :
            self._authorization = token
        if self._authorization:
            return requests.get(
                url=url,
                timeout=self._timeout,
                headers={"Authorization": self._authorization,'Accept-Encoding':'gzip, deflate, br'}
            )
        else:
            return requests.get(url=url, params=params, timeout=self._timeout)

    def _request_put(self, url, json_data, token):
        if token.strip() != "" :
            self._authorization = token
        if self._authorization:
            return requests.put(
                url=url,
                data=json.dumps(json_data),
                timeout=self._timeout,
                headers={"Authorization": self._authorization,"Content-Type":"application/json;charset=UTF-8"}
            )
        else:
            return requests.put(url=url, data=json.dumps(json_data), timeout=self._timeout, headers={"Content-Type":"application/json;charset=UTF-8"})

    def _request_delete(self, url, params={}, token=""):
        if token.strip() != "" :
            self._authorization = token
        if self._authorization:
            return requests.delete(
                url=url,
                params=params,
                timeout=self._timeout,
                headers={"Authorization": self._authorization,"Content-Type":"application/json;charset=UTF-8"}
            )
        else:
            return requests.delete(url=url, params=params, timeout=self._timeout, headers={"Content-Type":"application/json;charset=UTF-8"})

    def _request_post(self, url, json_data, token):
        if token.strip() != "" :
            self._authorization = token
        if self._authorization:
            return requests.post(
                url=url,
                data=json.dumps(json_data),
                timeout=self._timeout,
                headers={"Authorization": self._authorization,"Content-Type":"application/json;charset=UTF-8"}
            )
        else:
            return requests.post(url=url, data=json.dumps(json_data), timeout=self._timeout, headers={"Content-Type":"application/json;charset=UTF-8"})


class PrivateApolloClient(RequestClient):
    def __init__(self, portal_address, user, authorization, app_id='dsp', env='DEV', timeout=60):
        '''
        :param portal_address: apollo????????????
        :param app_id: ??????????????????AppId
        :param authorization: ????????????
        :param env: ????????????????????????
        :param timeout:
        '''
        RequestClient.__init__(self, timeout=timeout, authorization=authorization)
        self._portal_address = portal_address
        self._appid = app_id
        self._env = env
        self._user = user
        self._commentlimit = 64

    def get_cluster(self, appid='dsp', clusterName='dsp_ali_vg', token=""):
        '''
        ??????cluster
        :param appid: Cluster?????????AppId
        :param clusterName: Cluster?????????
        :return:
        '''
        __url = '{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters/{clusterName}'.format(
            portal_address=self._portal_address, env=self._env, appId=appid, clusterName=clusterName
        )
        print("%s: %s" %(sys._getframe().f_code.co_name, __url))
        try:
            resp = self._request_get(url=__url, token=token)
            if resp.status_code is 200 :
                return resp.json()
            else :
                print("%s: response code is %d, response detail: %s" %(sys._getframe().f_code.co_name, resp.status_code,resp.json()))
                return {}
        except BaseException as e:
            print("get_cluster err", e)
            return {}

    def create_cluster(self, appid='dsp', clusterName='dsp_ali_vg', dataChangeCreatedBy="", token=""):
        '''
        ??????cluster
        :param appid: Cluster?????????AppId
        :param clusterName: Cluster?????????
        :param dataChangeCreatedBy: item?????????????????????????????????????????????sso?????????User ID
        :return:
        '''
        if dataChangeCreatedBy == "" :
            dataChangeCreatedBy = self._user
        __url = '{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters'.format(
            portal_address=self._portal_address, env=self._env, appId=appid
        )
        __data = {
                "name":clusterName,
                "appId":appid,
                "dataChangeCreatedBy":dataChangeCreatedBy
            }
        print("%s: %s %s" %(sys._getframe().f_code.co_name, __url,__data))
        try:
            resp = self._request_post(url=__url, json_data=__data, token=token)
            if resp.status_code is 200 :
                return resp.json()
            else :
                print("%s: response code is %d, response detail: %s" %(sys._getframe().f_code.co_name, resp.status_code,resp.json()))
                return {}
        except BaseException as e:
            print("creat_cluster err", e)
            return {}

    def get_namespace(self, appid='dsp', clusterName='dsp_ali_vg', namespaceName='application', token=""):
        '''
        ??????namespace
        :param appid: Namespace?????????AppId
        :param namespaceName: Namespace?????????
        :return:
        '''
        __url = '{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}'.format(
            portal_address=self._portal_address, env=self._env, appId=appid, clusterName=clusterName, namespaceName=namespaceName
        )
        print("%s: %s" %(sys._getframe().f_code.co_name, __url))
        try:
            resp = self._request_get(url=__url, token=token)
            if resp.status_code is 200 :
                return resp.json()
            else :
                print("%s: response code is %d, response detail: %s" %(sys._getframe().f_code.co_name, resp.status_code,resp.json()))
                return {}
        except BaseException as e:
            print("get_namespace err", e)
            return {}

    def create_namespace(self, appid='dsp', namespaceName='application', format='properties', isPublic=False, dataChangeCreatedBy="", comment="", token=""):
        '''
        ??????namespace
        :param appid: Namespace?????????AppId
        :param namespaceName: Namespace?????????
        :param format: Namespace???????????????????????????????????? properties???xml???json???yml???yaml
        :param isPublic: ?????????????????????
        :param dataChangeCreatedBy: item?????????????????????????????????????????????sso?????????User ID
        :param comment: ???????????????,??????????????????1024?????????
        :return:
        '''
        if len(comment) > self._commentlimit :
            comment = comment[0:63]
        if dataChangeCreatedBy == "" :
            dataChangeCreatedBy = self._user
        __url = '{portal_address}/openapi/v1/apps/{appId}/appnamespaces'.format(
            portal_address=self._portal_address, appId=appid
        )
        __data = {
                "name":namespaceName,
                "appId":appid,
                "format":format,
                "isPublic":isPublic,
                "comment":comment,
                "dataChangeCreatedBy":dataChangeCreatedBy
            }
        print("%s: %s %s" %(sys._getframe().f_code.co_name, __url,__data))
        try:
            resp = self._request_post(url=__url, json_data=__data, token=token)
            if resp.status_code is 200 :
                return resp.json()
            else :
                print("%s: response code is %d, response detail: %s" %(sys._getframe().f_code.co_name, resp.status_code,resp.json()))
                return {}
        except BaseException as e:
            print("creat_namespace err", e)
            return {}

    def delete_namespace_items_key(self, key, appid='dsp', clusterName='dsp_ali_vg', namespaceName='application', dataChangeLastModifiedBy="", token=""):
        '''
        ??????????????????
        :param namespaceName: ????????????Namespace????????????????????????properties????????????????????????????????????sample.yml
        :param key: ?????????key?????????url??????key???????????????properties?????????key?????????content
        :param dataChangeLastModifiedBy: item?????????????????????????????????????????????sso?????????User ID
        :return:
        '''
        __url = '{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/items/{key}?operator={operator}'.format(
            portal_address=self._portal_address, env=self._env, appId=appid, clusterName=clusterName, namespaceName=namespaceName, key=key, operator=dataChangeLastModifiedBy
        )
        __data = {
                "key":"content",
                "operator":dataChangeLastModifiedBy
            }
        print("%s: %s %s" %(sys._getframe().f_code.co_name, __url,__data))
        try:
            resp = self._request_delete(url=__url, params=__data, token=token)
            if resp.status_code is 200 :
                return {"status_code":200}
            else :
                print("%s: response code is %d" %(sys._getframe().f_code.co_name, resp.status_code))
                return {}
        except BaseException as e:
            print("delete_namespace_items_key err", e)
            return {}


    def get_namespace_items_key(self, key, appid='dsp', clusterName='dsp_ali_vg', namespaceName='application', token=""):
        '''
        ??????????????????
        :param namespaceName: ????????????Namespace?????????
        :param clusterName: ?????????????????????????????? ????????????????????? default ??????????????????????????????????????????????????????????????????
        :param key: ???????????????key??????
        :return:
        '''
        #??????apollo?????????key?????????/?????????
        namespace_json = {}
        if "/" in key :
            #self.delete_namespace_items_key(key, appid=appid, clusterName=clusterName, namespaceName=namespaceName, dataChangeLastModifiedBy="apollo", token=token)
            namespace_json = self.get_namespace(appid=appid, clusterName=clusterName, namespaceName=namespaceName, token=token)
        if "items" in namespace_json:
            for i, val in enumerate(namespace_json["items"]):
                if "key" in val and val["key"] == key :
                    return val

        __url = '{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/items/{key}'.format(
            portal_address=self._portal_address, env=self._env, appId=appid, clusterName=clusterName, namespaceName=namespaceName, key=key.replace("/","\/")
        )
        print("%s: %s" %(sys._getframe().f_code.co_name, __url))
        try:
            resp = self._request_get(url=__url, token=token)
            if resp.status_code is 200 :
                return resp.json()
            else :
                print("%s: response code is %d, response detail: %s" %(sys._getframe().f_code.co_name, resp.status_code,resp.json()))
                return {}
        except BaseException as e:
            print("get_namespace_items_key err", e)
            return {}

    def update_namespace_items_key(self, key, value, appid='dsp', clusterName='dsp_ali_vg', namespaceName='application', dataChangeLastModifiedBy="", comment="", token=""):
        '''
        ??????????????????
        :param namespaceName: ????????????Namespace????????????????????????properties????????????????????????????????????sample.yml
        :param key: ?????????key?????????url??????key???????????????properties?????????key?????????content
        :param value: ?????????value?????????????????????20000???????????????properties?????????value?????????????????????
        :param comment: ???????????????,??????????????????1024?????????
        :param dataChangeLastModifiedBy: item?????????????????????????????????????????????sso?????????User ID
        :return:
        '''
        if len(comment) > self._commentlimit :
            comment = comment[0:63]
        if dataChangeLastModifiedBy == "" :
            dataChangeLastModifiedBy = self._user
        __url = '{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/items/{key}'.format(
            portal_address=self._portal_address, env=self._env, appId=appid, clusterName=clusterName, namespaceName=namespaceName, key=key
        )
        __data = {
                "key":key,
                "value":value,
                "comment":comment,
                "dataChangeLastModifiedBy":dataChangeLastModifiedBy
            }
        print("%s: %s %s" %(sys._getframe().f_code.co_name, __url,__data))
        try:
            resp = self._request_put(url=__url, json_data=__data, token=token)
            if resp.status_code is 200 :
                return {"status_code":200}
            else :
                print("%s: response code is %d" %(sys._getframe().f_code.co_name, resp.status_code))
                return {}
        except BaseException as e:
            print("update_namespace_items_key err", e)
            return {}


    def create_namespace_items_key_json(self, key, value, appid='dsp', clusterName='dsp_ali_vg', namespaceName='application', dataChangeCreatedBy="", comment="", token=""):
        '''
        ??????abtest/abtest_info????????????
        :param namespaceName: ????????????Namespace????????????????????????properties????????????????????????????????????sample.yml
        :param key: ?????????key?????????url??????key???????????????properties?????????key?????????content
        :param value: ?????????value?????????????????????20000???????????????properties?????????value?????????????????????
        :param comment: ???????????????,??????????????????1024?????????
        :param dataChangeCreatedBy: item?????????????????????????????????????????????sso?????????User ID
        :return:
        '''
        if len(comment) > self._commentlimit :
            comment = comment[0:63]
        create_abtest_fail = False
        if dataChangeCreatedBy == "" :
            dataChangeCreatedBy = self._user
        __url = '{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/items'.format(
            portal_address=self._portal_address, env=self._env, appId=appid, clusterName=clusterName, namespaceName=namespaceName
        )
        
        __data = {
                "key":"consul_key",
                "value":key.replace(".","/"),
                "comment":comment,
                "dataChangeCreatedBy":dataChangeCreatedBy
            }
        print("%s: %s %s" %(sys._getframe().f_code.co_name, __url,__data))
        try:
            resp = self._request_post(url=__url, json_data=__data, token=token)
            if resp.status_code is 200 :
                print(resp.json())
            else :
                print("%s: response code is %d, response detail: %s" %(sys._getframe().f_code.co_name, resp.status_code,resp.json()))
                create_abtest_fail = True
        except BaseException as e:
            print("create_namespace_items_key err", e)
            create_abtest_fail = True
        
        for real_value in literal_eval(value.strip()) :
            if "experiment" in real_value and "name" in  real_value["experiment"] :
                if "layer" in real_value :
                    real_key = "%s_%s" % (real_value["layer"], real_value["experiment"]["name"])
                else :
                    real_key = real_value["experiment"]["name"]
            else :
                print("invalid abtest_info config: ",real_value)
                continue
            __data = {
                    "key":real_key,
                    "value":json.dumps(real_value, sort_keys=True, indent=4, separators=(',', ':'),ensure_ascii=False),
                    "comment":comment,
                    "dataChangeCreatedBy":dataChangeCreatedBy
                }
            print("%s: %s %s" %(sys._getframe().f_code.co_name, __url,__data))
            try:
                resp = self._request_post(url=__url, json_data=__data, token=token)
                if resp.status_code is 200 :
                    print(resp.json())
                else :
                    print("%s: response code is %d, response detail: %s" %(sys._getframe().f_code.co_name, resp.status_code,resp.json()))
                    create_abtest_fail = True
            except BaseException as e:
                print("create_namespace_items_key err", e)
                create_abtest_fail = True
        if create_abtest_fail :
            return {}
        else :
            return {"status_code":200}

    def create_namespace_items_key_toml(self, key, value, appid='dsp', clusterName='dsp_ali_vg', namespaceName='application', dataChangeCreatedBy="", comment="", token=""):
        '''
        ??????bidforce????????????
        :param namespaceName: ????????????Namespace????????????????????????properties????????????????????????????????????sample.yml
        :param key: ?????????key?????????url??????key???????????????properties?????????key?????????content
        :param value: ?????????value?????????????????????20000???????????????properties?????????value?????????????????????
        :param comment: ???????????????,??????????????????1024?????????
        :param dataChangeCreatedBy: item?????????????????????????????????????????????sso?????????User ID
        :return:
        '''
        if len(comment) > self._commentlimit :
            comment = comment[0:63]
        create_abtest_fail = False
        if dataChangeCreatedBy == "" :
            dataChangeCreatedBy = self._user
        __url = '{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/items'.format(
            portal_address=self._portal_address, env=self._env, appId=appid, clusterName=clusterName, namespaceName=namespaceName
        )
        
        __data = {
                "key":"consul_key",
                "value":key.replace(".","/"),
                "comment":comment,
                "dataChangeCreatedBy":dataChangeCreatedBy
            }
        print("%s: %s %s" %(sys._getframe().f_code.co_name, __url,__data))
        try:
            resp = self._request_post(url=__url, json_data=__data, token=token)
            if resp.status_code is 200 :
                print(resp.json())
            else :
                print("%s: response code is %d, response detail: %s" %(sys._getframe().f_code.co_name, resp.status_code,resp.json()))
                create_abtest_fail = True
        except BaseException as e:
            print("create_namespace_items_key err", e)
            create_abtest_fail = True
        bidforce_config = toml.loads(value.strip())
        if "BidForceDeviceType" in bidforce_config :
            for real_key, value in bidforce_config["BidForceDeviceType"].items() :
                real_value = {}
                real_value["BidForceDeviceType"]={}
                real_value["BidForceDeviceType"][real_key] = value
                __data = {
                        "key":real_key,
                        "value":toml.dumps(real_value),
                        "comment":comment,
                        "dataChangeCreatedBy":dataChangeCreatedBy
                    }
                print("%s: %s %s" %(sys._getframe().f_code.co_name, __url,__data))
                try:
                    resp = self._request_post(url=__url, json_data=__data, token=token)
                    if resp.status_code is 200 :
                        print(resp.json())
                    else :
                        print("%s: response code is %d, response detail: %s" %(sys._getframe().f_code.co_name, resp.status_code,resp.json()))
                        create_abtest_fail &= True
                except BaseException as e:
                    print("create_namespace_items_key err", e)
                    create_abtest_fail &= True
        if create_abtest_fail :
            return {}
        else :
            return {"status_code":200}


    def create_namespace_items_key(self, key, value, appid='dsp', clusterName='dsp_ali_vg', namespaceName='application', dataChangeCreatedBy="", comment="", token=""):
        '''
        ??????????????????
        :param namespaceName: ????????????Namespace????????????????????????properties????????????????????????????????????sample.yml
        :param key: ?????????key?????????url??????key???????????????properties?????????key?????????content
        :param value: ?????????value?????????????????????20000???????????????properties?????????value?????????????????????
        :param comment: ???????????????,??????????????????1024?????????
        :param dataChangeCreatedBy: item?????????????????????????????????????????????sso?????????User ID
        :return:
        '''
        if len(comment) > self._commentlimit :
            comment = comment[0:63]
        if dataChangeCreatedBy == "" :
            dataChangeCreatedBy = self._user
        __url = '{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/items'.format(
            portal_address=self._portal_address, env=self._env, appId=appid, clusterName=clusterName, namespaceName=namespaceName
        )
        __data = {
                "key":key,
                "value":value,
                "comment":comment,
                "dataChangeCreatedBy":dataChangeCreatedBy
            }
        print("%s: %s %s" %(sys._getframe().f_code.co_name, __url,__data))
        try:
            resp = self._request_post(url=__url, json_data=__data, token=token)
            if resp.status_code is 200 :
                return resp.json()
            else :
                print("%s: response code is %d, response detail: %s" %(sys._getframe().f_code.co_name, resp.status_code,resp.json()))
                return {}
        except BaseException as e:
            print("create_namespace_items_key err", e)
            return {}

    def releases(self, releaseTitle, releaseComment, appid='dsp', clusterName='dsp_ali_vg', namespaceName='application', releasedBy="", token=""):
        '''
        ??????????????????
        :param releaseTitle: ??????????????????????????????????????????64?????????
        :param releaseComment: ????????????????????????????????????256?????????
        :param releasedBy: ???????????????????????????????????????ApolloConfigDB.ServerConfig??????namespace.lock.switch?????????true??????????????????false??????????????????????????????????????????????????????????????????????????????????????????zhanglea???????????????????????????zhanglea???
        :param namespaceName: ????????????Namespace????????????????????????properties????????????????????????????????????sample.yml
        :return:
        '''
        if len(releaseTitle) > self._commentlimit :
            releaseTitle = releaseTitle[0:63]       
        if len(releaseComment) > self._commentlimit :
            releaseComment = releaseComment[0:63]
        if releasedBy == "" :
            releasedBy = self._user
        __url = '{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/releases'.format(
            portal_address=self._portal_address, env=self._env, appId=appid, clusterName=clusterName, namespaceName=namespaceName)
        __data = {
                "releaseTitle":releaseTitle,
                "releaseComment":releaseComment,
                "releasedBy":releasedBy
            }
        print("%s: %s %s" %(sys._getframe().f_code.co_name, __url,__data))
        try:
            resp = self._request_post(url=__url, json_data=__data, token=token)
            if resp.status_code is 200 :
                return resp.json()
            else :
                print("%s: response code is %d, response detail: %s" %(sys._getframe().f_code.co_name, resp.status_code,resp.json()))
                return {}
        except BaseException as e:
            print("releases err", e)
            return {}

if __name__ == '__main__':
    pass