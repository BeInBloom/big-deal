use crate::domain::models::{PayOrderCommand, TransactionId};

#[derive(Debug, Default)]
pub(crate) struct PaymentProcessor;

impl PaymentProcessor {
    pub(crate) fn new() -> Self {
        Self
    }

    pub(crate) fn pay(&self, command: PayOrderCommand) -> TransactionId {
        let _ = (
            command.user_id.as_uuid(),
            command.order_id.as_uuid(),
            command.payment_method,
        );

        let transaction_id = TransactionId::new();
        println!("{transaction_id}");
        transaction_id
    }
}
