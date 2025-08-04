


insert into role_actions(role_id,action_id) values(9,1);
insert into role_actions(role_id,action_id) values(9,2);
insert into role_actions(role_id,action_id) values(9,3);
insert into role_actions(role_id,action_id) values(9,3);

INSERT INTO company_entities (company_id, entity_id) VALUES(16, 25);

INSERT INTO company_entities (company_id, entity_id) VALUES
    (1, 1),
    (1, 2),
    (1, 3),
    (1, 4),
    (1, 5),
    (1, 6),
    (1, 7),
    (1, 8),
    (1, 9),
    (1, 10),
    (1, 11),
    (1, 12),
    (1, 13),
    (1, 14),
    (1, 15),
    (1, 16),
    (1, 17),
    (1, 18),
    (1, 19),
    (1, 20),
    (1, 21);



-- company
insert into actions(name,entity_id) values('view',1);
insert into actions(name,entity_id) values('create',1);
insert into actions(name,entity_id) values('edit',1);

-- users
insert into actions(name,entity_id) values('view',11);
insert into actions(name,entity_id) values('create',11);
insert into actions(name,entity_id) values('edit',11);
insert into actions(name,entity_id) values('delete',11);

-- item
insert into actions(name,entity_id) values('view',2);
insert into actions(name,entity_id) values('create',2);
insert into actions(name,entity_id) values('edit',2);
insert into actions(name,entity_id) values('delete',2);

-- item-price
insert into actions(name,entity_id) values('view',3);
insert into actions(name,entity_id) values('create',3);
insert into actions(name,entity_id) values('edit',3);
insert into actions(name,entity_id) values('delete',3);

-- item-group
insert into actions(name,entity_id) values('view',4);
insert into actions(name,entity_id) values('create',4);
insert into actions(name,entity_id) values('edit',4);
insert into actions(name,entity_id) values('delete',4);

-- item-stock
insert into actions(name,entity_id) values('view',5);
insert into actions(name,entity_id) values('create',5);
insert into actions(name,entity_id) values('edit',5);
insert into actions(name,entity_id) values('delete',5);

-- item-attributes
insert into actions(name,entity_id) values('view',6);
insert into actions(name,entity_id) values('create',6);
insert into actions(name,entity_id) values('edit',6);
insert into actions(name,entity_id) values('delete',6);

-- item-warehouse
insert into actions(name,entity_id) values('view',7);
insert into actions(name,entity_id) values('create',7);
insert into actions(name,entity_id) values('edit',7);
insert into actions(name,entity_id) values('delete',7);

-- item-tax
insert into actions(name,entity_id) values('view',8);
insert into actions(name,entity_id) values('create',8);
insert into actions(name,entity_id) values('edit',8);
insert into actions(name,entity_id) values('delete',8);

-- price-list
insert into actions(name,entity_id) values('view',9);
insert into actions(name,entity_id) values('create',9);
insert into actions(name,entity_id) values('edit',9);
insert into actions(name,entity_id) values('delete',9);

-- supplier
insert into actions(name,entity_id) values('view',12);
insert into actions(name,entity_id) values('create',12);
insert into actions(name,entity_id) values('edit',12);
insert into actions(name,entity_id) values('delete',12);

-- purchase order
insert into actions(name,entity_id) values('view',14);
insert into actions(name,entity_id) values('create',14);
insert into actions(name,entity_id) values('edit',14);
insert into actions(name,entity_id) values('delete',14);

-- customer
insert into actions(name,entity_id) values('view',15);
insert into actions(name,entity_id) values('create',15);
insert into actions(name,entity_id) values('edit',15);
insert into actions(name,entity_id) values('delete',15);

-- address
insert into actions(name,entity_id) values('view',16);
insert into actions(name,entity_id) values('create',16);
insert into actions(name,entity_id) values('edit',16);
insert into actions(name,entity_id) values('delete',16);

