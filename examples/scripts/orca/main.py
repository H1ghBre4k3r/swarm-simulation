#!/usr/bin/env python3 -u

import sys
from typing import List

from halfplane import Halfplane
from mathutils import dist, norm, normalize
from obstacle import Obstacle
from orca import halfplane_intersection, obstacle_collision, orca
from participant import Participant
from simulation import Simulation


def is_static(p: Participant):
    return norm(p.velocity) < 1e-3


def main():

    # warnings.filterwarnings("error")

    simulation = Simulation()

    def callback(we: Participant, participants: List[Participant], obstacles: List[Obstacle]) -> None:
        # calculate halfplanes for each participant
        obstacle_planes = []
        halfplanes = []
        for i, p in enumerate(participants):
            in_obstacle = False
            if is_static(p):
                p.is_obstacle = True
                # TODO lome: maybe move out of if-statement
                for other in participants[i + 1:]:
                    # if they are too
                    if is_static(other) and dist(p.position, other.position) < p.radius + other.radius + p.safezone + other.safezone + (we.safezone + we.radius) * 2:
                        # if dist(p.position, other.position) < p.radius + other.radius + p.safezone + other.safezone + (we.safezone + we.radius) * 2:
                        obstacle_coords = {"start": {"x": p.position[0], "y": p.position[1]}, "end": {
                            "x": other.position[0], "y": other.position[1]}}
                        u, n = obstacle_collision(we, Obstacle(obstacle_coords, max(
                            p.radius + p.safezone, other.radius + other.safezone)))
                        obstacle_planes.append(Halfplane(u, n))
                        other.is_obstacle = True
                        in_obstacle = True

            # if used in obstacle or is static, we calculate half-planes like it is an obstacle
            if not in_obstacle and p.is_obstacle:
                obstacle_coords = {"start": {"x": p.position[0], "y": p.position[1]}, "end": {
                    "x": p.position[0], "y": p.position[1]}}
                u, n = obstacle_collision(we, Obstacle(
                    obstacle_coords, p.radius + p.safezone))
                obstacle_planes.append(Halfplane(u, n))
            # if it is not used in an obstacle, we calculate half-planes like it is a participant
            if not p.is_obstacle:
                u, n = orca(we, p)
                halfplanes.append(Halfplane(u, n))

        for o in obstacles:
            u, n = obstacle_collision(we, o)
            obstacle_planes.append(Halfplane(u, n))
        # try to find new velocity
        new_vel = None
        while new_vel is None:
            new_vel = halfplane_intersection(
                halfplanes + obstacle_planes, we.velocity, we.velocity)
            # move halfplanes outward at equal speed
            new_halfplanes = []
            for l in halfplanes:
                new_halfplanes.append(Halfplane(l.u - l.n * 0.0001, l.n))
            halfplanes = new_halfplanes
        if norm(new_vel) > we.vmax:
            new_vel = normalize(new_vel) * we.vmax
        return new_vel

    simulation.start(callback)


if __name__ == "__main__":
    main()
