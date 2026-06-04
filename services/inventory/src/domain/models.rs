#![allow(dead_code)]
use std::{
    borrow::Borrow,
    collections::{HashMap, HashSet},
    str::FromStr,
    time::SystemTime,
};

use anyhow::Result;
use uuid::Uuid;

use crate::domain::errors::{MeasurementError, PartIdError};

#[derive(Debug, Clone, PartialEq, Eq, Hash)]
pub(crate) struct Name(String);

impl Borrow<str> for Name {
    fn borrow(&self) -> &str {
        &self.0
    }
}

impl From<String> for Name {
    fn from(value: String) -> Self {
        Self(value)
    }
}

#[derive(Debug, Default, Clone, PartialEq, Eq)]
pub(crate) struct Names(HashSet<Name>);

impl Names {
    pub(crate) fn matches(&self, name: &str) -> bool {
        self.0.is_empty() || self.0.contains(name)
    }
}

impl From<Vec<String>> for Names {
    fn from(values: Vec<String>) -> Self {
        Self(values.into_iter().map(Name::from).collect())
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub(crate) struct PartId(Uuid);

impl FromStr for PartId {
    type Err = PartIdError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        if s.is_empty() {
            return Err(PartIdError::Missing);
        }

        Ok(Self(s.parse().map_err(|_| PartIdError::Invalid)?))
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Default)]
pub(crate) struct PartIds(HashSet<PartId>);

impl PartIds {
    pub(crate) fn matches(&self, id: &PartId) -> bool {
        self.0.is_empty() || self.0.contains(id)
    }
}

impl TryFrom<Vec<String>> for PartIds {
    type Error = PartIdError;

    fn try_from(value: Vec<String>) -> Result<Self, Self::Error> {
        Ok(Self(
            value
                .into_iter()
                .map(|v| v.parse())
                .collect::<Result<HashSet<_>, _>>()?,
        ))
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Hash)]
pub(crate) struct Tag(String);

impl From<String> for Tag {
    fn from(value: String) -> Self {
        Self(value)
    }
}

#[derive(Debug, Default, Clone, PartialEq, Eq)]
pub(crate) struct Tags(HashSet<Tag>);

impl Tags {
    pub(crate) fn matches(&self, tags: &Tags) -> bool {
        self.0.is_empty() || self.0.is_subset(&tags.0)
    }
}

impl From<Vec<String>> for Tags {
    fn from(values: Vec<String>) -> Self {
        Self(values.into_iter().map(Tag::from).collect())
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Hash)]
pub(crate) struct MoneyCents(u64);

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Hash)]
pub(crate) struct StockQuantity(u64);

#[derive(Debug, Clone, PartialEq, Eq, Hash)]
pub(crate) struct CountryCode(String);

impl From<String> for CountryCode {
    fn from(value: String) -> Self {
        Self(value)
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Default)]
pub(crate) struct CountryCodes(HashSet<CountryCode>);

impl CountryCodes {
    pub(crate) fn matches(&self, country: &CountryCode) -> bool {
        self.0.is_empty() || self.0.contains(country)
    }
}

impl From<Vec<String>> for CountryCodes {
    fn from(values: Vec<String>) -> Self {
        Self(values.into_iter().map(CountryCode::from).collect())
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub(crate) enum PartCategory {
    Engine,
    Fuel,
    Porthole,
    Wing,
}

#[derive(Debug, Clone, PartialEq, Eq, Default)]
pub(crate) struct PartCategories(HashSet<PartCategory>);

impl PartCategories {
    pub(crate) fn matches(&self, category: &PartCategory) -> bool {
        self.0.is_empty() || self.0.contains(category)
    }
}

impl FromIterator<PartCategory> for PartCategories {
    fn from_iter<T>(iter: T) -> Self
    where
        T: IntoIterator<Item = PartCategory>,
    {
        Self(iter.into_iter().collect())
    }
}

#[derive(Debug, Clone, Copy, PartialEq)]
pub(crate) struct PositiveMeasurement(f64);

impl TryFrom<f64> for PositiveMeasurement {
    type Error = MeasurementError;

    fn try_from(value: f64) -> Result<Self, Self::Error> {
        if !value.is_finite() || value <= 0.0 {
            return Err(MeasurementError::Invalid);
        }

        Ok(Self(value))
    }
}

impl From<PositiveMeasurement> for f64 {
    fn from(value: PositiveMeasurement) -> Self {
        value.0
    }
}

#[derive(Debug, Clone, Copy, PartialEq)]
pub(crate) struct Dimensions {
    pub(crate) length: PositiveMeasurement,
    pub(crate) width: PositiveMeasurement,
    pub(crate) height: PositiveMeasurement,
    pub(crate) weight: PositiveMeasurement,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub(crate) struct Manufacturer {
    pub(crate) name: String,
    pub(crate) country: CountryCode,
    pub(crate) website: String,
}

#[derive(Debug, Clone, PartialEq)]
pub(crate) enum MetadataValue {
    String(String),
    Int64(i64),
    Double(f64),
    Bool(bool),
}

#[derive(Debug, Clone, PartialEq)]
pub(crate) struct Part {
    pub(crate) id: PartId,
    pub(crate) name: String,
    pub(crate) price: MoneyCents,
    pub(crate) stock_quantity: StockQuantity,
    pub(crate) category: PartCategory,
    pub(crate) dimensions: Dimensions,
    pub(crate) manufacturer: Manufacturer,
    pub(crate) tags: Tags,
    pub(crate) metadata: HashMap<String, MetadataValue>,
    pub(crate) created_at: SystemTime,
    pub(crate) updated_at: SystemTime,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub(crate) struct GetPartQuery {
    pub(crate) id: PartId,
}

#[derive(Debug, Default, Clone, PartialEq, Eq)]
pub(crate) struct ListPartsQuery {
    pub(crate) ids: PartIds,
    pub(crate) names: Names,
    pub(crate) categories: PartCategories,
    pub(crate) manufacturer_countries: CountryCodes,
    pub(crate) tags: Tags,
}

impl ListPartsQuery {
    pub(crate) fn matches(&self, part: &Part) -> bool {
        self.ids.matches(&part.id)
            && self.names.matches(&part.name)
            && self.categories.matches(&part.category)
            && self
                .manufacturer_countries
                .matches(&part.manufacturer.country)
            && self.tags.matches(&part.tags)
    }
}
