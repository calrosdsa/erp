-- Asset
insert into ledgers(id,account_type,name,company_id,description) values(1,'Asset','Asset',1,'Group for asset accounts');
insert into ledgers(ledger_no,ledger_parent,account_type,name,company_id,description) values(1010,1,'Asset','Cash on hands',1,'Physical cash available.');
insert into ledgers(ledger_no,ledger_parent,account_type,name,company_id,description) values(1020,1,'Asset','Cash in bank',1,'Cash held in the business bank accounts.');
insert into ledgers(ledger_no,ledger_parent,account_type,name,company_id,description) values(1030,1,'Asset','Accounts receivable',1,'Amounts owed to the business by customers.');
insert into ledgers(ledger_no,ledger_parent,account_type,name,company_id,description) values(1040,1,'Asset','Inventory',1,'Value of goods available for sale.');

insert into ledger_accounts(ledger_no,can_debit,can_credit) values(1010,true,true);
insert into ledger_accounts(ledger_no,can_debit,can_credit) values(1020,true,true);
insert into ledger_accounts(ledger_no,can_debit,can_credit) values(1030,true,true);
insert into ledger_accounts(ledger_no,can_debit,can_credit) values(1040,true,true);


insert into ledgers(id,account_type,name,company_id,description) values(100,'Liabilities','Liabilities',1,'Group for Liabilities accounts');
insert into ledgers(ledger_no,ledger_parent,account_type,name,company_id,description) values(2010,100,'Liabilities','Accounts payable',1,'Amounts owed by the business to suppliers.');
insert into ledgers(ledger_no,ledger_parent,account_type,name,company_id,description) values(2020,100,'Liabilities','Short-term loans payable',1,'Any short-term loans the business has taken.');
insert into ledgers(ledger_no,ledger_parent,account_type,name,company_id,description) values(2030,100,'Liabilities','Credit card payable	',1,'Outstanding credit card balances.');

insert into ledger_accounts(ledger_no,can_debit,can_credit) values(2010,true,true);
insert into ledger_accounts(ledger_no,can_debit,can_credit) values(2020,true,true);
insert into ledger_accounts(ledger_no,can_debit,can_credit) values(2030,true,true);


insert into ledgers(id,account_type,name,company_id,description) values(101,'Equity','Equity',1,'Group for Equity accounts');
insert into ledgers(ledger_no,ledger_parent,account_type,name,company_id,description) values(3010,101,'Equity','Owners equity',1,'Owners investment in the business.');

insert into ledger_accounts(ledger_no,can_debit,can_credit) values(3010,true,true);


insert into ledgers(id,account_type,name,company_id,description) values(102,'Revenue','Revenue',1,'Group for Revenue accounts');
insert into ledgers(ledger_no,ledger_parent,account_type,name,company_id,description) values(4010,102,'Revenue','Product sales',1,'Income from selling products.');
insert into ledgers(ledger_no,ledger_parent,account_type,name,company_id,description) values(4020,102,'Revenue','Shipping revenue',1,'Income from shipping services.');

insert into ledger_accounts(ledger_no,can_debit,can_credit) values(4010,true,true);
insert into ledger_accounts(ledger_no,can_debit,can_credit) values(4020,true,true);

insert into ledgers(id,account_type,name,company_id,description) values(103,'Expenses','Expenses',1,'Group for Expenses accounts');
insert into ledgers(ledger_no,ledger_parent,account_type,name,company_id,description) values(5010,103,'Expenses','Cost of goodc sold',1,'Cost directly associated with producing goods.');
insert into ledgers(ledger_no,ledger_parent,account_type,name,company_id,description) values(5020,103,'Expenses','Shipping costs',1,'Costs associated with shipping.');
insert into ledgers(ledger_no,ledger_parent,account_type,name,company_id,description) values(5030,103,'Expenses','Website maintenance expenses	',1,'Expenses related to maintaining the online retail platform.');

insert into ledger_accounts(ledger_no,can_debit,can_credit) values(5010,true,true);
insert into ledger_accounts(ledger_no,can_debit,can_credit) values(5020,true,true);
insert into ledger_accounts(ledger_no,can_debit,can_credit) values(5030,true,true);