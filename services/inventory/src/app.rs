use std::net::SocketAddr;

use tonic::transport::Server;

use crate::proto::inventory_v1::inventory_service_server::{
    InventoryService, InventoryServiceServer,
};

pub(crate) struct App<T> {
    addr: SocketAddr,
    inventory: T,
}

impl<T> App<T>
where
    T: InventoryService,
{
    pub(crate) fn new(addr: SocketAddr, service: T) -> Self {
        Self {
            addr,
            inventory: service,
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
