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
    x = a + 360
    x %= 360
    y = b + 360
    y %= 360
    return min(abs(x - y), abs(x - y - 360), abs(x - y + 360))


def orca(a: Participant, b: Participant) -> np.ndarray:
    x = b.position - a.position
    r = b.radius + a.radius
    # r *= 2
    positionAngle = angle(x)

    # information about the cone
    sideAngle = arcsin(r, norm(x))
    rightSideAngle = positionAngle - sideAngle
    leftSideAngle = positionAngle + sideAngle
    rightSide = angle2Vec(rightSideAngle)
    leftSide = angle2Vec(leftSideAngle)

    # information about the truncating disc
    discCenter = x / tau
    discRadius = r / tau

    # relative velocity calculations
    # if norm(we.velocity) <= 0 or norm(other.velocity) <= 0:
    #     return we.velocity
    vel = a.velocity - b.velocity
    velocityAngle = angle(vel)
    differenceAngle = angle_diff(positionAngle, velocityAngle)
    u_vec = np.array([0, 0])
    if differenceAngle < sideAngle:
        # we are colliding at some point in time
        if norm(vel) <= norm(discCenter) and dist(vel, discCenter) < discRadius:
            # we are colliding with the truncating disc
            sys.stderr.write("Collision with truncating disc\n")
            vec = vel - discCenter
            vecDist = norm(vec)
            u = (discRadius - vecDist)
            u_vec = (vec / vecDist) * u
            u_ang = angle(u_vec)
            u_ang += 10
            u_vec = angle2Vec(u_ang)
            u_vec = u_vec * u
        elif norm(vel) > norm(discCenter):
            leftPoint = closest_point_on_line(np.array([0, 0]), leftSide, vel)
            rightPoint = closest_point_on_line(
                np.array([0, 0]), rightSide, vel)
            left_u = leftPoint - vel
            right_u = rightPoint - vel
            leftDist = dist(leftPoint, vel)
            rightDist = dist(rightPoint, vel)

            if leftDist < rightDist:
                u_vec = left_u
            else:
                u_vec = right_u

            u_vec = u_vec / 2
    elif dist(vel, discCenter) < discRadius:
        vec = vel - discCenter
        vecDist = norm(vec)
        u = (discRadius - vecDist)
        u_vec = (vec / vecDist) * u

    return a.velocity + u_vec
