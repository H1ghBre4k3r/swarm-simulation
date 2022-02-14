use ndarray::{arr1, Array1};

use crate::{
    halfplane::Halfplane,
    math::{
        angle2vec, angle_diff, arcsin, closest_point_on_line, cross, norm, normalize, vec2angle,
    },
    obstacle::Obstacle,
    participant::Participant,
};

const TAU: f64 = 4.0;

pub fn obstacle_collision(a: &Participant, obstacle: &Obstacle) -> (Array1<f64>, Array1<f64>) {
    let r_vec = normalize(&a.velocity) * (a.radius + a.safezone + obstacle.radius);
    let start = &obstacle.start - &(&a.position + &r_vec);
    let end = &obstacle.end - &(&a.position - &r_vec);

    let dist_vec =
        closest_point_on_line(&obstacle.start, &obstacle.end, &a.position, 0.0, 1.0) - &a.position;

    let u: Array1<f64>;
    let n: Array1<f64>;

    if norm(&dist_vec) < &a.radius + &a.safezone + &obstacle.radius {
        u = -&dist_vec - &a.velocity;
        n = normalize(&-dist_vec)
    } else {
        let start_tau = &start / TAU;
        let end_tau = &end / TAU;
        let closest = closest_point_on_line(&start_tau, &end_tau, &arr1(&[0.0, 0.0]), 0.0, 1.0);
        let w = closest;
        u = &w - &a.velocity;
        n = -normalize(&w);
    }

    return (u, n);
}

fn out_of_disk(
    disk_center: &Array1<f64>,
    disk_r: f64,
    velocity: &Array1<f64>,
) -> (Array1<f64>, Array1<f64>) {
    let rel_vec = velocity - disk_center;
    let w_length = norm(&rel_vec);
    let n = angle2vec(vec2angle(&rel_vec) + 10.0);
    let u = &n * (disk_r - w_length);
    return (u, n);
}

pub fn orca(a: &Participant, b: &Participant) -> (Array1<f64>, Array1<f64>) {
    let x = &b.position - &a.position;
    let r = &a.radius + &a.safezone + &b.radius + &b.safezone;
    let v = &a.velocity - &b.velocity;

    // check, if we are currently colliding
    if norm(&x) < r {
        return out_of_disk(&x, r, &v);
    }

    let disk_center = &x / TAU;
    let disk_r = r / TAU;
    let adjusted_disk_center = &disk_center * (1.0 - (r / norm(&x)).powi(2));
    // check, if we will collide with front disk
    if norm(&v) < norm(&adjusted_disk_center) && norm(&(&v - &adjusted_disk_center)) < r {
        return out_of_disk(&disk_center, disk_r, &v);
    }

    let position_angle = vec2angle(&x);
    let velocity_angle = vec2angle(&v);

    let difference_angle = angle_diff(position_angle, velocity_angle);

    let side_angle = arcsin(r, norm(&x));
    let right_side_angle = position_angle - side_angle;
    let left_side_angle = position_angle + side_angle;

    let right_side = angle2vec(right_side_angle);
    let left_side = angle2vec(left_side_angle);

    let left_point = closest_point_on_line(&arr1(&[0.0, 0.0]), &left_side, &v, 0.0, 1.0);
    let right_point = closest_point_on_line(&arr1(&[0.0, 0.0]), &right_side, &v, 0.0, 1.0);

    let left_u = &left_point - &v;
    let right_u = &right_point - &v;

    // let left_dist = dist(&left_point, &v);
    // let right_dist = dist(&right_point, &v);
    let left_dist = norm(&left_u);
    let right_dist = norm(&right_u);

    let u: Array1<f64>;
    if left_dist < right_dist {
        u = left_u;
    } else {
        u = right_u;
    }

    let mut n = normalize(&u);
    if difference_angle > side_angle {
        n *= -1.0
    }

    return (u, n);
}

pub fn halfplane_intersection(
    halfplanes_u: &Vec<Halfplane>,
    current_velocity: &Array1<f64>,
    optimal_point: &Array1<f64>,
) -> Option<Array1<f64>> {
    let mut halfplanes = Vec::new();
    for plane in halfplanes_u {
        halfplanes.push(Halfplane {
            u: current_velocity.clone() + plane.u.clone(),
            n: plane.n.clone(),
        });
    }
    let mut new_point = optimal_point.clone();

    for (i, plane) in halfplanes.iter().enumerate() {
        if (&new_point - &plane.u).dot(&plane.n) < 0.0 {
            let (left, right) = intersect_halfplane_with_other_halfplanes(plane, &halfplanes[..i]);
            if left.is_none() || right.is_none() {
                return None;
            }
            new_point = closest_point_on_line(
                &plane.u,
                &(&plane.u + &arr1(&[plane.n[1], -plane.n[0]])),
                optimal_point,
                left.unwrap(),
                right.unwrap(),
            )
        }
    }

    return Some(new_point);
}

fn intersect_halfplane_with_other_halfplanes(
    plane: &Halfplane,
    other_planes: &[Halfplane],
) -> (Option<f64>, Option<f64>) {
    let mut left = -f64::INFINITY;
    let mut right = f64::INFINITY;

    let direction = arr1(&[plane.n[1], -plane.n[0]]);

    for other_plane in other_planes {
        let other_dir = arr1(&[other_plane.n[1], -other_plane.n[0]]);
        let num = cross(&(&other_plane.u - &plane.u), &other_dir);
        let den = cross(&direction, &other_dir);

        if den == 0.0 {
            if num == 0.0 {
                return (None, None);
            }
            continue;
        }

        let t = &num / &den;
        if den > 0.0 {
            right = right.min(t);
        } else {
            left = left.max(t);
        }

        if left > right {
            return (None, None);
        }
    }

    return (Some(left), Some(right));
}
