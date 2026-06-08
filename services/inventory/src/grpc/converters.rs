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

#[cfg(test)]
mod tests {
    use crate::domain::errors::PartIdError;
    use crate::domain::models::{CountryCode, PartId, Tags};

    use super::*;

    #[test]
    fn part_category_maps_proto_values() {
        let cases = [
            (inventory_v1::PartCategory::Engine, PartCategory::Engine),
            (inventory_v1::PartCategory::Fuel, PartCategory::Fuel),
            (inventory_v1::PartCategory::Porthole, PartCategory::Porthole),
            (inventory_v1::PartCategory::Wing, PartCategory::Wing),
        ];

        for (proto_category, expected_category) in cases {
            let category = PartCategory::try_from(proto_category).unwrap();
            assert_eq!(category, expected_category);
        }
    }

    #[test]
    fn part_category_rejects_unspecified() {
        let err = PartCategory::try_from(inventory_v1::PartCategory::Unspecified).unwrap_err();
        assert!(matches!(err, InventoryRequestError::InvalidPartCategory));
    }

    #[test]
    fn get_part_request_converts_to_query() {
        let raw_uuid = uuid::Uuid::new_v4().to_string();
        let req = GetPartRequest {
            uuid: raw_uuid.clone(),
        };
        let query = GetPartQuery::try_from(req).unwrap();
        assert_eq!(String::from(query.id), raw_uuid);
    }

    #[test]
    fn get_part_request_rejects_empty_uuid() {
        let req = GetPartRequest {
            uuid: String::new(),
        };
        let err = GetPartQuery::try_from(req).unwrap_err();
        assert!(matches!(
            err,
            InventoryRequestError::PartId(PartIdError::Missing),
        ));
    }

    #[test]
    fn get_part_request_rejects_invalid_uuid() {
        let req = GetPartRequest {
            uuid: String::from("work harder, comrade"),
        };
        let err = GetPartQuery::try_from(req).unwrap_err();
        assert!(matches!(
            err,
            InventoryRequestError::PartId(PartIdError::Invalid(_))
        ));
    }

    #[test]
    fn list_parts_request_without_filter_converts_to_default_query() {
        let req = ListPartsRequest { filter: None };
        let query = ListPartsQuery::try_from(req).unwrap();
        assert_eq!(query, ListPartsQuery::default());
    }

    #[test]
    fn list_parts_request_converts_filter_to_query() {
        let raw_uuid = "11111111-1111-4111-8111-111111111111".to_string();

        let req = ListPartsRequest {
            filter: Some(inventory_v1::InventoryPartsFilter {
                uuids: vec![raw_uuid.clone()],
                names: vec!["Main engine".to_string()],
                categories: vec![inventory_v1::PartCategory::Engine as i32],
                manufacturer_countries: vec!["US".to_string()],
                tags: vec!["engine".to_string(), "critical".to_string()],
            }),
        };

        let query = ListPartsQuery::try_from(req).unwrap();

        let expected_id: PartId = raw_uuid.parse().unwrap();
        assert!(query.ids.matches(&expected_id));
        assert!(query.names.matches("Main engine"));
        assert!(query.categories.matches(&PartCategory::Engine));
        assert!(
            query
                .manufacturer_countries
                .matches(&CountryCode::from("US".to_string()))
        );
        assert!(query.tags.matches(&Tags::from(vec![
            "engine".to_string(),
            "critical".to_string(),
            "extra".to_string(),
        ])));

        let other_id: PartId = "22222222-2222-4222-8222-222222222222".parse().unwrap();
        assert!(!query.ids.matches(&other_id));
        assert!(!query.names.matches("Wrong name"));
        assert!(!query.categories.matches(&PartCategory::Wing));
        assert!(
            !query
                .manufacturer_countries
                .matches(&CountryCode::from("DE".to_string()))
        );
        assert!(!query.tags.matches(&Tags::from(vec!["engine".to_string()])));
    }

    #[test]
    fn list_parts_request_rejects_invalid_category() {
        let req = ListPartsRequest {
            filter: Some(inventory_v1::InventoryPartsFilter {
                uuids: Vec::new(),
                names: Vec::new(),
                categories: vec![999],
                manufacturer_countries: Vec::new(),
                tags: Vec::new(),
            }),
        };

        let err = ListPartsQuery::try_from(req).unwrap_err();
        assert!(matches!(err, InventoryRequestError::UnknownEnumValue(_)));
    }

    #[test]
    fn list_parts_request_rejects_empty_uuid() {
        let req = ListPartsRequest {
            filter: Some(inventory_v1::InventoryPartsFilter {
                uuids: vec![String::new()],
                names: Vec::new(),
                categories: Vec::new(),
                manufacturer_countries: Vec::new(),
                tags: Vec::new(),
            }),
        };

        let err = ListPartsQuery::try_from(req).unwrap_err();
        assert!(matches!(
            err,
            InventoryRequestError::PartId(PartIdError::Missing)
        ));
    }

    #[test]
    fn list_parts_request_rejects_invalid_uuid() {
        let req = ListPartsRequest {
            filter: Some(inventory_v1::InventoryPartsFilter {
                uuids: vec!["not-a-uuid".to_string()],
                names: Vec::new(),
                categories: Vec::new(),
                manufacturer_countries: Vec::new(),
                tags: Vec::new(),
            }),
        };

        let err = ListPartsQuery::try_from(req).unwrap_err();
        assert!(matches!(
            err,
            InventoryRequestError::PartId(PartIdError::Invalid(_))
        ));
    }

    #[test]
    fn list_parts_request_rejects_unspecified_category() {
        let req = ListPartsRequest {
            filter: Some(inventory_v1::InventoryPartsFilter {
                uuids: Vec::new(),
                names: Vec::new(),
                categories: vec![inventory_v1::PartCategory::Unspecified as i32],
                manufacturer_countries: Vec::new(),
                tags: Vec::new(),
            }),
        };

        let err = ListPartsQuery::try_from(req).unwrap_err();
        assert!(matches!(err, InventoryRequestError::InvalidPartCategory));
    }
}
