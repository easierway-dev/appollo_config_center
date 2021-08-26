#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os,json,toml,sys
from os import walk
 
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
        defmap["dsp"]["namespace"]["application"]=filelist

    if clustn.startswith("as_"):
      if not "as" in defmap.keys() :
        defmap["as"] = {}
        defmap["as"]["namespace"] = {}
      if len(clusterlist) > 0:
        defmap["as"]["cluster"] = clusterlist
      if len(filelist) > 0:
        defmap["as"]["namespace"]["application"]=filelist
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
                    input_json.append(join_json[index])
    else :
        print("%s:object type error %r %r %r %r" % (sys._getframe().f_code.co_name, input_json, type(input_json), join_json, type(join_json)))
        sys.exit(-1)

#根据映射规则将dsp/as的配置拆分成dsp/rtdsp(pioneer)/juno/dmp/drs(rs)
def split_map_conf(source_map, mapping_file):
  return source_map

if __name__ == "__main__":
  util_path = os.path.realpath(__file__)
  util_dir = os.path.dirname(util_path)
  mapping_conf_path = "%s/apollo_mapping.toml"%util_dir
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
  for task in tasklist :
    if not task in source_conf_map :
      merge_json = {}
      for _,value in source_conf_map.items() :
        json_merge_update(merge_json,deepcopy(value))
        source_conf_map[task] = deepcopy(merge_json)
  #根据mapping结果，对各业务线的配置进行瘦身（从全集中去掉不属于该业务线的内容）
  final_conf_map = split_map_conf(source_conf_map, mapping_conf_path)
  with open(gen_conf_path, "w") as fw: 
    #file.write(json.dumps(defmap, sort_keys=True, indent=4, separators=(',', ':'),ensure_ascii=False))
    toml.dump(final_conf_map,fw)
