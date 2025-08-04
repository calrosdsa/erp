package squaretypes

import "time"

type SquareCredentials struct {
	AccessToken   string `json:"accessToken"`
	ApplicationId string `json:"applicationId"`
	LocationId    string `json:"locationId"`
	ApiVersion    string `json:"apiVersion"`
}

type SubscriptionPlanResponse struct {
	CatalogObject struct {
		Type                  string    `json:"type"`
		ID                    string    `json:"id"`
		UpdatedAt             time.Time `json:"updated_at"`
		CreatedAt             time.Time `json:"created_at"`
		Version               int64     `json:"version"`
		IsDeleted             bool      `json:"is_deleted"`
		PresentAtAllLocations bool      `json:"present_at_all_locations"`
		SubscriptionPlanData  struct {
			Name     string `json:"name"`
			AllItems bool   `json:"all_items"`
		} `json:"subscription_plan_data"`
	} `json:"catalog_object"`
	IDMappings []struct {
		ClientObjectID string `json:"client_object_id"`
		ObjectID       string `json:"object_id"`
	} `json:"id_mappings"`
}

type CreateSubscriptionPlanResponse struct {
	CatalogObject struct {
		Type                          string    `json:"type"`
		ID                            string    `json:"id"`
		UpdatedAt                     time.Time `json:"updated_at"`
		CreatedAt                     time.Time `json:"created_at"`
		Version                       int64     `json:"version"`
		IsDeleted                     bool      `json:"is_deleted"`
		PresentAtAllLocations         bool      `json:"present_at_all_locations"`
		SubscriptionPlanVariationData struct {
			Name   string `json:"name"`
			Phases []struct {
				UID     string `json:"uid"`
				Cadence string `json:"cadence"`
				Periods int    `json:"periods"`
				Ordinal int    `json:"ordinal"`
				Pricing struct {
					Type  string `json:"type"`
					Price struct {
						Amount   int    `json:"amount"`
						Currency string `json:"currency"`
					} `json:"price"`
					PriceMoney struct {
						Amount   int    `json:"amount"`
						Currency string `json:"currency"`
					} `json:"price_money"`
				} `json:"pricing"`
			} `json:"phases"`
			SubscriptionPlanID string `json:"subscription_plan_id"`
		} `json:"subscription_plan_variation_data"`
	} `json:"catalog_object"`
	IDMappings []struct {
		ClientObjectID string `json:"client_object_id"`
		ObjectID       string `json:"object_id"`
	} `json:"id_mappings"`
}

type RetrieveCatalogRequest struct {
	Object struct {
		ObjectInfo
		SubscriptionPlanData struct {
			Name                       string `json:"name"`
			SubscriptionPlanVariations []SubscriptionPlanVariation `json:"subscription_plan_variations"`
			AllItems bool `json:"all_items"`
		} `json:"subscription_plan_data"`
	} `json:"object"`
}

type SubscriptionPlanVariation struct {
	ObjectInfo
	SubscriptionPlanVariationData struct {
		Name string `json:"name"`

		Phases []Phase `json:"phases"`

		SubscriptionPlanID string `json:"subscription_plan_id"`
	} `json:"subscription_plan_variation_data"`
}

