from participant import Participant
import numpy as np
import sys

tau = 2


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


def angle(vector: np.ndarray, deg=True) -> float:
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


def angle2Vec(angle: float) -> np.ndarray:
    """
    Convert an angle to a vector.
    """
    return np.array([np.cos(d2r(angle)), np.sin(d2r(angle))])


def arcsin(gegenKat: float, hypo: float, deg=True) -> float:
    """
    Calculate the arcsin of a value.
    """
    val = np.arctan2(gegenKat, hypo)
    if deg:
        val = r2d(val)
    return val


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


def closest_point_on_line(l0: np.ndarray, l1: np.ndarray, tar: np.ndarray) -> np.ndarray:
    """
    Find the closest point on a line to a target point.
    """
    c = tar - l0
    v = l1 - l0
    v = v / np.linalg.norm(v)
    d = norm(l0 - l1)
    t = np.dot(c, v) / d
    return mix(l0, l1, np.clip(t, 0, 1))


def angle_diff(a, b):
    if a < 0:
        a += 360
    if b < 0:
        b += 360
    return np.abs(a - b)


def orca(a: Participant, b: Participant) -> np.ndarray:
    x = b.position - a.position
    r = b.radius + a.radius
    v = a.velocity - b.velocity

    v_angle = angle(v)
    v_angle += 180
    v_angle %= 360
    v_angle -= 180

    x_angle = angle(x)
    x_angle -= v_angle
    x_angle += 180
    x_angle %= 360
    x_angle -= 180

    side_angle = arcsin(r, norm(x))
    right_angle = x_angle + side_angle
    left_angle = x_angle - side_angle
    right_side = angle2Vec(right_angle)
    left_side = angle2Vec(left_angle)

    disc_center = x / tau
    disc_r = r / tau

    u = np.array([0, 0])

    if np.abs(x_angle) < side_angle:
        if norm(v) <= norm(disc_center) and dist(v, disc_center) < disc_r:
            sys.stderr.write("ORCA: v is inside the disc\n")
        else:
            left_point = closest_point_on_line(np.array([0, 0]), left_side, v)
            right_point = closest_point_on_line(
                np.array([0, 0]), right_side, v)
            left_u = left_point - v
            right_u = right_point - v
            left_dist = norm(left_u)
            right_dist = norm(right_u)

            if left_dist < right_dist:
                u = left_u
            else:
                u = right_u

            l = norm(u)
            u_angle = angle(u)
            u_angle += v_angle
            u = angle2Vec(u_angle)
            u = u / norm(u)
            u *= l
            sys.stderr.write("ORCA: u is {}\n".format(l))

    return a.velocity + u / 2
