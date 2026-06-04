use crate::domain::{
    errors::{InventoryUseCaseError, PartRepoError},
    models::{GetPartQuery, ListPartsQuery, Part, PartId},
};

pub(crate) trait PartRepo: Send + Sync + 'static {
    fn get(&self, id: PartId) -> impl Future<Output = Result<Option<Part>, PartRepoError>> + Send;

    fn list(
        &self,
        query: ListPartsQuery,
    ) -> impl Future<Output = Result<Vec<Part>, PartRepoError>> + Send;
}

pub(crate) trait InventoryUseCases: Send + Sync + 'static {
    fn get_part(
        &self,
        query: GetPartQuery,
    ) -> impl Future<Output = Result<Option<Part>, InventoryUseCaseError>> + Send;

    fn list_parts(
        &self,
        query: ListPartsQuery,
    ) -> impl Future<Output = Result<Vec<Part>, InventoryUseCaseError>> + Send;
}
