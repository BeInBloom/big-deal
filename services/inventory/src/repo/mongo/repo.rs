use mongodb::Collection;

use crate::repo::mongo::models::PartDocument;

pub(crate) struct MongoPartRepo {
    collection: Collection<PartDocument>,
}
