import sys

import numpy as np


class Obstacle(object):
    def __init__(self, coords):
        self.coords = []
        for coord in coords:
            self.coords.append(np.array([coord["x"], coord["y"]]))
