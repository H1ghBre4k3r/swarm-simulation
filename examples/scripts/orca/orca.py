import sys

import numpy as np
from halfplane import Halfplane
from mathutils import (angle2Vec, angle_diff, arcsin, closest_point_on_line,
                       dist, norm, normalize, vec2angle)
from participant import Participant

tau = 1000


def out_of_disc(disc_center, disc_r, v):
    # calculate vector & distance from center of disc to our velocity
    rel_vec = v - disc_center
    w_length = norm(rel_vec)
    # rotate vector by a certain degree, so we do not deadlock
    rel_vec = angle2Vec(vec2angle(rel_vec) + 10)
    # calculate u (the velocity that will get us out of collision)
    u_vec = rel_vec * (disc_r - w_length)
    return u_vec, rel_vec


def orca(a: Participant, b: Participant) -> np.ndarray:
    x = b.position - a.position
    r = b.radius + a.radius
    r *= 2
    v = a.velocity - b.velocity

    if norm(x) < r:
        return out_of_disc(x, r, v)
    else:
        # calculate "tau-disc" center and radius
        disc_center = x / tau
        disc_r = r / tau

        # center of disc needs to be adjusted, sinec sides of the cone are not parallel
        adjusted_disc_center = disc_center * (1 - (r / norm(x))**2)

        # check, if we are colliding with the truncating disc
        if norm(v) < norm(adjusted_disc_center) and norm(v - adjusted_disc_center) < r:
            # we are colliding with truncating disc
            return out_of_disc(disc_center, disc_r, v)
        else:
            # we are colliding with cone
            side_len = np.sqrt(norm(x)**2 - r**2)
            # # determine the side of the cone we have to project on
            signed_radius = np.copysign(r, np.linalg.det([v, x]))
            rot = np.array([[side_len, signed_radius],
                            [-signed_radius, side_len]])
            x_rot = rot.dot(x) / norm(x)**2
            n = np.array([x_rot[1], -x_rot[0]]) * np.sign(signed_radius)
            return x_rot * np.dot(v, x_rot) - v, n


def halfplane_intersection(halfplanes_u: list[Halfplane], current_velocity: np.ndarray, optimal_point: np.ndarray) -> np.ndarray:
    halfplanes = []
    for halfplane in halfplanes_u:
        halfplanes.append(
            Halfplane(current_velocity + halfplane.u, halfplane.n))
    new_point = optimal_point
    for i in range(len(halfplanes)):
        plane = halfplanes[i]
        # check, if we are out of the halfplane
        if np.dot(new_point - plane.u, plane.n) < 0:
            left, right = intersect_halfplane_with_other_halfplanes(
                plane, halfplanes[:i])
            if left is None or right is None:
                return None
            new_point = closest_point_on_line(
                plane.u, plane.u + np.array([plane.n[1], -plane.n[0]]), optimal_point, left, right)
    return new_point


def intersect_halfplane_with_other_halfplanes(plane: Halfplane, other_planes: list[Halfplane]) -> tuple[np.float64, np.float64]:
    left = float('-inf')
    right = float('inf')

    direction = np.array([plane.n[1], -plane.n[0]])

    # see https://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect/565282#565282 for refernce
    for other_plane in other_planes:
        other_dir = np.array([other_plane.n[1], -other_plane.n[0]])
        num = np.cross(other_plane.u - plane.u, other_dir)
        den = np.cross(direction, other_dir)

        if den == 0:
            # planes are parallel
            if num == 0:
                # planes are coincident
                return None, None
            else:
                # planes are parallel, but not coincident
                continue

        t = num / den
        if den > 0:
            # intersection is to the left of the line
            right = min(right, t)
        else:
            left = max(left, t)

        if left > right:
            return None, None
            # intersection is to the right of the line
    return left, right
