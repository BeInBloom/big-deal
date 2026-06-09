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
