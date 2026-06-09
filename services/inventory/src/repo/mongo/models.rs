use std::collections::HashMap;

use mongodb::bson;
use serde::{Deserialize, Serialize};

use crate::{
    domain::models::{Dimensions, Manufacturer, MetadataValue, Part, PositiveMeasurement},
    repo::mongo::error::PartDocumentError,
};

const PART_SCHEMA_VERSION: i32 = 1;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub(crate) struct PartDocument {
    #[serde(rename = "_id")]
    pub(crate) id: String,
    pub(crate) schema_version: i32,
    pub(crate) name: String,
    pub(crate) description: String,
    pub(crate) price_cents: i64,
    pub(crate) stock_quantity: i64,
    pub(crate) category: String,
    pub(crate) dimensions: DimensionsDocument,
    pub(crate) manufacturer: ManufacturerDocument,
    pub(crate) tags: Vec<String>,
    pub(crate) metadata: HashMap<String, MetadataValueDocument>,
    pub(crate) created_at: bson::DateTime,
    pub(crate) updated_at: bson::DateTime,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub(crate) struct DimensionsDocument {
    pub(crate) length: f64,
    pub(crate) width: f64,
    pub(crate) height: f64,
    pub(crate) weight: f64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub(crate) struct ManufacturerDocument {
    pub(crate) name: String,
    pub(crate) country: String,
    pub(crate) website: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(tag = "kind", content = "value", rename_all = "snake_case")]
pub(crate) enum MetadataValueDocument {
    String(String),
    Int64(i64),
    Double(f64),
    Bool(bool),
}

impl From<Part> for PartDocument {
    fn from(part: Part) -> Self {
        Self {
            id: part.id.into(),
            schema_version: PART_SCHEMA_VERSION,
            name: part.name,
            description: part.description,
            price_cents: u64::from(part.price) as i64,
            stock_quantity: u64::from(part.stock_quantity) as i64,
            category: part.category.to_string(),
            dimensions: part.dimensions.into(),
            manufacturer: part.manufacturer.into(),
            tags: part.tags.into(),
            metadata: part
                .metadata
                .into_iter()
                .map(|(key, value)| (key, value.into()))
                .collect(),
            created_at: bson::DateTime::from_system_time(part.created_at),
            updated_at: bson::DateTime::from_system_time(part.updated_at),
        }
    }
}

impl From<Dimensions> for DimensionsDocument {
    fn from(value: Dimensions) -> Self {
        Self {
            length: value.length.into(),
            width: value.width.into(),
            height: value.height.into(),
            weight: value.weight.into(),
        }
    }
}

impl From<Manufacturer> for ManufacturerDocument {
    fn from(value: Manufacturer) -> Self {
        Self {
            name: value.name,
            country: value.country.into(),
            website: value.website,
        }
    }
}

impl From<MetadataValue> for MetadataValueDocument {
    fn from(value: MetadataValue) -> Self {
        match value {
            MetadataValue::String(value) => Self::String(value),
            MetadataValue::Int64(value) => Self::Int64(value),
            MetadataValue::Double(value) => Self::Double(value),
            MetadataValue::Bool(value) => Self::Bool(value),
        }
    }
}

impl TryFrom<PartDocument> for Part {
    type Error = PartDocumentError;

    fn try_from(document: PartDocument) -> Result<Self, Self::Error> {
        if document.schema_version != PART_SCHEMA_VERSION {
            return Err(PartDocumentError::UnsupportedSchemaVersion(
                document.schema_version,
            ));
        }

        Ok(Self {
            id: document.id.parse()?,
            name: document.name,
            description: document.description,
            price: document.price_cents.try_into()?,
            stock_quantity: document.stock_quantity.try_into()?,
            category: document.category.parse()?,
            dimensions: document.dimensions.try_into()?,
            manufacturer: document.manufacturer.into(),
            tags: document.tags.into(),
            metadata: document
                .metadata
                .into_iter()
                .map(|(key, value)| (key, value.into()))
                .collect(),
            created_at: document.created_at.to_system_time(),
            updated_at: document.updated_at.to_system_time(),
        })
    }
}

impl TryFrom<DimensionsDocument> for Dimensions {
    type Error = PartDocumentError;

    fn try_from(value: DimensionsDocument) -> Result<Self, Self::Error> {
        Ok(Self {
            length: positive_measurement("dimensions.length", value.length)?,
            width: positive_measurement("dimensions.width", value.width)?,
            height: positive_measurement("dimensions.height", value.height)?,
            weight: positive_measurement("dimensions.weight", value.weight)?,
        })
    }
}

fn positive_measurement(
    field: &'static str,
    value: f64,
) -> Result<PositiveMeasurement, PartDocumentError> {
    PositiveMeasurement::try_from(value)
        .map_err(|source| PartDocumentError::InvalidDimension { field, source })
}

impl From<ManufacturerDocument> for Manufacturer {
    fn from(value: ManufacturerDocument) -> Self {
        Self {
            name: value.name,
            country: value.country.into(),
            website: value.website,
        }
    }
}

impl From<MetadataValueDocument> for MetadataValue {
    fn from(value: MetadataValueDocument) -> Self {
        match value {
            MetadataValueDocument::String(value) => Self::String(value),
            MetadataValueDocument::Int64(value) => Self::Int64(value),
            MetadataValueDocument::Double(value) => Self::Double(value),
            MetadataValueDocument::Bool(value) => Self::Bool(value),
        }
    }
}

#[cfg(test)]
mod tests {
    use std::{
        collections::{HashMap, HashSet},
        time::{Duration, SystemTime},
    };

    use crate::domain::models::{CountryCode, MoneyCents, PartCategory, StockQuantity};

    use super::*;

    const RAW_PART_ID: &str = "11111111-1111-4111-8111-111111111111";

    #[test]
    fn part_converts_to_document() {
        let document = PartDocument::from(sample_part());

        assert_eq!(document.id, RAW_PART_ID);
        assert_eq!(document.schema_version, PART_SCHEMA_VERSION);
        assert_eq!(document.name, "Main engine");
        assert_eq!(document.description, "Primary propulsion unit");
        assert_eq!(document.price_cents, 12_500);
        assert_eq!(document.stock_quantity, 4);
        assert_eq!(document.category, "ENGINE");
        assert_eq!(document.dimensions.length, 1.0);
        assert_eq!(document.dimensions.width, 2.0);
        assert_eq!(document.dimensions.height, 3.0);
        assert_eq!(document.dimensions.weight, 4.0);
        assert_eq!(document.manufacturer.name, "ACME");
        assert_eq!(document.manufacturer.country, "US");
        assert_eq!(document.manufacturer.website, "https://acme.example");
        assert_eq!(
            string_set(&document.tags),
            HashSet::from(["critical", "engine"])
        );
        assert_eq!(document.created_at.to_system_time(), created_at());
        assert_eq!(document.updated_at.to_system_time(), updated_at());

        assert!(matches!(
            document.metadata.get("power_kw"),
            Some(MetadataValueDocument::Int64(900))
        ));
        assert!(matches!(
            document.metadata.get("reusable"),
            Some(MetadataValueDocument::Bool(true))
        ));
        assert!(matches!(
            document.metadata.get("material"),
            Some(MetadataValueDocument::String(value)) if value == "titanium"
        ));
        assert!(matches!(
            document.metadata.get("efficiency"),
            Some(MetadataValueDocument::Double(0.98))
        ));
    }

    #[test]
    fn part_document_serializes_id_as_mongo_id() {
        let raw_document = bson::to_document(&PartDocument::from(sample_part())).unwrap();

        assert_eq!(raw_document.get_str("_id").unwrap(), RAW_PART_ID);
        assert!(!raw_document.contains_key("id"));
    }

    #[test]
    fn document_converts_to_part() {
        let part = Part::try_from(sample_document()).unwrap();

        assert_eq!(String::from(part.id), RAW_PART_ID);
        assert_eq!(part.name, "Main engine");
        assert_eq!(part.description, "Primary propulsion unit");
        assert_eq!(u64::from(part.price), 12_500);
        assert_eq!(u64::from(part.stock_quantity), 4);
        assert_eq!(part.category, PartCategory::Engine);
        assert_eq!(f64::from(part.dimensions.length), 1.0);
        assert_eq!(f64::from(part.dimensions.width), 2.0);
        assert_eq!(f64::from(part.dimensions.height), 3.0);
        assert_eq!(f64::from(part.dimensions.weight), 4.0);
        assert_eq!(part.manufacturer.name, "ACME");
        assert_eq!(String::from(part.manufacturer.country), "US");
        assert_eq!(part.manufacturer.website, "https://acme.example");
        assert_eq!(part.created_at, created_at());
        assert_eq!(part.updated_at, updated_at());

        let tags: Vec<String> = part.tags.into();
        assert_eq!(string_set(&tags), HashSet::from(["critical", "engine"]));
        assert!(matches!(
            part.metadata.get("power_kw"),
            Some(MetadataValue::Int64(900))
        ));
        assert!(matches!(
            part.metadata.get("reusable"),
            Some(MetadataValue::Bool(true))
        ));
        assert!(matches!(
            part.metadata.get("material"),
            Some(MetadataValue::String(value)) if value == "titanium"
        ));
        assert!(matches!(
            part.metadata.get("efficiency"),
            Some(MetadataValue::Double(0.98))
        ));
    }

    #[test]
    fn document_rejects_unknown_schema_version() {
        let mut document = sample_document();
        document.schema_version = 2;

        assert!(matches!(
            Part::try_from(document),
            Err(PartDocumentError::UnsupportedSchemaVersion(2))
        ));
    }

    #[test]
    fn document_rejects_invalid_category() {
        let mut document = sample_document();
        document.category = "LASER".to_string();

        assert!(matches!(
            Part::try_from(document),
            Err(PartDocumentError::PartCategory(_))
        ));
    }

    #[test]
    fn document_rejects_negative_price_cents() {
        let mut document = sample_document();
        document.price_cents = -1;

        assert!(matches!(
            Part::try_from(document),
            Err(PartDocumentError::MoneyCents(_))
        ));
    }

    #[test]
    fn document_rejects_negative_stock_quantity() {
        let mut document = sample_document();
        document.stock_quantity = -1;

        assert!(matches!(
            Part::try_from(document),
            Err(PartDocumentError::StockQuantity(_))
        ));
    }

    #[test]
    fn document_rejects_invalid_dimension() {
        let mut document = sample_document();
        document.dimensions.length = 0.0;

        assert!(matches!(
            Part::try_from(document),
            Err(PartDocumentError::InvalidDimension {
                field: "dimensions.length",
                ..
            })
        ));
    }

    fn sample_part() -> Part {
        let mut metadata = HashMap::new();
        metadata.insert("power_kw".to_string(), MetadataValue::Int64(900));
        metadata.insert("reusable".to_string(), MetadataValue::Bool(true));
        metadata.insert(
            "material".to_string(),
            MetadataValue::String("titanium".to_string()),
        );
        metadata.insert("efficiency".to_string(), MetadataValue::Double(0.98));

        Part {
            id: RAW_PART_ID.parse().unwrap(),
            name: "Main engine".to_string(),
            description: "Primary propulsion unit".to_string(),
            price: MoneyCents::new(12_500),
            stock_quantity: StockQuantity::new(4),
            category: PartCategory::Engine,
            dimensions: sample_dimensions(),
            manufacturer: Manufacturer {
                name: "ACME".to_string(),
                country: CountryCode::from("US".to_string()),
                website: "https://acme.example".to_string(),
            },
            tags: vec!["engine".to_string(), "critical".to_string()].into(),
            metadata,
            created_at: created_at(),
            updated_at: updated_at(),
        }
    }

    fn sample_document() -> PartDocument {
        let mut metadata = HashMap::new();
        metadata.insert("power_kw".to_string(), MetadataValueDocument::Int64(900));
        metadata.insert("reusable".to_string(), MetadataValueDocument::Bool(true));
        metadata.insert(
            "material".to_string(),
            MetadataValueDocument::String("titanium".to_string()),
        );
        metadata.insert(
            "efficiency".to_string(),
            MetadataValueDocument::Double(0.98),
        );

        PartDocument {
            id: RAW_PART_ID.to_string(),
            schema_version: PART_SCHEMA_VERSION,
            name: "Main engine".to_string(),
            description: "Primary propulsion unit".to_string(),
            price_cents: 12_500,
            stock_quantity: 4,
            category: "ENGINE".to_string(),
            dimensions: DimensionsDocument {
                length: 1.0,
                width: 2.0,
                height: 3.0,
                weight: 4.0,
            },
            manufacturer: ManufacturerDocument {
                name: "ACME".to_string(),
                country: "US".to_string(),
                website: "https://acme.example".to_string(),
            },
            tags: vec!["engine".to_string(), "critical".to_string()],
            metadata,
            created_at: bson::DateTime::from_system_time(created_at()),
            updated_at: bson::DateTime::from_system_time(updated_at()),
        }
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

    fn created_at() -> SystemTime {
        SystemTime::UNIX_EPOCH
    }

    fn updated_at() -> SystemTime {
        SystemTime::UNIX_EPOCH + Duration::from_secs(1)
    }

    fn string_set(values: &[String]) -> HashSet<&str> {
        values.iter().map(String::as_str).collect()
    }
}