-- contact
insert into actions(name,entity_id) values('view',17);
insert into actions(name,entity_id) values('create',17);
insert into actions(name,entity_id) values('edit',17);
insert into actions(name,entity_id) values('delete',17);

-- purchase-invoice
insert into actions(name,entity_id) values('view',18);
insert into actions(name,entity_id) values('create',18);
insert into actions(name,entity_id) values('edit',18);
insert into actions(name,entity_id) values('delete',18);

-- payment
insert into actions(name,entity_id) values('view',19);
insert into actions(name,entity_id) values('create',19);
insert into actions(name,entity_id) values('edit',19);
insert into actions(name,entity_id) values('delete',19);

-- ledger 
insert into actions(name,entity_id) values('view',20);
insert into actions(name,entity_id) values('create',20);
insert into actions(name,entity_id) values('edit',20);
insert into actions(name,entity_id) values('delete',20);

-- receipt

insert into actions(name,entity_id) values('view',21);
insert into actions(name,entity_id) values('create',21);
insert into actions(name,entity_id) values('edit',21);
insert into actions(name,entity_id) values('delete',21);

--court (Regate)
insert into actions(name,entity_id) values('view',22);
insert into actions(name,entity_id) values('create',22);
insert into actions(name,entity_id) values('edit',22);
insert into actions(name,entity_id) values('delete',22);

-- booking (Regate)
insert into actions(name,entity_id) values('view',23);
insert into actions(name,entity_id) values('create',23);
insert into actions(name,entity_id) values('edit',23);
insert into actions(name,entity_id) values('delete',23);

-- booking (Reserva)
insert into actions(name,entity_id) values('view',24);
insert into actions(name,entity_id) values('create',24);
insert into actions(name,entity_id) values('edit',24);
insert into actions(name,entity_id) values('delete',24);

--sale order (Sale Order)
insert into actions(name,entity_id) values('view',25);
insert into actions(name,entity_id) values('create',25);
insert into actions(name,entity_id) values('edit',25);
insert into actions(name,entity_id) values('delete',25);

--sale order (Sale Invoice)
insert into actions(name,entity_id) values('view',26);
insert into actions(name,entity_id) values('create',26);
insert into actions(name,entity_id) values('edit',26);
insert into actions(name,entity_id) values('delete',26);


--piano forms (Piano Forms)
insert into actions(name,entity_id) values('view',27);
insert into actions(name,entity_id) values('create',27);
insert into actions(name,entity_id) values('edit',27);
insert into actions(name,entity_id) values('delete',27);

--regatechart (Regate Chart)
insert into actions(name,entity_id) values('view',28);


-- get parent company
WITH RECURSIVE companies_cte(id,parent_id,depth) as (
	SELECT 
    	companies.id, 
    	companies.parent_id, 
    	2
        FROM companies
        WHERE companies.id = 7
	UNION ALL 
	SELECT 
	    companies.id, 
    	companies.parent_id, 
    	depth - 1
        FROM companies_cte,
             companies
        WHERE companies.id = companies_cte.parent_id
)
SELECT *
FROM companies_cte;

-- get-group-decendents
WITH RECURSIVE groups_cte(id,parent_id,uuid,parent_uuid,name,is_group,enabled,depth) as (
	SELECT 
    	groups.id, 
    	groups.parent_id,
    	groups.uuid,
		null,
    	groups.name,
    	groups.is_group, 
    	groups.enabled, 
    	0
        FROM groups
        WHERE id = 2
	UNION ALL 
	SELECT 
	    groups.id, 
    	groups.parent_id, 
		groups.uuid,
    	(select uuid from groups where id = groups_cte.id),
    	groups.name, 
    	groups.is_group, 
    	groups.enabled, 
    	depth +1
        FROM groups_cte,
             groups
        WHERE groups.parent_id = groups_cte.id
)
SELECT g.uuid,g.parent_uuid,g.name,g.is_group,g.enabled,g.depth
FROM groups_cte as g;



