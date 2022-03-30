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
parser.add_argument("-f", nargs="+", required=True, help="Path to file")
parser.add_argument("-o", type=str, help="Path to output file")
parser.add_argument("-p", type=str, required=True, help="Participants")
parser.add_argument("-c", action=argparse.BooleanOptionalAction, default=False, help="Consensus")
parser.add_argument("-m", type=Mode, required=True, choices=list(Mode), help="Mode: runtime or collisions")
parser.add_argument("-s", type=float, default=1.0, help="Scale for Y-axis")
parser.add_argument("-ci", type=float, default=0.95, help="Confidence interval")
parser.add_argument("-d", type=str, required=True, help="Name of detail")
parser.add_argument("-t", type=str, required=True, help="Value of tau")
parser.add_argument("-l", nargs="+", default="", help="Labels for legend")
parser.add_argument("-xlabel", type=str, default="", help="Label for x-axis")
parser.add_argument("-ylabel", type=str, default="", help="Label for y-axis")

args = parser.parse_args()

summaries = {}

for f in args.f:
    file = json.loads(open(f).read())
    for (n, p) in file.items():
        if n != args.p:
            continue
        if args.t not in p:
            continue
        tau = p[args.t]
        for noise in tau:
            if noise not in summaries:
                summaries[noise] = {}
            noised = tau[noise]
            if str(args.c).lower() not in tau[noise]:
                continue
            consensus = noised[str(args.c).lower()]
            mode = consensus[args.m.value]
            detail = []
            for i in range (len(mode)):
                detail.append(mode[i][args.d]) 
            summaries[noise][f] = {
                "mean": np.mean(detail),
                "ci": st.t.interval(alpha=args.ci, df=len(detail)-1, loc=np.mean(detail), scale=st.sem(detail))
            }

x = [float(k) for k in summaries.keys()]
x.sort()
x = [str(k) if k != 0 else "0" for k in x]
ys = {}
for n in x:
    for noise in args.f:
        if noise not in ys:
            ys[noise] = {
                "mean": [],
                "lower": [],
                "upper": []
            }
        if noise not in summaries[n]: 
            ys[noise]["mean"].append(nan)
            ys[noise]["lower"].append(nan)
            ys[noise]["upper"].append(nan)
        else:
            ys[noise]["mean"].append(summaries[n][noise]["mean"]*args.s)
            if np.isnan(summaries[n][noise]["ci"][0]):
                ys[noise]["lower"].append(summaries[n][noise]["mean"]*args.s)
            else:
                ys[noise]["lower"].append(summaries[n][noise]["ci"][0]*args.s)
            if np.isnan(summaries[n][noise]["ci"][1]):
                
                ys[noise]["upper"].append(summaries[n][noise]["mean"]*args.s)
            else:
                ys[noise]["upper"].append(summaries[n][noise]["ci"][1]*args.s)

                
x = [str(float(k) * 1000) if k != "0" else "0" for k in x]
ci = args.ci * 100
i = 0
for (n, y) in ys.items():
    plt.plot(x, y["mean"], linestyle=linestyle_tuple[i], label=args.l[i%len(args.l)])
    plt.fill_between(x, y["lower"], y["upper"], alpha=.4, label=f"{ci if int(ci) != ci else int(ci)}% CI")
    i = (i+1)%len(linestyle_tuple)

plt.xlabel(args.xlabel)
plt.ylabel(args.ylabel)
plt.ylim(bottom=0)

plt.legend(loc='upper left', ncol=1)

if args.o is None:
    plt.show()
else:
    plt.savefig(args.o, format="pdf")
