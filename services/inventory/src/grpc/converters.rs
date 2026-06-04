use crate::{
    domain::{
        errors::InventoryRequestError,
        models::{GetPartQuery, ListPartsQuery, PartCategories, PartCategory},
    },
    proto::inventory_v1::{self, GetPartRequest, ListPartsRequest},
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
