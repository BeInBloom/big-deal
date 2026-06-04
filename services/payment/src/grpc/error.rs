use tonic::Status;

use crate::domain::error::PaymentError;

impl From<PaymentError> for Status {
    fn from(error: PaymentError) -> Self {
        Status::invalid_argument(error.to_string())
    }
}
