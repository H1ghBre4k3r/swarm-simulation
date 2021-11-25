#!/usr/bin/env python3 -u
import json
import math
import time


def main(): 
    i = 0
    oldX = math.sin(i * (math.pi / 180)) * 0.3 + 0.5
    oldY = math.cos(i * (math.pi / 180)) * 0.3 + 0.5
    while True:
        x = math.sin(i * (math.pi / 180)) * 0.3 + 0.5
        y = math.cos(i * (math.pi / 180)) * 0.3 + 0.5
        val = {
            "action": "move",
            "payload": {
                "x": x - oldX,
                "y": y - oldY,
            }
        }
        oldX = x
        oldY = y
        print(json.dumps(val))
        time.sleep(0.02)
        i += 2

if __name__ == "__main__":
    main()
