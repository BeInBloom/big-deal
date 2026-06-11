use mongodb::Collection;

use crate::{
    domain::{
        errors::PartRepoError,
        models::{ListPartsQuery, Part, PartId},
        traits::PartRepo,
    },
    repo::mongo::{error::MongoPartRepoError, models::PartDocument, query::IntoFilter},
};

pub(crate) struct MongoPartRepo {
    collection: Collection<PartDocument>,
}

impl MongoPartRepo {
    pub(crate) fn new(collection: Collection<PartDocument>) -> Self {
        Self { collection }
    }

    async fn try_get(&self, id: PartId) -> Result<Option<Part>, MongoPartRepoError> {
        let filter = id.into_filter()?.into_inner();

        Ok(self
            .collection
            .find_one(filter)
            .await?
            .map(Part::try_from)
            .transpose()?)
    }

    async fn try_list(&self, query: ListPartsQuery) -> Result<Vec<Part>, MongoPartRepoError> {
        let filter = query.into_filter()?.into_inner();
        let mut cursor = self.collection.find(filter).await?;
        let mut parts = Vec::new();

        while cursor.advance().await? {
            let doc = cursor.deserialize_current()?;
            parts.push(doc.try_into()?);
        }

        Ok(parts)
    }
}

impl PartRepo for MongoPartRepo {
    async fn get(&self, id: PartId) -> Result<Option<Part>, PartRepoError> {
        Ok(self.try_get(id).await?)
    }

    async fn list(&self, query: ListPartsQuery) -> Result<Vec<Part>, PartRepoError> {
        Ok(self.try_list(query).await?)
    }
}
