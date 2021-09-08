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
        
    #conn.put('the/key', 'the_value')
    def _putconsul(self, key, value):
        try:
            self.conn.put(key, value)
            resp = self.conn.get(key)
            if key in resp :
                return resp[key]
            else :
                print("%s: response code is %d" %(sys._getframe().f_code.co_name, resp.status_code))
                return ""
        except BaseException as e:
            print("_putconsul err", e)
            return ""
        
        
    def _getconsul(self, key):
        try:
            resp = self.conn.get(key)
            if key in resp :
                return resp[key]
            else :
                print("%s: response code is %d" %(sys._getframe().f_code.co_name, resp.status_code))
                return ""
        except BaseException as e:
            print("_getconsul err", e)
            return "error"
    """
    mapping = {
    'a/key': 'a_value',
    'another/k': 'another_value'
    }
    """
    def _putmapconsul(self, consulmap):
        if isinstance(consulmap, dict) :
            try:
                self.conn.put_mapping(consulmap)
                return True
            except BaseException as e:
                print("_putmapconsul err", e)
                return False
        else :
            print("need map, real:",type(consulmap))
            return False
    """
    dictionary = {
        'a': {
            'key': 'a_value'
        },
        'another': {
            'k': 'another_value'
        }
    }
    """
    def _putdictconsul(self, consuldict):
        if isinstance(consuldict, dict) :
            try:
                self.conn.put_dict(consuldict)
                return True
            except BaseException as e:
                print("_putdictconsul err", e)
                return False
        else :
            print("need dict, real:",type(consuldict))
            return False

if __name__ == '__main__':
    opconsul = OpConsul(consuladdr="http://47.252.4.203:8500/v1/")
    print(opconsul._getconsul("cnmultins"))
    print(opconsul._putconsul("cnmultins11",11))
    pass
