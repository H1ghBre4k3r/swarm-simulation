#!/usr/bin/env python3 -u
import math
import time

i = 0
while True:
    x = math.sin(i * (math.pi / 180)) * 300 + 512
    y = math.cos(i * (math.pi / 180)) * 300 + 512
    print("{} {}".format(x, y))
    time.sleep(0.01)
    i += 1
