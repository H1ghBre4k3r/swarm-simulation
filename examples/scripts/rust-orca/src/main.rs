mod math;
mod simulation;

use liborca::orca;
use simulation::Simulation;

fn main() {
    // create a new simulation and start it with our orca callback
    let mut sim = Simulation::new();
    sim.start(orca);
}
