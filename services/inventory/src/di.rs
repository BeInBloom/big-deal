use std::net::SocketAddr;

use crate::{
    app::App,
    config::models::{Config, MongoConfig, RepoConfig},
    domain::traits::PartRepo,
    grpc::inventory::InventoryGrpcHandler,
    repo::{
        map::map_repo::MapPartRepo,
        mongo::{models::PartDocument, repo::MongoPartRepo},
    },
    service::inventory::InventoryManager,
};

pub(crate) async fn run(config: Config) -> anyhow::Result<()> {
    let addr = config.grpc.addr;

    match config.repo {
        RepoConfig::HashMap => run_app(addr, MapPartRepo::new()).await,
        RepoConfig::Mongo(mongo_config) => {
            let repo = mongo_repo_from_config(mongo_config).await?;
            run_app(addr, repo).await
        }
    }
}

async fn run_app<R>(addr: SocketAddr, repo: R) -> anyhow::Result<()>
where
    R: PartRepo,
{
    let use_cases = InventoryManager::new(repo);
    let service = InventoryGrpcHandler::new(use_cases);
    App::new(addr, service).run().await
}

async fn mongo_repo_from_config(config: MongoConfig) -> anyhow::Result<MongoPartRepo> {
    let client = mongodb::Client::with_uri_str(config.url.as_str()).await?;
    let collection = client
        .database(&config.database)
        .collection::<PartDocument>(&config.parts_collection);

    Ok(MongoPartRepo::new(collection))
}
