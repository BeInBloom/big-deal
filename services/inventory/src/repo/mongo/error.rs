use thiserror::Error;

use crate::domain::errors::{
    MeasurementError, MoneyCentsError, PartCategoryError, PartIdError, PartRepoError,
    StockQuantityError,
};

#[derive(Debug, Error)]
pub(crate) enum MongoPartRepoError {
    #[error("mongo operation failed: {0}")]
    Driver(#[from] mongodb::error::Error),

    #[error("part document is invalid: {0}")]
    Document(#[from] PartDocumentError),
}

#[derive(Debug, Error)]
pub(crate) enum PartDocumentError {
    #[error("unsupported part schema version: {0}")]
    UnsupportedSchemaVersion(i32),

    #[error("invalid part id: {0}")]
    PartId(#[from] PartIdError),

    #[error("invalid part category: {0}")]
    PartCategory(#[from] PartCategoryError),

    #[error("invalid price cents: {0}")]
    MoneyCents(#[from] MoneyCentsError),

    #[error("invalid stock quantity: {0}")]
    StockQuantity(#[from] StockQuantityError),

    #[error("invalid dimension `{field}`: {source}")]
    InvalidDimension {
        field: &'static str,
        source: MeasurementError,
    },
}

impl From<MongoPartRepoError> for PartRepoError {
    fn from(error: MongoPartRepoError) -> Self {
        match error {
            MongoPartRepoError::Driver(_) => Self::Storage,
            MongoPartRepoError::Document(error) => Self::InvalidData(error.to_string()),
        }
    }
}
