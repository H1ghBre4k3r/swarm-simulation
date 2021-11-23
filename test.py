#!/usr/bin/env python3 -u
import math
import time

i = 0
oldX = math.sin(i * (math.pi / 180)) * 0.3 + 0.5
oldY = math.cos(i * (math.pi / 180)) * 0.3 + 0.5
while True:
    x = math.sin(i * (math.pi / 180)) * 0.3 + 0.5
    y = math.cos(i * (math.pi / 180)) * 0.3 + 0.5
    print("{} {}".format(x - oldX, y - oldY))
    oldX = x
    oldY = y
    time.sleep(0.01)
    i += 1
