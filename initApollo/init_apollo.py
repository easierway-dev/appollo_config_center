#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# @Time    : 2021/08/16 18:52
# @Author  : Xuexia.li
# @Site    : 
# @File    : apollo_client_api.py
# @Description:
import requests
import json,toml

from apollo_client_api import *
from opconsul import *

class InitApollo(object):
    def __init__(self, user="", appid="", base_config = "./baseconfig.toml", apollo_to_consul_config="./apollo_to_consul.toml"):
        self.base_config_data = {}
        self.apollo_to_consul_config_data = {}
        with open(base_config, "r") as fs: 
            self.base_config_data = toml.load(fs)   
        with open(apollo_to_consul_config, "r") as fs: 
            self.apollo_to_consul_config_data = toml.load(fs)  
        if "portaddr" in self.base_config_data:
            _portaddr = self.base_config_data["portaddr"]
        else :
            _portaddr = "http://localhost:80"
        if "env" in self.base_config_data:
             _env = self.base_config_data["env"]
        else :
            _env = "DEV"
        if appid != "":
            _appid = appid
        else :
            _appid = "dsp"
        if user != "":
            _user = user    
        else :
            _user = "apollo"    
        if "usermap" in self.base_config_data and _user in self.base_config_data["usermap"]:
             _token = self.base_config_data["usermap"][_user]["token"]
        else :
            _token = "0bcbd744e2c08203a384a740f5aa9ab13f7cc24c" 
        if "timeout" in self.base_config_data:
            _timeout = self.base_config_data["timeout"]
        else :
            _timeout = 60       
        self.PrivateApolloClient = PrivateApolloClient( _portaddr, _user, _token, _appid, _env, _timeout)
        
    def setup_apollo(self):
        for appid, clusters in self.apollo_to_consul_config_data.items() :
            if "cluster" in clusters : 
                for cluster in clusters["cluster"] :
                    if cluster in self.base_config_data :
                        consul_addr = "%s/v1/" %(self.base_config_data[cluster])
                        operateConsul = OpConsul(consul_addr)
                    else :
                        print("init consul:%s connection failed" %(consul_addr))
                        continue
                    #cluster已经创建
                    if bool(self.PrivateApolloClient.get_cluster(appid, cluster)) :
                        #不使用默认namespace=application
                        if "namespace" in clusters : 
                            for namespace, consulkeylist in clusters["namespace"].items():
                                #namespace已经创建
                                get_namespace_resp = self.PrivateApolloClient.get_namespace(appid, cluster,namespace)
                                if bool(get_namespace_resp) :
                                    for key in consulkeylist :
                                        value = operateConsul._getconsul(key)
                                        if value == "" :
                                            print("consul key=%s value=%s, will skip" %(key, value))
                                            continue
                                        get_namespace_key_fail  = True
                                        getResp = self.PrivateApolloClient.get_namespace_items_key(key,appid,cluster,namespace)
                                        old_value = ""
                                        if not bool(getResp) :
                                            for item_map in get_namespace_resp["items"] :
                                                if "key" in item_map :
                                                    get_namespace_key_fail  = False
                                                    if "value" in item_map :
                                                        old_value = item_map["value"]
                                                    break
                                        else :
                                            get_namespace_key_fail  = False
                                            if "value" in getResp :
                                                old_value = getResp["value"] 
                                        if not get_namespace_key_fail :
                                            if old_value != value :
                                                self.PrivateApolloClient.update_namespace_items_key(key, value,appid,cluster,namespace, comment="update k=%s ov=%s nv=%s " %(key,old_value,value))
                                            else :
                                                print("noNeed to update !!!",getResp)
                                        else :
                                            self.PrivateApolloClient.create_namespace_items_key(key, value,appid,cluster,namespace, comment="insert k=%s v=%s" %(key,value))
                                else :
                                    #namespace创建
                                    if bool(self.PrivateApolloClient.create_namespace(appid, namespace)) :
                                        for key in consulkeylist :
                                            value = operateConsul._getconsul(key)
                                            if value == "" :
                                                print("consul key=%s value=%s, will skip" %(key, value))
                                                continue
                                            get_namespace_key_fail  = True
                                            getResp = self.PrivateApolloClient.get_namespace_items_key(key,appid,cluster,namespace)
                                            old_value = ""
                                            if not bool(getResp) :
                                                for item_map in get_namespace_resp["items"] :
                                                    if "key" in item_map :
                                                        get_namespace_key_fail  = False
                                                        if "value" in item_map :
                                                            old_value = item_map["value"]
                                                        break
                                            else :
                                                get_namespace_key_fail  = False
                                                if "value" in getResp :
                                                    old_value = getResp["value"] 
                                            if not get_namespace_key_fail :
                                                if old_value != value :
                                                    self.PrivateApolloClient.update_namespace_items_key(key, value,appid,cluster,namespace, comment="update k=%s ov=%s nv=%s " %(key,old_value,value))
                                                else :
                                                    print("noNeed to update !!!",getResp)
                                            else :
                                                self.PrivateApolloClient.create_namespace_items_key(key, value,appid,cluster,namespace, comment="insert k=%s v=%s" %(key,value))
                                    else :
                                        print("create_namespace appid_%s cluster_%s namespace_%s failed" %(appid, cluster,namespace))
                                self.PrivateApolloClient.releases("release", "release appid:%s cluster:%s"%(appid,cluster),appid, cluster,namespace)
                        else :
                            print("appid:%s namespace config not find failed" %(appid))
                    else :
                        #cluster创建
                        if bool(self.PrivateApolloClient.creat_cluster(appid, cluster)) :
                            #不使用默认namespace=application
                            if "namespace" in clusters : 
                                for namespace, consulkeylist in clusters["namespace"].items():
                                    #namespace已经创建
                                    if bool(self.PrivateApolloClient.get_namespace(appid, cluster,namespace)) :
                                        for key in consulkeylist :
                                            value = operateConsul._getconsul(key)
                                            if value == "" :
                                                print("consul key=%s value=%s, will skip" %(key, value))
                                                continue
                                            get_namespace_key_fail  = True
                                            getResp = self.PrivateApolloClient.get_namespace_items_key(key,appid,cluster,namespace)
                                            old_value = ""
                                            if not bool(getResp) :
                                                for item_map in get_namespace_resp["items"] :
                                                    if "key" in item_map :
                                                        get_namespace_key_fail  = False
                                                        if "value" in item_map :
                                                            old_value = item_map["value"]
                                                        break
                                            else :
                                                get_namespace_key_fail  = False
                                                if "value" in getResp :
                                                    old_value = getResp["value"] 
                                            if not get_namespace_key_fail :
                                                if old_value != value :
                                                    self.PrivateApolloClient.update_namespace_items_key(key, value,appid,cluster,namespace, comment="update k=%s ov=%s nv=%s " %(key,old_value,value))
                                                else :
                                                    print("noNeed to update !!!",getResp)
                                            else :
                                                self.PrivateApolloClient.create_namespace_items_key(key, value,appid,cluster,namespace, comment="insert k=%s v=%s" %(key,value))
                                    else :
                                        #namespace创建
                                        if bool(self.PrivateApolloClient.create_namespace(appid, namespace)) :
                                            for key in consulkeylist :
                                                value = operateConsul._getconsul(key)
                                                if value == "" :
                                                    print("consul key=%s value=%s, will skip" %(key, value))
                                                    continue
                                                get_namespace_key_fail  = True
                                                getResp = self.PrivateApolloClient.get_namespace_items_key(key,appid,cluster,namespace)
                                                old_value = ""
                                                if not bool(getResp) :
                                                    for item_map in get_namespace_resp["items"] :
                                                        if "key" in item_map :
                                                            get_namespace_key_fail  = False
                                                            if "value" in item_map :
                                                                old_value = item_map["value"]
                                                            break
                                                else :
                                                    get_namespace_key_fail  = False
                                                    if "value" in getResp :
                                                        old_value = getResp["value"] 
                                                if not get_namespace_key_fail :
                                                    if old_value != value :
                                                        self.PrivateApolloClient.update_namespace_items_key(key, value,appid,cluster,namespace, comment="update k=%s ov=%s nv=%s " %(key,old_value,value))
                                                    else :
                                                        print("noNeed to update !!!",getResp)
                                                else :
                                                    self.PrivateApolloClient.create_namespace_items_key(key, value,appid,cluster,namespace, comment="insert k=%s v=%s" %(key,value))
                                        else :
                                            print("create_namespace appid_%s cluster_%s namespace_%s failed" %(appid, cluster,namespace))
                                    self.PrivateApolloClient.releases("release", "release appid:%s cluster:%s"%(appid,cluster),appid, cluster,namespace)
                            else :
                                print("appid:%s namespace config not find" %(appid))
                                continue        
                        else :
                            print("create_cluster appid_%s cluster_%s failed" %(appid, cluster))
                            continue                

if __name__ == '__main__':
    appid = "dsp"
    inputuser = "apollo"
    initClient = InitApollo(user=inputuser,appid=appid)
    initClient.setup_apollo()

