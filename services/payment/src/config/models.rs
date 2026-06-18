use std::net::SocketAddr;

use crate::config::error::ConfigError;

pub(crate) struct RawConfig {
    pub(crate) raw_grpc_config: RawGrpcConfig,
}

pub(crate) struct RawGrpcConfig {
    pub(crate) addr: String,
}

pub(crate) struct Config {
    pub(crate) grpc_config: GrpcConfig,
}

pub(crate) struct GrpcConfig {
    pub(crate) addr: SocketAddr,
}

impl TryInto<Config> for RawConfig {
    type Error = ConfigError;

    fn try_into(self) -> Result<Config, Self::Error> {
        Ok(Config {
            grpc_config: self.raw_grpc_config.try_into()?,
        })
    }
}

impl TryInto<GrpcConfig> for RawGrpcConfig {
    type Error = ConfigError;

    fn try_into(self) -> Result<GrpcConfig, Self::Error> {
        Ok(GrpcConfig {
            addr: self.addr.parse()?,
        })
    }
}
