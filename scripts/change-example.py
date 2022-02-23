#!/usr/bin/env -S python3

import argparse
import json

parser = argparse.ArgumentParser(description='Change examples')

parser.add_argument("-r", type=float, help="Radius of the participants")
parser.add_argument("-v", type=float, help="Velocity of the participants")
parser.add_argument("-z", type=float, help="Safezone for participants")

parser.add_argument("-t", type=float, help="TAU for simulation")
parser.add_argument("-s", type=str, required=True, help="Path to script")

args = parser.parse_args()
configuration = json.load(open(args.s, "r"));

if args.t:
    configuration["settings"]["tau"] = args.t

for p in configuration["participants"]:
    if args.r:
        p["radius"] = args.r
    if args.v:
        p["vmax"] = args.v
    if args.z:
        p["safezone"] = args.z

json.dump(configuration, open(args.s, "w"), indent=2)