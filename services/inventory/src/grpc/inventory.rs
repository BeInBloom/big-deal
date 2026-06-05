use tonic::{Request, Response, Status};

use crate::{
    domain::{
        errors::InventoryUseCaseError,
        models::{GetPartQuery, ListPartsQuery, Part},
        traits::InventoryUseCases,
    },
    grpc::error::InventoryRequestError,
    proto::inventory_v1::{
        GetPartRequest, GetPartResponse, InventoryPart, ListPartsRequest, ListPartsResponse,
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

    async fn get_part_or_not_found(&self, query: GetPartQuery) -> Result<InventoryPart, Status> {
        Ok(self
            .inventory_manager
            .get_part(query)
            .await
            .map_err(map_use_case_error)?
            .ok_or_else(|| Status::not_found("part not found"))?
            .into())
    }

    async fn get_list_parts(&self, query: ListPartsQuery) -> Result<Vec<InventoryPart>, Status> {
        Ok(self
            .inventory_manager
            .list_parts(query)
            .await
            .map_err(map_use_case_error)?
            .into_iter()
            .map(InventoryPart::from)
            .collect())
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
        let query = request.into_inner().try_into().map_err(map_request_error)?;
        let part = self.get_part_or_not_found(query).await?;
        Ok(Response::new(GetPartResponse { part: Some(part) }))
    }

    async fn list_parts(
        &self,
        request: Request<ListPartsRequest>,
    ) -> Result<Response<ListPartsResponse>, Status> {
        let query = request.into_inner().try_into().map_err(map_request_error)?;
        let parts = self.get_list_parts(query).await?;
        Ok(Response::new(ListPartsResponse { parts }))
    }
}

fn map_request_error(error: InventoryRequestError) -> Status {
    Status::invalid_argument(error.to_string())
}

fn map_use_case_error(error: InventoryUseCaseError) -> Status {
    match error {
        InventoryUseCaseError::Storage(_) => Status::internal("inventory storage failed"),
    }
}
