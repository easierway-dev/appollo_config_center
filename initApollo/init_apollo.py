#!/usr/bin/env python
# -*- coding: utf-8 -*-
# @Time    : 2021/08/16 18:52
# @Author  : Xuexia.li
# @Site    : 
# @File    : apollo_client_api.py
# @Description:
import requests
import json

class InitApollo(object):
    def __init__(self, user="", appid="", base_config = "./baseconfig.toml", apollo_to_consul_config="./apollo_to_consul.toml"):
		self.base_config_data = {}
		self.apollo_to_consul_config_data = {}
		with open(base_config, "r") as fs: 
			self.base_config_data = toml.load(fs)	
		with open(apollo_to_consul_config, "r") as fs: 
			self._apollo_to_consul_config_data = toml.load(fs)	
		if "portaddr" in self.base_config_data:
			_portaddr = self.base_config_data["portaddr"]
		else :
			_portaddr = "http://localhost:80"
		if "env" in self.base_config_data:
			_env = self.base_config_data["env"]
		else :
			_env = "DEV"
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
		PrivateApolloClient.__init__(self, _portaddr,_user,_token,_appid,_dev,_timeout)
		
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
					if bool(self.get_cluster(appid, cluster)) :
						#不使用默认namespace=application
						if "namespace" in clusters : 
							for namespace, consulkeylist in clusters["namespace"].items():
								#namespace已经创建
								if bool(get_namespace(appid, cluster,namespace)) :
									for key in consulkeylist :
										value = operateConsul._getconsul(key)
										getResp = self.get_namespace_items_key(key,appid,cluster,namespace)
										if bool(getResp) :
											if getResp["value"] != value :
												self.update_namespace_items_key(key, value,appid,cluster,namespace, comment="update k=%s ov=%s nv=%s " %(key,getResp.json()["value"],value))
												if bool(putResp) and enRelease :
													print("update:",self.releases("release", "release %s:%s"%(key,value)))
											else :
												print("noNeed to update !!!",getResp)
										else :
											self.create_namespace_items_key(key, value,appid,cluster,namespace, comment="insert k=%s v=%s" %(key,value))
											if bool(postResp) and enRelease :
												print("insert:",self.releases("release", "release %s:%s"%(key,value)))
								else :
									#namespace创建
									if bool(create_namespace(appid, cluster,namespace)) :
										for key in consulkeylist :
											value = operateConsul._getconsul(key)
											getResp = self.get_namespace_items_key(key,appid,cluster,namespace)
											if bool(getResp) :
												if getResp["value"] != value :
													self.update_namespace_items_key(key, value,appid,cluster,namespace, comment="update k=%s ov=%s nv=%s " %(key,getResp.json()["value"],value))
													if bool(putResp) and enRelease :
														print("update:",self.releases("release", "release %s:%s"%(key,value)))
												else :
													print("noNeed to update !!!",getResp)
											else :
												self.create_namespace_items_key(key, value,appid,cluster,namespace, comment="insert k=%s v=%s" %(key,value))
												if bool(postResp) and enRelease :
													print("insert:",self.releases("release", "release %s:%s"%(key,value)))
									else :
										print("create_namespace appid_%s cluster_%s namespace_%s failed" %(appid, cluster,namespace))
										continue						
						else :
							print("appid:%s namespace config not find failed" %(appid))
							continue
					else :
						#cluster创建
						if bool(self.creat_cluster(appid, cluster)) :
							#不使用默认namespace=application
							if "namespace" in clusters :
								for namespace, consulkeylist in clusters["namespace"].items():
									#namespace已经创建
									if bool(get_namespace(appid, cluster,namespace)) :
										for key in consulkeylist :
											value = operateConsul._getconsul(key)
											getResp = self.get_namespace_items_key(key,appid,cluster,namespace)
											if bool(getResp) :
												if getResp["value"] != value :
													self.update_namespace_items_key(key, value,appid,cluster,namespace, comment="update k=%s ov=%s nv=%s " %(key,getResp.json()["value"],value))
													if bool(putResp) and enRelease :
														print("update:",self.releases("release", "release %s:%s"%(key,value)))
												else :
													print("noNeed to update !!!",getResp)
											else :
												self.create_namespace_items_key(key, value,appid,cluster,namespace, comment="insert k=%s v=%s" %(key,value))
												if bool(postResp) and enRelease :
													print("insert:",self.releases("release", "release %s:%s"%(key,value)))
									else :
										#namespace创建
										if bool(create_namespace(appid, cluster,namespace)) :
											for key in consulkeylist :
												value = operateConsul._getconsul(key)
												getResp = self.get_namespace_items_key(key,appid,cluster,namespace)
												if bool(getResp) :
													if getResp["value"] != value :
														self.update_namespace_items_key(key, value,appid,cluster,namespace, comment="update k=%s ov=%s nv=%s " %(key,getResp.json()["value"],value))
														if bool(putResp) and enRelease :
															print("update:",self.releases("release", "release %s:%s"%(key,value)))
													else :
														print("noNeed to update !!!",getResp)
												else :
													self.create_namespace_items_key(key, value,appid,cluster,namespace, comment="insert k=%s v=%s" %(key,value))
													if bool(postResp) and enRelease :
														print("insert:",self.releases("release", "release %s:%s"%(key,value)))
										else :
											print("create_namespace appid_%s cluster_%s namespace_%s failed" %(appid, cluster,namespace))
											continue					
								
							else :
								print("appid:%s namespace config not find failed" %(appid))
								continue			
						else :
							print("create_cluster appid_%s cluster_%s failed" %(appid, cluster))
							continue				
				self.releases("release", "release appid:%s cluster:%s"%(appid,cluster))

