#!/usr/bin/env python
import logging
from sys import argv, exit
from apifs import ApiFs
from fuse import FUSE
from var_dump import var_dump
import requests
import json
from tree import Tree
from defaultenv import env
from errno import ENOENT
import pickle
import var_dump

if not hasattr(__builtins__, 'bytes'):
    bytes = str

if __name__ == '__main__':
    if len(argv) != 2:
        print('usage: %s <mountpoint>' % argv[0])
        exit(1)

    logging.basicConfig(level=logging.DEBUG)

    def loader(path, tree):
        logging.debug('load: %s', path)
        if '/' == path:
            tree.setDir('/', loaded=True)
            tree.setDir('/jobs_queue')
            tree.setFile('/features')
            tree.setFile('/status')
        else:
            resp = requests.get("http://korchasa.host/api/v1" + path.rstrip('/'))
            tree.setFile(path+'/request', var_dump.var_export(resp.request).encode('utf8'))
            tree.setDir(path, loaded=True)
            data = resp.json().get('data', [])
            for i, value in data.items() if isinstance(data, dict) else enumerate(data):
                content = json.dumps(value, ensure_ascii=False, indent=2).encode('utf8')
                tree.setFile(path+'/'+str(i).replace('/', '_'), content)
    fuse = FUSE(ApiFs(loader), argv[1], foreground=True)
