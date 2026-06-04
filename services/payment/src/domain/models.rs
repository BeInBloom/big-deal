use std::{fmt::Display, str::FromStr};

use uuid::Uuid;

use crate::domain::error::PaymentError;

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub(crate) struct TransactionId(Uuid);

impl TransactionId {
    pub(crate) fn new() -> Self {
        Self(Uuid::new_v4())
    }
}

impl Display for TransactionId {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        self.0.fmt(f)
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub(crate) struct UserId(Uuid);

impl UserId {
    pub(crate) fn as_uuid(&self) -> Uuid {
        self.0
    }
}

impl TryFrom<&str> for UserId {
    type Error = PaymentError;

    fn try_from(value: &str) -> Result<Self, Self::Error> {
        if value.is_empty() {
            return Err(PaymentError::MissingUserId);
        }

        let uuid = value.parse()?;

        Ok(Self(uuid))
    }
}

impl FromStr for UserId {
    type Err = PaymentError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        UserId::try_from(s)
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub(crate) struct OrderId(Uuid);

impl OrderId {
    pub(crate) fn as_uuid(&self) -> Uuid {
        self.0
    }
}

impl TryFrom<&str> for OrderId {
    type Error = PaymentError;

    fn try_from(value: &str) -> Result<Self, Self::Error> {
        if value.is_empty() {
            return Err(PaymentError::MissingOrderId);
        }

        let uuid = value.parse()?;

        Ok(Self(uuid))
    }
}

impl FromStr for OrderId {
    type Err = PaymentError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        OrderId::try_from(s)
    }
}

#[derive(Debug, Clone, Copy)]
pub(crate) enum PaymentMethod {
    Card,
    Sbp,
    CreditCard,
    InvestorMoney,
}

#[derive(Debug, Clone, Copy)]
pub(crate) struct PayOrderCommand {
    pub(crate) user_id: UserId,
    pub(crate) order_id: OrderId,
    pub(crate) payment_method: PaymentMethod,
}
