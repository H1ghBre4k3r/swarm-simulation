mod halfplane;
mod math;
mod obstacle;
mod orca;
mod participant;

use halfplane::Halfplane;
use math::{dist, max, norm, normalize};
use ndarray::{arr1, Array1};
use obstacle::Obstacle;
use orca::{halfplane_intersection, obstacle_collision, orca};
use participant::Participant;
use serde_json::{json, Value};
use std::io;

const CONF: f64 = 3.0;

fn main() {
    let mut buffer = String::new();
    io::stdin()
        .read_line(&mut buffer)
        .expect("Error reading stdio");

    let setup: Value = serde_json::from_str(buffer.as_str()).expect("Error decoding setup message");
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

    let mut we = Participant {
        velocity: normalize(&(&target - &position)) * vmax,
        position,
        target,
        vmax,
        radius: setup["radius"].as_f64().unwrap(),
        confidence: 0.0,
        in_obstacle: false,
    };

    loop {
        buffer = String::new();
        io::stdin()
            .read_line(&mut buffer)
            .expect("Error reading stdio");
        let inp: Value =
            serde_json::from_str(buffer.trim()).expect("Error decoding message from simulation.");

        let position = arr1(&[
            inp["position"]["x"].as_f64().unwrap(),
            inp["position"]["y"].as_f64().unwrap(),
        ]);
        we.update_position(&position);
        we.confidence = inp["stddev"].as_f64().unwrap() * CONF + 0.001;

        let mut participants: Vec<Participant> = Vec::new();
        for p in inp["participants"].as_array().unwrap() {
            participants.push(Participant {
                position: arr1(&[
                    p["position"]["x"].as_f64().unwrap(),
                    p["position"]["y"].as_f64().unwrap(),
                ]),
                velocity: arr1(&[
                    p["velocity"]["x"].as_f64().unwrap(),
                    p["velocity"]["y"].as_f64().unwrap(),
                ]),
                radius: p["radius"].as_f64().unwrap(),
                confidence: p["stddev"].as_f64().unwrap() * CONF + 0.001,
                target: arr1(&[0.0, 0.0]),
                vmax: 0.0,
                in_obstacle: false,
            });
        }

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

        let vel = callback(&we, &mut participants, &obstacles, tau);
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

fn is_static(p: &Participant) -> bool {
    return norm(&p.velocity) == 0.0;
}

/// Callback for "orca"
fn callback(
    we: &Participant,
    participants: &mut [Participant],
    obstacles: &[Obstacle],
    tau: f64,
) -> Array1<f64> {
    let (mut halfplanes, obstacle_planes) = generate_halfplanes(we, participants, obstacles, tau);

    let mut new_vel: Option<Array1<f64>> = None;
    while new_vel.is_none() {
        // combine halfplanes
        let mut hp = halfplanes.to_vec();
        let mut op = obstacle_planes.to_vec();
        hp.append(&mut op);
        // calculate new velocity
        new_vel = halfplane_intersection(&hp, &we.velocity, &we.velocity);
        // adjust halfplanes (move them outward)
        let mut new_halfplanes: Vec<Halfplane> = Vec::new();
        for l in halfplanes {
            new_halfplanes.push(Halfplane {
                u: &l.u - &(&l.n * 0.0001),
                n: l.n.clone(),
            });
        }
        halfplanes = new_halfplanes
    }
    let mut vel = new_vel.unwrap();
    if norm(&vel) > we.vmax {
        vel = normalize(&vel) * we.vmax;
    }
    return vel;
}

/// Generate the halfplanes for other participants and static obstacles.
///
/// Participants, which are too close to each other will be combined into a static obstacle aswell.
fn generate_halfplanes(
    we: &Participant,
    participants: &mut [Participant],
    obstacles: &[Obstacle],
    tau: f64,
) -> (Vec<Halfplane>, Vec<Halfplane>) {
    let mut obstacle_planes = Vec::new();
    let mut halfplanes = Vec::new();
    // some cloning, since rust does not like multiple mutable borrows
    let parts = participants.to_vec().clone();
    for (i, p) in parts.iter().enumerate() {
        let mut in_obstacle = false;
        if is_static(p) {
            for (j, other) in participants.iter_mut().enumerate() {
                // check, if "p" is already part of an obstacle
                if i == j {
                    in_obstacle = other.in_obstacle
                }
                // ignore all participants before "other"
                if j < i + 1 {
                    continue;
                }
                // check, if "other" is static _and_ too close
                if is_static(other)
                    && dist(&p.position, &other.position)
                        < p.radius
                            + other.radius
                            + p.confidence
                            + other.confidence
                            + (we.confidence + we.radius) * 2.0
                {
                    let (u, n) = obstacle_collision(
                        we,
                        &Obstacle {
                            start: p.position.clone(),
                            end: other.position.clone(),
                            radius: max(p.radius + p.confidence, other.radius + other.confidence),
                        },
                        tau,
                    );
                    obstacle_planes.push(Halfplane { u, n });
                    in_obstacle = true;
                    other.in_obstacle();
                }
            }
            // if this participant is not part of a static obstacle, we treat it like a "normal" participant
            if !in_obstacle {
                let (u, n) = orca(we, p, tau);
                halfplanes.push(Halfplane { u, n });
            }
        } else {
            let (u, n) = orca(we, p, tau);
            halfplanes.push(Halfplane { u, n });
        }
    }

    for o in obstacles {
        let (u, n) = obstacle_collision(we, o, tau);
        obstacle_planes.push(Halfplane { u, n });
    }

    return (halfplanes, obstacle_planes);
}
