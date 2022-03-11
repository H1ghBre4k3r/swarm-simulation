#!/usr/bin/env python3

import argparse
import json
from os import listdir, walk
from os.path import isfile, join

import matplotlib.pyplot as plt
import numpy as np

parser = argparse.ArgumentParser(description='Perform stuff')
parser.add_argument("-d", type=str, required=True, help="Path to directory")
parser.add_argument("-o", type=str, required=True, help="Path to output file")
args = parser.parse_args()

folders = listdir(args.d)

summaryByPartsByNoise = {}

for folderName in folders:
    if folderName.startswith("."):
        continue
    folderPath = join(args.d, folderName)
    # parse the folder name
    vals = folderName.split("-")
    p = vals[0]
    # tau = vals[3]
    noise = vals[4]
    consensus = vals[5]
    if p not in summaryByPartsByNoise:
        summaryByPartsByNoise[p] = {}
    if noise not in summaryByPartsByNoise[p]:
        summaryByPartsByNoise[p][noise] = {}
    if consensus not in summaryByPartsByNoise[p][noise]:
        summaryByPartsByNoise[p][noise][consensus] = {}
    if "summary" not in summaryByPartsByNoise[p][noise][consensus]:
        summaryByPartsByNoise[p][noise][consensus]["summary"] = {}
    # aggregate information
    files = [f for f in listdir(folderPath) if isfile(join(folderPath, f))]
    for fileName in files:
        summary = json.loads(open(join(folderPath, fileName)).read())
        if "collisions" not in summaryByPartsByNoise[p][noise][consensus]["summary"]:
            summaryByPartsByNoise[p][noise][consensus]["summary"]["collisions"] = []
        summaryByPartsByNoise[p][noise][consensus]["summary"]["collisions"].append(summary["collisions"])
        if "runtime" not in summaryByPartsByNoise[p][noise][consensus]["summary"]:
            summaryByPartsByNoise[p][noise][consensus]["summary"]["runtime"] = []
        summaryByPartsByNoise[p][noise][consensus]["summary"]["runtime"].append(summary["runtime"])

json.dump(summaryByPartsByNoise, open(args.o, "w"))
















# x = np.linspace(0, 2 * np.pi, 200)
# y = np.sin(x)

# fig, ax = plt.subplots()
# ax.plot(x, y)
# plt.savefig('sin.svg', format="svg", transparent=True)
