use anyhow::Result;
use tonic::{Request, Response, Status};

use crate::{
    domain::traits::InventoryUseCases,
    proto::inventory_v1::{
        GetPartRequest, GetPartResponse, ListPartsRequest, ListPartsResponse,
        inventory_service_server::InventoryService,
    },
};

pub(crate) struct InventoryGrpcHandler<R>
where
    R: InventoryUseCases,
{
    inventory_manager: R,
}

impl<M> InventoryGrpcHandler<M>
where
    M: InventoryUseCases,
{
    pub(crate) fn new(inventory_manager: M) -> Self {
        Self { inventory_manager }
    }
}

#[tonic::async_trait]
impl<R> InventoryService for InventoryGrpcHandler<R>
where
    R: InventoryUseCases,
{
    async fn get_part(
        &self,
        request: Request<GetPartRequest>,
    ) -> Result<Response<GetPartResponse>, Status> {
        todo!()
    }

    async fn list_parts(
        &self,
        request: Request<ListPartsRequest>,
    ) -> Result<Response<ListPartsResponse>, Status> {
        todo!()
    }
}
