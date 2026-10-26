use tonic::Status;

use crate::domain::error::{PayOrderCommandError, PaymentError};

impl From<PayOrderCommandError> for Status {
    fn from(error: PayOrderCommandError) -> Self {
        Status::invalid_argument(error.to_string())
    }
}

impl From<PaymentError> for Status {
    fn from(error: PaymentError) -> Self {
        match error {
            PaymentError::AlreadyPaidWithDifferentParameters => {
                Status::already_exists(error.to_string())
            }
        }
    }
}
