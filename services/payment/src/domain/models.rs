use std::{fmt::Display, str::FromStr};

use uuid::Uuid;

use crate::domain::error::{OrderIdError, UserIdError};

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
    type Error = UserIdError;

    fn try_from(value: &str) -> Result<Self, Self::Error> {
        if value.is_empty() {
            return Err(UserIdError::Missing);
        }

        let uuid = value.parse()?;

        Ok(Self(uuid))
    }
}

impl FromStr for UserId {
    type Err = UserIdError;

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
    type Error = OrderIdError;

    fn try_from(value: &str) -> Result<Self, Self::Error> {
        if value.is_empty() {
            return Err(OrderIdError::Missing);
        }

        let uuid = value.parse()?;

        Ok(Self(uuid))
    }
}

impl FromStr for OrderId {
    type Err = OrderIdError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        OrderId::try_from(s)
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
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

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn user_id_parses_uuid() {
        let raw_uuid = uuid::Uuid::new_v4().to_string();
        let user_id = UserId::try_from(raw_uuid.as_str()).unwrap();
        assert_eq!(user_id.as_uuid().to_string(), raw_uuid);
    }

    #[test]
    fn user_id_rejects_empty_value() {
        let err = UserId::try_from("").unwrap_err();
        assert!(matches!(err, UserIdError::Missing));
    }

    #[test]
    fn user_id_rejects_invalid_uuid() {
        let err = UserId::try_from("wanna snu snu?").unwrap_err();
        assert!(matches!(err, UserIdError::Invalid(_)));
    }

    #[test]
    fn order_id_parses_uuid() {
        let raw_uuid = uuid::Uuid::new_v4().to_string();
        let order_id = OrderId::try_from(raw_uuid.as_str()).unwrap();
        assert_eq!(order_id.as_uuid().to_string(), raw_uuid);
    }

    #[test]
    fn order_id_rejects_empty_value() {
        let err = OrderId::try_from("").unwrap_err();
        assert!(matches!(err, OrderIdError::Missing));
    }

    #[test]
    fn order_id_rejects_invalid_uuid() {
        let err = OrderId::try_from("lol kek cheburek").unwrap_err();
        assert!(matches!(err, OrderIdError::Invalid(_)));
    }
}
