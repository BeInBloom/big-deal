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

#[cfg(test)]
mod tests {
    use crate::domain::{
        errors::{InventoryUseCaseError, PartRepoError},
        traits::MockPartRepo,
    };

    use super::*;

    #[tokio::test]
    async fn get_part_delegates_to_repo() {
        let part_id = "11111111-1111-4111-8111-111111111111".parse().unwrap();

        let mut repo = MockPartRepo::new();
        repo.expect_get()
            .with(mockall::predicate::eq(part_id))
            .times(1)
            .returning(|_| Box::pin(std::future::ready(Ok(None))));

        let manager = InventoryManager::new(repo);

        let res = manager
            .get_part(GetPartQuery { id: part_id })
            .await
            .unwrap();

        assert!(res.is_none());
    }

    #[tokio::test]
    async fn get_part_maps_repo_error() {
        let part_id = "11111111-1111-4111-8111-111111111111".parse().unwrap();

        let mut repo = MockPartRepo::new();
        repo.expect_get()
            .with(mockall::predicate::eq(part_id))
            .times(1)
            .returning(|_| Box::pin(std::future::ready(Err(PartRepoError::Failed))));

        let manager = InventoryManager::new(repo);

        let err = manager
            .get_part(GetPartQuery { id: part_id })
            .await
            .unwrap_err();

        assert!(matches!(err, InventoryUseCaseError::Storage(_)));
    }

    #[tokio::test]
    async fn list_parts_delegates_to_repo() {
        let query = ListPartsQuery::default();

        let mut repo = MockPartRepo::new();
        repo.expect_list()
            .with(mockall::predicate::eq(query.clone()))
            .times(1)
            .returning(|_| Box::pin(std::future::ready(Ok(Vec::new()))));

        let manager = InventoryManager::new(repo);
        let parts = manager.list_parts(query).await.unwrap();
        assert!(parts.is_empty());
    }

    #[tokio::test]
    async fn list_parts_maps_repo_error() {
        let query = ListPartsQuery::default();

        let mut repo = MockPartRepo::new();
        repo.expect_list()
            .with(mockall::predicate::eq(query.clone()))
            .times(1)
            .returning(|_| Box::pin(std::future::ready(Err(PartRepoError::Failed))));

        let manager = InventoryManager::new(repo);
        let err = manager.list_parts(query).await.unwrap_err();
        assert!(matches!(err, InventoryUseCaseError::Storage(_)));
    }
}
