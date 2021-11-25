#!/usr/bin/env python3 -u
import json
import math
import sys
import time


def main(): 
    i = 0
    while True:
        inp = json.loads(sys.stdin.readline())
        pos = inp["position"]
        oldX = pos["x"]
        oldY = pos["y"]
        x = math.sin(i * (math.pi / 180)) * 0.3 + 0.5
        y = math.cos(i * (math.pi / 180)) * 0.3 + 0.5
        val = {
            "action": "move",
            "payload": {
                "x": x - oldX,
                "y": y - oldY,
            }
        }
        print(json.dumps(val))
        i += 4

if __name__ == "__main__":
    main()
