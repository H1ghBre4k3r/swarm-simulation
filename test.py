#!/usr/bin/env python3 -u
import json
import sys
import numpy as np


def main(): 
    setup = json.loads(sys.stdin.readline())
    position = np.array([setup["position"]["x"], setup["position"]["y"]])
    target = np.array([setup["target"]["x"], setup["target"]["y"]])
    vmax = setup["vmax"]
    while True:
        inp = json.loads(sys.stdin.readline())
        position = np.array([inp["position"]["x"], inp["position"]["y"]])
        targetDir = target - position
        if np.linalg.norm(targetDir) > vmax:
            targetDir = targetDir / np.linalg.norm(targetDir)
            targetDir = targetDir * vmax
        
        val = {
            "action": "move",
            "payload": {
                "x": targetDir[0],
                "y": targetDir[1],
            }
        }
        print(json.dumps(val))

if __name__ == "__main__":
    main()
