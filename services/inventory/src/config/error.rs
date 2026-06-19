use std::net::AddrParseError;

use thiserror::Error;

#[derive(Debug, Error)]
pub(crate) enum ConfigError {
    #[error("Addr is invalid: {0}")]
    Addr(#[from] AddrParseError),

    #[error("MongoUrl is invalid: {0}")]
    MongoUrl(#[from] url::ParseError),

    #[error("Unknown repo kind: {0}")]
    UnknownRepoKind(String),

    #[error("Empty value")]
    EmptyValue,
}
