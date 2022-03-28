#!/usr/bin/python3

import re
from collections import defaultdict


metrics = defaultdict(int)

with open("test.log") as f:
    for l in f:
        res = re.findall(".*(main.go:\d+):", l)
        metrics[res[0]] += 1

print(metrics)