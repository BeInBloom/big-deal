use crate::{ArrayOp, ConditionMergeError, Op, ScalarOp, array_operand::ArrayOperand, bson};

#[derive(Debug, Clone, PartialEq)]
pub struct Condition {
    operator: Op,
    operand: bson::Bson,
}

impl Condition {
    pub fn scalar(operator: ScalarOp, operand: bson::Bson) -> Self {
        Self::new(operator.into(), operand)
    }

    pub fn array(operator: ArrayOp, values: Vec<bson::Bson>) -> Self {
        Self::new(operator.into(), bson::Bson::Array(values))
    }

    pub(crate) fn new(operator: Op, operand: bson::Bson) -> Self {
        Self { operator, operand }
    }

    pub(crate) fn operator(&self) -> Op {
        self.operator
    }

    pub(crate) fn into_entry(self) -> (String, bson::Bson) {
        (self.operator.as_str().to_owned(), self.operand)
    }

    pub(crate) fn merge(&mut self, incoming: Self) -> Result<(), ConditionMergeError> {
        debug_assert_eq!(self.operator, incoming.operator);

        match self.operator {
            Op::In => self.intersect_array(incoming),
            Op::All => self.union_array(incoming),
            _ if self.operand == incoming.operand => Ok(()),
            _ => Err(ConditionMergeError::DuplicateOperator {
                operator: self.operator,
            }),
        }
    }

    fn intersect_array(&mut self, incoming: Self) -> Result<(), ConditionMergeError> {
        let existing_values =
            ArrayOperand::from_bson(self.operator, std::mem::take(&mut self.operand))?;
        let incoming_values = ArrayOperand::from_bson(self.operator, incoming.operand)?;

        self.operand = existing_values.intersection(incoming_values).into_bson();

        Ok(())
    }

    fn union_array(&mut self, incoming: Self) -> Result<(), ConditionMergeError> {
        let existing_values =
            ArrayOperand::from_bson(self.operator, std::mem::take(&mut self.operand))?;
        let incoming_values = ArrayOperand::from_bson(self.operator, incoming.operand)?;

        self.operand = existing_values.union(incoming_values).into_bson();

        Ok(())
    }
}
