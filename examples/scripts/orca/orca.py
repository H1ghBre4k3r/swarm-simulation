import sys

import numpy as np
from halfplane import Halfplane
from mathutils import (angle2vec, angle_diff, arcsin, closest_point_on_line,
                       dist, norm, normalize, vec2angle)
from obstacle import Obstacle
from participant import Participant

tau = 4


def obstacle_collision(a: Participant, obstacle: Obstacle) -> tuple[np.ndarray, np.ndarray]:
    r_vec = normalize(a.velocity) * (a.radius + a.safezone)
    start = obstacle.start - (a.position + r_vec)
    end = obstacle.end - (a.position + r_vec)

    v = a.velocity
    v_norm = norm(v)

    start_tau = start / tau
    end_tau = end / tau

    # determine the side to project on
    closest_front = closest_point_on_line(start_tau, end_tau, a.velocity)
    closest_side_start = closest_point_on_line(
        start_tau, normalize(start_tau) * v_norm, a.velocity)
    closest_side_end = closest_point_on_line(
        end_tau, normalize(end_tau) * v_norm, a.velocity)
    closest = closest_front
    if dist(closest_side_start, a.velocity) < dist(closest, a.velocity):
        closest = closest_side_start
    if dist(closest_side_end, a.velocity) < dist(closest, a.velocity):
        closest = closest_side_end

    w = closest - v
    u = w
    n = normalize(w)
    if np.dot(v - closest, closest) < 0:
        n *= -1
    return u, n


def out_of_disc(disc_center, disc_r, v):
    # calculate vector & distance from center of disc to our velocity
    rel_vec = v - disc_center
    w_length = norm(rel_vec)
    # rotate vector by a certain degree, so we do not deadlock
    rel_vec = angle2vec(vec2angle(rel_vec) + 10)
    # calculate u (the velocity that will get us out of collision)
    u_vec = rel_vec * (disc_r - w_length)
    return u_vec, rel_vec


def orca(a: Participant, b: Participant) -> tuple[np.ndarray, np.ndarray]:
    x = b.position - a.position
    r = a.radius + a.safezone + b.radius + b.safezone
    # r *= 2
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
            # calculate angles of relative position and velocity
            positionAngle = vec2angle(x)
            velocityAngle = vec2angle(v)

            # calculate angle between velocity and position vector
            differenceAngle = angle_diff(positionAngle, velocityAngle)

            # calculate angles of cone
            sideAngle = arcsin(r, norm(x))
            rightSideAngle = positionAngle - sideAngle
            leftSideAngle = positionAngle + sideAngle

            # calculate vectors for sides of code
            rightSide = angle2vec(rightSideAngle)
            leftSide = angle2vec(leftSideAngle)

            # calculate closest points on both side
            leftPoint = closest_point_on_line(
                np.array([0, 0]), leftSide, v)
            rightPoint = closest_point_on_line(
                np.array([0, 0]), rightSide, v)

            # calculate vectors for those closest points
            left_u = leftPoint - v
            right_u = rightPoint - v

            # calculate distance to closest points on lines
            leftDist = dist(leftPoint, v)
            rightDist = dist(rightPoint, v)

            # decide, which side is closer
            if leftDist < rightDist:
                u = left_u
            else:
                u = right_u

            # calculate n (the normal vector at the closest point pointing OUTWARDS)
            n = normalize(u)
            if differenceAngle > sideAngle:
                n *= -1

            # adjust u in relation to percentage of own radius
            # perc = 1
            # if norm(a.velocity) + norm(b.velocity) != 0:
            #     perc = 1 - (norm(a.velocity) /
            #                 (norm(a.velocity) + norm(b.velocity)))
            # perc = 1 - (a.radius / r)
            return u, n


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
            # check, if feasible solution was found
            if left is None or right is None:
                return None
            new_point = closest_point_on_line(
                plane.u, plane.u + np.array([plane.n[1], -plane.n[0]]), optimal_point, left, right)
    return new_point


def intersect_halfplane_with_other_halfplanes(plane: Halfplane, other_planes: list[Halfplane]) -> tuple[np.float64, np.float64]:
    left = float('-inf')
    right = float('inf')

    direction = np.array([plane.n[1], -plane.n[0]])

    # see https://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect/565282#565282 for reference
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
