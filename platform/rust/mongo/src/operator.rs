#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Op {
    Eq,
    In,
    All,
    Gt,
    Gte,
    Lt,
    Lte,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ArrayOp {
    In,
    All,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ScalarOp {
    Eq,
    Gt,
    Gte,
    Lt,
    Lte,
}

impl Op {
    pub fn as_str(self) -> &'static str {
        match self {
            Self::Eq => "$eq",
            Self::In => "$in",
            Self::All => "$all",
            Self::Gt => "$gt",
            Self::Gte => "$gte",
            Self::Lt => "$lt",
            Self::Lte => "$lte",
        }
    }
}

impl From<ArrayOp> for Op {
    fn from(value: ArrayOp) -> Self {
        match value {
            ArrayOp::In => Self::In,
            ArrayOp::All => Self::All,
        }
    }
}

impl From<ScalarOp> for Op {
    fn from(value: ScalarOp) -> Self {
        match value {
            ScalarOp::Eq => Self::Eq,
            ScalarOp::Gt => Self::Gt,
            ScalarOp::Gte => Self::Gte,
            ScalarOp::Lt => Self::Lt,
            ScalarOp::Lte => Self::Lte,
        }
    }
}
