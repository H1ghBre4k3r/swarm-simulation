#!/usr/bin/env python3

import argparse
import subprocess

parser = argparse.ArgumentParser(description='Perform stuff')
parser.add_argument("-e", type=str, required=True, help="Path to exectuable")
parser.add_argument("-s", nargs="+", required=True, help="Path to script")
parser.add_argument("-o", type=str, required=True, help="Path to output")
parser.add_argument("-c", action=argparse.BooleanOptionalAction, default=False, help="Consensus")

noises = [
    0, 0.0001, 0.0002, 0.0005, 0.001, 0.002, 0.005, 
0.01, 
# 0.02, 0.05
]

args = parser.parse_args()

for script in args.s:
    for noise in noises:
        for i in range(10):
            subprocess.call([args.e, "-c", script, "-o", args.o, "-n", str(noise), f"-consensus={args.c}"])
