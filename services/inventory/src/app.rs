use std::net::{Ipv4Addr, SocketAddr, SocketAddrV4};

use tonic::transport::Server;

use crate::{
    grpc::inventory::InventoryGrpcHandler,
    proto::inventory_v1::inventory_service_server::InventoryServiceServer,
};

const ADDR: SocketAddr = SocketAddr::V4(SocketAddrV4::new(Ipv4Addr::LOCALHOST, 50052));

pub(crate) struct App {
    addr: SocketAddr,
    inventory: InventoryGrpcHandler,
}

impl App {
    pub(crate) fn new() -> Self {
        Self {
            addr: ADDR,
            inventory: InventoryGrpcHandler::new(),
        }
    }

    pub(crate) async fn run(self) -> anyhow::Result<()> {
        let inventory_server = InventoryServiceServer::new(self.inventory);

        Server::builder()
            .add_service(inventory_server)
            .serve(self.addr)
            .await?;

        Ok(())
    }
}
