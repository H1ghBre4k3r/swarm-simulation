use std::io;

use orca_rs::{ndarray::arr1, obstacle::Obstacle, participant::Participant};
use serde_json::{json, Value};

pub struct Simulation {
    we: Participant,
    tau: f64,
}

#[cfg(feature = "confidence")]
const CONF: f64 = 2.0;

#[cfg(not(feature = "confidence"))]
const CONF: f64 = 0.0;

impl Simulation {
    pub fn new() -> Simulation {
        let mut buffer = String::new();
        io::stdin()
            .read_line(&mut buffer)
            .expect("Error reading stdio");

        // Read setup information from stdin
        let setup: Value =
            serde_json::from_str(buffer.as_str()).expect("Error decoding setup message");
        let position = arr1(&[
            setup["position"]["x"].as_f64().unwrap(),
            setup["position"]["y"].as_f64().unwrap(),
        ]);
        let target = arr1(&[
            setup["target"]["x"].as_f64().unwrap(),
            setup["target"]["y"].as_f64().unwrap(),
        ]);
        let vmax = setup["vmax"].as_f64().unwrap();
        let tau = setup["tau"].as_f64().unwrap();
        let radius = setup["radius"].as_f64().unwrap();
        let we = Participant::new(position, arr1(&[0.0, 0.0]), radius, 0.0)
            .with_inner_state(vmax, target);

        Simulation { we, tau }
    }

    /// Start this simulation with a callback to call during each tick
    pub fn start(&mut self) {
        // continuously perform ticks
        loop {
            let mut buffer = String::new();
            // read information from stdin
            io::stdin()
                .read_line(&mut buffer)
                .expect("Error reading stdin");
            let inp: Value = serde_json::from_str(buffer.trim())
                .expect("Error decoding message from simulation.");

            // update our information
            let position = arr1(&[
                inp["position"]["x"].as_f64().unwrap(),
                inp["position"]["y"].as_f64().unwrap(),
            ]);
            self.we.update_position(&position);
            self.we.confidence = inp["stddev"].as_f64().unwrap() * CONF + 0.001;

            // convert the information about all other participants
            let mut participants: Vec<Participant> = Vec::new();
            for p in inp["participants"].as_array().unwrap() {
                let position = arr1(&[
                    p["position"]["x"].as_f64().unwrap(),
                    p["position"]["y"].as_f64().unwrap(),
                ]);
                let velocity = arr1(&[
                    p["velocity"]["x"].as_f64().unwrap(),
                    p["velocity"]["y"].as_f64().unwrap(),
                ]);
                let radius = p["radius"].as_f64().unwrap();
                let confidence = p["stddev"].as_f64().unwrap() * CONF + 0.001;
                participants.push(Participant::new(position, velocity, radius, confidence))
            }

            // get information about static obstacles
            let mut obstacles: Vec<Obstacle> = Vec::new();
            for o in inp["obstacles"].as_array().unwrap() {
                obstacles.push(Obstacle {
                    start: arr1(&[
                        o["start"]["x"].as_f64().unwrap(),
                        o["start"]["y"].as_f64().unwrap(),
                    ]),
                    end: arr1(&[
                        o["end"]["x"].as_f64().unwrap(),
                        o["end"]["y"].as_f64().unwrap(),
                    ]),
                    radius: 0.0,
                });
            }
            // get new velocity from callback & send it via stdout
            let vel = self.we.orca(&mut participants, &obstacles, self.tau);
            let val = json!({
                "action": "move",
                "payload": {
                    "x": vel[0],
                    "y": vel[1]
                }
            });
            println!("{}", val);
        }
    }
}
