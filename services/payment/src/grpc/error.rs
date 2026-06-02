use tonic::Status;

use crate::domain::error::PaymentError;

impl From<PaymentError> for Status {
    fn from(error: PaymentError) -> Self {
        match error {
            PaymentError::MissingUserId
            | PaymentError::InvalidUserId
            | PaymentError::MissingOrderId
            | PaymentError::InvalidOrderId
            | PaymentError::MissingPaymentMethod
            | PaymentError::InvalidPaymentMethod => Status::invalid_argument(error.to_string()),
        }
    }
}
