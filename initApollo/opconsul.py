#!/usr/bin/env python
# -*- coding: utf-8 -*-
# @Time    : 2021/08/11 18:52
# @Author  : xuexia.li
# @Site    : 
# @File    : opconsul.py
# @Description:
from consul_kv import Connection

class OpConsul(object):
    def __init__(self, consuladdr='http://localhost:8500/v1/'):
        self.conn = Connection(endpoint=consuladdr)
        
    def _putconsul(self, key, value):
        conn.put(key, value)
        
    def _getconsul(self, key):
        conn.gut(key)
    
    def _putmapconsul(self, consulmap):
        if isinstance(consulmap, dict) :
            conn.put_mapping(consulmap)
        else :
            print("need dict, real:",type(consulmap))

if __name__ == '__main__':
    pass
