use prost::UnknownEnumValue;
use thiserror::Error;

#[derive(Debug, Error)]
pub(crate) enum UserIdError {
    #[error("user_uuid is required")]
    Missing,

    #[error("user_uuid is invalid: {0}")]
    Invalid(#[from] uuid::Error),
}

#[derive(Debug, Error)]
pub(crate) enum OrderIdError {
    #[error("order_uuid is required")]
    Missing,

    #[error("order_uuid is invalid: {0}")]
    Invalid(#[from] uuid::Error),
}

#[derive(Debug, Error, PartialEq, Eq)]
pub(crate) enum PaymentMethodError {
    #[error("payment_method is required")]
    Missing,

    #[error("payment_method is invalid: {0}")]
    Invalid(#[from] UnknownEnumValue),
}

#[derive(Debug, Error)]
pub(crate) enum PayOrderCommandError {
    #[error(transparent)]
    UserId(#[from] UserIdError),

    #[error(transparent)]
    OrderId(#[from] OrderIdError),

    #[error(transparent)]
    PaymentMethod(#[from] PaymentMethodError),
}
