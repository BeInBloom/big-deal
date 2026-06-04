use std::collections::HashMap;

use crate::{
    domain::models::{
        Dimensions, GetPartQuery, ListPartsQuery, Manufacturer, MetadataValue, MoneyCents, Part,
        PartCategories, PartCategory, StockQuantity,
    },
    grpc::error::InventoryRequestError,
    proto::inventory_v1::{
        self, GetPartRequest, InventoryDimensions, InventoryMetadataValue, InventoryPart,
        ListPartsRequest, inventory_metadata_value::Kind,
    },
};

impl TryFrom<GetPartRequest> for GetPartQuery {
    type Error = InventoryRequestError;

    fn try_from(value: GetPartRequest) -> Result<Self, Self::Error> {
        Ok(Self {
            id: value.uuid.parse()?,
        })
    }
}

impl TryFrom<ListPartsRequest> for ListPartsQuery {
    type Error = InventoryRequestError;

    fn try_from(value: ListPartsRequest) -> Result<Self, Self::Error> {
        let Some(filter) = value.filter else {
            return Ok(Self::default());
        };

        Ok(Self {
            ids: filter.uuids.try_into()?,
            names: filter.names.into(),
            categories: filter.categories.try_into()?,
            manufacturer_countries: filter.manufacturer_countries.into(),
            tags: filter.tags.into(),
        })
    }
}

impl TryFrom<i32> for PartCategory {
    type Error = InventoryRequestError;

    fn try_from(value: i32) -> Result<Self, Self::Error> {
        inventory_v1::PartCategory::try_from(value)?.try_into()
    }
}

impl TryFrom<inventory_v1::PartCategory> for PartCategory {
    type Error = InventoryRequestError;

    fn try_from(value: inventory_v1::PartCategory) -> Result<Self, Self::Error> {
        match value {
            inventory_v1::PartCategory::Engine => Ok(Self::Engine),
            inventory_v1::PartCategory::Fuel => Ok(Self::Fuel),
            inventory_v1::PartCategory::Porthole => Ok(Self::Porthole),
            inventory_v1::PartCategory::Wing => Ok(Self::Wing),
            inventory_v1::PartCategory::Unspecified => {
                Err(InventoryRequestError::InvalidPartCategory)
            }
        }
    }
}

impl TryFrom<Vec<i32>> for PartCategories {
    type Error = InventoryRequestError;

    fn try_from(values: Vec<i32>) -> Result<Self, Self::Error> {
        values.into_iter().map(PartCategory::try_from).collect()
    }
}

impl From<PartCategory> for inventory_v1::PartCategory {
    fn from(value: PartCategory) -> Self {
        match value {
            PartCategory::Engine => Self::Engine,
            PartCategory::Fuel => Self::Fuel,
            PartCategory::Porthole => Self::Porthole,
            PartCategory::Wing => Self::Wing,
        }
    }
}

impl From<Part> for InventoryPart {
    fn from(part: Part) -> Self {
        Self {
            uuid: part.id.into(),
            name: part.name,
            description: String::new(),
            price: money_cents_to_proto(part.price),
            stock_quantity: stock_quantity_to_proto(part.stock_quantity),
            category: inventory_v1::PartCategory::from(part.category) as i32,
            dimensions: Some(part.dimensions.into()),
            manufacturer: Some(part.manufacturer.into()),
            tags: part.tags.into(),
            metadata: convert_metadata(part.metadata),
            created_at: Some(part.created_at.into()),
            updated_at: Some(part.updated_at.into()),
        }
    }
}

impl From<Dimensions> for InventoryDimensions {
    fn from(value: Dimensions) -> Self {
        Self {
            length: value.length.into(),
            width: value.width.into(),
            height: value.height.into(),
            weight: value.weight.into(),
        }
    }
}

impl From<Manufacturer> for inventory_v1::InventoryManufacturer {
    fn from(value: Manufacturer) -> Self {
        Self {
            name: value.name,
            country: value.country.into(),
            website: value.website,
        }
    }
}

impl From<MetadataValue> for InventoryMetadataValue {
    fn from(value: MetadataValue) -> Self {
        let kind = match value {
            MetadataValue::String(value) => Kind::StringValue(value),
            MetadataValue::Int64(value) => Kind::Int64Value(value),
            MetadataValue::Double(value) => Kind::DoubleValue(value),
            MetadataValue::Bool(value) => Kind::BoolValue(value),
        };

        Self { kind: Some(kind) }
    }
}

fn convert_metadata(
    metadata: HashMap<String, MetadataValue>,
) -> HashMap<String, InventoryMetadataValue> {
    metadata
        .into_iter()
        .map(|(key, value)| (key, value.into()))
        .collect()
}

fn money_cents_to_proto(value: MoneyCents) -> f64 {
    u64::from(value) as f64 / 100.0
}

fn stock_quantity_to_proto(value: StockQuantity) -> i64 {
    u64::from(value) as i64
}
