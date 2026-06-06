use crate::{
    domain::{
        error::{PayOrderCommandError, PaymentMethodError},
        models::{PayOrderCommand, PaymentMethod},
    },
    proto::payment_v1::{self, PayOrderRequest},
};

impl TryFrom<PayOrderRequest> for PayOrderCommand {
    type Error = PayOrderCommandError;

    fn try_from(value: PayOrderRequest) -> Result<Self, Self::Error> {
        let payment_method = payment_v1::PaymentMethod::try_from(value.payment_method)
            .map_err(PaymentMethodError::Invalid)?;

        Ok(Self {
            user_id: value.user_uuid.parse()?,
            order_id: value.order_uuid.parse()?,
            payment_method: payment_method.try_into()?,
        })
    }
}

impl TryFrom<payment_v1::PaymentMethod> for PaymentMethod {
    type Error = PaymentMethodError;

    fn try_from(value: payment_v1::PaymentMethod) -> Result<Self, Self::Error> {
        match value {
            payment_v1::PaymentMethod::Card => Ok(Self::Card),
            payment_v1::PaymentMethod::Sbp => Ok(Self::Sbp),
            payment_v1::PaymentMethod::CreditCard => Ok(Self::CreditCard),
            payment_v1::PaymentMethod::InvestorMoney => Ok(Self::InvestorMoney),
            payment_v1::PaymentMethod::Unspecified => Err(PaymentMethodError::Missing),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::domain::error::OrderIdError;

    #[test]
    fn payment_method_maps_proto_values() {
        let cases = [
            (payment_v1::PaymentMethod::Card, PaymentMethod::Card),
            (payment_v1::PaymentMethod::Sbp, PaymentMethod::Sbp),
            (
                payment_v1::PaymentMethod::CreditCard,
                PaymentMethod::CreditCard,
            ),
            (
                payment_v1::PaymentMethod::InvestorMoney,
                PaymentMethod::InvestorMoney,
            ),
        ];

        for (proto_method, expected_method) in cases {
            let method = PaymentMethod::try_from(proto_method).unwrap();
            assert_eq!(method, expected_method);
        }
    }

    #[test]
    fn payment_method_rejects_unspecified() {
        let err = PaymentMethod::try_from(payment_v1::PaymentMethod::Unspecified).unwrap_err();
        assert!(matches!(err, PaymentMethodError::Missing))
    }

    #[test]
    fn pay_order_request_converts_to_command() {
        let user_uuid = uuid::Uuid::new_v4();
        let order_uuid = uuid::Uuid::new_v4();

        let req = PayOrderRequest {
            user_uuid: user_uuid.to_string(),
            order_uuid: order_uuid.to_string(),
            payment_method: payment_v1::PaymentMethod::Card as i32,
        };

        let command = PayOrderCommand::try_from(req).unwrap();

        assert_eq!(command.user_id.as_uuid(), user_uuid);
        assert_eq!(command.order_id.as_uuid(), order_uuid);
        assert_eq!(command.payment_method, PaymentMethod::Card);
    }

    #[test]
    fn pay_order_request_rejects_missing_user_id() {
        let order_uuid = uuid::Uuid::new_v4();

        let req = PayOrderRequest {
            user_uuid: String::new(),
            order_uuid: order_uuid.to_string(),
            payment_method: payment_v1::PaymentMethod::Card as i32,
        };

        let err = PayOrderCommand::try_from(req).unwrap_err();
        assert!(matches!(err, PayOrderCommandError::UserId(_)));
    }

    #[test]
    fn pay_order_request_rejects_missing_order_id() {
        let user_uuid = uuid::Uuid::new_v4();

        let request = PayOrderRequest {
            user_uuid: user_uuid.to_string(),
            order_uuid: String::new(),
            payment_method: payment_v1::PaymentMethod::Card as i32,
        };

        let err = PayOrderCommand::try_from(request).unwrap_err();
        assert!(matches!(
            err,
            PayOrderCommandError::OrderId(OrderIdError::Missing)
        ));
    }

    #[test]
    fn pay_order_request_rejects_invalid_order_id() {
        let user_uuid = uuid::Uuid::new_v4();

        let request = PayOrderRequest {
            user_uuid: user_uuid.to_string(),
            order_uuid: "best_kitty_ever".to_string(),
            payment_method: payment_v1::PaymentMethod::Card as i32,
        };

        let err = PayOrderCommand::try_from(request).unwrap_err();
        assert!(matches!(
            err,
            PayOrderCommandError::OrderId(OrderIdError::Invalid(_))
        ));
    }

    #[test]
    fn pay_order_request_rejects_missing_payment_method() {
        let user_uuid = uuid::Uuid::new_v4();
        let order_uuid = uuid::Uuid::new_v4();

        let request = PayOrderRequest {
            user_uuid: user_uuid.to_string(),
            order_uuid: order_uuid.to_string(),
            payment_method: payment_v1::PaymentMethod::Unspecified as i32,
        };

        let err = PayOrderCommand::try_from(request).unwrap_err();

        assert!(matches!(
            err,
            PayOrderCommandError::PaymentMethod(PaymentMethodError::Missing)
        ));
    }

    #[test]
    fn pay_order_request_rejects_invalid_payment_method() {
        let user_uuid = uuid::Uuid::new_v4();
        let order_uuid = uuid::Uuid::new_v4();

        let request = PayOrderRequest {
            user_uuid: user_uuid.to_string(),
            order_uuid: order_uuid.to_string(),
            payment_method: 999,
        };

        let err = PayOrderCommand::try_from(request).unwrap_err();

        assert!(matches!(
            err,
            PayOrderCommandError::PaymentMethod(PaymentMethodError::Invalid(_))
        ));
    }
}
