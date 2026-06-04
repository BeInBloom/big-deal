use thiserror::Error;

#[allow(dead_code)]
#[derive(Debug, Error)]
pub(crate) enum PartRepoError {
    #[error("part repository failed")]
    Failed,
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

#[derive(Debug, thiserror::Error)]
pub(crate) enum InventoryUseCaseError {
    #[error("inventory storage failed: {0}")]
    Storage(#[from] PartRepoError),
}
