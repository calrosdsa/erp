package internal

import "gorm.io/gen"

type Querier interface {
	// INSERT INTO party_references (party_id,reference_id) VALUES (@party_id,@reference_id)
	InsertReference(party_id int64,reference_id int64)(error)
	// INSERT INTO activities (party_id,profile_id,type,data) VALUES (@party_id,@profile_id,@typeActivity,@data)
	InsertActivity(party_id int64,profile_id int64,typeActivity string,data *string)(error)
	// INSERT INTO parties (party_type_code) VALUES (@party_type_code) returning id
	InsertParty(party_type_code string) (int64, error)
	// SELECT r.* FROM @@table as r
	// inner join parties as p on p.id = r.id
	//  {{where}}
	//    @@column=@value and
	//    p.party_type_code = @partyType
	//    {{if enabled != "" }}
	//      and enabled = @enabled
	//    {{end}}
	//    {{if query != "" }}
	//      and name ILIKE = %@query%
	//    {{end}}
	//    {{if orderColumn != "" && order != ""}}
	//       ORDER BY @orderColumn @order
	//    {{end}}
	// {{end}}
	// LIMIT @limit OFFSET @offset;
	PaginateParty(partyType, column string, value int64, query string, limit, offset int, orderColumn, order,enabled string) ([]gen.T, error)
	// SELECT count(*) FROM @@table as r
	// inner join parties as p on p.id = r.id
	//  {{where}}
	//    @@column=@value and
	//    p.party_type_code = @partyType
	//    {{if enabled != "" }}
	//      and enabled = @enabled
	//    {{end}}
	//    {{if query != "" }}
	//      and name ILIKE = %@query%
	//    {{end}}
	// {{end}}
	CountPaginateParty(partyType, column string, value int64, query,enabled string) (int64, error)

	// SELECT r.* FROM @@table as r
	//  {{where}}
	//   {{if column != "" }}
	//    @@column=@value
	//    {{end}}
	//    {{if query != "" }}
	//      and name ILIKE = %@query%
	//    {{end}}
	//    {{if enabled != "" }}
	//      and enabled = @enabled
	//    {{end}}
	//    {{if orderColumn != "" && order != ""}}
	//       ORDER BY @orderColumn @order
	//    {{end}}
	// {{end}}
	// LIMIT @limit OFFSET @offset;
	Paginate(column string, value int64, query string, limit, offset int, orderColumn, order,enabled string,) ([]gen.T, error)

	// SELECT count(*) FROM @@table as r
	//  {{where}}
	//   {{if column != "" }}
	//    @@column=@value
	//    {{end}}
	//    {{if query != "" }}
	//      and name ILIKE = %@query%
	//    {{end}}
	//    {{if enabled != "" }}
	//      and enabled = @enabled
	//    {{end}}
	// {{end}}
	CountPaginate(column string, value int64, query,enabled string) (int64, error)
}
