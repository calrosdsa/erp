package dtosquare

import (
	entitysquare "erp/internal/app/plugin/square/entitiy_square"
	squaretypes "erp/internal/app/plugin/square/square_types"
)

type SquareCatalogRequest struct {
	ItemGroupUuid string `path:"uuid"`
}

type SquareCatalogResponse struct {
	Body struct {
		Catalog squaretypes.RetrieveCatalogRequest `json:"catalog"`
		Objects []entitysquare.SquareObject        `json:"objects"`
	}
}

type SquareObjectRequest struct {
	ObjectId      string `path:"object_id"`
	ItemGroupUuid string `path:"uuid"`
}

type SquareObjectResponse struct {
	Body struct {
		PlanVariation squaretypes.RetrieveObjectRequest `json:"plan_variation"`
		SquareObject  entitysquare.SquareObject         `json:"square_object"`
	}
}
