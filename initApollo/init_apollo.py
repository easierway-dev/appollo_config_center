#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# @Time    : 2021/08/16 18:52
# @Author  : Xuexia.li
# @Site    : 
# @File    : apollo_client_api.py
# @Description:
import requests
import json,toml
import time
import argparse

from apollo_client_api import *
from opconsul import *

skipNull = False
cluster_num = 0
namespace_num = 0
key_num = 0

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
        if "usermap" in self.base_config_data and _appid in self.base_config_data["usermap"]:
             _token = self.base_config_data["usermap"][_appid]["token"]
        else :
            _token = "0bcbd744e2c08203a384a740f5aa9ab13f7cc24c" 
        if "timeout" in self.base_config_data:
            _timeout = self.base_config_data["timeout"]
        else :
            _timeout = 60  

        self._waittime = 5    
        self.PrivateApolloClient = PrivateApolloClient( _portaddr, _user, _token, _appid, _env, _timeout)

    def wait_update(self):
        time.sleep(self._waittime)
        
    def setup_apollo(self):
        global cluster_num, namespace_num, key_num
        for appid, clusters in self.apollo_to_consul_config_data.items() :
            if "cluster" in clusters : 
                for cluster in clusters["cluster"] :
                    if cluster in self.base_config_data :
                        consul_addr = "%s/v1/" %(self.base_config_data[cluster])
                        operateConsul = OpConsul(consul_addr)
                    else :
                        print("init consul:%s connection failed" %(cluster))
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
                                        #value = operateConsul._getconsul(key.replace(".","/"))
                                        value = operateConsul._getconsul(key)
                                        if value == None :
                                            value = ""
                                        if value == "" and skipNull or value == "error":
                                            print("consul key=%s, will skip" %(key))
                                            continue
                                        getResp = self.PrivateApolloClient.get_namespace_items_key(key,appid,cluster,namespace)
                                        if bool(getResp) :
                                            if "value" in getResp :
                                                old_value = getResp["value"] 
                                            if old_value != value :
                                                update_namespace_items_key_resp = self.PrivateApolloClient.update_namespace_items_key(key, value,appid,cluster,namespace,comment="update %s_%s_%s_%s" %(appid,cluster,namespace,key))
                                                if bool (update_namespace_items_key_resp) :
                                                    key_num += 1
                                                else :
                                                    print("update %s_%s_%s_%s failed" %(appid,cluster,namespace,key))
                                            else :
                                                print("noNeed to update !!!",getResp)
                                        else :
                                            if namespace == "abtesting" :
                                                create_namespace_items_key_resp = self.PrivateApolloClient.create_namespace_items_key_json(key, value,appid,cluster,namespace,comment="create %s_%s_%s_%s" %(appid,cluster,namespace,key))
                                            elif appid == "bidforce" :
                                                create_namespace_items_key_resp = self.PrivateApolloClient.create_namespace_items_key_toml(key, value,appid,cluster,namespace,comment="create %s_%s_%s_%s" %(appid,cluster,namespace,key))
                                            else :
                                                create_namespace_items_key_resp = self.PrivateApolloClient.create_namespace_items_key(key, value,appid,cluster,namespace,comment="create %s_%s_%s_%s" %(appid,cluster,namespace,key))
                                            if bool (create_namespace_items_key_resp) :
                                                key_num += 1
                                            else :
                                                print("create %s_%s_%s_%s failed" %(appid,cluster,namespace,key))
                                else :
                                    #namespace创建
                                    if bool(self.PrivateApolloClient.create_namespace(appid, namespace,comment="create %s_%s" %(appid,namespace))) :
                                        self.PrivateApolloClient.releases("release", "release appid:%s cluster:%s"%(appid,cluster),appid, cluster,namespace)
                                        self.wait_update()
                                        namespace_num += 1
                                        for key in consulkeylist :
                                            value = operateConsul._getconsul(key)
                                            if value == None :
                                                value = ""
                                            if value == "" and skipNull or value == "error":
                                                print("consul key=%s, will skip" %(key))
                                                continue
                                            getResp = self.PrivateApolloClient.get_namespace_items_key(key,appid,cluster,namespace)
                                            old_value = ""
                                            if bool(getResp) :
                                                if "value" in getResp :
                                                    old_value = getResp["value"]
                                                if old_value != value :
                                                    update_namespace_items_key_resp = self.PrivateApolloClient.update_namespace_items_key(key, value,appid,cluster,namespace,comment="update %s_%s_%s_%s" %(appid,cluster,namespace,key))
                                                    if bool(update_namespace_items_key_resp) :
                                                        key_num += 1
                                                    else :
                                                        print("update %s_%s_%s_%s failed" %(appid,cluster,namespace,key))
                                                else :
                                                    print("noNeed to update !!!",getResp)
                                            else :
                                                if namespace == "abtesting" :
                                                    create_namespace_items_key_resp = self.PrivateApolloClient.create_namespace_items_key_json(key, value,appid,cluster,namespace,comment="create %s_%s_%s_%s" %(appid,cluster,namespace,key))
                                                elif appid == "bidforce" :
                                                    create_namespace_items_key_resp = self.PrivateApolloClient.create_namespace_items_key_toml(key, value,appid,cluster,namespace,comment="create %s_%s_%s_%s" %(appid,cluster,namespace,key))
                                                else :
                                                    create_namespace_items_key_resp = self.PrivateApolloClient.create_namespace_items_key(key, value,appid,cluster,namespace,comment="create %s_%s_%s_%s" %(appid,cluster,namespace,key))
                                                if bool (create_namespace_items_key_resp) :
                                                    key_num += 1
                                                else :
                                                    print("create %s_%s_%s_%s failed" %(appid,cluster,namespace,key))                                                    
                                    else :
                                        print("create_namespace appid_%s cluster_%s namespace_%s failed" %(appid, cluster,namespace))
                                self.PrivateApolloClient.releases("release", "release appid:%s cluster:%s"%(appid,cluster),appid, cluster,namespace)
                        else :
                            print("appid:%s namespace config not find failed" %(appid))
                    else :
                        #cluster创建
                        if bool(self.PrivateApolloClient.create_cluster(appid, cluster)) :
                            cluster_num += 1
                            #不使用默认namespace=application
                            if "namespace" in clusters : 
                                for namespace, consulkeylist in clusters["namespace"].items():
                                    #namespace已经创建
                                    get_namespace_resp = self.PrivateApolloClient.get_namespace(appid, cluster,namespace)
                                    if bool(get_namespace_resp) :
                                        for key in consulkeylist :
                                            value = operateConsul._getconsul(key)
                                            if value == None :
                                                value = ""
                                            if value == "" and skipNull or value == "error":
                                                print("consul key=%s, will skip" %(key))
                                                continue
                                            getResp = self.PrivateApolloClient.get_namespace_items_key(key,appid,cluster,namespace)
                                            old_value = ""
                                            if bool(getResp) :
                                                if "value" in getResp :
                                                    old_value = getResp["value"]
                                                if old_value != value :
                                                    update_namespace_items_key_resp = self.PrivateApolloClient.update_namespace_items_key(key, value,appid,cluster,namespace,comment="update %s_%s_%s_%s" %(appid,cluster,namespace,key))
                                                    if bool(update_namespace_items_key_resp) :
                                                        key_num += 1
                                                    else :
                                                        print("update %s_%s_%s_%s failed" %(appid,cluster,namespace,key))                                                        
                                                else :
                                                    print("noNeed to update !!!",getResp)
                                            else :
                                                if namespace == "abtesting" :
                                                    create_namespace_items_key_resp = self.PrivateApolloClient.create_namespace_items_key_json(key, value,appid,cluster,namespace,comment="create %s_%s_%s_%s" %(appid,cluster,namespace,key))
                                                elif appid == "bidforce" :
                                                    create_namespace_items_key_resp = self.PrivateApolloClient.create_namespace_items_key_toml(key, value,appid,cluster,namespace,comment="create %s_%s_%s_%s" %(appid,cluster,namespace,key))
                                                else :
                                                    create_namespace_items_key_resp = self.PrivateApolloClient.create_namespace_items_key(key, value,appid,cluster,namespace,comment="create %s_%s_%s_%s" %(appid,cluster,namespace,key))
                                                if bool (create_namespace_items_key_resp) :
                                                    key_num += 1
                                                else :
                                                    print("create %s_%s_%s_%s failed" %(appid,cluster,namespace,key))
                                    else :
                                        #namespace创建
                                        if bool(self.PrivateApolloClient.create_namespace(appid, namespace,comment="create %s_%s" %(appid,namespace))) :
                                            self.PrivateApolloClient.releases("release", "release appid:%s cluster:%s"%(appid,cluster),appid, cluster,namespace)
                                            self.wait_update()
                                            namespace_num += 1
                                            for key in consulkeylist :
                                                value = operateConsul._getconsul(key)
                                                if value == None :
                                                    value = ""                                                
                                                if value == "" and skipNull or value == "error":
                                                    print("consul key=%s, will skip" %(key))
                                                    continue
                                                getResp = self.PrivateApolloClient.get_namespace_items_key(key,appid,cluster,namespace)
                                                old_value = ""
                                                if bool(getResp) :
                                                    if "value" in getResp :
                                                        old_value = getResp["value"]
                                                    if old_value != value :
                                                        update_namespace_items_key_resp = self.PrivateApolloClient.update_namespace_items_key(key, value,appid,cluster,namespace,comment="update %s_%s_%s_%s" %(appid,cluster,namespace,key))
                                                        if bool (update_namespace_items_key_resp) :
                                                            key_num += 1
                                                        else :
                                                            print("update_namespace_items_key %s_%s_%s_%s failed" %(appid,cluster,namespace,key))
                                                    else :
                                                        print("noNeed to update !!!",getResp)
                                                else :
                                                    if namespace == "abtesting" :
                                                        create_namespace_items_key_resp = self.PrivateApolloClient.create_namespace_items_key_json(key, value,appid,cluster,namespace,comment="create %s_%s_%s_%s" %(appid,cluster,namespace,key))
                                                    elif appid == "bidforce" :
                                                        create_namespace_items_key_resp = self.PrivateApolloClient.create_namespace_items_key_toml(key, value,appid,cluster,namespace,comment="create %s_%s_%s_%s" %(appid,cluster,namespace,key))
                                                    else :
                                                        create_namespace_items_key_resp = self.PrivateApolloClient.create_namespace_items_key(key, value,appid,cluster,namespace,comment="create %s_%s_%s_%s" %(appid,cluster,namespace,key))
                                                    if bool (create_namespace_items_key_resp) :
                                                        key_num += 1
                                                    else :
                                                        print("create %s_%s_%s_%s failed" %(appid,cluster,namespace,key))                                                        
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
    parser = argparse.ArgumentParser(description='manual to this script')
    parser.add_argument('--appid', type=str, default = "dsp")
    parser.add_argument('--user', type=str, default = "apollo")

    args = parser.parse_args()
    appid = args.appid
    inputuser = args.user

    initClient = InitApollo(user=inputuser,appid=appid)
    initClient.setup_apollo()
    print("cluster_num=%d namespace_num=%d key_num=%d" %(cluster_num,namespace_num,key_num))

