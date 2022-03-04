mod math;
mod simulation;

use liborca::orca;
use simulation::Simulation;

fn main() {
    let mut sim = Simulation::new();
    sim.start(orca);
}
