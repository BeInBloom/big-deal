use platform_mongo::{ArrayOp, ArrayPredicate, Field, Filter, FilterError, IntoBson, bson};

use crate::domain::models::{
    CountryCode, CountryCodes, ListPartsQuery, Name, Names, PartCategories, PartCategory, PartId,
    PartIds, Tag, Tags,
};

pub(in crate::repo::mongo) trait IntoFilter {
    fn into_filter(self) -> Result<Filter, FilterError>;
}

#[derive(Debug, Clone, Copy)]
pub(crate) enum PartField {
    Id,
    Name,
    Category,
    ManufacturerCountry,
    Tags,
}

impl Field for PartField {
    fn path(&self) -> &'static str {
        match self {
            Self::Id => "_id",
            Self::Name => "name",
            Self::Category => "category",
            Self::ManufacturerCountry => "manufacturer.country",
            Self::Tags => "tags",
        }
    }
}

impl IntoBson for Name {
    fn into_bson(self) -> bson::Bson {
        bson::Bson::String(self.into())
    }
}

impl IntoBson for Tag {
    fn into_bson(self) -> bson::Bson {
        bson::Bson::String(self.into())
    }
}

impl IntoBson for CountryCode {
    fn into_bson(self) -> bson::Bson {
        bson::Bson::String(self.into())
    }
}

impl IntoBson for PartId {
    fn into_bson(self) -> bson::Bson {
        bson::Bson::String(self.into())
    }
}

impl IntoBson for PartCategory {
    fn into_bson(self) -> bson::Bson {
        bson::Bson::String(self.to_string())
    }
}

impl ArrayPredicate for PartIds {
    type Field = PartField;

    const FIELD: Self::Field = PartField::Id;
    const OPERATOR: ArrayOp = ArrayOp::In;
}

impl ArrayPredicate for Names {
    type Field = PartField;

    const FIELD: Self::Field = PartField::Name;
    const OPERATOR: ArrayOp = ArrayOp::In;
}

impl ArrayPredicate for PartCategories {
    type Field = PartField;

    const FIELD: Self::Field = PartField::Category;
    const OPERATOR: ArrayOp = ArrayOp::In;
}

impl ArrayPredicate for CountryCodes {
    type Field = PartField;

    const FIELD: Self::Field = PartField::ManufacturerCountry;
    const OPERATOR: ArrayOp = ArrayOp::In;
}

impl ArrayPredicate for Tags {
    type Field = PartField;

    const FIELD: Self::Field = PartField::Tags;
    const OPERATOR: ArrayOp = ArrayOp::All;
}

impl IntoFilter for ListPartsQuery {
    fn into_filter(self) -> Result<Filter, FilterError> {
        let predicates = [
            self.ids.into_predicate(),
            self.names.into_predicate(),
            self.categories.into_predicate(),
            self.manufacturer_countries.into_predicate(),
            self.tags.into_predicate(),
        ]
        .into_iter()
        .flatten();

        Filter::try_from_predicates(predicates)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    const RAW_PART_ID: &str = "11111111-1111-4111-8111-111111111111";

    #[test]
    fn empty_list_parts_query_converts_to_empty_filter() {
        let filter = ListPartsQuery::default()
            .into_filter()
            .unwrap()
            .into_inner();

        assert_eq!(filter, bson::Document::new());
    }

    #[test]
    fn list_parts_query_converts_ids_to_id_in_filter() {
        let filter = ListPartsQuery {
            ids: PartIds::try_from(vec![RAW_PART_ID.to_string()]).unwrap(),
            ..ListPartsQuery::default()
        }
        .into_filter()
        .unwrap()
        .into_inner();

        assert_eq!(
            filter,
            bson::doc! {
                "_id": { "$in": [RAW_PART_ID] },
            }
        );
    }

    #[test]
    fn list_parts_query_converts_names_to_name_in_filter() {
        let filter = ListPartsQuery {
            names: Names::from(vec!["Main engine".to_string()]),
            ..ListPartsQuery::default()
        }
        .into_filter()
        .unwrap()
        .into_inner();

        assert_eq!(
            filter,
            bson::doc! {
                "name": { "$in": ["Main engine"] },
            }
        );
    }

    #[test]
    fn list_parts_query_converts_categories_to_category_in_filter() {
        let filter = ListPartsQuery {
            categories: [PartCategory::Engine].into_iter().collect(),
            ..ListPartsQuery::default()
        }
        .into_filter()
        .unwrap()
        .into_inner();

        assert_eq!(
            filter,
            bson::doc! {
                "category": { "$in": ["ENGINE"] },
            }
        );
    }

    #[test]
    fn list_parts_query_converts_manufacturer_countries_to_country_in_filter() {
        let filter = ListPartsQuery {
            manufacturer_countries: CountryCodes::from(vec!["US".to_string()]),
            ..ListPartsQuery::default()
        }
        .into_filter()
        .unwrap()
        .into_inner();

        assert_eq!(
            filter,
            bson::doc! {
                "manufacturer.country": { "$in": ["US"] },
            }
        );
    }

    #[test]
    fn list_parts_query_converts_tags_to_tags_all_filter() {
        let filter = ListPartsQuery {
            tags: Tags::from(vec!["critical".to_string(), "engine".to_string()]),
            ..ListPartsQuery::default()
        }
        .into_filter()
        .unwrap()
        .into_inner();

        assert_field_array_contains(
            &filter,
            "tags",
            "$all",
            &[
                bson::Bson::String("critical".to_string()),
                bson::Bson::String("engine".to_string()),
            ],
        );
    }

    #[test]
    fn list_parts_query_converts_all_filters() {
        let filter = ListPartsQuery {
            ids: PartIds::try_from(vec![RAW_PART_ID.to_string()]).unwrap(),
            names: Names::from(vec!["Main engine".to_string()]),
            categories: [PartCategory::Engine].into_iter().collect(),
            manufacturer_countries: CountryCodes::from(vec!["US".to_string()]),
            tags: Tags::from(vec!["critical".to_string(), "engine".to_string()]),
        }
        .into_filter()
        .unwrap()
        .into_inner();

        assert_field_array_contains(
            &filter,
            "_id",
            "$in",
            &[bson::Bson::String(RAW_PART_ID.to_string())],
        );
        assert_field_array_contains(
            &filter,
            "name",
            "$in",
            &[bson::Bson::String("Main engine".to_string())],
        );
        assert_field_array_contains(
            &filter,
            "category",
            "$in",
            &[bson::Bson::String("ENGINE".to_string())],
        );
        assert_field_array_contains(
            &filter,
            "manufacturer.country",
            "$in",
            &[bson::Bson::String("US".to_string())],
        );
        assert_field_array_contains(
            &filter,
            "tags",
            "$all",
            &[
                bson::Bson::String("critical".to_string()),
                bson::Bson::String("engine".to_string()),
            ],
        );
    }

    fn assert_field_array_contains(
        filter: &bson::Document,
        field: &str,
        operator: &str,
        expected_values: &[bson::Bson],
    ) {
        let values = filter
            .get_document(field)
            .unwrap()
            .get_array(operator)
            .unwrap();

        assert_eq!(values.len(), expected_values.len());

        for expected_value in expected_values {
            assert!(values.contains(expected_value));
        }
    }
}
