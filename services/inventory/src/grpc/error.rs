use thiserror::Error;

use crate::domain::errors::PartIdError;

#[derive(Debug, Error)]
pub(crate) enum InventoryRequestError {
    #[error(transparent)]
    PartId(#[from] PartIdError),

    #[error("part category is invalid")]
    InvalidPartCategory,

    #[error(transparent)]
    UnknownEnumValue(#[from] prost::UnknownEnumValue),
}
