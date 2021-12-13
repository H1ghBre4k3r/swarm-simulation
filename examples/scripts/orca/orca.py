from agent import Participant
import numpy as np
import sys

tau = 100


def r2d(rad: float) -> float:
    return np.rad2deg(rad)


def d2r(deg: float) -> float:
    return np.deg2rad(deg)


def angle(x: np.ndarray, deg=True) -> float:
    ref = np.array([1, 0])
    ref_unit = ref / np.linalg.norm(ref)
    x_unit = x / np.linalg.norm(x)
    ang = np.arccos(np.dot(ref_unit, x_unit))
    if deg:
        ang = r2d(ang)
        if x[1] < 0:
            ang = -ang
    return ang


def angle2Vec(angle: float) -> np.ndarray:
    return np.array([np.cos(d2r(angle)), np.sin(d2r(angle))])


def arcsin(gegenKat: float, hypo: float, deg=True) -> float:
    val = np.arctan2(gegenKat, hypo)
    if deg:
        val = r2d(val)
    return val


def norm(x: np.ndarray) -> float:
    return np.linalg.norm(x)


def dist(x: np.ndarray, y: np.ndarray) -> float:
    return np.linalg.norm(x-y)


def orca(we: Participant, other: Participant) -> np.ndarray:
    x = other.position - we.position
    r = other.radius + we.radius
    # r *= 2
    positionAngle = angle(x)

    # information about the cone
    sideAngle = arcsin(norm(r), norm(x))
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
    vel = we.velocity - other.velocity
    velocityAngle = angle(vel)
    differenceAngle = np.abs(velocityAngle % 360 - positionAngle % 360)
    u_vec = np.array([0, 0])
    if differenceAngle < sideAngle:
        # we are colliding at some point in time
        if norm(vel) <= norm(discCenter) and dist(vel, discCenter) < discRadius:
            # we are colliding with the truncating disc
            sys.stderr.write("Collision with truncating disc\n")
        if norm(vel) >= norm(discCenter):
            # we are colliding with the cone
            distLeft = norm(
                np.cross(leftSide, -vel)) / norm(leftSide)
            distRight = norm(
                np.cross(rightSide, -vel)) / norm(rightSide)
            u = 0
            if distLeft < distRight:
                u = distLeft
                u_vec = np.array([-leftSide[1], leftSide[0]])
            else:
                u = distRight
                u_vec = np.array([leftSide[1], -leftSide[0]])
            u_vec = u_vec / norm(u_vec)
            u_vec *= u / 2

    return we.velocity + u_vec
