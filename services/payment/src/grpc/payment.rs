use tonic::{Request, Response, Status};

use crate::{
    proto::payment_v1::{
        PayOrderRequest, PayOrderResponse, payment_service_server::PaymentService,
    },
    service::payment::PaymentProcessor,
};

#[derive(Debug, Default)]
pub(crate) struct PaymentGrpcHandler {
    payment_processor: PaymentProcessor,
}

impl PaymentGrpcHandler {
    pub(crate) fn new() -> Self {
        let payment_processor = PaymentProcessor::new();
        Self { payment_processor }
    }
}

#[tonic::async_trait]
impl PaymentService for PaymentGrpcHandler {
    async fn pay_order(
        &self,
        request: Request<PayOrderRequest>,
    ) -> Result<Response<PayOrderResponse>, Status> {
        let command = request.into_inner().try_into()?;
        let transaction_id = self.payment_processor.pay(command);
        let response = PayOrderResponse {
            transaction_uuid: transaction_id.to_string(),
        };
        Ok(Response::new(response))
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::proto::payment_v1::PaymentMethod;

    #[tokio::test]
    async fn payorder_reutrns_transaction_uuid() {
        let handler = PaymentGrpcHandler::new();

        let user_uuid = uuid::Uuid::new_v4().to_string();
        let order_uuid = uuid::Uuid::new_v4().to_string();

        let req = Request::new(PayOrderRequest {
            user_uuid,
            order_uuid,
            payment_method: PaymentMethod::Card as i32,
        });

        let res = handler.pay_order(req).await.unwrap().into_inner();
        let transaction_uuid = uuid::Uuid::parse_str(&res.transaction_uuid);

        assert!(transaction_uuid.is_ok());
    }

    #[tokio::test]
    async fn pay_order_returns_invalid_argument_for_invalid_request() {
        let handler = PaymentGrpcHandler::new();

        let request = Request::new(PayOrderRequest {
            user_uuid: String::new(),
            order_uuid: uuid::Uuid::new_v4().to_string(),
            payment_method: PaymentMethod::Card as i32,
        });

        let err = handler.pay_order(request).await.unwrap_err();

        assert_eq!(err.code(), tonic::Code::InvalidArgument);
    }
}
