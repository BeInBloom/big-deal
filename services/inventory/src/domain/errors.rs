#![allow(dead_code)]

use thiserror::Error;

#[derive(Debug, Error)]
pub(crate) enum InventoryError {
    #[error("measurement must be finite and positive")]
    InvalidMeasurement,

    #[error("wrong uuid: {0}")]
    WrongUuid(#[from] PartIdError),
}

#[derive(Debug, Error)]
pub(crate) enum InventoryServiceError {
    #[error("part not found")]
    PartNotFound,

    #[error("part category is invalid")]
    InvalidPartCategory,

    #[error("part repository error: {0}")]
    PartRepo(#[from] PartRepoError),
}

#[derive(Debug, Error)]
pub(crate) enum PartRepoError {
    #[error("part repository failed")]
    Failed,
}

#[derive(Debug, Error)]
pub(crate) enum PartIdError {
    #[error("part id is reuired")]
    Missing,

    #[error("part id is invalid")]
    Invalid,
}

#[derive(Debug, Error)]
pub(crate) enum MeasurementError {
    #[error("meaurement must be finite and positive")]
    Invalid,
}

#[derive(Debug, Error)]
pub(crate) enum InventoryRequestError {
    #[error(transparent)]
    PartId(#[from] PartIdError),

    #[error("part category is invalid")]
    InvalidPartCategory,

    #[error(transparent)]
    UnknownEnumValue(#[from] prost::UnknownEnumValue),
}
