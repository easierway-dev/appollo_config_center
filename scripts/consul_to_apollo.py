#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os,json,toml
from os import walk
 
# 获取文件夹的中的文件夹和文件夹里文件
def do_file(save_filepath,o_filepath): #定义函数 传入写入文档保存的位置和要操作的任意电脑路劲
  file=open(save_filepath,"w+")
  #with open(save_filepath, "w") as fw: 
  # 遍历文件路径
  defmap={}
  for parent,dirnames,filenames in walk(o_filepath):
    clustn=parent.replace(o_filepath+"/","")
    if ".git" in parent or "script" in parent or parent.endswith("consul_backup"):
      continue
    print(("根目录为：{0}\n\n").format(parent))
    #parentn={0}.format(parent)a
    clusterlist = []
    if "dsp" in defmap.keys() and "cluster" in  defmap["dsp"].keys() and clustn.startswith("dsp_"):
       clusterlist = defmap["dsp"]["cluster"]
    if "as" in defmap.keys() and "cluster" in  defmap["as"].keys() and clustn.startswith("as_"):
       clusterlist = defmap["as"]["cluster"]
    if clustn not in clusterlist :
       clusterlist.append(clustn)
    #for dirname in dirnames:
    #  file.write(("  里面的文件夹有：{0}\n\n").format(dirname))
    filelist=[]
    if "dsp" in defmap.keys() and "consulkey" in  defmap["dsp"].keys() and clustn.startswith("dsp_"):
       filelist = defmap["dsp"]["consulkey"]
    if "as" in defmap.keys() and "consulkey" in  defmap["as"].keys() and clustn.startswith("as_"):
       filelist = defmap["as"]["consulkey"]
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
      #file.write(("  里面的文件有：{0}\n\n").format(filename))
  #toml.dump(defmap,fw)
  #file.write(json.dumps(defmap, sort_keys=True, indent=4, separators=(',', ':'),ensure_ascii=False))
  toml.dump(defmap,file)
  file.close()

util_path = os.path.realpath(__file__)
util_dir = os.path.dirname(util_path)  
do_file("%s/test.file"%util_dir,"%s/consul_backup"%util_dir)#传入相关的参数即可
