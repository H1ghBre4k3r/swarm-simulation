import numpy as np


class Participant(object):
    def __init__(self, position: np.ndarray, velocity: np.ndarray, radius: np.float64) -> None:
        super().__init__()
        self.position = position
        self.velocity = velocity
        self.radius = radius
