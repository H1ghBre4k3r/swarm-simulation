import sys

import numpy as np


class Obstacle(object):
    def __init__(self, coords):
        self.start = coords["start"]
        self.end = coords["end"]
