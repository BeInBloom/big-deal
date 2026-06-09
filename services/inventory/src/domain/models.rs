use std::{
    borrow::Borrow,
    collections::{HashMap, HashSet},
    fmt,
    str::FromStr,
    time::SystemTime,
};

use uuid::Uuid;

use crate::domain::errors::{
    MeasurementError, MoneyCentsError, PartCategoryError, PartIdError, StockQuantityError,
};

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

impl From<Name> for String {
    fn from(value: Name) -> Self {
        value.0
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

impl From<Names> for Vec<String> {
    fn from(value: Names) -> Self {
        value.0.into_iter().map(String::from).collect()
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

        Ok(Self(s.parse()?))
    }
}

impl From<PartId> for String {
    fn from(value: PartId) -> Self {
        value.0.to_string()
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

impl From<PartIds> for Vec<String> {
    fn from(value: PartIds) -> Self {
        value.0.into_iter().map(String::from).collect()
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Hash)]
pub(crate) struct Tag(String);

impl From<String> for Tag {
    fn from(value: String) -> Self {
        Self(value)
    }
}

impl From<Tag> for String {
    fn from(value: Tag) -> Self {
        value.0
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

impl From<Tags> for Vec<String> {
    fn from(value: Tags) -> Self {
        value.0.into_iter().map(String::from).collect()
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Hash)]
pub(crate) struct MoneyCents(u64);

impl MoneyCents {
    pub(crate) fn new(num: u64) -> Self {
        Self(num)
    }
}

impl From<MoneyCents> for u64 {
    fn from(value: MoneyCents) -> Self {
        value.0
    }
}

impl From<u64> for MoneyCents {
    fn from(value: u64) -> Self {
        Self::new(value)
    }
}

impl TryFrom<i64> for MoneyCents {
    type Error = MoneyCentsError;

    fn try_from(value: i64) -> Result<Self, Self::Error> {
        let value = u64::try_from(value).map_err(|_| MoneyCentsError::Negative(value))?;
        Ok(Self::from(value))
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Hash)]
pub(crate) struct StockQuantity(u64);

impl StockQuantity {
    pub(crate) fn new(num: u64) -> Self {
        Self(num)
    }
}

impl From<StockQuantity> for u64 {
    fn from(value: StockQuantity) -> Self {
        value.0
    }
}

impl From<u64> for StockQuantity {
    fn from(value: u64) -> Self {
        Self::new(value)
    }
}

impl TryFrom<i64> for StockQuantity {
    type Error = StockQuantityError;

    fn try_from(value: i64) -> Result<Self, Self::Error> {
        let value = u64::try_from(value).map_err(|_| StockQuantityError::Negative(value))?;
        Ok(Self::from(value))
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Hash)]
pub(crate) struct CountryCode(String);

impl From<String> for CountryCode {
    fn from(value: String) -> Self {
        Self(value)
    }
}

impl From<CountryCode> for String {
    fn from(value: CountryCode) -> Self {
        value.0
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

impl From<CountryCodes> for Vec<String> {
    fn from(value: CountryCodes) -> Self {
        value.0.into_iter().map(String::from).collect()
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub(crate) enum PartCategory {
    Engine,
    Fuel,
    Porthole,
    Wing,
}

impl fmt::Display for PartCategory {
    fn fmt(&self, formatter: &mut fmt::Formatter<'_>) -> fmt::Result {
        let value = match self {
            Self::Engine => "ENGINE",
            Self::Fuel => "FUEL",
            Self::Porthole => "PORTHOLE",
            Self::Wing => "WING",
        };

        formatter.write_str(value)
    }
}

impl FromStr for PartCategory {
    type Err = PartCategoryError;

    fn from_str(value: &str) -> Result<Self, Self::Err> {
        match value {
            "ENGINE" => Ok(Self::Engine),
            "FUEL" => Ok(Self::Fuel),
            "PORTHOLE" => Ok(Self::Porthole),
            "WING" => Ok(Self::Wing),
            _ => Err(PartCategoryError::Unknown(value.to_string())),
        }
    }
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
    //TODO: Подумать над тем, чтобы заменить на Name
    //пока смысла в этом не много
    pub(crate) name: String,
    pub(crate) price: MoneyCents,
    pub(crate) stock_quantity: StockQuantity,
    pub(crate) category: PartCategory,
    pub(crate) description: String,
    pub(crate) dimensions: Dimensions,
    pub(crate) manufacturer: Manufacturer,
    pub(crate) tags: Tags,
    pub(crate) metadata: HashMap<String, MetadataValue>,
    pub(crate) created_at: SystemTime,
    //Возможно стоит более явно различать было обновление или нет
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

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn part_id_parses_valid_uuid() {
        let id = uuid::Uuid::new_v4();
        let part_id: PartId = id.to_string().parse().unwrap();
        assert_eq!(part_id.0, id);
    }

    #[test]
    fn part_id_rejects_empty_value() {
        let err = "".parse::<PartId>().unwrap_err();
        assert!(matches!(err, PartIdError::Missing));
    }

    #[test]
    fn part_id_rejects_invalid_uuid() {
        let err = "some pretty kitty".parse::<PartId>().unwrap_err();
        assert!(matches!(err, PartIdError::Invalid(_)));
    }

    #[test]
    fn positive_measurement_accepts_positive_finite_value() {
        let num = 10.10_f64;
        let positive_measurement: PositiveMeasurement = num.try_into().unwrap();
        assert_eq!(num, positive_measurement.0);
    }

    #[test]
    fn positive_measurement_rejects_zero() {
        let num = 0_f64;
        let err = PositiveMeasurement::try_from(num).unwrap_err();
        assert!(matches!(err, MeasurementError::Invalid));
    }

    #[test]
    fn positive_measurement_rejects_negative_value() {
        let num = -1.123_f64;
        let err = PositiveMeasurement::try_from(num).unwrap_err();
        assert!(matches!(err, MeasurementError::Invalid));
    }

    #[test]
    fn positive_measurement_rejects_nan() {
        let num = f64::NAN;
        let err = PositiveMeasurement::try_from(num).unwrap_err();
        assert!(matches!(err, MeasurementError::Invalid));
    }

    #[test]
    fn positive_measurement_rejects_infinity() {
        let num = f64::INFINITY;
        let err = PositiveMeasurement::try_from(num).unwrap_err();
        assert!(matches!(err, MeasurementError::Invalid));
    }

    #[test]
    fn list_parts_query_matches_any_part_by_default() {
        let sample_parts = sample_parts();
        let default_query = ListPartsQuery::default();

        for part in sample_parts {
            assert!(default_query.matches(&part));
        }
    }

    #[test]
    fn list_parts_query_matches_part_when_all_filters_match() {
        let part = sample_part();

        let query = ListPartsQuery {
            ids: vec![String::from(part.id)].try_into().unwrap(),
            names: vec![part.name.clone()].into(),
            categories: [part.category].into_iter().collect(),
            manufacturer_countries: vec![String::from(part.manufacturer.country.clone())].into(),
            tags: vec!["engine".to_string(), "critical".to_string()].into(),
        };

        assert!(query.matches(&part));
    }

    #[test]
    fn list_parts_query_rejects_part_when_filter_does_not_match() {
        let part = sample_part();

        let cases = [
            ListPartsQuery {
                ids: vec!["99999999-9999-4999-8999-999999999999".to_string()]
                    .try_into()
                    .unwrap(),
                ..ListPartsQuery::default()
            },
            ListPartsQuery {
                names: vec!["Wrong name".to_string()].into(),
                ..ListPartsQuery::default()
            },
            ListPartsQuery {
                categories: [PartCategory::Porthole].into_iter().collect(),
                ..ListPartsQuery::default()
            },
            ListPartsQuery {
                manufacturer_countries: vec!["FR".to_string()].into(),
                ..ListPartsQuery::default()
            },
            ListPartsQuery {
                tags: vec!["missing-tag".to_string()].into(),
                ..ListPartsQuery::default()
            },
        ];

        for query in cases {
            assert!(!query.matches(&part));
        }
    }

    fn sample_parts() -> Vec<Part> {
        vec![
            Part {
                id: "11111111-1111-4111-8111-111111111111".parse().unwrap(),
                name: "Main engine".to_string(),
                price: MoneyCents(12_500),
                stock_quantity: StockQuantity(4),
                category: PartCategory::Engine,
                description: String::default(),
                dimensions: sample_dimensions(),
                manufacturer: Manufacturer {
                    name: "ACME".to_string(),
                    country: CountryCode("US".to_string()),
                    website: "https://acme.example".to_string(),
                },
                tags: vec!["engine".to_string(), "critical".to_string()].into(),
                metadata: HashMap::new(),
                created_at: SystemTime::UNIX_EPOCH,
                updated_at: SystemTime::UNIX_EPOCH,
            },
            Part {
                id: "22222222-2222-4222-8222-222222222222".parse().unwrap(),
                name: "Fuel tank".to_string(),
                price: MoneyCents(8_000),
                stock_quantity: StockQuantity(12),
                category: PartCategory::Fuel,
                description: String::default(),
                dimensions: sample_dimensions(),
                manufacturer: Manufacturer {
                    name: "Orbital Parts".to_string(),
                    country: CountryCode("DE".to_string()),
                    website: "https://orbital.example".to_string(),
                },
                tags: vec!["fuel".to_string(), "storage".to_string()].into(),
                metadata: HashMap::new(),
                created_at: SystemTime::UNIX_EPOCH,
                updated_at: SystemTime::UNIX_EPOCH,
            },
            Part {
                id: "33333333-3333-4333-8333-333333333333".parse().unwrap(),
                name: "Left wing".to_string(),
                price: MoneyCents(20_000),
                stock_quantity: StockQuantity(2),
                category: PartCategory::Wing,
                description: String::default(),
                dimensions: sample_dimensions(),
                manufacturer: Manufacturer {
                    name: "Sky Forge".to_string(),
                    country: CountryCode("JP".to_string()),
                    website: "https://skyforge.example".to_string(),
                },
                tags: vec!["wing".to_string(), "aero".to_string()].into(),
                metadata: HashMap::new(),
                created_at: SystemTime::UNIX_EPOCH,
                updated_at: SystemTime::UNIX_EPOCH,
            },
        ]
    }

    fn sample_part() -> Part {
        sample_parts().into_iter().next().unwrap()
    }

    fn sample_dimensions() -> Dimensions {
        Dimensions {
            length: measurement(1.0),
            width: measurement(2.0),
            height: measurement(3.0),
            weight: measurement(4.0),
        }
    }

    fn measurement(value: f64) -> PositiveMeasurement {
        PositiveMeasurement::try_from(value).unwrap()
    }
}
