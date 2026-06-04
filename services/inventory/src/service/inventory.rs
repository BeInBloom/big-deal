use anyhow::Result;

use crate::domain::{
    errors::InventoryUseCaseError,
    models::{GetPartQuery, ListPartsQuery, Part},
    traits::{InventoryUseCases, PartRepo},
};

pub(crate) struct InventoryManager<T> {
    part_repo: T,
}

impl<R> InventoryManager<R>
where
    R: PartRepo,
{
    pub(crate) fn new(part_repo: R) -> Self {
        Self { part_repo }
    }
}

impl<R> InventoryUseCases for InventoryManager<R>
where
    R: PartRepo,
{
    async fn get_part(&self, query: GetPartQuery) -> Result<Option<Part>, InventoryUseCaseError> {
        Ok(self.part_repo.get(query.id).await?)
    }

    async fn list_parts(&self, query: ListPartsQuery) -> Result<Vec<Part>, InventoryUseCaseError> {
        Ok(self.part_repo.list(query).await?)
    }
}
