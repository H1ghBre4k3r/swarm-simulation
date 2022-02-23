import json
import sys
from threading import Thread

import numpy as np
from mathutils import normalize
from obstacle import Obstacle
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
        safezone = setup["safezone"]
        target = np.array([setup["target"]["x"], setup["target"]["y"]])
        vmax = setup["vmax"]
        # create participant for this unit
        self.we = Participant(position, normalize(
            target - position) * vmax, radius, safezone, vmax, target)

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

    def __loop(self, callback):
        """
        The main loop of the simulation
        """
        while self.running:
            # get current tick information
            inp = json.loads(sys.stdin.readline())
            position = np.array([inp["position"]["x"], inp["position"]["y"]])
            self.we.update_position(position)

            # read information about all other participants
            participants = []
            for p in inp["participants"]:
                participant = Participant(np.array([p["position"]["x"], p["position"]["y"]]), np.array(
                    [p["velocity"]["x"], p["velocity"]["y"]]), p["radius"], p["safezone"])
                participants.append(participant)
            obstacles = []
            for o in inp["obstacles"]:
                obstacles.append(Obstacle(o))

            # call the callback function
            velocity = callback(self.we, participants, obstacles)
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