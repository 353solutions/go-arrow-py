#!/usr/bin/env python

import pyarrow as pa

size = 1 << 20
m = pa.memory_map('/dev/shm/shm-allocator')
buf = m.read_buffer(size)
