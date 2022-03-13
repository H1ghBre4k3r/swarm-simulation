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
parser.add_argument("-p", nargs="+", required=True, help="Participants")
parser.add_argument("-c", action=argparse.BooleanOptionalAction, default=False, help="Consensus")
parser.add_argument("-m", type=Mode, required=True, choices=list(Mode), help="Mode: runtime or collisions")
parser.add_argument("-d", type=str, required=True, help="Name of detail")
parser.add_argument("-t", type=str, required=True, help="Value of tau")
parser.add_argument("-l", type=str, default="", help="Label for legend")
parser.add_argument("-xlabel", type=str, default="", help="Label for x-axis")
parser.add_argument("-ylabel", type=str, default="", help="Label for y-axis")

args = parser.parse_args()

file = json.loads(open(args.f).read())

summaries = {}

# participant -> tau -> noise -> consensus

for (n, p) in file.items():
    if n not in args.p:
        continue
    if args.t not in p:
        continue
    tau = p[args.t]
    for participant in tau:
        if participant not in summaries:
            summaries[participant] = {}
        noised = tau[participant]
        if str(args.c).lower() not in tau[participant]:
            continue
        consensus = noised[str(args.c).lower()]
        mode = consensus[args.m.value]
        detail = []
        for i in range (len(mode)):
            detail.append(mode[i][args.d]) 
        summaries[participant][n] = {
            "mean": np.mean(detail),
            "ci": st.t.interval(alpha=0.95, df=len(detail)-1, loc=np.mean(detail), scale=st.sem(detail))
        }

x = [float(k) for k in summaries.keys()]
x.sort()
x = [str(k) if k != 0 else "0" for k in x]
ys = {}
for n in x:
    for participant in args.p:
        if participant not in ys:
            ys[participant] = {
                "mean": [],
                "lower": [],
                "upper": []
            }
        if participant not in summaries[n]: 
            ys[participant]["mean"].append(nan)
            ys[participant]["lower"].append(nan)
            ys[participant]["upper"].append(nan)
        else:
            ys[participant]["mean"].append(summaries[n][participant]["mean"])
            if np.isnan(summaries[n][participant]["ci"][0]):
                ys[participant]["lower"].append(summaries[n][participant]["mean"])
            else:
                ys[participant]["lower"].append(summaries[n][participant]["ci"][0])
            if np.isnan(summaries[n][participant]["ci"][1]):
                
                ys[participant]["upper"].append(summaries[n][participant]["mean"])
            else:
                ys[participant]["upper"].append(summaries[n][participant]["ci"][1])

i = 0
for (n, y) in ys.items():
    plt.plot(x, y["mean"], linestyle=linestyle_tuple[i], label=f"{n} {args.l}")
    plt.fill_between(x, y["lower"], y["upper"], alpha=.4, label="95% CI")
    i = (i+1)%len(linestyle_tuple)

plt.xlabel(args.xlabel)
plt.ylabel(args.ylabel)

plt.legend(loc='upper left', ncol=2)

if args.o is None:
    plt.show()
else:
    plt.savefig(args.o, format="pdf")
