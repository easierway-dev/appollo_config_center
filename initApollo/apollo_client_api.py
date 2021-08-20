#!/usr/bin/env python
# -*- coding: utf-8 -*-
# @Time    : 2021/08/16 18:52
# @Author  : Xuexia.li
# @Site    : 
# @File    : apollo_client_api.py
# @Description:
import requests
import json
import sys

class RequestClient(object):
    def __init__(self, timeout=60, authorization="0bcbd744e2c08203a384a740f5aa9ab13f7cc24c"):
        self._timeout = timeout
        self._authorization = authorization

    def _request_get(self, url):
        if self._authorization:
            return requests.get(
                url=url,
                timeout=self._timeout,
                headers={"Authorization": self._authorization}
            )
        else:
            return requests.get(url=url, params=params, timeout=self._timeout)

    def _request_put(self, url, json_data):
        if self._authorization:
            return requests.put(
                url=url,
                data=json.dumps(json_data),
                timeout=self._timeout,
                headers={"Authorization": self._authorization,"Content-Type":"application/json;charset=UTF-8"}
            )
        else:
            return requests.put(url=url, data=json.dumps(json_data), timeout=self._timeout, headers={"Content-Type":"application/json;charset=UTF-8"})

    def _request_delete(self, url, params={}):
        if self._authorization:
            return requests.delete(
                url=url,
                params=params,
                timeout=self._timeout,
                headers={"Authorization": self._authorization,"Content-Type":"application/json;charset=UTF-8"}
            )
        else:
            return requests.delete(url=url, params=params, timeout=self._timeout, headers={"Content-Type":"application/json;charset=UTF-8"})

    def _request_post(self, url, json_data):
        print(self._authorization)
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
        :param portal_address: apollo接口地址
        :param app_id: 所管理的配置AppId
        :param authorization: 鉴权参数
        :param env: 所管理的配置环境
        :param timeout:
        '''
        RequestClient.__init__(self, timeout=timeout, authorization=authorization)
        self._portal_address = portal_address
        self._appid = app_id
        self._env = env
        self._user = user

    def get_cluster(self, appid='dsp', clusterName='dsp_ali_vg'):
        '''
        读取cluster
        :param appid: Cluster所属的AppId
        :param clusterName: Cluster的名字
        :return:
        '''
        __url = '{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters/{clusterName}'.format(
            portal_address=self._portal_address, env=self._env, appId=appid, clusterName=clusterName
        )
        print("%s: %s" %(sys._getframe().f_code.co_name, __url))
        try:
            resp = self._request_get(url=__url)
            if resp.status_code is 200 :
                return resp.json()
            else :
                print("%s: response code is %d" %(sys._getframe().f_code.co_name, resp.status_code))
                return {}
        except BaseException as e:
            print("get_cluster err", e)
            return {}

    def create_cluster(self, appid='dsp', clusterName='dsp_ali_vg', dataChangeCreatedBy=""):
        '''
        新增cluster
        :param appid: Cluster所属的AppId
        :param clusterName: Cluster的名字
        :param dataChangeCreatedBy: item的创建人，格式为域账号，也就是sso系统的User ID
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
            resp = self._request_post(url=__url, json_data=__data)
            if resp.status_code is 200 :
                return resp.json()
            else :
                print("%s: response code is %d" %(sys._getframe().f_code.co_name, resp.status_code))
                return {}
        except BaseException as e:
            print("creat_cluster err", e)
            return {}

    def get_namespace(self, appid='dsp', clusterName='dsp_ali_vg', namespaceName='application'):
        '''
        新增namespace
        :param appid: Namespace所属的AppId
        :param namespaceName: Namespace的名字
        :return:
        '''
        __url = '{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}'.format(
            portal_address=self._portal_address, env=self._env, appId=appid, clusterName=clusterName, namespaceName=namespaceName
        )
        print("%s: %s" %(sys._getframe().f_code.co_name, __url))
        try:
            resp = self._request_get(url=__url)
            if resp.status_code is 200 :
                return resp.json()
            else :
                print("%s: response code is %d" %(sys._getframe().f_code.co_name, resp.status_code))
                return {}
        except BaseException as e:
            print("get_namespace err", e)
            return {}

    def create_namespace(self, appid='dsp', namespaceName='application', format='properties', isPublic=False, dataChangeCreatedBy="", comment=None):
        '''
        新增namespace
        :param appid: Namespace所属的AppId
        :param namespaceName: Namespace的名字
        :param format: Namespace的格式，只能是以下类型： properties、xml、json、yml、yaml
        :param isPublic: 是否是公共文件
        :param dataChangeCreatedBy: item的创建人，格式为域账号，也就是sso系统的User ID
        :param comment: 配置的备注,长度不能超过1024个字符
        :return:
        '''
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
            resp = self._request_post(url=__url, json_data=__data)
            if resp.status_code is 200 :
                return resp.json()
            else :
                print("%s: response code is %d" %(sys._getframe().f_code.co_name, resp.status_code))
                return {}
        except BaseException as e:
            print("creat_namespace err", e)
            return {}

    def get_namespace_items_key(self, key, appid='dsp', clusterName='dsp_ali_vg', namespaceName='application'):
        '''
        读取配置接口
        :param namespaceName: 所管理的Namespace的名称，如果是非properties格式，需要加上后缀名，如sample.yml
        :param clusterName: 所管理的配置集群名， 一般情况下传入 default 即可。如果是特殊集群，传入相应集群的名称即可
        :param key: 配置对应的key名称
        :return:
        '''
        __url = '{portal_address}/openapi/v1/envs/{env}/apps/{appId}/clusters/{clusterName}/namespaces/{namespaceName}/items/{key}'.format(
            portal_address=self._portal_address, env=self._env, appId=appid, clusterName=clusterName, namespaceName=namespaceName, key=key
        )
        print("%s: %s" %(sys._getframe().f_code.co_name, __url))
        try:
            resp = self._request_get(url=__url)
            if resp.status_code is 200 :
                return resp.json()
            else :
                print("%s: response code is %d" %(sys._getframe().f_code.co_name, resp.status_code))
                return {}
        except BaseException as e:
            print("get_namespace_items_key err", e)
            return {}

    def update_namespace_items_key(self, key, value, appid='dsp', clusterName='dsp_ali_vg', namespaceName='application', dataChangeLastModifiedBy="", comment=None):
        '''
        修改配置接口
        :param namespaceName: 所管理的Namespace的名称，如果是非properties格式，需要加上后缀名，如sample.yml
        :param key: 配置的key，需和url中的key值一致。非properties格式，key固定为content
        :param value: 配置的value，长度不能超过20000个字符，非properties格式，value为文件全部内容
        :param comment: 配置的备注,长度不能超过1024个字符
        :param dataChangeLastModifiedBy: item的修改人，格式为域账号，也就是sso系统的User ID
        :return:
        '''
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
            resp = self._request_put(url=__url, json_data=__data)
            if resp.status_code is 200 :
                return {"status_code":200}
            else :
                print("%s: response code is %d" %(sys._getframe().f_code.co_name, resp.status_code))
                return {}
        except BaseException as e:
            print("update_namespace_items_key err", e)
            return {}

    def create_namespace_items_key(self, key, value, appid='dsp', clusterName='dsp_ali_vg', namespaceName='application', dataChangeCreatedBy="", comment=None):
        '''
        新增配置接口
        :param namespaceName: 所管理的Namespace的名称，如果是非properties格式，需要加上后缀名，如sample.yml
        :param key: 配置的key，需和url中的key值一致。非properties格式，key固定为content
        :param value: 配置的value，长度不能超过20000个字符，非properties格式，value为文件全部内容
        :param comment: 配置的备注,长度不能超过1024个字符
        :param dataChangeCreatedBy: item的创建人，格式为域账号，也就是sso系统的User ID
        :return:
        '''
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
            resp = self._request_post(url=__url, json_data=__data)
            if resp.status_code is 200 :
                return resp.json()
            else :
                print("%s: response code is %d" %(sys._getframe().f_code.co_name, resp.status_code))
                return {}
        except BaseException as e:
            print("create_namespace_items_key err", e)
            return {}

    def releases(self, releaseTitle, releaseComment, appid='dsp', clusterName='dsp_ali_vg', namespaceName='application', releasedBy=""):
        '''
        发布配置接口
        :param releaseTitle: 此次发布的标题，长度不能超过64个字符
        :param releaseComment: 发布的备注，长度不能超过256个字符
        :param releasedBy: 发布人，域账号，注意：如果ApolloConfigDB.ServerConfig中的namespace.lock.switch设置为true的话（默认是false），那么该环境不允许发布人和编辑人为同一人。所以如果编辑人是zhanglea，发布人就不能再是zhanglea。
        :param namespaceName: 所管理的Namespace的名称，如果是非properties格式，需要加上后缀名，如sample.yml
        :return:
        '''
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
            resp = self._request_post(url=__url, json_data=__data)
            if resp.status_code is 200 :
                return resp.json()
            else :
                print("%s: response code is %d" %(sys._getframe().f_code.co_name, resp.status_code))
                return {}
        except BaseException as e:
            print("releases err", e)
            return {}

