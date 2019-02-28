#!/usr/bin/env python

import os
from mmap import mmap


file = open('/dev/shm/lassie', 'r+b')
size = file.seek(0, os.SEEK_END)
file.seek(0, os.SEEK_SET)
data = mmap(file.fileno(), size)
print(data)
