#!/usr/bin/env python3 -u

import sys

from halfplane import Halfplane
from mathutils import norm, normalize
from obstacle import Obstacle
from orca import halfplane_intersection, obstacle_collision, orca
from participant import Participant
from simulation import Simulation


def main():

    # warnings.filterwarnings("error")

    simulation = Simulation()

    def callback(we: Participant, participants: list[Participant], obstacles: list[Obstacle]) -> None:
        # calculate halfplanes for each participant
        halfplanes = []
        for p in participants:
            u, n = orca(we, p)
            halfplanes.append(Halfplane(u, n))

        obstacle_planes = []
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
