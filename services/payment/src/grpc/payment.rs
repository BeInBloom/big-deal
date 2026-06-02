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
