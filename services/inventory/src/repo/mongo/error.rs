use platform_mongo::FilterError;
use thiserror::Error;

use crate::domain::errors::{
    MeasurementError, MoneyCentsError, PartCategoryError, PartIdError, PartRepoError,
    StockQuantityError,
};

#[derive(Debug, Error)]
pub(crate) enum MongoPartRepoError {
    #[error("mongo operation failed: {0}")]
    Driver(mongodb::error::Error),

    #[error("part document deserialization failed: {0}")]
    DocumentDecode(mongodb::error::Error),

    #[error("part document is invalid: {0}")]
    Document(#[from] PartDocumentError),

    #[error("mongo filter is invalid: {0}")]
    Filter(#[from] FilterError),
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

impl From<mongodb::error::Error> for MongoPartRepoError {
    fn from(error: mongodb::error::Error) -> Self {
        if matches!(
            error.kind.as_ref(),
            mongodb::error::ErrorKind::BsonDeserialization(_)
        ) {
            Self::DocumentDecode(error)
        } else {
            Self::Driver(error)
        }
    }
}

impl From<MongoPartRepoError> for PartRepoError {
    fn from(error: MongoPartRepoError) -> Self {
        match error {
            MongoPartRepoError::Driver(_) | MongoPartRepoError::Filter(_) => Self::Storage,
            MongoPartRepoError::Document(error) => Self::InvalidData(error.to_string()),
            MongoPartRepoError::DocumentDecode(error) => Self::InvalidData(error.to_string()),
        }
    }
}
