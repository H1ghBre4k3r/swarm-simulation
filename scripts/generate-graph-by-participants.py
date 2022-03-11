#!/usr/bin/env python3

import argparse
import json
from enum import Enum

import matplotlib.pyplot as plt
import numpy as np
import scipy.stats as st


class Mode(Enum):
    runtime="runtime"
    collisions="collisions"

    def __str__(self):
        return self.value

parser = argparse.ArgumentParser(description='Perform stuff')
parser.add_argument("-f", type=str, required=True, help="Path to file")
parser.add_argument("-o", type=str, help="Path to output file")
parser.add_argument("-n", nargs="+", required=True, help="Noise")
parser.add_argument("-c", action=argparse.BooleanOptionalAction, default=False, help="Consensus")
parser.add_argument("-m", type=Mode, required=True, choices=list(Mode), help="Mode: runtime or collisions")
parser.add_argument("-d", type=str, required=True, help="Name of detail")

args = parser.parse_args()

file = json.loads(open(args.f).read())

summaries = {}

for (n, p) in file.items():
    summaries[n] = {}
    for noise in args.n:
        if noise not in p:
            summaries[n][noise] = None
            continue
        noised = p[noise]
        if str(args.c).lower() not in noised:
            summaries[n][noise] = None
            continue
        consensus = noised[str(args.c).lower()]["summary"]
        mode = consensus[args.m.value]
        detail = []
        for i in range(len(mode)):
            detail.append(mode[i][args.d])
        summaries[n][noise] = {
            "mean": np.mean(detail),
            "ci": st.t.interval(alpha=0.95, df=len(detail)-1, loc=np.mean(detail), scale=st.sem(detail))
        }

x = [int(k) for k in summaries.keys()]
x.sort()
x = [str(k) for k in x]
ys = {}
for n in x:
    for noise in args.n:
        if noise not in ys:
            ys[noise] = []
        if summaries[n][noise] is None:
            ys[noise].append(None)
        else:
            ys[noise].append(summaries[n][noise]["mean"])

for (n, y) in ys.items():
    plt.plot(x, y, label=n)

plt.legend()

if args.o is None:
    plt.show()
else:
    plt.savefig(args.o, format="pdf")
    
