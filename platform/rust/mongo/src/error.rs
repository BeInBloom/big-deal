use crate::Op;
use thiserror::Error;

#[derive(Debug, Clone, Copy, PartialEq, Eq, Error)]
pub enum FilterError {
    #[error("failed to merge mongo filter conditions for `{path}`")]
    FieldMerge {
        path: &'static str,
        #[source]
        source: ConditionMergeError,
    },
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Error)]
pub enum ConditionMergeError {
    #[error("expected array operand for `{}`", .operator.as_str())]
    ExpectedArrayOperand { operator: Op },
    #[error("duplicate mongo operator `{}`", .operator.as_str())]
    DuplicateOperator { operator: Op },
}
