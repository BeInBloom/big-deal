use std::net::SocketAddr;

use url::Url;

use crate::config::error::ConfigError;

const DEFAULT_GRPC_ADDR: &str = "0.0.0.0:50052";

pub(in crate::config) struct RawConfig {
    pub(in crate::config) raw_grpc_config: RawGrpcConfig,
    pub(in crate::config) raw_repo_config: RawRepoConfig,
}

pub(in crate::config) struct RawGrpcConfig {
    pub(in crate::config) addr: String,
}

pub(in crate::config) struct RawRepoConfig {
    pub(in crate::config) kind: String,
    pub(in crate::config) raw_mongo_config: RawMongoConfig,
}

pub(in crate::config) struct RawMongoConfig {
    pub(in crate::config) url: String,
    pub(in crate::config) database: String,
    pub(in crate::config) parts_collection: String,
}

pub(crate) struct Config {
    pub(crate) grpc: GrpcConfig,
    pub(crate) repo: RepoConfig,
}

pub(crate) struct GrpcConfig {
    pub(crate) addr: SocketAddr,
}

pub(crate) enum RepoConfig {
    HashMap,
    Mongo(MongoConfig),
}

pub(crate) struct MongoConfig {
    pub(crate) url: Url,
    pub(crate) database: String,
    pub(crate) parts_collection: String,
}

impl TryFrom<RawConfig> for Config {
    type Error = ConfigError;

    fn try_from(value: RawConfig) -> Result<Self, Self::Error> {
        Ok(Self {
            grpc: value.raw_grpc_config.try_into()?,
            repo: value.raw_repo_config.try_into()?,
        })
    }
}

impl TryFrom<RawGrpcConfig> for GrpcConfig {
    type Error = ConfigError;

    fn try_from(value: RawGrpcConfig) -> Result<Self, Self::Error> {
        let addr = if value.addr.is_empty() {
            DEFAULT_GRPC_ADDR.to_string()
        } else {
            value.addr
        };

        Ok(Self {
            addr: addr.parse()?,
        })
    }
}

impl TryFrom<RawRepoConfig> for RepoConfig {
    type Error = ConfigError;

    fn try_from(value: RawRepoConfig) -> Result<Self, Self::Error> {
        match value.kind.as_str() {
            "hash_map" => Ok(Self::HashMap),
            "mongo" => Ok(Self::Mongo(value.raw_mongo_config.try_into()?)),
            kind => Err(ConfigError::UnknownRepoKind(kind.to_string())),
        }
    }
}

impl TryFrom<RawMongoConfig> for MongoConfig {
    type Error = ConfigError;

    fn try_from(value: RawMongoConfig) -> Result<Self, Self::Error> {
        Ok(Self {
            url: value.url.parse()?,
            database: non_empty(value.database)?,
            parts_collection: non_empty(value.parts_collection)?,
        })
    }
}

fn non_empty(value: String) -> Result<String, ConfigError> {
    if value.is_empty() {
        return Err(ConfigError::EmptyValue);
    }

    Ok(value)
}
