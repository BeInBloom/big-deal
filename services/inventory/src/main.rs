use crate::app::App;

mod app;
mod domain;
mod grpc;
mod proto;
mod service;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let app = App::new();
    app.run().await?;
    Ok(())
}
