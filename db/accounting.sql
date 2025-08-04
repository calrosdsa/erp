create table if not exists payments (
    id bigint not null,
    uuid uuid not null default uuid_generate_v4(),
    created_at timestamp not null default now(),
    deleted_at timestamp,
    updated_at timestamp,
    posting_date date not null,
    payment_type text not null,
    mode_of_payment text not null,
    amount integer not null
);

create table if not exists party_payments (
    party_id bigint not null,
    payment_id bigint not null,
    party_bank_account text,
    company_bank_account text,
    contact_id bigint,
);


create table if not exists mode_of_payments(
    id bigint not null,
    uuid uuid not null default uuid_generate_v4(),
    name text not null,
    enabled boolean not null default false,
    
)


WITH invoice_data AS (
    -- Materialize the invoice lookup to avoid repeating it per row in transaction_ledgers
    SELECT id, code
    FROM invoices
    WHERE code IN (SELECT DISTINCT voucher_code FROM transaction_ledgers WHERE party_id = any(array[14,1]))
)
SELECT 
    tx.created_at as posting_date,
	('supplier') as party_type,
	suppl.name as party_name, 
	suppl.uuid as party_uuid,
	lg.name as receivable_account,
	lg.uuid as receivable_account_uuid,
	tx.voucher_type,
	tx.voucher_code as voucher_no,
    coalesce(t.total, 0) as invoiced_amount,
    coalesce(t.paid,tx.balance) as paid_amount
FROM 
    transaction_ledgers AS tx
LEFT JOIN LATERAL (
    SELECT 
        SUM(ii.paid_amount) AS paid,
        SUM(il.quantity * il.rate) AS total
    FROM 
        invoiced_item_lines ii
    JOIN 
        item_lines il ON ii.item_line = il.id
    WHERE 
        ii.invoice_id = (SELECT id FROM invoice_data WHERE code = tx.voucher_code)
) AS t ON true
join suppliers as suppl on suppl.id = tx.party_id
join ledgers as lg on lg.id = tx.ledger
WHERE 
    tx.party_id = any(array[14,1]);