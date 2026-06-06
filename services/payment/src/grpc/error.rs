use tonic::Status;

use crate::domain::error::PayOrderCommandError;

impl From<PayOrderCommandError> for Status {
    fn from(error: PayOrderCommandError) -> Self {
        Status::invalid_argument(error.to_string())
    }
}
