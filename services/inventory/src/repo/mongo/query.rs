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

impl IntoMongoPredicate for PartIds {
    fn into_mongo_predicate(self) -> Option<MongoPredicate> {
        in_predicate(PartField::Id, self)
    }
}

impl IntoMongoPredicate for Names {
    fn into_mongo_predicate(self) -> Option<MongoPredicate> {
        in_predicate(PartField::Name, self)
    }
}

impl IntoMongoPredicate for PartCategories {
    fn into_mongo_predicate(self) -> Option<MongoPredicate> {
        in_predicate(PartField::Category, self)
    }
}

impl IntoMongoPredicate for CountryCodes {
    fn into_mongo_predicate(self) -> Option<MongoPredicate> {
        in_predicate(PartField::ManufacturerCountry, self)
    }
}

impl IntoMongoPredicate for Tags {
    fn into_mongo_predicate(self) -> Option<MongoPredicate> {
        all_predicate(PartField::Tags, self)
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

fn in_predicate<I, T>(field: PartField, values: I) -> Option<MongoPredicate>
where
    T: IntoMongoValue,
    I: IntoIterator<Item = T>,
{
    let values = values
        .into_iter()
        .map(IntoMongoValue::into_mongo_value)
        .collect::<Vec<_>>();

    if values.is_empty() {
        return None;
    }

    Some(MongoPredicate {
        field,
        filter: bson::doc! { "$in": values}.into(),
    })
}

fn all_predicate<I, T>(field: PartField, values: I) -> Option<MongoPredicate>
where
    T: IntoMongoValue,
    I: IntoIterator<Item = T>,
{
    let values = values
        .into_iter()
        .map(IntoMongoValue::into_mongo_value)
        .collect::<Vec<_>>();

    if values.is_empty() {
        return None;
    }

    Some(MongoPredicate {
        field,
        filter: bson::doc! { "$all": values}.into(),
    })
}
