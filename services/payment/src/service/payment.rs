use std::{
    collections::{HashMap, hash_map::Entry},
    sync::{Mutex, PoisonError},
};

use crate::domain::{
    error::PaymentError,
    models::{OrderId, PayOrderCommand, PaymentMethod, TransactionId, UserId},
};

#[derive(Debug, Clone, Copy)]
struct Payment {
    user_id: UserId,
    payment_method: PaymentMethod,
    transaction_id: TransactionId,
}

/// Payments are idempotent by order: a retry gets the original transaction.
///
/// Payments live in memory: the map is never pruned and is lost on restart,
/// so idempotency only holds for the lifetime of the process.
#[derive(Debug, Default)]
pub(crate) struct PaymentProcessor {
    payments: Mutex<HashMap<OrderId, Payment>>,
}

impl PaymentProcessor {
    pub(crate) fn new() -> Self {
        Self::default()
    }

    pub(crate) fn pay(&self, command: PayOrderCommand) -> Result<TransactionId, PaymentError> {
        let mut payments = self.payments.lock().unwrap_or_else(PoisonError::into_inner);

        match payments.entry(command.order_id) {
            Entry::Occupied(entry) => {
                let payment = entry.get();

                if payment.user_id != command.user_id
                    || payment.payment_method != command.payment_method
                {
                    return Err(PaymentError::AlreadyPaidWithDifferentParameters);
                }

                Ok(payment.transaction_id)
            }
            Entry::Vacant(entry) => {
                let transaction_id = TransactionId::new();

                entry.insert(Payment {
                    user_id: command.user_id,
                    payment_method: command.payment_method,
                    transaction_id,
                });

                Ok(transaction_id)
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn pay_returns_same_transaction_for_retry() {
        let processor = PaymentProcessor::new();
        let command = command(PaymentMethod::Card);

        let first = processor.pay(command).unwrap();
        let second = processor.pay(command).unwrap();

        assert_eq!(first, second);
    }

    #[test]
    fn pay_returns_different_transactions_for_different_orders() {
        let processor = PaymentProcessor::new();

        let first = processor.pay(command(PaymentMethod::Card)).unwrap();
        let second = processor.pay(command(PaymentMethod::Card)).unwrap();

        assert_ne!(first, second);
    }

    #[test]
    fn pay_rejects_retry_with_another_payment_method() {
        let processor = PaymentProcessor::new();
        let command = command(PaymentMethod::Card);
        processor.pay(command).unwrap();

        let err = processor
            .pay(PayOrderCommand {
                payment_method: PaymentMethod::Sbp,
                ..command
            })
            .unwrap_err();

        assert_eq!(err, PaymentError::AlreadyPaidWithDifferentParameters);
    }

    fn command(payment_method: PaymentMethod) -> PayOrderCommand {
        PayOrderCommand {
            user_id: uuid::Uuid::new_v4().to_string().parse().unwrap(),
            order_id: uuid::Uuid::new_v4().to_string().parse().unwrap(),
            payment_method,
        }
    }
}
