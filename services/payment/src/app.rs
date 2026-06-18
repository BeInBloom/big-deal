use std::net::{Ipv4Addr, SocketAddr, SocketAddrV4};

use tonic::transport::Server;

use crate::{
    config::models::Config, grpc::payment::PaymentGrpcHandler,
    proto::payment_v1::payment_service_server::PaymentServiceServer,
};

const ADDR: SocketAddr = SocketAddr::V4(SocketAddrV4::new(Ipv4Addr::LOCALHOST, 50051));

pub(crate) struct App {
    addr: SocketAddr,
    payment: PaymentGrpcHandler,
}

impl App {
    pub(crate) fn new(config: Config) -> Self {
        Self {
            addr: config.grpc_config.addr,
            payment: PaymentGrpcHandler::new(),
        }
    }

    pub(crate) async fn run(self) -> anyhow::Result<()> {
        let payment_server = PaymentServiceServer::new(self.payment);

        Server::builder()
            .add_service(payment_server)
            .serve(self.addr)
            .await?;

        Ok(())
    }
}
