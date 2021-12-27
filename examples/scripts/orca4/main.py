#!/usr/bin/env python3 -u
import json
import sys

import numpy as np
from orca import orca
from participant import Participant


def main():
    setup = json.loads(sys.stdin.readline())
    position = np.array([setup["position"]["x"], setup["position"]["y"]])
    radius = setup["radius"]
    target = np.array([setup["target"]["x"], setup["target"]["y"]])
    vmax = setup["vmax"]
    while True:
        inp = json.loads(sys.stdin.readline())
        position = np.array([inp["position"]["x"], inp["position"]["y"]])
        velocity = target - position

        if np.linalg.norm(velocity) > vmax:
            velocity = velocity / np.linalg.norm(velocity)
            velocity = velocity * vmax

        we = Participant(position, velocity, radius)

        for p in inp["participants"]:
            participant = Participant(np.array([p["position"]["x"], p["position"]["y"]]), np.array(
                [p["velocity"]["x"], p["velocity"]["y"]]), p["radius"])
            u, n = orca(we, participant)
            velocity += u

        val = {
            "action": "move",
            "payload": {
                "x": velocity[0],
                "y": velocity[1],
            },
            "debug": inp["participants"]
        }
        print(json.dumps(val))


if __name__ == "__main__":
    main()
