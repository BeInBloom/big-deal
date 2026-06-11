use mongodb::bson;

use crate::domain::models::{
    CountryCode, CountryCodes, ListPartsQuery, Name, Names, PartCategories, PartCategory, PartId,
    PartIds, Tag, Tags,
};

pub(in crate::repo::mongo) trait IntoMongoFilter {
    fn into_mongo_filter(self) -> PartFilterDocument;
}

pub(in crate::repo::mongo) trait IntoMongoPredicate {
    fn into_mongo_predicate(self) -> Option<MongoPredicate>;
}

pub(in crate::repo::mongo) trait IntoMongoValue {
    fn into_mongo_value(self) -> bson::Bson;
}

trait PartFilterSpec: IntoIterator
where
    Self::Item: IntoMongoValue,
{
    const FIELD: PartField;
    const OPERATOR: MongoOperator;
}

pub(in crate::repo::mongo) struct PartFilterDocument {
    inner: bson::Document,
}

pub(in crate::repo::mongo) struct MongoPredicate {
    field: PartField,
    filter: bson::Bson,
}

#[derive(Debug, Clone, Copy)]
pub(in crate::repo::mongo) enum PartField {
    Id,
    Name,
    Category,
    ManufacturerCountry,
    Tags,
}

#[derive(Debug, Clone, Copy)]
enum MongoOperator {
    In,
    All,
}

impl PartFilterDocument {
    pub(in crate::repo::mongo) fn into_inner(self) -> bson::Document {
        self.inner
    }
}

impl PartField {
    fn path(self) -> &'static str {
        match self {
            Self::Id => "_id",
            Self::Name => "name",
            Self::Category => "category",
            Self::ManufacturerCountry => "manufacturer.country",
            Self::Tags => "tags",
        }
    }
}

impl MongoOperator {
    fn as_str(self) -> &'static str {
        match self {
            Self::In => "$in",
            Self::All => "$all",
        }
    }

    fn into_filter(self, values: Vec<bson::Bson>) -> bson::Bson {
        bson::doc! { self.as_str(): values }.into()
    }
}

impl FromIterator<MongoPredicate> for PartFilterDocument {
    fn from_iter<T>(predicates: T) -> Self
    where
        T: IntoIterator<Item = MongoPredicate>,
    {
        let inner = predicates
            .into_iter()
            .map(|predicate| (predicate.field.path().to_owned(), predicate.filter))
            .collect();

        Self { inner }
    }
}

impl IntoMongoValue for String {
    fn into_mongo_value(self) -> bson::Bson {
        bson::Bson::String(self)
    }
}

impl IntoMongoValue for Name {
    fn into_mongo_value(self) -> bson::Bson {
        bson::Bson::String(self.into())
    }
}

impl IntoMongoValue for Tag {
    fn into_mongo_value(self) -> bson::Bson {
        bson::Bson::String(self.into())
    }
}

impl IntoMongoValue for CountryCode {
    fn into_mongo_value(self) -> bson::Bson {
        bson::Bson::String(self.into())
    }
}

impl IntoMongoValue for PartId {
    fn into_mongo_value(self) -> bson::Bson {
        bson::Bson::String(self.into())
    }
}

impl IntoMongoValue for PartCategory {
    fn into_mongo_value(self) -> bson::Bson {
        bson::Bson::String(self.to_string())
    }
}

impl PartFilterSpec for PartIds {
    const FIELD: PartField = PartField::Id;
    const OPERATOR: MongoOperator = MongoOperator::In;
}

impl PartFilterSpec for Names {
    const FIELD: PartField = PartField::Name;
    const OPERATOR: MongoOperator = MongoOperator::In;
}

impl PartFilterSpec for PartCategories {
    const FIELD: PartField = PartField::Category;
    const OPERATOR: MongoOperator = MongoOperator::In;
}

impl PartFilterSpec for CountryCodes {
    const FIELD: PartField = PartField::ManufacturerCountry;
    const OPERATOR: MongoOperator = MongoOperator::In;
}

impl PartFilterSpec for Tags {
    const FIELD: PartField = PartField::Tags;
    const OPERATOR: MongoOperator = MongoOperator::All;
}

impl<T> IntoMongoPredicate for T
where
    T: PartFilterSpec,
    T::Item: IntoMongoValue,
{
    fn into_mongo_predicate(self) -> Option<MongoPredicate> {
        let values = self
            .into_iter()
            .map(IntoMongoValue::into_mongo_value)
            .collect::<Vec<_>>();

        if values.is_empty() {
            return None;
        }

        Some(MongoPredicate {
            field: T::FIELD,
            filter: T::OPERATOR.into_filter(values),
        })
    }
}

impl IntoMongoFilter for ListPartsQuery {
    fn into_mongo_filter(self) -> PartFilterDocument {
        [
            self.ids.into_mongo_predicate(),
            self.names.into_mongo_predicate(),
            self.categories.into_mongo_predicate(),
            self.manufacturer_countries.into_mongo_predicate(),
            self.tags.into_mongo_predicate(),
        ]
        .into_iter()
        .flatten()
        .collect()
    }
}
