import numpy as np
from mathutils import normalize


class Participant(object):
    def __init__(self, position: np.ndarray, velocity: np.ndarray, radius: np.float64, vmax=None, target=None) -> None:
        super().__init__()
        self.position = position
        self.velocity = velocity
        self.radius = radius
        self.vmax = vmax
        self.target = target

    def update_position(self, position: np.ndarray) -> None:
        self.position = position
        self.velocity = self.target - self.position
        if np.linalg.norm(self.velocity) > self.vmax:
            self.velocity = normalize(self.velocity) * self.vmax
