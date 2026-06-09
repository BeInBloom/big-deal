use tonic::{Request, Response, Status};

use crate::{
    domain::{
        errors::InventoryUseCaseError,
        models::{GetPartQuery, ListPartsQuery},
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

#[cfg(test)]
mod tests {
    use crate::{
        domain::{
            errors::{InventoryUseCaseError, PartRepoError},
            traits::MockInventoryUseCases,
        },
        proto::inventory_v1,
    };

    use super::*;

    #[tokio::test]
    async fn get_part_returns_invalid_argument_for_invalid_request() {
        let use_cases = MockInventoryUseCases::new();
        let handler = InventoryGrpcHandler::new(use_cases);

        let req = Request::new(GetPartRequest {
            uuid: String::new(),
        });

        let err = handler.get_part(req).await.unwrap_err();
        assert_eq!(err.code(), tonic::Code::InvalidArgument);
    }

    #[tokio::test]
    async fn get_part_returns_not_found_when_part_is_missing() {
        let raw_uuid = "11111111-1111-4111-8111-111111111111".to_string();
        let part_id = raw_uuid.parse().unwrap();

        let mut use_cases = MockInventoryUseCases::new();
        use_cases
            .expect_get_part()
            .with(mockall::predicate::eq(GetPartQuery { id: part_id }))
            .times(1)
            .returning(|_| Box::pin(std::future::ready(Ok(None))));

        let handler = InventoryGrpcHandler::new(use_cases);

        let req = Request::new(GetPartRequest { uuid: raw_uuid });

        let err = handler.get_part(req).await.unwrap_err();
        assert_eq!(err.code(), tonic::Code::NotFound);
    }

    #[tokio::test]
    async fn get_part_returns_internal_when_use_case_fails() {
        let raw_uuid = "11111111-1111-4111-8111-111111111111".to_string();
        let part_id = raw_uuid.parse().unwrap();

        let mut use_cases = MockInventoryUseCases::new();
        use_cases
            .expect_get_part()
            .with(mockall::predicate::eq(GetPartQuery { id: part_id }))
            .times(1)
            .returning(|_| {
                Box::pin(std::future::ready(Err(InventoryUseCaseError::Storage(
                    PartRepoError::Storage,
                ))))
            });

        let handler = InventoryGrpcHandler::new(use_cases);

        let req = Request::new(GetPartRequest { uuid: raw_uuid });

        let err = handler.get_part(req).await.unwrap_err();
        assert_eq!(err.code(), tonic::Code::Internal);
    }

    #[tokio::test]
    async fn list_parts_returns_empty_list() {
        let query = ListPartsQuery::default();

        let mut use_cases = MockInventoryUseCases::new();
        use_cases
            .expect_list_parts()
            .with(mockall::predicate::eq(query))
            .times(1)
            .returning(|_| Box::pin(std::future::ready(Ok(Vec::new()))));

        let handler = InventoryGrpcHandler::new(use_cases);

        let req = Request::new(ListPartsRequest { filter: None });

        let res = handler.list_parts(req).await.unwrap().into_inner();
        assert!(res.parts.is_empty());
    }

    #[tokio::test]
    async fn list_parts_returns_invalid_argument_for_invalid_request() {
        let use_cases = MockInventoryUseCases::new();
        let handler = InventoryGrpcHandler::new(use_cases);

        let req = Request::new(ListPartsRequest {
            filter: Some(inventory_v1::InventoryPartsFilter {
                uuids: vec!["not-a-uuid".to_string()],
                names: Vec::new(),
                categories: Vec::new(),
                manufacturer_countries: Vec::new(),
                tags: Vec::new(),
            }),
        });

        let err = handler.list_parts(req).await.unwrap_err();
        assert_eq!(err.code(), tonic::Code::InvalidArgument);
    }

    #[tokio::test]
    async fn list_parts_returns_internal_when_use_case_fails() {
        let query = ListPartsQuery::default();

        let mut use_cases = MockInventoryUseCases::new();
        use_cases
            .expect_list_parts()
            .with(mockall::predicate::eq(query))
            .times(1)
            .returning(|_| {
                Box::pin(std::future::ready(Err(InventoryUseCaseError::Storage(
                    PartRepoError::Storage,
                ))))
            });

        let handler = InventoryGrpcHandler::new(use_cases);

        let req = Request::new(ListPartsRequest { filter: None });

        let err = handler.list_parts(req).await.unwrap_err();
        assert_eq!(err.code(), tonic::Code::Internal);
    }
}
