use crate::config::from_env::config_from_env;

mod app;
mod config;
mod di;
mod domain;
mod grpc;
mod proto;
mod repo;
mod service;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let config = config_from_env()?;
    di::run(config).await
}
