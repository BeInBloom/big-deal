use thiserror::Error;

#[allow(dead_code)]
#[derive(Debug, Error)]
pub(crate) enum PartRepoError {
    #[error("part repository storage failed")]
    Storage,

    #[error("part repository returned invalid data: {0}")]
    InvalidData(String),
}

#[derive(Debug, Error)]
pub(crate) enum PartIdError {
    #[error("part id is reuired")]
    Missing,

    #[error("part id is invalid: {0}")]
    Invalid(#[from] uuid::Error),
}

#[derive(Debug, Error)]
pub(crate) enum MeasurementError {
    #[error("meaurement must be finite and positive")]
    Invalid,
}

#[derive(Debug, Error)]
pub(crate) enum PartCategoryError {
    #[error("unknown value `{0}`")]
    Unknown(String),
}

#[derive(Debug, Error)]
pub(crate) enum MoneyCentsError {
    #[error("value must be non-negative: {0}")]
    Negative(i64),
}

#[derive(Debug, Error)]
pub(crate) enum StockQuantityError {
    #[error("value must be non-negative: {0}")]
    Negative(i64),
}

#[derive(Debug, thiserror::Error)]
pub(crate) enum InventoryUseCaseError {
    #[error("inventory storage failed: {0}")]
    Storage(#[from] PartRepoError),
}
