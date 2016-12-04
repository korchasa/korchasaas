import string
from stat import S_IFDIR, S_IFLNK, S_IFREG
from time import time
import logging

# from fuse import FUSE, FuseOSError, Operations, LoggingMixIn

class Tree():

    def __init__(self):
        self.tree = {}
        self.setDir('/')

    def nearest(self, path):
        if self.node(path):
            return self.node(path)
        elif '/' not in path:
            return self.node('/')
        else:
            return self.nearest(path.split('/', 1)[-1])

    def node(self, path):
        return self.tree.get(path, {})

    def children(self, path):
        if "/" == path:
            return {k:v for k,v in iter(self.tree.items()) if 1 == k.count('/')}
        else:
            return {k.rpartition('/')[-1]:v for k,v in iter(self.tree.items()) if k.startswith(path) and path != k}

    def setFile(self, path, data=None, mode=0o644):
        self.tree[path] = dict(dir=False, data=data, loaded= data != None)

    def setDir(self, path, loaded=False):
        self.tree[path] = dict(dir=True, loaded=loaded)

    def load(self, path, loader, attempts=10):
        if not attempts:
            return None
        node = self.node(path)
        if node:
            if node.get('loaded'):
                return self.node(path)
            elif attempts > 0:
                loader(path, self)
                self.tree[path]['loaded'] = True
                return self.load(path, loader, attempts - 1)
        else:
            node = self.nearest(path)
            loader(path, self)
            return self.load(path, loader, attempts - 1)
