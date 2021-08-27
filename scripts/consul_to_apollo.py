#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os,json,toml,sys
from os import walk
from copy import deepcopy
 
# 获取文件夹的中的文件夹和文件夹里文件
def do_file(o_filepath): #定义函数 传入写入文档保存的位置和要操作的任意电脑路劲
  # 遍历文件路径
  defmap={}
  for parent,dirnames,filenames in walk(o_filepath):
    clustn=parent.replace(o_filepath+"/","")
    #if ".git" in parent or "script" in parent or parent.endswith("consul_backup"):
    if not parent.endswith(("dsp_ali_cn","dsp_ali_vg","dsp_hw_hk","as_ali_vg","as_aws_fk","as_aws_se","as_aws_sg","as_aws_vg")):
      continue
    print(("根目录为：{0}\n\n").format(parent))

    #####获取cluster
    clusterlist = []
    if "dsp" in defmap.keys() and "cluster" in  defmap["dsp"].keys() and clustn.startswith("dsp_"):
       clusterlist = defmap["dsp"]["cluster"]
    if "as" in defmap.keys() and "cluster" in  defmap["as"].keys() and clustn.startswith("as_"):
       clusterlist = defmap["as"]["cluster"]
    if clustn not in clusterlist :
       clusterlist.append(clustn)

    #####获取consul keylist
    filelist=[]
    if "dsp" in defmap.keys() and "namespace" in  defmap["dsp"].keys() and "application" in defmap["dsp"]["namespace"].keys() and clustn.startswith("dsp_"):
       filelist = defmap["dsp"]["namespace"]["application"]
    if "as" in defmap.keys() and "namespace" in  defmap["as"].keys() and "application" in defmap["as"]["namespace"].keys() and clustn.startswith("as_"):
       filelist = defmap["as"]["namespace"]["application"]
    for filename in filenames:
      #fileddn={0}.format(filename)
      if filename.endswith("#"):
        continue
      filekey=filename.replace("#","/")
      if filekey not in filelist :
        filelist.append(filekey)

    if clustn.startswith("dsp_"):
      if not "dsp" in defmap.keys() :
        defmap["dsp"] = {}
        defmap["dsp"]["namespace"] = {}
      if len(clusterlist) > 0:
        defmap["dsp"]["cluster"] = clusterlist
      if len(filelist) > 0:
        defmap["dsp"]["namespace"]["application"]=sorted(filelist)

    if clustn.startswith("as_"):
      if not "as" in defmap.keys() :
        defmap["as"] = {}
        defmap["as"]["namespace"] = {}
      if len(clusterlist) > 0:
        defmap["as"]["cluster"] = clusterlist
      if len(filelist) > 0:
        defmap["as"]["namespace"]["application"]=sorted(filelist)
  return defmap

def json_merge_update(input_json, join_json) :
    if isinstance(join_json, dict) and isinstance(input_json, dict):
        for k,v in join_json.items() :
            if k not in input_json :
                input_json[k] = v
                continue
            else :
                if isinstance(v, (dict)) or isinstance(v, (tuple, list)):
                    json_merge_update(input_json[k], v)
                else :
                    input_json[k] = v
    elif isinstance(input_json, (tuple, list)) and isinstance(join_json, (tuple, list)):
        simpleFlag = True
        for index in range(len(join_json)) :
            if isinstance(join_json[index], dict) or isinstance(join_json[index], (tuple, list)) :
                simpleFlag = False
        if simpleFlag == True :
            for index in range(len(join_json)) :
                if not join_json[index] in input_json :
                    input_json.append(join_json[index])
        else :
            for index in range(len(join_json)) :
                if index < len(input_json) :
                    json_merge_update(input_json[index], join_json[index])
                else :
                    if not join_json[index] in input_json :
                        input_json.append(join_json[index])
    else :
        print("%s:object type error %r %r %r %r" % (sys._getframe().f_code.co_name, input_json, type(input_json), join_json, type(join_json)))
        sys.exit(-1)

