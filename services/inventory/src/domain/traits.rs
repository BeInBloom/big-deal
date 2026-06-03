#![allow(dead_code)]
use crate::domain::{
    errors::PartRepoError,
    models::{GetPartQuery, ListPartsQuery, Part, PartId},
};

pub(crate) trait PartRepo: Send + Sync + 'static {
    async fn get(&self, id: PartId) -> Result<Option<Part>, PartRepoError>;
    async fn list(&self, query: ListPartsQuery) -> Result<Vec<Part>, PartRepoError>;
}

pub(crate) trait InventoryUseCases: Send + Sync + 'static {
    async fn get_part(&self, query: GetPartQuery) -> Result<Option<Part>, PartRepoError>;
    async fn list_parts(&self, query: ListPartsQuery) -> Result<Vec<Part>, PartRepoError>;
}
