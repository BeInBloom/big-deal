use crate::{Field, FilterError, Predicate, bson, field_filter::FieldFilters};

#[derive(Debug, Clone, PartialEq)]
pub struct Filter {
    inner: bson::Document,
}

impl Filter {
    pub fn into_inner(self) -> bson::Document {
        self.inner
    }

    pub fn try_from_predicates<F, I>(predicates: I) -> Result<Self, FilterError>
    where
        F: Field,
        I: IntoIterator<Item = Predicate<F>>,
    {
        let mut filters = FieldFilters::default();

        for predicate in predicates {
            let (field, condition) = predicate.into_parts();
            filters.merge(field.path(), condition)?;
        }

        Ok(Self {
            inner: filters.into_document(),
        })
    }
}

#[cfg(test)]
mod tests {
    use crate::{
        ArrayOp, Condition, ConditionMergeError, Field, Filter, FilterError, Op, Predicate,
        ScalarOp, bson,
    };

    #[derive(Debug, Clone, Copy)]
    enum TestField {
        Name,
        Price,
        Tags,
    }

    impl Field for TestField {
        fn path(&self) -> &'static str {
            match self {
                Self::Name => "name",
                Self::Price => "price",
                Self::Tags => "tags",
            }
        }
    }

    #[test]
    fn builds_filter_for_distinct_fields() {
        let filter = Filter::try_from_predicates([
            array_predicate(TestField::Name, ArrayOp::In, vec![string("engine")]),
            array_predicate(TestField::Tags, ArrayOp::All, vec![string("critical")]),
        ])
        .unwrap()
        .into_inner();

        assert_eq!(
            filter,
            bson::doc! {
                "name": { "$in": ["engine"] },
                "tags": { "$all": ["critical"] },
            }
        );
    }

    #[test]
    fn merges_different_operators_for_same_field() {
        let filter = Filter::try_from_predicates([
            scalar_predicate(TestField::Price, ScalarOp::Gte, bson::Bson::Int64(0)),
            scalar_predicate(TestField::Price, ScalarOp::Lte, bson::Bson::Int64(10)),
        ])
        .unwrap()
        .into_inner();

        assert_eq!(
            filter,
            bson::doc! {
                "price": { "$gte": 0_i64, "$lte": 10_i64 },
            }
        );
    }

    #[test]
    fn intersects_duplicate_in_operator_for_same_field() {
        let filter = Filter::try_from_predicates([
            array_predicate(
                TestField::Name,
                ArrayOp::In,
                vec![string("engine"), string("wing")],
            ),
            array_predicate(
                TestField::Name,
                ArrayOp::In,
                vec![string("wing"), string("fuel")],
            ),
        ])
        .unwrap()
        .into_inner();

        assert_eq!(
            filter,
            bson::doc! {
                "name": { "$in": ["wing"] },
            }
        );
    }

    #[test]
    fn unions_duplicate_all_operator_for_same_field() {
        let filter = Filter::try_from_predicates([
            array_predicate(
                TestField::Tags,
                ArrayOp::All,
                vec![string("critical"), string("engine")],
            ),
            array_predicate(
                TestField::Tags,
                ArrayOp::All,
                vec![string("engine"), string("reusable")],
            ),
        ])
        .unwrap()
        .into_inner();

        assert_eq!(
            filter,
            bson::doc! {
                "tags": { "$all": ["critical", "engine", "reusable"] },
            }
        );
    }

    #[test]
    fn keeps_duplicate_scalar_operator_when_operand_is_same() {
        let filter = Filter::try_from_predicates([
            scalar_predicate(TestField::Price, ScalarOp::Gte, bson::Bson::Int64(0)),
            scalar_predicate(TestField::Price, ScalarOp::Gte, bson::Bson::Int64(0)),
        ])
        .unwrap()
        .into_inner();

        assert_eq!(
            filter,
            bson::doc! {
                "price": { "$gte": 0_i64 },
            }
        );
    }

    #[test]
    fn rejects_duplicate_scalar_operator_when_operand_differs() {
        let error = Filter::try_from_predicates([
            scalar_predicate(TestField::Price, ScalarOp::Gte, bson::Bson::Int64(0)),
            scalar_predicate(TestField::Price, ScalarOp::Gte, bson::Bson::Int64(1)),
        ])
        .unwrap_err();

        assert_eq!(
            error,
            FilterError::FieldMerge {
                path: "price",
                source: ConditionMergeError::DuplicateOperator { operator: Op::Gte },
            }
        );
    }

    #[test]
    fn rejects_non_array_operand_when_merging_in_operator() {
        let error = Filter::try_from_predicates([
            raw_predicate(TestField::Name, Op::In, string("engine")),
            array_predicate(TestField::Name, ArrayOp::In, vec![string("engine")]),
        ])
        .unwrap_err();

        assert_eq!(
            error,
            FilterError::FieldMerge {
                path: "name",
                source: ConditionMergeError::ExpectedArrayOperand { operator: Op::In },
            }
        );
    }

    fn array_predicate(
        field: TestField,
        operator: ArrayOp,
        values: Vec<bson::Bson>,
    ) -> Predicate<TestField> {
        Predicate::new(field, Condition::array(operator, values))
    }

    fn scalar_predicate(
        field: TestField,
        operator: ScalarOp,
        operand: bson::Bson,
    ) -> Predicate<TestField> {
        Predicate::new(field, Condition::scalar(operator, operand))
    }

    fn raw_predicate(field: TestField, operator: Op, operand: bson::Bson) -> Predicate<TestField> {
        Predicate::new(field, Condition::new(operator, operand))
    }

    fn string(value: &str) -> bson::Bson {
        bson::Bson::String(value.to_owned())
    }
}
