#!/usr/bin/env python3

import argparse
import json
from cmath import nan
from enum import Enum

import matplotlib.pyplot as plt
import numpy as np
import scipy.stats as st

linestyle_tuple = [
     (0, (1, 1)),
     (0, (5, 10)),
     (0, (3, 5, 1, 5)),
     (0, (5, 5)),
     (0, (5, 1)),
     (0, (3, 10, 1, 10)),
     (0, (3, 1, 1, 1)),
     (0, (3, 5, 1, 5, 1, 5)),
     (0, (3, 10, 1, 10, 1, 10)),
     (0, (3, 1, 1, 1, 1, 1))]

class Mode(Enum):
    runtime="runtime"
    collisions="collisions"

    def __str__(self):
        return self.value

parser = argparse.ArgumentParser(description='Perform stuff')
parser.add_argument("-f", type=str, required=True, help="Path to file")
parser.add_argument("-o", type=str, help="Path to output file")
parser.add_argument("-n", type=str, required=True, help="Noise")
parser.add_argument("-c", action=argparse.BooleanOptionalAction, default=False, help="Consensus")
parser.add_argument("-m", type=Mode, required=True, choices=list(Mode), help="Mode: runtime or collisions")
parser.add_argument("-d", type=str, required=True, help="Name of detail")
parser.add_argument("-p", nargs="+", required=True, help="Participants")
parser.add_argument("-l", type=str, default="", help="Label for legend")
args = parser.parse_args()

file = json.loads(open(args.f).read())

summaries = {}

for (n, p) in file.items():
    if n not in args.p:
        continue
    for (tau, ns) in p.items():
        if args.n not in ns:
            continue
        noised = ns[args.n]
        if str(args.c).lower() not in noised:
            continue
        if tau not in summaries:
            summaries[tau] = {}
        consensus = noised[str(args.c).lower()]
        mode = consensus[args.m.value]
        detail = []
        for i in range(len(mode)):
            detail.append(mode[i][args.d])
        summaries[tau][n] = {
            "mean": np.mean(detail),
            "ci": st.t.interval(alpha=0.95, df=len(detail)-1, loc=np.mean(detail), scale=st.sem(detail))
        }

x = [int(k) for k in summaries.keys()]
x.sort()
x = [str(k) for k in x]
ys = {}

for tau in x:
    for p in args.p:
        if p not in ys:
            ys[p] = {
                "mean": [],
                "lower": [],
                "upper": []
            }
        if summaries[tau][p] is None:
            ys[p]["mean"].append(nan)
            ys[p]["lower"].append(nan)
            ys[p]["upper"].append(nan)
        else:
            ys[p]["mean"].append(summaries[tau][p]["mean"])
            if np.isnan(summaries[tau][p]["ci"][0]):
                ys[p]["lower"].append(summaries[tau][p]["mean"])
            else:
                ys[p]["lower"].append(summaries[tau][p]["ci"][0])
            if np.isnan(summaries[tau][p]["ci"][1]):
                ys[p]["upper"].append(summaries[tau][p]["mean"])
            else:
                ys[p]["upper"].append(summaries[tau][p]["ci"][1])

i = 0
for (n, y) in ys.items():
    plt.plot(x, y["mean"],linestyle=linestyle_tuple[i], label=f"{n} {args.l}")
    plt.fill_between(x, y["lower"], y["upper"], alpha=.4, label="95% CI")
    i = (i+1) % len(linestyle_tuple)

plt.xlabel("N° of participants")
plt.ylabel(args.m.value)

plt.legend()

if args.o is None:
    plt.show()
else:
    plt.savefig(args.o, format="pdf")