if __name__ == '__main__':
    appid = "dsp"
    inputuser = "apollo"
	initClient = InitApollo(user=inputuser,appid=appid)
	initClient.setup_apollo()
	"""
    inputdata = {
        "dsp":{
            "cluster": [ "dsp_ali_cn", "dsp_hw_hk", "dsp_ali_vg",],
            "consulkey": [ "test", "xiongjia/aa",]
			"namespace":{
				"application": [ "test",],
				"abtesting": ["xiongjia/aa",]
			}
        }
    }
    for appid, clusters in inputdata.items() :
        if "cluster" in clusters and for cluster in clusters["cluster"] :
			operateConsul = OpConsul(cluster_map[cluster])
			if bool(apolloClient.get_cluster(appid, cluster)) :
				if "namespace" in clusters and for namespace, consulkeylist in clusters["namespace"].items():
					if bool(get_namespace(appid, cluster,namespace)) :
						for key in consulkeylist :
							value = operateConsul._getconsul(key)
							getResp = apolloClient.get_namespace_items_key(key,appid,cluster,namespace)
							if bool(getResp) :
								if getResp["value"] != value :
									apolloClient.update_namespace_items_key(key, value,appid,cluster,namespace, comment="update k=%s ov=%s nv=%s " %(key,getResp.json()["value"],value))
									if bool(putResp) and enRelease :
										print("update:",apolloClient.releases("release", "release %s:%s"%(key,value)))
								else :
									print("noNeed to update !!!",getResp)
							else :
								apolloClient.create_namespace_items_key(key, value,appid,cluster,namespace, comment="insert k=%s v=%s" %(key,value))
								if bool(postResp) and enRelease :
									print("insert:",apolloClient.releases("release", "release %s:%s"%(key,value)))
					else :
						if bool(create_namespace(appid, cluster,namespace)) :
							if "consulkey" in clusters and for key in clusters["consulkey"] :
								value = operateConsul._getconsul(key)
								getResp = apolloClient.get_namespace_items_key(key,appid,cluster,namespace)
								if bool(getResp) :
									if getResp["value"] != value :
										apolloClient.update_namespace_items_key(key, value,appid,cluster,namespace, comment="update k=%s ov=%s nv=%s " %(key,getResp.json()["value"],value))
										if bool(putResp) and enRelease :
											print("update:",apolloClient.releases("release", "release %s:%s"%(key,value)))
									else :
										print("noNeed to update !!!",getResp)
								else :
									apolloClient.create_namespace_items_key(key, value,appid,cluster,namespace, comment="insert k=%s v=%s" %(key,value))
									if bool(postResp) and enRelease :
										print("insert:",apolloClient.releases("release", "release %s:%s"%(key,value)))
						else :
							print("create_namespace appid_%s cluster_%s namespace_%s failed" %(appid, cluster,namespace))
							continue					
					
				else :
					namespace = "application"
					if bool(get_namespace(appid, cluster,namespace)) :
						if "consulkey" in clusters and for key in clusters["consulkey"] :
							value = operateConsul._getconsul(key)
							getResp = apolloClient.get_namespace_items_key(key,appid,cluster,namespace)
							if bool(getResp) :
								if getResp["value"] != value :
									apolloClient.update_namespace_items_key(key, value,appid,cluster,namespace, comment="update k=%s ov=%s nv=%s " %(key,getResp.json()["value"],value))
									if bool(putResp) and enRelease :
										print("update:",apolloClient.releases("release", "release %s:%s"%(key,value)))
								else :
									print("noNeed to update !!!",getResp)
							else :
								apolloClient.create_namespace_items_key(key, value,appid,cluster,namespace, comment="insert k=%s v=%s" %(key,value))
								if bool(postResp) and enRelease :
									print("insert:",apolloClient.releases("release", "release %s:%s"%(key,value)))
					else :
						if bool(create_namespace(appid, cluster,namespace)) :
							if "consulkey" in clusters and for key in clusters["consulkey"] :
								value = operateConsul._getconsul(key)
								getResp = apolloClient.get_namespace_items_key(key,appid,cluster,namespace)
								if bool(getResp) :
									if getResp["value"] != value :
										apolloClient.update_namespace_items_key(key, value,appid,cluster,namespace, comment="update k=%s ov=%s nv=%s " %(key,getResp.json()["value"],value))
										if bool(putResp) and enRelease :
											print("update:",apolloClient.releases("release", "release %s:%s"%(key,value)))
									else :
										print("noNeed to update !!!",getResp)
								else :
									apolloClient.create_namespace_items_key(key, value,appid,cluster,namespace, comment="insert k=%s v=%s" %(key,value))
									if bool(postResp) and enRelease :
										print("insert:",apolloClient.releases("release", "release %s:%s"%(key,value)))
						else :
							print("create_namespace appid_%s cluster_%s namespace_%s failed" %(appid, cluster,namespace))
							continue
			else :
				if bool(apolloClient.creat_cluster(appid, cluster))
				else :
					print("create_cluster appid_%s cluster_%s failed" %(appid, cluster))
					continue
				
			apolloClient.releases("release", "release appid:%s cluster:%s"%(appid,cluster))
			"""

