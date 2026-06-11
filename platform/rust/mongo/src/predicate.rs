use crate::{ArrayOp, Condition, Field, IntoBson, ScalarOp};

pub trait ArrayPredicate: IntoIterator
where
    Self::Item: IntoBson,
{
    type Field: Field;

    const FIELD: Self::Field;
    const OPERATOR: ArrayOp;

    fn into_predicate(self) -> Option<Predicate<Self::Field>>
    where
        Self: Sized,
    {
        let values = self
            .into_iter()
            .map(IntoBson::into_bson)
            .collect::<Vec<_>>();

        if values.is_empty() {
            return None;
        }

        Some(Predicate::new(
            Self::FIELD,
            Condition::array(Self::OPERATOR, values),
        ))
    }
}

pub trait ScalarPredicate: IntoBson {
    type Field: Field;

    const FIELD: Self::Field;
    const OPERATOR: ScalarOp;

    fn into_predicate(self) -> Predicate<Self::Field>
    where
        Self: Sized,
    {
        Predicate::new(
            Self::FIELD,
            Condition::scalar(Self::OPERATOR, self.into_bson()),
        )
    }
}

#[derive(Debug, Clone, PartialEq)]
pub struct Predicate<F> {
    field: F,
    condition: Condition,
}

impl<F> Predicate<F> {
    pub fn new(field: F, condition: Condition) -> Self {
        Self { field, condition }
    }

    pub(crate) fn into_parts(self) -> (F, Condition) {
        (self.field, self.condition)
    }
}

#[cfg(test)]
mod tests {
    use crate::{
        ArrayOp, ArrayPredicate, Field, Filter, IntoBson, ScalarOp, ScalarPredicate, bson,
    };

    #[derive(Debug, Clone, Copy)]
    enum TestField {
        ManufacturerCountry,
        Name,
        Price,
        Sku,
    }

    impl Field for TestField {
        fn path(&self) -> &'static str {
            match self {
                Self::ManufacturerCountry => "manufacturer.country",
                Self::Name => "name",
                Self::Price => "price",
                Self::Sku => "sku",
            }
        }
    }

    struct TestNames(Vec<String>);

    impl IntoIterator for TestNames {
        type Item = String;
        type IntoIter = std::vec::IntoIter<String>;

        fn into_iter(self) -> Self::IntoIter {
            self.0.into_iter()
        }
    }

    impl ArrayPredicate for TestNames {
        type Field = TestField;

        const FIELD: Self::Field = TestField::Name;
        const OPERATOR: ArrayOp = ArrayOp::In;
    }

    struct MinPrice(i64);

    impl IntoBson for MinPrice {
        fn into_bson(self) -> bson::Bson {
            bson::Bson::Int64(self.0)
        }
    }

    impl ScalarPredicate for MinPrice {
        type Field = TestField;

        const FIELD: Self::Field = TestField::Price;
        const OPERATOR: ScalarOp = ScalarOp::Gte;
    }

    struct Skus(Vec<Sku>);

    impl IntoIterator for Skus {
        type Item = Sku;
        type IntoIter = std::vec::IntoIter<Sku>;

        fn into_iter(self) -> Self::IntoIter {
            self.0.into_iter()
        }
    }

    impl ArrayPredicate for Skus {
        type Field = TestField;

        const FIELD: Self::Field = TestField::Sku;
        const OPERATOR: ArrayOp = ArrayOp::In;
    }

    struct Sku(String);

    impl IntoBson for Sku {
        fn into_bson(self) -> bson::Bson {
            bson::Bson::String(format!("sku:{}", self.0))
        }
    }

    struct ManufacturerCountry(String);

    impl IntoBson for ManufacturerCountry {
        fn into_bson(self) -> bson::Bson {
            bson::Bson::String(self.0)
        }
    }

    impl ScalarPredicate for ManufacturerCountry {
        type Field = TestField;

        const FIELD: Self::Field = TestField::ManufacturerCountry;
        const OPERATOR: ScalarOp = ScalarOp::Eq;
    }

    #[test]
    fn array_predicate_returns_none_for_empty_values() {
        assert!(TestNames(Vec::new()).into_predicate().is_none());
    }

    #[test]
    fn array_predicate_builds_filter_for_non_empty_values() {
        let predicate = TestNames(vec!["engine".to_string(), "wing".to_string()])
            .into_predicate()
            .unwrap();
        let filter = Filter::try_from_predicates([predicate])
            .unwrap()
            .into_inner();

        assert_eq!(
            filter,
            bson::doc! {
                "name": { "$in": ["engine", "wing"] },
            }
        );
    }

    #[test]
    fn array_predicate_uses_custom_into_bson_values() {
        let predicate = Skus(vec![Sku("a".to_string()), Sku("b".to_string())])
            .into_predicate()
            .unwrap();
        let filter = Filter::try_from_predicates([predicate])
            .unwrap()
            .into_inner();

        assert_eq!(
            filter,
            bson::doc! {
                "sku": { "$in": ["sku:a", "sku:b"] },
            }
        );
    }

    #[test]
    fn scalar_predicate_uses_custom_field_path() {
        let filter =
            Filter::try_from_predicates([ManufacturerCountry("US".to_string()).into_predicate()])
                .unwrap()
                .into_inner();

        assert_eq!(
            filter,
            bson::doc! {
                "manufacturer.country": { "$eq": "US" },
            }
        );
    }

    #[test]
    fn scalar_predicate_always_builds_predicate() {
        let filter = Filter::try_from_predicates([MinPrice(100).into_predicate()])
            .unwrap()
            .into_inner();

        assert_eq!(
            filter,
            bson::doc! {
                "price": { "$gte": 100_i64 },
            }
        );
    }
}
