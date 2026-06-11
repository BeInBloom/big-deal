use crate::{Condition, FilterError, bson};

#[derive(Debug, Clone, PartialEq)]
pub(crate) struct FieldFilter {
    path: &'static str,
    conditions: Vec<Condition>,
}

#[derive(Debug, Default, Clone, PartialEq)]
pub(crate) struct FieldFilters {
    entries: Vec<FieldFilter>,
}

impl FieldFilter {
    pub(crate) fn new(path: &'static str, condition: Condition) -> Self {
        Self {
            path,
            conditions: vec![condition],
        }
    }

    pub(crate) fn merge(&mut self, condition: Condition) -> Result<(), FilterError> {
        match self
            .conditions
            .iter_mut()
            .find(|existing| existing.operator() == condition.operator())
        {
            Some(existing) => existing
                .merge(condition)
                .map_err(|source| FilterError::FieldMerge {
                    path: self.path,
                    source,
                }),
            None => {
                self.conditions.push(condition);
                Ok(())
            }
        }
    }

    pub(crate) fn path(&self) -> &'static str {
        self.path
    }

    pub(crate) fn into_entry(self) -> (String, bson::Bson) {
        (
            self.path.to_owned(),
            bson::Bson::Document(
                self.conditions
                    .into_iter()
                    .map(Condition::into_entry)
                    .collect(),
            ),
        )
    }
}

impl FieldFilters {
    pub(crate) fn merge(
        &mut self,
        path: &'static str,
        condition: Condition,
    ) -> Result<(), FilterError> {
        match self.entries.iter_mut().find(|filter| filter.path() == path) {
            Some(filter) => filter.merge(condition),
            None => {
                self.entries.push(FieldFilter::new(path, condition));
                Ok(())
            }
        }
    }

    pub(crate) fn into_document(self) -> bson::Document {
        self.entries
            .into_iter()
            .map(FieldFilter::into_entry)
            .collect()
    }
}
