use std::env;

use crate::config::{
    error::ConfigError,
    models::{Config, RawConfig, RawGrpcConfig},
};

const ENV_GRPC_ADDR: &str = "PAYMENT_GRPC_ADDR";
const DEFAULT_GRPC_ADDR: &str = "127.0.0.1:50051";

pub(crate) fn config_from_env() -> Result<Config, ConfigError> {
    let raw_config = get_raw_config();
    raw_config.try_into()
}

fn get_raw_config() -> RawConfig {
    let mut addr = env::var(ENV_GRPC_ADDR).unwrap_or(DEFAULT_GRPC_ADDR.to_string());

    if addr.is_empty() {
        addr = DEFAULT_GRPC_ADDR.to_string();
    }

    RawConfig {
        raw_grpc_config: RawGrpcConfig { addr },
    }
}
