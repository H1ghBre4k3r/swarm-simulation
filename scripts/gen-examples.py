#!/usr/bin/env -S python3

import argparse
import json
import os
from enum import Enum

import numpy as np


def angle2vec(angle: float) -> np.ndarray:
    """
    Convert an angle to a vector.
    """
    return np.array([np.cos(d2r(angle)), np.sin(d2r(angle))])

def d2r(deg: float) -> float:
    """
    Convert between degrees and radians.
    """
    return np.deg2rad(deg)

def dist(p1: np.ndarray, p2: np.ndarray) -> float:
    """
    Calculate the distance between two points.
    """
    return np.linalg.norm(p2 - p1)

class Mode(Enum):
    circle="circle"
    random="random"

    def __str__(self):
        return self.value

parser = argparse.ArgumentParser(description='Generate examples')
parser.add_argument("-n", type=int, required=True, help="Number of participants")
parser.add_argument("-r", type=float, required=True, help="Radius of the participants")
parser.add_argument("-d", type=float, default=0, help="Additional distance between participants")
parser.add_argument("-v", type=float, required=True, help="Velocity of the participants")
parser.add_argument("-t", type=float, default=1, help="TAU for simulation")
parser.add_argument("-s", type=str, required=True, help="Path to script")
parser.add_argument("-o", type=str, default="", help="Path to output")
parser.add_argument("-l", type=float, default=0.1, help="tick length")
parser.add_argument("-m", type=Mode, default=Mode.circle, choices=list(Mode), help="mode")

args = parser.parse_args()

def some(list_, pred):
    return any(pred(i) for i in list_) #booleanize the values, and pass them to any

participants = []
for i in range(args.n):
    if args.m == Mode.circle:
        start = angle2vec((360 / args.n) * i) * 0.4 + np.array([0.5, 0.5])
        target = start + (np.array([0.5, 0.5]) - start) * 2
    elif args.m == Mode.random:
        start = np.random.uniform(0, 1, 2)
        while some(participants, lambda p: dist(np.array([p["start"]["x"], p["start"]["y"]]), start) < args.r * 2 + args.d):
            start = np.random.uniform(0, 1, 2)
        target = np.random.uniform(0, 1, 2)
        while some(participants, lambda p: dist(np.array([p["target"]["x"], p["target"]["y"]]), target) < args.r * 2 + args.d):
            target = np.random.uniform(0, 1, 2)
    else:
        raise Exception("Unknown mode")
    participants.append({
        "start": {
            "x": start[0],
            "y": start[1]
        },
        "target": {
            "x": target[0],
            "y": target[1]
        },
        "vmax": args.v,
        "radius": args.r,
        "script": args.s
    })

settings = {
    "tickLength": args.l,
    "tau": args.t
}

configuration = {
    "settings": settings,
    "participants": participants,
}

json.dump(configuration, open(os.path.join(args.o, f"{args.n}-participants-{args.m}.json"), "w"), indent=2)
