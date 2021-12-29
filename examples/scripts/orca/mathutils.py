import numpy as np


def r2d(rad: float) -> float:
    """
    Convert radians to degrees.
    """
    return np.rad2deg(rad)


def d2r(deg: float) -> float:
    """
    Convert between degrees and radians.
    """
    return np.deg2rad(deg)


def vec2angle(vector: np.ndarray, deg=True) -> float:
    """
    Calculate the angle of a vector.
    """
    ref = np.array([1, 0])
    ref_unit = ref / np.linalg.norm(ref)
    x_unit = vector / np.linalg.norm(vector)
    ang = np.arccos(np.dot(ref_unit, x_unit))
    if deg:
        ang = r2d(ang)
        if vector[1] < 0:
            ang = -ang
    return ang


def angle2vec(angle: float) -> np.ndarray:
    """
    Convert an angle to a vector.
    """
    return np.array([np.cos(d2r(angle)), np.sin(d2r(angle))])


def arcsin(gegenKat: float, hypo: float, deg=True) -> float:
    """
    Calculate the arcsin of a value.
    """
    val = np.arcsin(gegenKat / hypo)
    if deg:
        val = r2d(val)
    return val


def angle_diff(a, b):
    x = a + 360
    x %= 360
    y = b + 360
    y %= 360
    return min(abs(x - y), abs(x - y - 360), abs(x - y + 360))


def norm(x: np.ndarray) -> float:
    """"
    Calculate the norm of a vector.
    """
    return np.linalg.norm(x)


def dist(x: np.ndarray, y: np.ndarray) -> float:
    """
    Calculate the distance between two points.
    """
    return np.linalg.norm(x-y)


def mix(a, b, amount):
    return a + (b - a) * amount


def closest_point_on_line(l0: np.ndarray, l1: np.ndarray, tar: np.ndarray, left=0, right=1) -> np.ndarray:
    """
    Find the closest point on a line(-segment) to a target point.
    """
    c = tar - l0
    v = l1 - l0
    v = v / np.linalg.norm(v)
    d = norm(l0 - l1)
    t = np.dot(c, v) / d
    return mix(l0, l1, np.clip(t, left, right))


def normalize(n):
    return n / np.linalg.norm(n)
