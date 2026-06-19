use std::env;

use crate::config::{
    error::ConfigError,
    models::{Config, RawConfig, RawGrpcConfig, RawMongoConfig, RawRepoConfig},
};

const ENV_GRPC_ADDR: &str = "INVENTORY_GRPC_ADDR";

const ENV_REPO_KIND: &str = "INVENTORY_REPO_KIND";
const ENV_REPO_MONGO_URL: &str = "INVENTORY_REPO_MONGO_URL";
const ENV_REPO_MONGO_DATABASE: &str = "INVENTORY_REPO_MONGO_DATABASE";
const ENV_REPO_MONGO_PARTS_COLLECTION: &str = "INVENTORY_REPO_MONGO_PARTS_COLLECTION";

pub(crate) fn config_from_env() -> Result<Config, ConfigError> {
    let raw_config = get_raw_config();
    raw_config.try_into()
}

fn get_raw_config() -> RawConfig {
    let raw_grpc_config = get_raw_grpc_config();
    let raw_repo_config = get_raw_repo_config();

    RawConfig {
        raw_grpc_config,
        raw_repo_config,
    }
}

fn get_raw_grpc_config() -> RawGrpcConfig {
    let addr = env::var(ENV_GRPC_ADDR).unwrap();
    RawGrpcConfig { addr }
}

fn get_raw_repo_config() -> RawRepoConfig {
    let kind = env::var(ENV_REPO_KIND).unwrap();
    let raw_mongo_config = get_raw_mongo_config();

    RawRepoConfig {
        kind,
        raw_mongo_config,
    }
}

fn get_raw_mongo_config() -> RawMongoConfig {
    RawMongoConfig {
        url: env::var(ENV_REPO_MONGO_URL).unwrap(),
        database: env::var(ENV_REPO_MONGO_DATABASE).unwrap(),
        parts_collection: env::var(ENV_REPO_MONGO_PARTS_COLLECTION).unwrap(),
    }
}
