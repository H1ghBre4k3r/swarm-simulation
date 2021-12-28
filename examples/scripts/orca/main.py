#!/usr/bin/env python3 -u
import json
import sys

import numpy as np
from halfplane import Halfplane
from mathutils import norm, normalize
from orca import halfplane_intersection, orca
from participant import Participant

# optional factor for changing "lookup range"
FPS = 120


def main():
    setup = json.loads(sys.stdin.readline())
    position = np.array([setup["position"]["x"], setup["position"]["y"]])
    radius = setup["radius"]
    target = np.array([setup["target"]["x"], setup["target"]["y"]])
    vmax = setup["vmax"] * FPS

    while True:
        inp = json.loads(sys.stdin.readline())
        position = np.array([inp["position"]["x"], inp["position"]["y"]])
        velocity = target - position

        if np.linalg.norm(velocity) > vmax:
            velocity = velocity / np.linalg.norm(velocity)
            velocity = velocity * vmax

        we = Participant(position, velocity, radius)

        # get u for each other participant
        halfplanes = []
        for p in inp["participants"]:
            participant = Participant(np.array([p["position"]["x"], p["position"]["y"]]), np.array(
                [p["velocity"]["x"], p["velocity"]["y"]]) * FPS, p["radius"])
            u, n = orca(we, participant)
            if norm(u) > 0:
                halfplanes.append(Halfplane(u, n))

        # get new "perfect" velocity from halfplane intersections
        new_vel = None
        while new_vel is None:
            new_vel = halfplane_intersection(
                halfplanes, we.velocity, we.velocity)
            new_halfplanes = []
            for l in halfplanes:
                offset = normalize(l.u)
                new_halfplanes.append(Halfplane(l.u - offset + 0.01, l.n))
            halfplanes = new_halfplanes
        velocity = new_vel

        # if it's too fast, slow it down
        if np.linalg.norm(velocity) > vmax:
            velocity = velocity / np.linalg.norm(velocity)
            velocity = velocity * vmax

        velocity /= FPS

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
