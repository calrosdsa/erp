package dto

type (
	PaginationResult[T any] struct {
		Results       T                 `json:"results"`
		Total         int64             `json:"total"`
		FilterOptions []FilterOptionDto `json:"filters"`
	}

	ResultEntity[T any] struct {
		Entity     T             `json:"entity"`
		Addresses  []AddressDto  `json:"addresses"`
		Contacts   []ContactDto  `json:"contacts"`
		Activities []ActivityDto `json:"activities"`
	}

	ResponseData[T any] struct {
		Body struct {
			Result            T                   `json:"result"`
			Actions           []ActionDto         `json:"actions"`
			Message           string              `json:"message"`
			AssociatedActions map[int][]ActionDto `json:"associated_actions"`
		}
	}

	ResponseDataList[T any] struct {
		Body struct {
			Result            T                   `json:"result"`
			Actions           []ActionDto         `json:"actions"`
			Message           string              `json:"message"`
			AssociatedActions map[int][]ActionDto `json:"associated_actions"`
			FilterOptions     []FilterOptionDto   `json:"filters"`
		}
	}

	DeleteRequest struct {
		ID        string `query:"id" required:"true"`
		PartyType string `query:"party_type" required:"false"`
		Type      string `query:" required:"false"`
	}

	AcceptLanguageHeader struct {
		AcceptLanguage string `header:"Accept-Language"`
	}

	RequestDataWithPartyType struct {
		PartyType string `path:"party" required:"true"`
	}

	RequestEntityWithParty struct {
		AuthParams
		//can be ID Uuid Code
		ID        string `path:"id" required:"tre"`
		PartyType string `query:"party" required:"true"`
	}

	RequestEntity struct {
		AuthParams
		OptionalQueryParams
		//can be ID Uuid Code
		ID string `path:"id"`
	}

	RequestData struct {
		ID string `query:"id"`
	}

	ExportDocumentRequest struct {
		Body ExportDocumentData
	}

	ExportDocumentData struct {
		PartyType string `json:"party_type" required:"false"`
		ID        string `json:"id"`
	}

	RequestEntityPartyWithEvent struct {
		AuthParams
		ID        string `path:"id"`
		Event     string `query:"event" required:"true"`
		PartyType string `query:"party" required:"true"`
	}

	UpdateStatusWithEvent struct {
		AuthParams
		Body struct {
			// PartyID represents a unique identifier for a party, which can be either a UUID or a code.
			PartyID      string  `json:"party_id"`
			PartyType    string  `json:"party_type" required:"false"`
			Events       []int32 `json:"events" minItems:"1"`
			CurrentState string  `json:"current_state" minLength:"1"`
		}
	}

	RequestEntityWithEvent struct {
		AuthParams
		ID    string `path:"id"`
		Event string `event:"event" required:"true"`
	}

	UuidRequest struct {
		Uuid string `path:"uuid"`
		AcceptLanguageHeader
		// AcceptLanguage string `header:"Accept-Language"`
	}

	UpsertRequestEntity[T any] struct {
		AuthParams
		Body struct {
			Entity T `json:"entity"`
		}
	}

	AuthParams struct {
		Auth string `header:"Authorization"`
		// ActiveCompany   string `header:"Active-Company"`
		UserSessionUuid string `header:"User-Session-Uuid"`
		// Role            string `header:"Role"`
	}

	// ClientData struct {
	// 	ClientUuid string `header:"Client-Uuid" required:"true"`
	// }

	ResponseMessage struct {
		Body struct {
			Message string            `json:"message"`
			Errors  map[string]string `json:"errors"`
		}
	}

	// ActiveCompanyHeader struct {
	// 	ActiveCompany string `header:"Active-Company"`
	// }

	PaginationParams struct {
		Page    string `query:"page" required:"false"`
		Size    string `query:"size" required:"true"`
		Enabled string `query:"enabled" required:"false"`
		Status  string `query:"status" required:"false"`

		IsGroup string `query:"is_group" required:"false"`
	}

	DefaultListParams struct {
		Size        string `query:"size" required:"false"`
		Status      string `query:"status" required:"false"`
		Order       string `query:"orientation" required:"false"`
		OrderColumn string `query:"column" required:"false"`
	}

	OptionalQueryParams struct {
		Query string `query:"query" required:"false"`

		//FOR ORDERING
		Order       string `query:"orientation" required:"false"`
		OrderColumn string `query:"column" required:"false"`
		// FOR PARENT ENTITIES
		FilterID string `query:"parentId" required:"false"`
	}

	RequestPaginationData struct {
		PaginationParams
		OptionalQueryParams
		AuthParams
	}

	// For groups
	RequestPaginationPartyData struct {
		PartyType string `path:"party" required:"true"`
		PaginationParams
		OptionalQueryParams
		AuthParams
	}

	PaginationResponse[T any] struct {
		Body struct {
			PaginationResult T           `json:"pagination_result"`
			Actions          []ActionDto `json:"actions"`
		}
	}

	EntityResponse[T any] struct {
		Body struct {
			Message string      `json:"message"`
			Result  T           `json:"result"`
			Actions []ActionDto `json:"actions"`
			//ASSOCIATED ACTIONS
			AssociatedActions map[int][]ActionDto `json:"associated_actions"`
		}
	}

	TreeEntryDto struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		UUID     string `json:"uuid"`
		IsGroup  bool   `json:"is_group"`
		ParentID int64  `json:"parent"`
	}

	ItemActionData struct {
		ID     int64  `json:"id" required:"true"`
		Action string `json:"action" required:"true"`
	}
)