def split_map_conf(source_map, merge_map, mapping_rule):
    merged_consul_list = []
    mapping_conf_map = {}
    if "namespace" in merge_map and "application" in merge_map["namespace"] :
        merged_consul_list.extend(merge_map["namespace"]["application"])
    for rulemap in mapping_rule["mappingrules"] :
        for appid, matchlist in rulemap.items() :
            if not appid in mapping_conf_map:
                mapping_conf_map[appid] = []
            print("xxdebug matchlist=",appid,appid)
            print("before:", merged_consul_list)
            needremove = []
            if len(matchlist) > 0 :
                for consulkey in merged_consul_list :
                    print("ing:", merged_consul_list)
                    print("ing:", consulkey, matchlist,findcheck(consulkey, matchlist))
                    if findcheck(consulkey, matchlist) :
                        if not consulkey in mapping_conf_map[appid] :
                            mapping_conf_map[appid].append(consulkey)
                        needremove.append(consulkey)
                for rmkey in needremove :
                    merged_consul_list.remove(rmkey)
                print("after:",mapping_conf_map[appid])
                print("after:",merged_consul_list)
                if appid in sourcemap :
                    sourcemap[appid]["namespace"]["application"] = mapping_conf_map[appid]
                else :
                    sourcemap[appid] = deepcopy(merge_map)
                    sourcemap[appid]["namespace"]["application"] = mapping_conf_map[appid]
            else :
                mapping_conf_map[appid] = list(set(mapping_conf_map[appid]+merged_consul_list))
                if appid in sourcemap :
                    sourcemap[appid]["namespace"]["application"] = mapping_conf_map[appid]
                else :
                    sourcemap[appid] = deepcopy(merge_map)
                    sourcemap[appid]["namespace"]["application"] = mapping_conf_map[appid]
    return source_map

if __name__ == "__main__":
  util_path = os.path.realpath(__file__)
  util_dir = os.path.dirname(util_path)
  mapping_rule_path = "%s/apollo_mapping.toml"%util_dir
  gen_conf_path = "%s/consul_to_apollo.toml"%util_dir
  watch_path = "%s/consul_backup"%util_dir
  tasklist = ["dsp","as","rtdsp","juno","dmp","drs"]
  if len(sys.argv) == 2 :
    watch_path = sys.argv[1]
  if len(sys.argv) == 3 :
    watch_path = sys.argv[1]
    gen_conf_path = sys.argv[2]
  if len(sys.argv) == 4 :
    watch_path = sys.argv[1]
    gen_conf_path = sys.argv[2]
    mapping_conf_path = sys.argv[3]

  source_conf_map = do_file(watch_path)#根据consul备份结果生成包含dsp/as的map

  #除dsp/as之外的服务的配置默认是dsp/as的并集
  merge_map = {}
  for _,value in source_conf_map.items() :
    json_merge_update(merge_map,deepcopy(value))
  
  if os.path.exists(mapping_rule_path) :
    mapping_rule = toml.load(mapping_rule_path, _dict=dict)
  else :
    mapping_rule = {"mappingrules":[]}
  #根据mapping结果，对各业务线的配置进行瘦身（从全集中去掉不属于该业务线的内容）
  final_conf_map = split_map_conf(source_conf_map, merge_map, mapping_rule)
  
  #abtest信息独立namespace存储
  final_conf_map["as"]["namespace"]["abtesting"] = ["abtest/abtest_info"]
  final_conf_map["as"]["namespace"]["application"].remove("abtest/abtest_info")
  final_conf_map["dsp"]["namespace"]["abtesting"] = ["abtest/abtest_info"]
  final_conf_map["dsp"]["namespace"]["application"].remove("abtest/abtest_info")

  with open(gen_conf_path, "w") as fw: 
    #file.write(json.dumps(defmap, sort_keys=True, indent=4, separators=(',', ':'),ensure_ascii=False))
    toml.dump(final_conf_map,fw)
