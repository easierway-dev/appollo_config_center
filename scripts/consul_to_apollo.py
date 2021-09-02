#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os,json,toml,sys,time
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

def find_check(a, alist) :
    find = False
    for checkstr in alist :
        if a.lower().find(checkstr.lower()) != -1 :
            find = True
            break
    return find

def not_find_check(a, alist) :
    nofind = True
    for checkstr in alist :
        if a.lower().find(checkstr.lower()) != -1 :
            nofind = False
            break
    return nofind

def list_minus(alist, blist) :
    for b in blist :
        if b in alist :
            alist.remove(b)
    return alist

def split_map_conf(source_map, merge_map, mapping_rule):
    merged_consul_list = []
    mapping_conf_map = {}
    if "namespace" in merge_map and "application" in merge_map["namespace"] :
        merged_consul_list.extend(merge_map["namespace"]["application"])
    needremove = []
    for rulemap in mapping_rule["mappingrules"] :
        for appid, matchlist in rulemap.items() :
            if not appid in mapping_conf_map:
                mapping_conf_map[appid] = []
            print("before: appid,matchlist=",appid,matchlist)
            print("before:merged_consul_list", merged_consul_list)
            if "white" in matchlist and len(matchlist["white"]) > 0 or "black" in matchlist and len(matchlist["black"]) > 0:
                for consulkey in merged_consul_list :
                    white_match = False
                    black_nomatch = False
                    if not "white" in matchlist or "white" in matchlist and find_check(consulkey, matchlist["white"]) :
                        white_match = True
                    if not "black" in matchlist or "black" in matchlist and not_find_check(consulkey, matchlist["black"]) :
                        black_nomatch = True
                    print("ing:merged_consul_list", merged_consul_list)
                    print("ing:consulkey, matchlist,white_match, black_nomatch", consulkey, matchlist, white_match, black_nomatch)
                    if white_match and black_nomatch :
                        if not consulkey in mapping_conf_map[appid] :
                            mapping_conf_map[appid].append(consulkey)
                        needremove.append(consulkey)
                merged_consul_list = list_minus(merged_consul_list, needremove)
                print("after:mapping_conf_map",mapping_conf_map[appid])
                print("after:merged_consul_list",merged_consul_list)
            else :               
                if appid in source_map :
                    mapping_conf_map[appid] =  = list_minus(source_map[appid]["namespace"]["application"], needremove)
                else :
                    mapping_conf_map[appid] = list(set(mapping_conf_map[appid]+merged_consul_list))
            if not appid in source_map :
                source_map[appid] = deepcopy(merge_map)
            source_map[appid]["namespace"]["application"] = mapping_conf_map[appid]

    #没有配置规则的默认不在apollo上同步
    for skey in source_map.keys() :
        if not skey in mapping_conf_map :
            del source_map[skey]
    return source_map

if __name__ == "__main__":
  print("start:",time.strftime("%Y-%m-%d %H:%M:%S", time.localtime()))
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

  #根据consul备份结果生成包含dsp/as的map
  source_conf_map = do_file(watch_path)

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
  abtest_key = "abtesting"
  abtest_value = "abtest/abtest_info"
  if "as" in final_conf_map :
    if abtest_value in final_conf_map["as"]["namespace"]["application"] :
      final_conf_map["as"]["namespace"]["application"].remove(abtest_value)
    final_conf_map["as"]["namespace"][abtest_key] = [abtest_value]
  if "dsp" in final_conf_map :
    if abtest_value in final_conf_map["dsp"]["namespace"]["application"] :
      final_conf_map["dsp"]["namespace"]["application"].remove(abtest_value)
    final_conf_map["dsp"]["namespace"][abtest_key] = [abtest_value]

  with open(gen_conf_path, "w") as fw: 
    #file.write(json.dumps(defmap, sort_keys=True, indent=4, separators=(',', ':'),ensure_ascii=False))
    toml.dump(final_conf_map,fw)
  print("end:",time.strftime("%Y-%m-%d %H:%M:%S", time.localtime()))