if __name__ == '__main__':
    reqClient = RequestClient()
    portaddr = "http://localhost:80"
    appid = "dsp"
    token = "0bcbd744e2c08203a384a740f5aa9ab13f7cc24c"
    inputuser = "apollo"
    apolloClient = PrivateApolloClient(portaddr,inputuser,token,appid)
    key= "test"
    value ="5"
    enRelease = True
    inputdata = {
        "dsp":{
            "cluster": [ "dsp_ali_cn", "dsp_hw_hk", "dsp_ali_vg",],
            "consulkey": [ "test", "xiongjia/aa",]
        }
    }
    #for appid, clusters in inputdata.items() :
    #    if "cluster" in clusters.keys() :
    #        for cluster in clusters["cluster"] :


    getResp = apolloClient.get_namespace_items_key(key)
    #if getResp.status_code is 200 and getResp.
    if bool(getResp) :
        if getResp["value"] != value :
            putResp = apolloClient.put_namespace_items_key(key, value, comment="update k=%s ov=%s nv=%s " %(key,getResp.json()["value"],value))
            if bool(putResp) and enRelease :
                print("update:",apolloClient.releases("release", "release %s:%s"%(key,value)))
        else :
            print("noNeed to update !!!",getResp)
    else :
        postResp = apolloClient.create_namespace_items_key(key, value, inputuser, comment="insert k=%s v=%s" %(key,value))
        if bool(postResp) and enRelease :
            print("insert:",apolloClient.releases("release", "release %s:%s"%(key,value)))

