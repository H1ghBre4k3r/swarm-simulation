#!/usr/bin/env python3

import argparse
import json
from fileinput import filename
from os import listdir, walk
from os.path import isfile, join

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
    tau = vals[3]
    noise = vals[4]
    consensus = vals[5]
    if p not in summaryByPartsByNoise:
        summaryByPartsByNoise[p] = {}
    if tau not in summaryByPartsByNoise[p]:
        summaryByPartsByNoise[p][tau] = {}
    if noise not in summaryByPartsByNoise[p][tau]:
        summaryByPartsByNoise[p][tau][noise] = {}
    if consensus not in summaryByPartsByNoise[p][tau][noise]:
        summaryByPartsByNoise[p][tau][noise][consensus] = {}
    if "summary" not in summaryByPartsByNoise[p][tau][noise][consensus]:
        summaryByPartsByNoise[p][tau][noise][consensus] = {}
    # aggregate information
    files = [f for f in listdir(folderPath) if isfile(join(folderPath, f))]
    for fileName in files:
        summary = json.loads(open(join(folderPath, fileName)).read())
        if "collisions" not in summaryByPartsByNoise[p][tau][noise][consensus]:
            summaryByPartsByNoise[p][tau][noise][consensus]["collisions"] = []
        summaryByPartsByNoise[p][tau][noise][consensus]["collisions"].append(summary["collisions"])
        if "runtime" not in summaryByPartsByNoise[p][tau][noise][consensus]:
            summaryByPartsByNoise[p][tau][noise][consensus]["runtime"] = []
        summaryByPartsByNoise[p][tau][noise][consensus]["runtime"].append(summary["runtime"])

json.dump(summaryByPartsByNoise, open(args.o, "w"))












