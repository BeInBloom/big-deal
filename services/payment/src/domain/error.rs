use prost::UnknownEnumValue;
use thiserror::Error;

#[derive(Debug, Error)]
pub(crate) enum PaymentError {
    #[error("user_uuid is required")]
    MissingUserId,

    #[error("user_uuid is invalid")]
    InvalidUserId,

    #[error("order_uuid is required")]
    MissingOrderId,

    #[error("order_uuid is invalid: {0}")]
    InvalidOrderId(#[from] uuid::Error),

    #[error("payment_method is required")]
    MissingPaymentMethod,

    #[error("payment_method is invalid: {0}")]
    InvalidPaymentMethod(#[from] UnknownEnumValue),
}
