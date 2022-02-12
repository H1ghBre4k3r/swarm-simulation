use halfplane::Halfplane;
use math::{dist, max, norm, normalize};
use ndarray::{arr1, Array1};
use obstacle::Obstacle;
use orca::{halfplane_intersection, obstacle_collision, orca};
use participant::Participant;
mod halfplane;
mod math;
mod obstacle;
mod orca;
mod participant;

fn main() {
    let we = Participant {
        position: arr1(&[0.0, 0.0]),
        radius: 0.005,
        safezone: 0.0,
        target: arr1(&[0.0, 0.0]),
        velocity: arr1(&[0.0, 0.0]),
        vmax: 0.015,
    };

    callback(&we, &mut [], &[]);
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
                n: l.n,
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
