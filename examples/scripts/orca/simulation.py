import json
import sys
from threading import Thread

import numpy as np
from mathutils import norm, normalize
from participant import Participant
from util import log


class Simulation(Thread):
    """
    Class for communicating with the simulation. It handles serialization and deserialization of the data.
    """

    running = False

    def __init__(self):
        """
        Constructor

        Paremeters:
            fps: The frames per second of the simulation. This does not refer to the simulation ticks, but rather the "lookup range" for each participant.
        """
        super().__init__()
        self.__setup()

    def __setup(self):
        # gather setup information
        setup = json.loads(sys.stdin.readline())
        position = np.array([setup["position"]["x"], setup["position"]["y"]])
        radius = setup["radius"]
        target = np.array([setup["target"]["x"], setup["target"]["y"]])
        self.fps = setup["FPS"]
        vmax = setup["vmax"] * self.fps
        # create participant for this unit
        self.we = Participant(position, normalize(
            target - position) * vmax, radius, vmax, target)

    def start(self, cb):
        """
        Starts the simulation and registers the callback function, which is called every tick.
        """
        if self.running:
            # don't be rude and start the simulation twice
            raise RuntimeError("Simulation already running")
        self.cb = cb
        self.running = True
        super().start()

    def run(self):
        self.__loop(self.cb)

    def __loop(self, cb):
        """
        The main loop of the simulation
        """
        while self.running:
            # get current tick information
            inp = json.loads(sys.stdin.readline())
            position = np.array([inp["position"]["x"], inp["position"]["y"]])
            self.we.update_position(position)

            # read information about all other participants
            participansts = []
            for p in inp["participants"]:
                participant = Participant(np.array([p["position"]["x"], p["position"]["y"]]), np.array(
                    [p["velocity"]["x"], p["velocity"]["y"]]) * self.fps, p["radius"])
                participansts.append(participant)

            # call the callback function
            velocity = cb(self.we, participansts)
            velocity /= self.fps
            # send the new velocity to the simulation
            val = {
                "action": "move",
                "payload": {
                    "x": velocity[0],
                    "y": velocity[1],
                },
                "debug": inp["participants"]
            }
            print(json.dumps(val))

    def log(self, msg):
        """
        Log a message to the simulation.
        """
        log("log", msg)

    def stop(self):
        """
        Stop the simulation.
        """
        # set flag and send stop message
        self.running = False
        val = {
            "action": "stop",
            "payload": {}
        }
        print(json.dumps(val))
