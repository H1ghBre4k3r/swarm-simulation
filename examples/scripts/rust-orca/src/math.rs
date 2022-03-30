/// Normalize a given vector.
pub fn normalize(vec: &ndarray::Array1<f64>) -> ndarray::Array1<f64> {
    let len = norm(vec);
    if len == 0.0 {
        return ndarray::arr1(&[0.0, 0.0]);
    }
    return vec / len;
}

fn norm(vector: &ndarray::Array1<f64>) -> f64 {
    return vector.dot(vector).sqrt();
}
