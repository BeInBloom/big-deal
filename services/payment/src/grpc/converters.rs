use crate::{
    domain::{
        error::PaymentError,
        models::{PayOrderCommand, PaymentMethod},
    },
    proto::payment_v1::{self, PayOrderRequest},
};

impl TryFrom<PayOrderRequest> for PayOrderCommand {
    type Error = PaymentError;

    fn try_from(value: PayOrderRequest) -> Result<Self, Self::Error> {
        let payment_method = payment_v1::PaymentMethod::try_from(value.payment_method)?;

        Ok(Self {
            user_id: value.user_uuid.parse()?,
            order_id: value.order_uuid.parse()?,
            payment_method: payment_method.try_into()?,
        })
    }
}

impl TryFrom<payment_v1::PaymentMethod> for PaymentMethod {
    type Error = PaymentError;

    fn try_from(value: payment_v1::PaymentMethod) -> Result<Self, Self::Error> {
        match value {
            payment_v1::PaymentMethod::Card => Ok(Self::Card),
            payment_v1::PaymentMethod::Sbp => Ok(Self::Sbp),
            payment_v1::PaymentMethod::CreditCard => Ok(Self::CreditCard),
            payment_v1::PaymentMethod::InvestorMoney => Ok(Self::InvestorMoney),
            payment_v1::PaymentMethod::Unspecified => Err(PaymentError::MissingPaymentMethod),
        }
    }
}
