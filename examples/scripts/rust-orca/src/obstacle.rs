use ndarray::Array1;

pub struct Obstacle {
    pub start: Array1<f64>,
    pub end: Array1<f64>,
    pub radius: f64,
}
