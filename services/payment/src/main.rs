mod app;
mod config;
mod domain;
mod grpc;
mod proto;
mod service;

use crate::{app::App, config::from_env::config_from_env};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let config = config_from_env()?;
    let app = App::new(config);
    app.run().await?;
    Ok(())
}
