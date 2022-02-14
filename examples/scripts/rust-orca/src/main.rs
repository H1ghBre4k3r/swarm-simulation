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

fn main() {
    let mut buffer = String::new();
    io::stdin()
        .read_line(&mut buffer)
        .expect("Error reading stdio");
    dbg!(&buffer);

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

    let mut we = Participant {
        velocity: normalize(&(&target - &position)) * vmax,
        position,
        target,
        vmax,
        radius: setup["radius"].as_f64().unwrap(),
        safezone: setup["safezone"].as_f64().unwrap(),
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
                safezone: p["safezone"].as_f64().unwrap(),
                target: arr1(&[0.0, 0.0]),
                vmax: 0.0,
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

        let vel = callback(&we, &mut participants, &obstacles);
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
    return norm(&p.velocity) < 1e-3;
}

fn callback(
    we: &Participant,
    participants: &mut [Participant],
    obstacles: &[Obstacle],
) -> Array1<f64> {
    let (mut halfplanes, obstacle_planes) = generate_halfplanes(we, participants, obstacles);

    let mut new_vel: Option<Array1<f64>> = None;
    while new_vel.is_none() {
        let mut hp = halfplanes.to_vec();
        let mut op = obstacle_planes.to_vec();
        hp.append(&mut op);
        new_vel = halfplane_intersection(&hp, &we.velocity, &we.velocity);
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

fn generate_halfplanes(
    we: &Participant,
    participants: &[Participant],
    obstacles: &[Obstacle],
) -> (Vec<Halfplane>, Vec<Halfplane>) {
    let mut obstacle_planes = Vec::new();
    let mut halfplanes = Vec::new();
    for (i, p) in participants.iter().enumerate() {
        let mut in_obstacle = false;
        if is_static(p) {
            for (j, other) in participants.iter().enumerate() {
                if j < i + 1 {
                    continue;
                }
                if is_static(other)
                    && dist(&p.position, &other.position)
                        < p.radius
                            + other.radius
                            + p.safezone
                            + other.safezone
                            + (we.safezone + we.radius) * 2.0
                {
                    let (u, n) = obstacle_collision(
                        we,
                        &Obstacle {
                            start: p.position.clone(),
                            end: other.position.clone(),
                            radius: max(p.radius + p.safezone, other.radius + other.safezone),
                        },
                    );
                    obstacle_planes.push(Halfplane { u, n });
                    in_obstacle = true;
                }
            }
            if !in_obstacle {
                let (u, n) = obstacle_collision(
                    we,
                    &Obstacle {
                        start: p.position.clone(),
                        end: p.position.clone(),
                        radius: p.radius,
                    },
                );
                obstacle_planes.push(Halfplane { u, n });
            }
        } else {
            let (u, n) = orca(we, p);
            halfplanes.push(Halfplane { u, n });
        }
    }

    for o in obstacles {
        let (u, n) = obstacle_collision(we, o);
        obstacle_planes.push(Halfplane { u, n });
    }

    return (halfplanes, obstacle_planes);
}
