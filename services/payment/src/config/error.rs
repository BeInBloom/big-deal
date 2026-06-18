use std::net::AddrParseError;

use thiserror::Error;

#[derive(Debug, Error)]
pub(crate) enum ConfigError {
    #[error("Addres is invalid: {0}")]
    Addr(#[from] AddrParseError),
}
