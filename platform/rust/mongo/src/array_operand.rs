use crate::{ConditionMergeError, Op, bson};

#[derive(Debug, Clone, PartialEq)]
pub(crate) struct ArrayOperand {
    values: Vec<bson::Bson>,
}

impl ArrayOperand {
    pub(crate) fn from_bson(
        operator: Op,
        operand: bson::Bson,
    ) -> Result<Self, ConditionMergeError> {
        match operand {
            bson::Bson::Array(values) => Ok(Self::from_values(values)),
            _ => Err(ConditionMergeError::ExpectedArrayOperand { operator }),
        }
    }

    pub(crate) fn intersection(self, other: Self) -> Self {
        self.values
            .into_iter()
            .filter(|value| other.contains(value))
            .collect()
    }

    pub(crate) fn union(mut self, other: Self) -> Self {
        for value in other.values {
            self.push_unique(value);
        }

        self
    }

    pub(crate) fn into_bson(self) -> bson::Bson {
        bson::Bson::Array(self.values)
    }

    fn from_values(values: Vec<bson::Bson>) -> Self {
        values
            .into_iter()
            .fold(Self { values: Vec::new() }, |mut operand, value| {
                operand.push_unique(value);
                operand
            })
    }

    fn contains(&self, value: &bson::Bson) -> bool {
        self.values.contains(value)
    }

    fn push_unique(&mut self, value: bson::Bson) {
        if !self.contains(&value) {
            self.values.push(value);
        }
    }
}

impl FromIterator<bson::Bson> for ArrayOperand {
    fn from_iter<T>(iter: T) -> Self
    where
        T: IntoIterator<Item = bson::Bson>,
    {
        Self::from_values(iter.into_iter().collect())
    }
}
