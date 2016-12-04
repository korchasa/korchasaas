from __future__ import print_function, absolute_import, division
from fuse import FuseOSError, Operations, LoggingMixIn
from collections import defaultdict
from errno import ENOENT
from stat import S_IFDIR, S_IFLNK, S_IFREG
from time import time
from tree import Tree
import json
import requests
import sys

class ApiFs(LoggingMixIn, Operations):
    'Api filesystem for http://korchasa.host'

    def __init__(self, loader):
        self.tree = Tree()
        self.loader = loader
        self.fd = 0

    def create(self, path, mode=0o644):
        self.files[path] = self._file(mode)
        self.fd += 1
        return self.fd

    getxattr = None

    def getattr(self, path, fh=None):
        if len(path) > 1 and path[1] == '.':
            raise FuseOSError(ENOENT)
        node = self.tree.load(path, self.loader)
        if node:
            return self._fs_node(node)
        else:
            raise FuseOSError(ENOENT)

    def open(self, path, flags):
        self.fd += 1
        return self.fd

    def read(self, path, size, offset, fh):
        return self.tree.node(path)['data'][offset:offset + size]

    def readdir(self, path, fh):
        node = self.tree.nearest(path)
        if not node.get('loaded'):
            self.loader(path, self.tree)
        children = ['.', '..'] + [url.lstrip('/') for url in list(self.tree.children(path).keys()) if url != '/']

        return children

    def statfs(self, path):
        return dict(f_bsize=512, f_blocks=4096, f_bavail=2048)

    def _fs_node(self, tree_node):
        if tree_node.get('dir'):
            return dict(
                st_mode=(S_IFDIR | 0o755),
                st_ctime=time(),
                st_mtime=time(),
                st_atime=time(),
                st_nlink=2
            )
        else:
            return dict(
                st_mode=(S_IFREG | 0o644),
                st_nlink=1,
                st_size=len(tree_node.get('data')),
                st_ctime=time(),
                st_mtime=time(),
                st_atime=time()
            )
