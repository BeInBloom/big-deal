use crate::{
    app::App, grpc::inventory::InventoryGrpcHandler, repo::map::map_repo::MapPartRepo,
    service::inventory::InventoryManager,
};

mod app;
mod domain;
mod grpc;
mod proto;
mod repo;
mod service;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let repo = MapPartRepo::new();
    let use_cases = InventoryManager::new(repo);
    let service = InventoryGrpcHandler::new(use_cases);
    let app = App::new(service);
    app.run().await?;
    Ok(())
}
