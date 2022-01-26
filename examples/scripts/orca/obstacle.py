import sys

import numpy as np


class Obstacle(object):
    def __init__(self, coords, radius=0):
        self.start = np.array([coords["start"]["x"], coords["start"]["y"]])
        self.end = np.array([coords["end"]["x"], coords["end"]["y"]])
        self.radius = radius
