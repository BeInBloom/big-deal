use std::collections::HashMap;

use tokio::sync::RwLock;

use crate::domain::{
    errors::PartRepoError,
    models::{ListPartsQuery, Part, PartId},
    traits::PartRepo,
};

pub(crate) struct MapPartRepo {
    map: RwLock<HashMap<PartId, Part>>,
}

impl MapPartRepo {
    pub(crate) fn new() -> Self {
        let map = RwLock::new(HashMap::new());
        Self { map }
    }
}

impl PartRepo for MapPartRepo {
    async fn get(&self, id: PartId) -> Result<Option<Part>, PartRepoError> {
        Ok(self.map.read().await.get(&id).cloned())
    }

    async fn list(&self, query: ListPartsQuery) -> Result<Vec<Part>, PartRepoError> {
        Ok(self
            .map
            .read()
            .await
            .values()
            .filter(|part| query.matches(part))
            .cloned()
            .collect())
    }
}
