class Halfplane(object):
    """A half-plane defined by a line in Hessian normal form."""

    def __init__(self, u, n):
        super(Halfplane, self).__init__()
        self.u = u
        self.n = n

    def __repr__(self):
        return f"Halfplane({self.u}, {self.n})"
