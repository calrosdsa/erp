--
-- PostgreSQL database dump
--

-- Dumped from database version 17.2 (Debian 17.2-1.pgdg120+1)
-- Dumped by pg_dump version 17.2

-- Started on 2025-01-21 12:28:37

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 7 (class 2615 OID 2200)
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--

-- *not* creating schema, since initdb creates it


ALTER SCHEMA public OWNER TO postgres;

--
-- TOC entry 4592 (class 0 OID 0)
-- Dependencies: 7
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON SCHEMA public IS '';


--
-- TOC entry 2 (class 3079 OID 20529)
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- TOC entry 4594 (class 0 OID 0)
-- Dependencies: 2
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- TOC entry 3 (class 3079 OID 20566)
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- TOC entry 4595 (class 0 OID 0)
-- Dependencies: 3
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 219 (class 1259 OID 20577)
-- Name: account_settings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.account_settings (
    id integer NOT NULL,
    company_id bigint NOT NULL,
    bank_account bigint NOT NULL,
    cash_accunt bigint NOT NULL,
    payable_account bigint NOT NULL,
    cost_of_good_sold_account bigint NOT NULL,
    receivable_account bigint NOT NULL,
    income_account bigint NOT NULL
);


ALTER TABLE public.account_settings OWNER TO postgres;

--
-- TOC entry 220 (class 1259 OID 20580)
-- Name: account_settings_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.account_settings ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.account_settings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 221 (class 1259 OID 20581)
-- Name: account_statements; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.account_statements (
    account_no integer NOT NULL,
    date date NOT NULL,
    closing_balance date NOT NULL,
    total_credit integer,
    total_debit integer
);


ALTER TABLE public.account_statements OWNER TO postgres;

--
-- TOC entry 222 (class 1259 OID 20584)
-- Name: account_type_exts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.account_type_exts (
    account_type text NOT NULL,
    description text NOT NULL,
    transaction_type_code text NOT NULL
);


ALTER TABLE public.account_type_exts OWNER TO postgres;

--
-- TOC entry 223 (class 1259 OID 20589)
-- Name: account_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.account_types (
    account_type text NOT NULL,
    description text NOT NULL
);


ALTER TABLE public.account_types OWNER TO postgres;

--
-- TOC entry 224 (class 1259 OID 20594)
-- Name: accounts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.accounts (
    account_no integer NOT NULL,
    entity_type text NOT NULL,
    establishment_date date
);


ALTER TABLE public.accounts OWNER TO postgres;

--
-- TOC entry 225 (class 1259 OID 20599)
-- Name: actions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.actions (
    id bigint NOT NULL,
    name text NOT NULL,
    entity_id bigint NOT NULL
);


ALTER TABLE public.actions OWNER TO postgres;

--
-- TOC entry 226 (class 1259 OID 20604)
-- Name: actions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.actions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.actions_id_seq OWNER TO postgres;

--
-- TOC entry 4596 (class 0 OID 0)
-- Dependencies: 226
-- Name: actions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.actions_id_seq OWNED BY public.actions.id;


--
-- TOC entry 227 (class 1259 OID 20605)
-- Name: activities; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.activities (
    profile_id bigint NOT NULL,
    type text NOT NULL,
    party_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    comment text,
    deleted_at timestamp without time zone,
    id integer NOT NULL,
    arg1 text
);


ALTER TABLE public.activities OWNER TO postgres;

--
-- TOC entry 228 (class 1259 OID 20611)
-- Name: activities_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.activities ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.activities_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 229 (class 1259 OID 20612)
-- Name: addresses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.addresses (
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    updated_at timestamp with time zone DEFAULT now(),
    company text,
    street_line1 text NOT NULL,
    street_line2 text NOT NULL,
    city text NOT NULL,
    province text,
    postal_code text,
    phone_number text,
    country_code text,
    identification_number text,
    email text,
    title text NOT NULL,
    company_id bigint NOT NULL,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    is_shipping_address boolean DEFAULT false NOT NULL,
    is_billing_address boolean DEFAULT false NOT NULL,
    enabled boolean DEFAULT false NOT NULL
);


ALTER TABLE public.addresses OWNER TO postgres;

--
-- TOC entry 230 (class 1259 OID 20623)
-- Name: addresses_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.addresses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.addresses_id_seq OWNER TO postgres;

--
-- TOC entry 4597 (class 0 OID 0)
-- Dependencies: 230
-- Name: addresses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.addresses_id_seq OWNED BY public.addresses.id;


--
-- TOC entry 231 (class 1259 OID 20624)
-- Name: batch_bundles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.batch_bundles (
    id bigint NOT NULL,
    item_id bigint NOT NULL,
    created_at timestamp without time zone,
    company_id bigint NOT NULL,
    warehouse_id bigint NOT NULL,
    batch_bundle_no text NOT NULL,
    voucher_type text NOT NULL,
    voucher_code text NOT NULL,
    posting_date date NOT NULL,
    posting_time time without time zone NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.batch_bundles OWNER TO postgres;

--
-- TOC entry 232 (class 1259 OID 20629)
-- Name: charges_template; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.charges_template (
    id bigint NOT NULL,
    deleted_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name text NOT NULL,
    company_id bigint NOT NULL,
    status text DEFAULT 'ENABLED'::text NOT NULL
);


ALTER TABLE public.charges_template OWNER TO postgres;

--
-- TOC entry 233 (class 1259 OID 20637)
-- Name: companies; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.companies (
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name text NOT NULL,
    is_group boolean DEFAULT false NOT NULL,
    parent_id bigint,
    ordinal integer NOT NULL,
    code text NOT NULL,
    logo text,
    site_url text
);


ALTER TABLE public.companies OWNER TO postgres;

--
-- TOC entry 234 (class 1259 OID 20646)
-- Name: companies_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.companies_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.companies_id_seq OWNER TO postgres;

--
-- TOC entry 4598 (class 0 OID 0)
-- Dependencies: 234
-- Name: companies_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.companies_id_seq OWNED BY public.companies.id;


--
-- TOC entry 235 (class 1259 OID 20647)
-- Name: company_defaults; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.company_defaults (
    id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    updated_at timestamp without time zone,
    country text,
    tax_id bigint,
    currency text NOT NULL,
    company_id bigint NOT NULL
);


ALTER TABLE public.company_defaults OWNER TO postgres;

--
-- TOC entry 236 (class 1259 OID 20653)
-- Name: company_defaults_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.company_defaults ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.company_defaults_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 237 (class 1259 OID 20654)
-- Name: company_entities; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.company_entities (
    company_id bigint NOT NULL,
    entity_id bigint NOT NULL,
    enabled boolean DEFAULT false NOT NULL
);


ALTER TABLE public.company_entities OWNER TO postgres;

--
-- TOC entry 238 (class 1259 OID 20658)
-- Name: contacts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.contacts (
    id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    family_name text,
    given_name text NOT NULL,
    gender text,
    email text,
    phone_number text,
    company_id bigint NOT NULL
);


ALTER TABLE public.contacts OWNER TO postgres;

--
-- TOC entry 239 (class 1259 OID 20665)
-- Name: cost_centers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cost_centers (
    id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    updated_at timestamp without time zone,
    name text NOT NULL,
    company_id bigint NOT NULL,
    status text NOT NULL,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL
);


ALTER TABLE public.cost_centers OWNER TO postgres;

--
-- TOC entry 240 (class 1259 OID 20672)
-- Name: currencies; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.currencies (
    "code " text NOT NULL,
    "name " text NOT NULL
);


ALTER TABLE public.currencies OWNER TO postgres;

--
-- TOC entry 241 (class 1259 OID 20677)
-- Name: currency_exchanges; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.currency_exchanges (
    id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    status text NOT NULL,
    company_id bigint NOT NULL,
    name text NOT NULL,
    from_currency text NOT NULL,
    to_currency text NOT NULL,
    exchange_rate integer NOT NULL,
    for_buying boolean DEFAULT false NOT NULL,
    for_selling boolean DEFAULT false NOT NULL,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL
);


ALTER TABLE public.currency_exchanges OWNER TO postgres;

--
-- TOC entry 242 (class 1259 OID 20686)
-- Name: customers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.customers (
    id bigint NOT NULL,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    updated_at timestamp without time zone,
    customer_type text NOT NULL,
    group_id bigint,
    name text NOT NULL,
    company_id bigint NOT NULL,
    status text DEFAULT 'ENABLED'::text NOT NULL
);


ALTER TABLE public.customers OWNER TO postgres;

--
-- TOC entry 243 (class 1259 OID 20694)
-- Name: delivery_line_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.delivery_line_items (
    item_line_id integer NOT NULL,
    source_warehouse_id bigint NOT NULL,
    billed_quantity integer NOT NULL
);


ALTER TABLE public.delivery_line_items OWNER TO postgres;

--
-- TOC entry 244 (class 1259 OID 20697)
-- Name: entities; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.entities (
    id bigint NOT NULL,
    name text NOT NULL,
    href text DEFAULT ''::text NOT NULL
);


ALTER TABLE public.entities OWNER TO postgres;

--
-- TOC entry 245 (class 1259 OID 20703)
-- Name: entities_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.entities_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.entities_id_seq OWNER TO postgres;

--
-- TOC entry 4599 (class 0 OID 0)
-- Dependencies: 245
-- Name: entities_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.entities_id_seq OWNED BY public.entities.id;


--
-- TOC entry 246 (class 1259 OID 20704)
-- Name: entity_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.entity_types (
    entity_type text NOT NULL,
    name text
);


ALTER TABLE public.entity_types OWNER TO postgres;

--
-- TOC entry 247 (class 1259 OID 20709)
-- Name: groups; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.groups (
    id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name text NOT NULL,
    ordinal smallint NOT NULL,
    company_id bigint NOT NULL,
    parent_id bigint,
    is_group boolean DEFAULT false NOT NULL,
    enabled boolean DEFAULT true NOT NULL
);


ALTER TABLE public.groups OWNER TO postgres;

--
-- TOC entry 248 (class 1259 OID 20718)
-- Name: invoiced_item_lines; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.invoiced_item_lines (
    item_line integer NOT NULL,
    receipt_quantity integer DEFAULT 0 NOT NULL,
    invoice_id bigint NOT NULL,
    paid_amount bigint DEFAULT 0 NOT NULL
);


ALTER TABLE public.invoiced_item_lines OWNER TO postgres;

--
-- TOC entry 249 (class 1259 OID 20723)
-- Name: invoices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.invoices (
    id bigint NOT NULL,
    code text NOT NULL,
    party_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    due_date date,
    status text NOT NULL,
    company_id bigint NOT NULL,
    currency text NOT NULL,
    bill_date date,
    update_stock boolean DEFAULT false NOT NULL,
    posting_date date DEFAULT CURRENT_DATE NOT NULL,
    posting_time time without time zone DEFAULT CURRENT_TIME NOT NULL,
    tz text DEFAULT ''::text NOT NULL,
    project_id bigint,
    cost_center_id bigint,
    doc_reference_id bigint,
    price_list_id bigint,
    warehouse_id bigint
);


ALTER TABLE public.invoices OWNER TO postgres;

--
-- TOC entry 250 (class 1259 OID 20733)
-- Name: item_attribute_values; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.item_attribute_values (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    updated_at timestamp with time zone DEFAULT now(),
    ordinal integer NOT NULL,
    value text NOT NULL,
    abbreviation text NOT NULL,
    item_attribute_id bigint NOT NULL
);


ALTER TABLE public.item_attribute_values OWNER TO postgres;

--
-- TOC entry 251 (class 1259 OID 20740)
-- Name: item_attribute_values_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.item_attribute_values_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.item_attribute_values_id_seq OWNER TO postgres;

--
-- TOC entry 4600 (class 0 OID 0)
-- Dependencies: 251
-- Name: item_attribute_values_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.item_attribute_values_id_seq OWNED BY public.item_attribute_values.id;


--
-- TOC entry 252 (class 1259 OID 20741)
-- Name: item_attributes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.item_attributes (
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    updated_at timestamp with time zone DEFAULT now(),
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name text NOT NULL,
    company_id bigint NOT NULL
);


ALTER TABLE public.item_attributes OWNER TO postgres;

--
-- TOC entry 253 (class 1259 OID 20749)
-- Name: item_attributes_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.item_attributes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.item_attributes_id_seq OWNER TO postgres;

--
-- TOC entry 4601 (class 0 OID 0)
-- Dependencies: 253
-- Name: item_attributes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.item_attributes_id_seq OWNED BY public.item_attributes.id;


--
-- TOC entry 254 (class 1259 OID 20750)
-- Name: item_inventory_settings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.item_inventory_settings (
    item_id bigint NOT NULL,
    shelf_life_in_days integer,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    warranty_period_in_days integer,
    has_serial_no boolean,
    serial_no_template text,
    weight_uom_id bigint,
    weight_per_unit integer
);


ALTER TABLE public.item_inventory_settings OWNER TO postgres;

--
-- TOC entry 255 (class 1259 OID 20756)
-- Name: item_line_receipts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.item_line_receipts (
    item_line integer NOT NULL,
    accepted_warehouse bigint NOT NULL,
    rejected_warehouse bigint,
    accepted_quantity integer DEFAULT 0 NOT NULL,
    rejected_quantity integer DEFAULT 0 NOT NULL,
    invoice_item_line_id integer,
    billed_quantity integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.item_line_receipts OWNER TO postgres;

--
-- TOC entry 256 (class 1259 OID 20762)
-- Name: item_line_stock_entries; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.item_line_stock_entries (
    id integer NOT NULL,
    source_warehouse_id bigint,
    target_warehouse_id bigint,
    item_line integer NOT NULL
);


ALTER TABLE public.item_line_stock_entries OWNER TO postgres;

--
-- TOC entry 257 (class 1259 OID 20765)
-- Name: item_line_stock_entries_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.item_line_stock_entries ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.item_line_stock_entries_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 258 (class 1259 OID 20766)
-- Name: item_line_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.item_line_types (
    type text NOT NULL
);


ALTER TABLE public.item_line_types OWNER TO postgres;

--
-- TOC entry 259 (class 1259 OID 20771)
-- Name: item_lines; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.item_lines (
    id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    deleted_at timestamp without time zone,
    rate bigint NOT NULL,
    quantity integer NOT NULL,
    party_id bigint,
    item_line_reference_id integer,
    "type " text DEFAULT 'ITEM_LINE_ORDER'::text NOT NULL,
    item_id bigint NOT NULL,
    unit_of_measure_id bigint NOT NULL
);


ALTER TABLE public.item_lines OWNER TO postgres;

--
-- TOC entry 260 (class 1259 OID 20778)
-- Name: item_lines_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.item_lines_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.item_lines_id_seq OWNER TO postgres;

--
-- TOC entry 4602 (class 0 OID 0)
-- Dependencies: 260
-- Name: item_lines_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.item_lines_id_seq OWNED BY public.item_lines.id;


--
-- TOC entry 261 (class 1259 OID 20779)
-- Name: price_lists; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.price_lists (
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    updated_at timestamp with time zone DEFAULT now(),
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name text NOT NULL,
    is_buying boolean DEFAULT false NOT NULL,
    is_selling boolean DEFAULT false NOT NULL,
    currency text NOT NULL,
    company_id bigint NOT NULL,
    status text DEFAULT 'DRAFT'::text NOT NULL
);


ALTER TABLE public.price_lists OWNER TO postgres;

--
-- TOC entry 262 (class 1259 OID 20790)
-- Name: item_price_lists_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.item_price_lists_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.item_price_lists_id_seq OWNER TO postgres;

--
-- TOC entry 4603 (class 0 OID 0)
-- Dependencies: 262
-- Name: item_price_lists_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.item_price_lists_id_seq OWNED BY public.price_lists.id;


--
-- TOC entry 263 (class 1259 OID 20791)
-- Name: item_price_plugins; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.item_price_plugins (
    plugin text NOT NULL,
    base_id bigint NOT NULL
);


ALTER TABLE public.item_price_plugins OWNER TO postgres;

--
-- TOC entry 264 (class 1259 OID 20796)
-- Name: item_prices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.item_prices (
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    updated_at timestamp with time zone DEFAULT now(),
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    item_id bigint NOT NULL,
    price_list_id bigint NOT NULL,
    company_id bigint NOT NULL,
    item_quantity integer DEFAULT 1 NOT NULL,
    rate bigint NOT NULL,
    unit_of_measure_id bigint NOT NULL
);


ALTER TABLE public.item_prices OWNER TO postgres;

--
-- TOC entry 265 (class 1259 OID 20803)
-- Name: item_prices_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.item_prices_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.item_prices_id_seq OWNER TO postgres;

--
-- TOC entry 4604 (class 0 OID 0)
-- Dependencies: 265
-- Name: item_prices_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.item_prices_id_seq OWNED BY public.item_prices.id;


--
-- TOC entry 266 (class 1259 OID 20804)
-- Name: item_variants; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.item_variants (
    item_attribute_value_id integer NOT NULL,
    item_id bigint NOT NULL,
    variant_id bigint NOT NULL
);


ALTER TABLE public.item_variants OWNER TO postgres;

--
-- TOC entry 267 (class 1259 OID 20807)
-- Name: items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.items (
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name text NOT NULL,
    code text NOT NULL,
    group_id bigint,
    company_id bigint NOT NULL,
    unit_of_measure_id bigint NOT NULL,
    parent_id bigint,
    item_type text DEFAULT 'ITEM'::text NOT NULL,
    maintain_stock boolean DEFAULT true NOT NULL,
    description text
);


ALTER TABLE public.items OWNER TO postgres;

--
-- TOC entry 268 (class 1259 OID 20817)
-- Name: items_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.items_id_seq OWNER TO postgres;

--
-- TOC entry 4605 (class 0 OID 0)
-- Dependencies: 268
-- Name: items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.items_id_seq OWNED BY public.items.id;


--
-- TOC entry 269 (class 1259 OID 20818)
-- Name: journal_entries; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.journal_entries (
    id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    posting_date date NOT NULL,
    deleted_at timestamp without time zone,
    updated_at timestamp without time zone,
    company_id bigint NOT NULL,
    code text NOT NULL,
    status text NOT NULL,
    entry_type text NOT NULL
);


ALTER TABLE public.journal_entries OWNER TO postgres;

--
-- TOC entry 270 (class 1259 OID 20824)
-- Name: journal_entry_lines; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.journal_entry_lines (
    id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    ledger_id bigint NOT NULL,
    cost_center_id bigint,
    debit bigint DEFAULT 0 NOT NULL,
    credit bigint DEFAULT 0 NOT NULL,
    currency text NOT NULL,
    project_id bigint,
    journal_entry_id bigint NOT NULL
);


ALTER TABLE public.journal_entry_lines OWNER TO postgres;

--
-- TOC entry 271 (class 1259 OID 20832)
-- Name: journal_entry_lines_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.journal_entry_lines ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.journal_entry_lines_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 272 (class 1259 OID 20833)
-- Name: key_values; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.key_values (
    id integer NOT NULL,
    key text NOT NULL,
    value text NOT NULL,
    party_id bigint NOT NULL
);


ALTER TABLE public.key_values OWNER TO postgres;

--
-- TOC entry 273 (class 1259 OID 20838)
-- Name: key_values_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.key_values ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.key_values_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 274 (class 1259 OID 20839)
-- Name: ledger_accounts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ledger_accounts (
    ledger_id bigint NOT NULL,
    can_credit boolean DEFAULT true NOT NULL,
    can_debit boolean DEFAULT true NOT NULL,
    "limit" integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    currency text NOT NULL
);


ALTER TABLE public.ledger_accounts OWNER TO postgres;

--
-- TOC entry 275 (class 1259 OID 20848)
-- Name: ledger_statements; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ledger_statements (
    ledger_no integer NOT NULL,
    date timestamp without time zone NOT NULL,
    closing_balance timestamp without time zone NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.ledger_statements OWNER TO postgres;

--
-- TOC entry 276 (class 1259 OID 20852)
-- Name: ledger_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ledger_types (
    ledger_type text NOT NULL,
    description text
);


ALTER TABLE public.ledger_types OWNER TO postgres;

--
-- TOC entry 277 (class 1259 OID 20857)
-- Name: ledgers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ledgers (
    ledger_parent bigint,
    account_type text,
    name text NOT NULL,
    company_id bigint NOT NULL,
    is_group boolean DEFAULT false NOT NULL,
    id bigint NOT NULL,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    ledger_no text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone,
    account_root_type text NOT NULL,
    report_type text,
    cash_flow_section text,
    status text DEFAULT 'ENABLED'::text NOT NULL,
    is_offset_account boolean DEFAULT false NOT NULL
);


ALTER TABLE public.ledgers OWNER TO postgres;

--
-- TOC entry 278 (class 1259 OID 20867)
-- Name: module_sections; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.module_sections (
    id integer NOT NULL,
    name text NOT NULL,
    module_id bigint NOT NULL,
    entity_id integer NOT NULL
);


ALTER TABLE public.module_sections OWNER TO postgres;

--
-- TOC entry 279 (class 1259 OID 20872)
-- Name: module_sections_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.module_sections ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.module_sections_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 280 (class 1259 OID 20873)
-- Name: modules; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.modules (
    id bigint NOT NULL,
    company_id bigint NOT NULL,
    deleted_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    icon_code text,
    label text NOT NULL,
    href text NOT NULL,
    status text DEFAULT 'ENABLED'::text NOT NULL,
    uuid text DEFAULT public.uuid_generate_v4() NOT NULL,
    icon_name text,
    has_direct_access boolean DEFAULT false NOT NULL,
    priority integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.modules OWNER TO postgres;

--
-- TOC entry 281 (class 1259 OID 20883)
-- Name: orders; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.orders (
    id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    updated_at timestamp without time zone,
    code text NOT NULL,
    delivery_date timestamp without time zone,
    company_id bigint NOT NULL,
    currency text NOT NULL,
    status text DEFAULT 'DRAFT'::text NOT NULL,
    posting_date date DEFAULT CURRENT_DATE NOT NULL,
    party_id bigint NOT NULL,
    posting_time text DEFAULT ''::text NOT NULL,
    tz text DEFAULT ''::text NOT NULL,
    project_id bigint,
    cost_center_id bigint,
    price_list_id bigint
);


ALTER TABLE public.orders OWNER TO postgres;

--
-- TOC entry 282 (class 1259 OID 20893)
-- Name: parties; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.parties (
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    updated_at timestamp with time zone DEFAULT now(),
    party_type_code text NOT NULL
);


ALTER TABLE public.parties OWNER TO postgres;

--
-- TOC entry 283 (class 1259 OID 20900)
-- Name: parties_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.parties_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.parties_id_seq OWNER TO postgres;

--
-- TOC entry 4606 (class 0 OID 0)
-- Dependencies: 283
-- Name: parties_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.parties_id_seq OWNED BY public.parties.id;


--
-- TOC entry 284 (class 1259 OID 20901)
-- Name: parties_id_seq1; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.parties ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.parties_id_seq1
    START WITH 1423
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 285 (class 1259 OID 20902)
-- Name: party_addresses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.party_addresses (
    party_id bigint NOT NULL,
    address_id bigint NOT NULL,
    is_shipping_address boolean DEFAULT false NOT NULL,
    is_billing_address boolean DEFAULT false NOT NULL,
    enabled boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.party_addresses OWNER TO postgres;

--
-- TOC entry 286 (class 1259 OID 20910)
-- Name: party_payments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.party_payments (
    party_id bigint NOT NULL,
    payment_id bigint NOT NULL,
    party_bank_account text,
    company_bank_account text,
    contact_id bigint
);


ALTER TABLE public.party_payments OWNER TO postgres;

--
-- TOC entry 287 (class 1259 OID 20915)
-- Name: party_references; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.party_references (
    party_id bigint NOT NULL,
    reference_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.party_references OWNER TO postgres;

--
-- TOC entry 288 (class 1259 OID 20919)
-- Name: party_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.party_types (
    code text NOT NULL,
    name text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    entity_id bigint
);


ALTER TABLE public.party_types OWNER TO postgres;

--
-- TOC entry 289 (class 1259 OID 20925)
-- Name: payment_references; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payment_references (
    payment_id bigint NOT NULL,
    party_id bigint NOT NULL,
    party_code text NOT NULL,
    total bigint NOT NULL,
    outstanding bigint NOT NULL,
    allocated bigint NOT NULL,
    id integer NOT NULL,
    currency text DEFAULT 'BOB'::text NOT NULL
);


ALTER TABLE public.payment_references OWNER TO postgres;

--
-- TOC entry 290 (class 1259 OID 20931)
-- Name: payment_references_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.payment_references ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.payment_references_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 291 (class 1259 OID 20932)
-- Name: payment_terms; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payment_terms (
    id bigint NOT NULL,
    uuid text DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    name text NOT NULL,
    invoice_portion integer NOT NULL,
    credit_days integer,
    due_date_base_on text NOT NULL,
    description text,
    company_id bigint NOT NULL,
    status text NOT NULL,
    discount_type text,
    discount bigint,
    discount_validity_base_on text,
    discount_validity integer
);


ALTER TABLE public.payment_terms OWNER TO postgres;

--
-- TOC entry 292 (class 1259 OID 20939)
-- Name: payment_terms_lines; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payment_terms_lines (
    id integer NOT NULL,
    create_at timestamp without time zone NOT NULL,
    deleted_at timestamp without time zone,
    payment_terms_id bigint NOT NULL,
    document_id bigint NOT NULL,
    description text,
    invoice_portion integer NOT NULL,
    due_date_base_on text NOT NULL,
    credit_days integer
);


ALTER TABLE public.payment_terms_lines OWNER TO postgres;

--
-- TOC entry 293 (class 1259 OID 20944)
-- Name: payment_terms_lines_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.payment_terms_lines ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.payment_terms_lines_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 294 (class 1259 OID 20945)
-- Name: payment_terms_templates; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payment_terms_templates (
    id bigint NOT NULL,
    uuid text DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    name text NOT NULL,
    company_id bigint NOT NULL,
    status text NOT NULL
);


ALTER TABLE public.payment_terms_templates OWNER TO postgres;

--
-- TOC entry 295 (class 1259 OID 20952)
-- Name: payments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payments (
    id bigint NOT NULL,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    updated_at timestamp without time zone,
    posting_date date NOT NULL,
    payment_type text NOT NULL,
    amount bigint NOT NULL,
    code text NOT NULL,
    company_id bigint NOT NULL,
    party_id bigint NOT NULL,
    status text DEFAULT 'DRAFT'::text NOT NULL,
    account_paid_from_id bigint NOT NULL,
    account_paid_to_id bigint NOT NULL,
    project_id bigint,
    cost_center_id bigint
);


ALTER TABLE public.payments OWNER TO postgres;

--
-- TOC entry 296 (class 1259 OID 20960)
-- Name: piano_form; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.piano_form (
    id bigint NOT NULL,
    moving_date date NOT NULL,
    piano_type text NOT NULL,
    pickup_flights text NOT NULL,
    pickup_street text NOT NULL,
    pickup_city text NOT NULL,
    pickup_state text NOT NULL,
    pickup_zip text NOT NULL,
    dropoff_flights text NOT NULL,
    dropoff_street text NOT NULL,
    dropoff_city text NOT NULL,
    dropoff_state text NOT NULL,
    dropoff_zip text NOT NULL,
    company_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    updated_at timestamp without time zone,
    first_name text NOT NULL,
    last_name text NOT NULL,
    email text NOT NULL,
    phone_number text NOT NULL,
    rent_piano boolean NOT NULL,
    stairs_dropoff boolean NOT NULL,
    stairs_pickup boolean NOT NULL
);


ALTER TABLE public.piano_form OWNER TO postgres;

--
-- TOC entry 297 (class 1259 OID 20966)
-- Name: plugins; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.plugins (
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    name text NOT NULL,
    code text NOT NULL
);


ALTER TABLE public.plugins OWNER TO postgres;

--
-- TOC entry 298 (class 1259 OID 20972)
-- Name: plugins_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.plugins_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.plugins_id_seq OWNER TO postgres;

--
-- TOC entry 4607 (class 0 OID 0)
-- Dependencies: 298
-- Name: plugins_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.plugins_id_seq OWNED BY public.plugins.id;


--
-- TOC entry 299 (class 1259 OID 20973)
-- Name: pricing_line_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pricing_line_items (
    id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    supplier_id bigint,
    part_number text,
    description text,
    quantity integer,
    pl_unit integer,
    pricing_id bigint NOT NULL,
    fob_unit_fn text,
    retention_fn text,
    cost_zf_fn text,
    cost_alm_fn text,
    tva_fn text,
    cantidad_fn text,
    precio_unitario_fn text,
    precio_total_fn text,
    precio_unitario_tc_fn text,
    precio_total_tc_fn text,
    fob_total_fn text,
    gpl_total_fn text,
    tva_total_fn text,
    is_title boolean,
    color text
);


ALTER TABLE public.pricing_line_items OWNER TO postgres;

--
-- TOC entry 300 (class 1259 OID 20979)
-- Name: priced_item_lines_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.pricing_line_items ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.priced_item_lines_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 301 (class 1259 OID 20980)
-- Name: pricing_charges; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pricing_charges (
    id integer NOT NULL,
    name text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    deleted_at timestamp without time zone,
    rate integer NOT NULL,
    pricing_id bigint NOT NULL
);


ALTER TABLE public.pricing_charges OWNER TO postgres;

--
-- TOC entry 302 (class 1259 OID 20985)
-- Name: pricing_charges_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.pricing_charges ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.pricing_charges_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 303 (class 1259 OID 20986)
-- Name: pricings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pricings (
    id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    code text NOT NULL,
    status text NOT NULL,
    company_id bigint NOT NULL,
    project_id bigint,
    cost_center_id bigint,
    customer_id bigint
);


ALTER TABLE public.pricings OWNER TO postgres;

--
-- TOC entry 304 (class 1259 OID 20992)
-- Name: profiles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.profiles (
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    updated_at timestamp with time zone DEFAULT now(),
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    given_name text NOT NULL,
    family_name text NOT NULL,
    email_address text NOT NULL,
    phone_number text,
    avatar text
);


ALTER TABLE public.profiles OWNER TO postgres;

--
-- TOC entry 305 (class 1259 OID 21000)
-- Name: profiles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.profiles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.profiles_id_seq OWNER TO postgres;

--
-- TOC entry 4608 (class 0 OID 0)
-- Dependencies: 305
-- Name: profiles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.profiles_id_seq OWNED BY public.profiles.id;


--
-- TOC entry 306 (class 1259 OID 21001)
-- Name: progress_invoices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.progress_invoices (
    invoice_id bigint NOT NULL,
    total_amount bigint NOT NULL,
    paid_amount bigint DEFAULT 0 NOT NULL
);


ALTER TABLE public.progress_invoices OWNER TO postgres;

--
-- TOC entry 307 (class 1259 OID 21005)
-- Name: progress_orders; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.progress_orders (
    order_id bigint NOT NULL,
    total_items integer DEFAULT 0 NOT NULL,
    received_items integer DEFAULT 0 NOT NULL,
    total_amount bigint DEFAULT 0 NOT NULL,
    billed_amount bigint DEFAULT 0 NOT NULL
);


ALTER TABLE public.progress_orders OWNER TO postgres;

--
-- TOC entry 308 (class 1259 OID 21012)
-- Name: projects; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.projects (
    id bigint NOT NULL,
    created_at timestamp without time zone NOT NULL,
    deleted_at timestamp without time zone,
    updated_at timestamp without time zone,
    name text NOT NULL,
    status text NOT NULL,
    company_id bigint NOT NULL,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL
);


ALTER TABLE public.projects OWNER TO postgres;

--
-- TOC entry 309 (class 1259 OID 21018)
-- Name: purchase_records; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.purchase_records (
    id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    supplier_nit text NOT NULL,
    supplier_business_name text NOT NULL,
    authorization_code text NOT NULL,
    invoice_no text NOT NULL,
    dui_dim_no text NOT NULL,
    invoice_dui_dim_date date NOT NULL,
    total_purchase_amount bigint NOT NULL,
    ice_amount integer NOT NULL,
    iehd_amount integer NOT NULL,
    ipj_amount integer NOT NULL,
    tax_rates integer NOT NULL,
    other_not_subject_to_tax_credit integer NOT NULL,
    exempt_amounts integer NOT NULL,
    zero_rate_taxable_purchases_amount integer NOT NULL,
    subtotal bigint NOT NULL,
    discounts_bonus_rebates_subject_to_vat integer NOT NULL,
    gift_card_amount integer NOT NULL,
    cf_base_amount integer NOT NULL,
    tax_credit integer NOT NULL,
    purchase_type text NOT NULL,
    control_code text NOT NULL,
    consolidation_status text NOT NULL,
    company_id bigint NOT NULL,
    supplier_id bigint NOT NULL,
    status text NOT NULL,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    with_tax_credit_right boolean NOT NULL,
    invoice_id bigint
);


ALTER TABLE public.purchase_records OWNER TO postgres;

--
-- TOC entry 310 (class 1259 OID 21025)
-- Name: quotations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.quotations (
    id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    updated_at timestamp without time zone,
    posting_date date NOT NULL,
    posting_time time without time zone NOT NULL,
    party_id bigint NOT NULL,
    status text DEFAULT 'DRAFT'::text NOT NULL,
    company_id bigint NOT NULL,
    currency text NOT NULL,
    valid_till date NOT NULL,
    code text NOT NULL,
    tz text NOT NULL,
    project_id bigint,
    cost_center_id bigint,
    price_list_id bigint
);


ALTER TABLE public.quotations OWNER TO postgres;

--
-- TOC entry 311 (class 1259 OID 21032)
-- Name: r_booking_events; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.r_booking_events (
    booking_id bigint NOT NULL,
    event_id bigint NOT NULL
);


ALTER TABLE public.r_booking_events OWNER TO postgres;

--
-- TOC entry 312 (class 1259 OID 21035)
-- Name: r_booking_prices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.r_booking_prices (
    id bigint NOT NULL,
    total_price integer NOT NULL,
    paid integer NOT NULL,
    discount integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.r_booking_prices OWNER TO postgres;

--
-- TOC entry 313 (class 1259 OID 21039)
-- Name: r_booking_slots; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.r_booking_slots (
    id integer NOT NULL,
    court_id bigint NOT NULL,
    paid_amount integer DEFAULT 0 NOT NULL,
    datetime timestamp with time zone NOT NULL,
    booking_id bigint NOT NULL,
    company_id bigint NOT NULL,
    type text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    total_price integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.r_booking_slots OWNER TO postgres;

--
-- TOC entry 314 (class 1259 OID 21047)
-- Name: r_booking_slots_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.r_booking_slots ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.r_booking_slots_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 315 (class 1259 OID 21048)
-- Name: r_bookings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.r_bookings (
    id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone,
    start_date timestamp without time zone NOT NULL,
    end_date timestamp without time zone NOT NULL,
    status text NOT NULL,
    party bigint NOT NULL,
    type text NOT NULL,
    company_id bigint NOT NULL,
    court_id bigint NOT NULL,
    code text NOT NULL
);


ALTER TABLE public.r_bookings OWNER TO postgres;

--
-- TOC entry 316 (class 1259 OID 21054)
-- Name: r_court_rates; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.r_court_rates (
    id integer NOT NULL,
    rate integer NOT NULL,
    enabled boolean DEFAULT false NOT NULL,
    court_id bigint NOT NULL,
    day_week integer NOT NULL,
    company_id bigint NOT NULL,
    "time" time without time zone NOT NULL,
    currency text NOT NULL,
    status text DEFAULT 'ENABLED'::text NOT NULL
);


ALTER TABLE public.r_court_rates OWNER TO postgres;

--
-- TOC entry 317 (class 1259 OID 21061)
-- Name: r_court_rates_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.r_court_rates ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.r_court_rates_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 318 (class 1259 OID 21062)
-- Name: r_courts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.r_courts (
    id bigint NOT NULL,
    name text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone,
    enabled boolean DEFAULT false NOT NULL,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    company_id bigint NOT NULL,
    status text DEFAULT 'ENABLED'::text NOT NULL
);


ALTER TABLE public.r_courts OWNER TO postgres;

--
-- TOC entry 319 (class 1259 OID 21071)
-- Name: r_events; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.r_events (
    id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    updated_at timestamp without time zone,
    name text NOT NULL,
    description text,
    status text DEFAULT 'ACTIVE'::text NOT NULL,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    company_id bigint NOT NULL
);


ALTER TABLE public.r_events OWNER TO postgres;

--
-- TOC entry 320 (class 1259 OID 21079)
-- Name: receipts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.receipts (
    id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    updated_at timestamp without time zone,
    code text NOT NULL,
    posting_date date NOT NULL,
    party_id bigint NOT NULL,
    currency text NOT NULL,
    status text NOT NULL,
    company_id bigint NOT NULL,
    posting_time time without time zone DEFAULT CURRENT_TIME NOT NULL,
    tz text DEFAULT ''::text NOT NULL,
    project_id bigint,
    cost_center_id bigint,
    doc_reference_id bigint,
    price_list_id bigint,
    warehouse_id bigint NOT NULL
);


ALTER TABLE public.receipts OWNER TO postgres;

--
-- TOC entry 321 (class 1259 OID 21087)
-- Name: role_actions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.role_actions (
    role_id bigint NOT NULL,
    action_id bigint NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.role_actions OWNER TO postgres;

--
-- TOC entry 322 (class 1259 OID 21090)
-- Name: role_templates; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.role_templates (
    id bigint NOT NULL,
    name text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.role_templates OWNER TO postgres;

--
-- TOC entry 323 (class 1259 OID 21096)
-- Name: roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.roles (
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    code text NOT NULL,
    description text,
    company_id bigint NOT NULL
);


ALTER TABLE public.roles OWNER TO postgres;

--
-- TOC entry 324 (class 1259 OID 21104)
-- Name: roles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.roles_id_seq OWNER TO postgres;

--
-- TOC entry 4609 (class 0 OID 0)
-- Dependencies: 324
-- Name: roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.roles_id_seq OWNED BY public.roles.id;


--
-- TOC entry 325 (class 1259 OID 21105)
-- Name: sales_records; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sales_records (
    id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    invoice_date date NOT NULL,
    invoice_no text NOT NULL,
    authorization_code text NOT NULL,
    customer_nit_ci text NOT NULL,
    supplement text NOT NULL,
    name_or_business_name text NOT NULL,
    total_sale_amount bigint NOT NULL,
    ice_amount integer NOT NULL,
    iehd_amount integer NOT NULL,
    ipj_amount integer NOT NULL,
    tax_rates integer NOT NULL,
    other_not_subject_to_vat integer NOT NULL,
    exports_and_exempt_operations integer NOT NULL,
    zero_rate_taxable_sales integer NOT NULL,
    subtotal bigint NOT NULL,
    discounts_bonus_and_rebates_subject_to_vat integer NOT NULL,
    gift_card_amount integer NOT NULL,
    base_amount_for_tax_debit integer NOT NULL,
    tax_debit integer NOT NULL,
    state text NOT NULL,
    control_code text NOT NULL,
    sale_type text NOT NULL,
    with_tax_credit_right boolean NOT NULL,
    consolidation_status text NOT NULL,
    company_id bigint NOT NULL,
    customer_id bigint NOT NULL,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    status text NOT NULL,
    invoice_id bigint
);


ALTER TABLE public.sales_records OWNER TO postgres;

--
-- TOC entry 326 (class 1259 OID 21112)
-- Name: serial_no_transactions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.serial_no_transactions (
    id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    serial_no_id bigint NOT NULL,
    qty smallint NOT NULL,
    status text NOT NULL,
    batch_bundle_id bigint NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.serial_no_transactions OWNER TO postgres;

--
-- TOC entry 327 (class 1259 OID 21118)
-- Name: serial_no_transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.serial_no_transactions ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.serial_no_transactions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 328 (class 1259 OID 21119)
-- Name: serial_nos; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.serial_nos (
    id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    status text NOT NULL,
    serial_no text NOT NULL,
    batch_bundle_id bigint NOT NULL,
    valuation_rate bigint NOT NULL,
    item_id bigint NOT NULL
);


ALTER TABLE public.serial_nos OWNER TO postgres;

--
-- TOC entry 329 (class 1259 OID 21125)
-- Name: states; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.states (
    state text NOT NULL
);


ALTER TABLE public.states OWNER TO postgres;

--
-- TOC entry 330 (class 1259 OID 21130)
-- Name: stock_defaults; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stock_defaults (
    id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone,
    company_id bigint NOT NULL,
    valuation_method text DEFAULT 'FIFO'::text NOT NULL,
    default_warehouse text
);


ALTER TABLE public.stock_defaults OWNER TO postgres;

--
-- TOC entry 331 (class 1259 OID 21137)
-- Name: stock_entries; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stock_entries (
    id bigint NOT NULL,
    code text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    updated_at timestamp without time zone,
    entry_type text NOT NULL,
    status text NOT NULL,
    company_id bigint NOT NULL,
    currency text NOT NULL,
    posting_date date NOT NULL,
    posting_time time without time zone NOT NULL,
    tz text DEFAULT 'America/La_Paz'::text NOT NULL,
    project_id bigint,
    cost_center_id bigint,
    source_warehouse_id bigint,
    target_warehouse_id bigint
);


ALTER TABLE public.stock_entries OWNER TO postgres;

--
-- TOC entry 332 (class 1259 OID 21144)
-- Name: stock_levels; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stock_levels (
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    updated_at timestamp with time zone DEFAULT now(),
    enabled boolean NOT NULL,
    stock integer NOT NULL,
    out_of_stock_threshold integer NOT NULL,
    ware_house_id bigint NOT NULL,
    item_id bigint NOT NULL,
    id bigint NOT NULL,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL
);


ALTER TABLE public.stock_levels OWNER TO postgres;

--
-- TOC entry 333 (class 1259 OID 21150)
-- Name: stock_movements; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stock_movements (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    updated_at timestamp with time zone DEFAULT now(),
    quantity integer NOT NULL,
    item_id bigint NOT NULL,
    ware_house_id bigint NOT NULL,
    stock_movement_type text NOT NULL
);


ALTER TABLE public.stock_movements OWNER TO postgres;

--
-- TOC entry 334 (class 1259 OID 21157)
-- Name: stock_movements_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.stock_movements_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.stock_movements_id_seq OWNER TO postgres;

--
-- TOC entry 4610 (class 0 OID 0)
-- Dependencies: 334
-- Name: stock_movements_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.stock_movements_id_seq OWNED BY public.stock_movements.id;


--
-- TOC entry 335 (class 1259 OID 21158)
-- Name: stock_settings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stock_settings (
    id integer NOT NULL,
    inventory_account bigint NOT NULL,
    stock_received_but_not_billed bigint NOT NULL,
    company_id bigint NOT NULL,
    stock_adjustment bigint NOT NULL
);


ALTER TABLE public.stock_settings OWNER TO postgres;

--
-- TOC entry 336 (class 1259 OID 21161)
-- Name: stock_settings_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.stock_settings ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.stock_settings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 337 (class 1259 OID 21162)
-- Name: stock_transactions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stock_transactions (
    id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    item_id bigint NOT NULL,
    uom_id bigint NOT NULL,
    in_qty integer DEFAULT 0 NOT NULL,
    out_qty integer DEFAULT 0 NOT NULL,
    balance_qty integer NOT NULL,
    warehouse_id bigint NOT NULL,
    incoming_rate bigint DEFAULT 0 NOT NULL,
    balance_value bigint DEFAULT 0 NOT NULL,
    valuation_rate bigint NOT NULL,
    voucher_type text NOT NULL,
    voucher_no text NOT NULL,
    currency text DEFAULT 'USD'::text NOT NULL,
    average_rate bigint DEFAULT 0 NOT NULL,
    available_qty integer DEFAULT 0 NOT NULL,
    deleted_at timestamp without time zone,
    posting_date date DEFAULT CURRENT_DATE NOT NULL
);


ALTER TABLE public.stock_transactions OWNER TO postgres;

--
-- TOC entry 338 (class 1259 OID 21176)
-- Name: stock_transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.stock_transactions ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.stock_transactions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 339 (class 1259 OID 21177)
-- Name: supplier_orders; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.supplier_orders (
    supplier_id bigint NOT NULL,
    order_id bigint NOT NULL
);


ALTER TABLE public.supplier_orders OWNER TO postgres;

--
-- TOC entry 340 (class 1259 OID 21180)
-- Name: suppliers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.suppliers (
    id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone,
    updated_at timestamp without time zone,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name text NOT NULL,
    company_id bigint NOT NULL,
    group_id bigint,
    status text DEFAULT 'ENABLED'::text NOT NULL
);


ALTER TABLE public.suppliers OWNER TO postgres;

--
-- TOC entry 341 (class 1259 OID 21188)
-- Name: tax_and_charge_lines; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tax_and_charge_lines (
    id integer NOT NULL,
    type text NOT NULL,
    account_head bigint NOT NULL,
    tax_rate smallint,
    doc_party_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    amount bigint DEFAULT 0 NOT NULL,
    is_deducted boolean DEFAULT false NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.tax_and_charge_lines OWNER TO postgres;

--
-- TOC entry 342 (class 1259 OID 21196)
-- Name: tax_and_charge_lines_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.tax_and_charge_lines ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.tax_and_charge_lines_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 343 (class 1259 OID 21197)
-- Name: taxes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.taxes (
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    updated_at timestamp with time zone DEFAULT now(),
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name text NOT NULL,
    enabled boolean DEFAULT false NOT NULL,
    value numeric(5,2) NOT NULL,
    company_id bigint NOT NULL
);


ALTER TABLE public.taxes OWNER TO postgres;

--
-- TOC entry 344 (class 1259 OID 21206)
-- Name: taxes_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.taxes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.taxes_id_seq OWNER TO postgres;

--
-- TOC entry 4611 (class 0 OID 0)
-- Dependencies: 344
-- Name: taxes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.taxes_id_seq OWNED BY public.taxes.id;


--
-- TOC entry 345 (class 1259 OID 21207)
-- Name: terms_and_conditions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.terms_and_conditions (
    id bigint NOT NULL,
    uuid text DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    name text NOT NULL,
    terms_and_conditions text NOT NULL,
    status text NOT NULL,
    company_id bigint NOT NULL
);


ALTER TABLE public.terms_and_conditions OWNER TO postgres;

--
-- TOC entry 346 (class 1259 OID 21214)
-- Name: transaction_accounts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transaction_accounts (
    ledger_no integer NOT NULL,
    datetime timestamp without time zone DEFAULT now() NOT NULL,
    xact_type_code text NOT NULL,
    xact_type_code_ext text NOT NULL,
    account_no integer NOT NULL
);


ALTER TABLE public.transaction_accounts OWNER TO postgres;

--
-- TOC entry 347 (class 1259 OID 21220)
-- Name: transaction_ledgers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transaction_ledgers (
    ledger bigint NOT NULL,
    ledger_against bigint,
    voucher_code text NOT NULL,
    voucher_type text NOT NULL,
    voucher_subtype text NOT NULL,
    party_id bigint,
    id integer NOT NULL,
    credit bigint DEFAULT 0 NOT NULL,
    debit bigint DEFAULT 0 NOT NULL,
    currency text DEFAULT 'USD'::text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    posting_date date DEFAULT CURRENT_DATE NOT NULL,
    cost_center_id bigint,
    project_id bigint,
    posting_time text DEFAULT CURRENT_TIME NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.transaction_ledgers OWNER TO postgres;

--
-- TOC entry 348 (class 1259 OID 21231)
-- Name: transaction_ledgers_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.transaction_ledgers ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.transaction_ledgers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 349 (class 1259 OID 21232)
-- Name: transaction_type_de; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transaction_type_de (
    type_code text NOT NULL,
    name text NOT NULL
);


ALTER TABLE public.transaction_type_de OWNER TO postgres;

--
-- TOC entry 350 (class 1259 OID 21237)
-- Name: transaction_type_ext; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transaction_type_ext (
    type_code_ext text NOT NULL,
    description text,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.transaction_type_ext OWNER TO postgres;

--
-- TOC entry 351 (class 1259 OID 21243)
-- Name: transaction_type_ledgers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transaction_type_ledgers (
    type_code_ledger text NOT NULL,
    name text NOT NULL
);


ALTER TABLE public.transaction_type_ledgers OWNER TO postgres;

--
-- TOC entry 352 (class 1259 OID 21248)
-- Name: unit_of_measure_translations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.unit_of_measure_translations (
    language_code text NOT NULL,
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    updated_at timestamp with time zone DEFAULT now(),
    name text NOT NULL,
    base_id bigint NOT NULL
);


ALTER TABLE public.unit_of_measure_translations OWNER TO postgres;

--
-- TOC entry 353 (class 1259 OID 21255)
-- Name: unit_of_measure_translations_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.unit_of_measure_translations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.unit_of_measure_translations_id_seq OWNER TO postgres;

--
-- TOC entry 4612 (class 0 OID 0)
-- Dependencies: 353
-- Name: unit_of_measure_translations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.unit_of_measure_translations_id_seq OWNED BY public.unit_of_measure_translations.id;


--
-- TOC entry 354 (class 1259 OID 21256)
-- Name: unit_of_measures; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.unit_of_measures (
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    updated_at timestamp with time zone DEFAULT now(),
    code text NOT NULL,
    enabled boolean,
    company_id bigint
);


ALTER TABLE public.unit_of_measures OWNER TO postgres;

--
-- TOC entry 355 (class 1259 OID 21263)
-- Name: unit_of_measures_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.unit_of_measures_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.unit_of_measures_id_seq OWNER TO postgres;

--
-- TOC entry 4613 (class 0 OID 0)
-- Dependencies: 355
-- Name: unit_of_measures_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.unit_of_measures_id_seq OWNED BY public.unit_of_measures.id;


--
-- TOC entry 356 (class 1259 OID 21264)
-- Name: user_relations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_relations (
    user_id bigint NOT NULL,
    profile_id bigint NOT NULL,
    company_id bigint NOT NULL,
    role_id bigint NOT NULL,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.user_relations OWNER TO postgres;

--
-- TOC entry 357 (class 1259 OID 21269)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    identifier text NOT NULL,
    last_login timestamp with time zone,
    password_hash text NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 358 (class 1259 OID 21277)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- TOC entry 4614 (class 0 OID 0)
-- Dependencies: 358
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 359 (class 1259 OID 21278)
-- Name: ware_houses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ware_houses (
    id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    updated_at timestamp with time zone DEFAULT now(),
    name text NOT NULL,
    parent_id bigint,
    ordinal bigint NOT NULL,
    company_id bigint NOT NULL,
    enabled boolean NOT NULL,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    is_group boolean DEFAULT false NOT NULL
);


ALTER TABLE public.ware_houses OWNER TO postgres;

--
-- TOC entry 360 (class 1259 OID 21287)
-- Name: ware_houses_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ware_houses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ware_houses_id_seq OWNER TO postgres;

--
-- TOC entry 4615 (class 0 OID 0)
-- Dependencies: 360
-- Name: ware_houses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ware_houses_id_seq OWNED BY public.ware_houses.id;


--
-- TOC entry 3704 (class 2604 OID 21288)
-- Name: actions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.actions ALTER COLUMN id SET DEFAULT nextval('public.actions_id_seq'::regclass);


--
-- TOC entry 3706 (class 2604 OID 21289)
-- Name: addresses id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.addresses ALTER COLUMN id SET DEFAULT nextval('public.addresses_id_seq'::regclass);


--
-- TOC entry 3716 (class 2604 OID 21290)
-- Name: companies id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.companies ALTER COLUMN id SET DEFAULT nextval('public.companies_id_seq'::regclass);


--
-- TOC entry 3734 (class 2604 OID 21291)
-- Name: entities id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.entities ALTER COLUMN id SET DEFAULT nextval('public.entities_id_seq'::regclass);


--
-- TOC entry 3747 (class 2604 OID 21292)
-- Name: item_attribute_values id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_attribute_values ALTER COLUMN id SET DEFAULT nextval('public.item_attribute_values_id_seq'::regclass);


--
-- TOC entry 3750 (class 2604 OID 21293)
-- Name: item_attributes id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_attributes ALTER COLUMN id SET DEFAULT nextval('public.item_attributes_id_seq'::regclass);


--
-- TOC entry 3758 (class 2604 OID 21294)
-- Name: item_lines id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_lines ALTER COLUMN id SET DEFAULT nextval('public.item_lines_id_seq'::regclass);


--
-- TOC entry 3768 (class 2604 OID 21295)
-- Name: item_prices id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_prices ALTER COLUMN id SET DEFAULT nextval('public.item_prices_id_seq'::regclass);


--
-- TOC entry 3773 (class 2604 OID 21296)
-- Name: items id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items ALTER COLUMN id SET DEFAULT nextval('public.items_id_seq'::regclass);


--
-- TOC entry 3821 (class 2604 OID 21297)
-- Name: plugins id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugins ALTER COLUMN id SET DEFAULT nextval('public.plugins_id_seq'::regclass);


--
-- TOC entry 3761 (class 2604 OID 21298)
-- Name: price_lists id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.price_lists ALTER COLUMN id SET DEFAULT nextval('public.item_price_lists_id_seq'::regclass);


--
-- TOC entry 3856 (class 2604 OID 21299)
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);


--
-- TOC entry 3871 (class 2604 OID 21300)
-- Name: stock_movements id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_movements ALTER COLUMN id SET DEFAULT nextval('public.stock_movements_id_seq'::regclass);


--
-- TOC entry 3889 (class 2604 OID 21301)
-- Name: taxes id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.taxes ALTER COLUMN id SET DEFAULT nextval('public.taxes_id_seq'::regclass);


--
-- TOC entry 3904 (class 2604 OID 21302)
-- Name: unit_of_measure_translations id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.unit_of_measure_translations ALTER COLUMN id SET DEFAULT nextval('public.unit_of_measure_translations_id_seq'::regclass);


--
-- TOC entry 3907 (class 2604 OID 21303)
-- Name: unit_of_measures id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.unit_of_measures ALTER COLUMN id SET DEFAULT nextval('public.unit_of_measures_id_seq'::regclass);


--
-- TOC entry 3912 (class 2604 OID 21304)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- TOC entry 3916 (class 2604 OID 21305)
-- Name: ware_houses id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ware_houses ALTER COLUMN id SET DEFAULT nextval('public.ware_houses_id_seq'::regclass);


--
-- TOC entry 3922 (class 2606 OID 21308)
-- Name: account_settings account_settings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.account_settings
    ADD CONSTRAINT account_settings_pkey PRIMARY KEY (id);


--
-- TOC entry 3925 (class 2606 OID 21310)
-- Name: account_statements account_statements_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.account_statements
    ADD CONSTRAINT account_statements_pkey PRIMARY KEY (account_no, date);


--
-- TOC entry 3927 (class 2606 OID 21312)
-- Name: account_type_exts account_type_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.account_type_exts
    ADD CONSTRAINT account_type_pkey PRIMARY KEY (account_type);


--
-- TOC entry 3929 (class 2606 OID 21314)
-- Name: account_types account_types_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.account_types
    ADD CONSTRAINT account_types_pkey PRIMARY KEY (account_type);


--
-- TOC entry 3931 (class 2606 OID 21316)
-- Name: actions actions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.actions
    ADD CONSTRAINT actions_pkey PRIMARY KEY (id);


--
-- TOC entry 3933 (class 2606 OID 21318)
-- Name: activities activities_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.activities
    ADD CONSTRAINT activities_pkey PRIMARY KEY (id);


--
-- TOC entry 3937 (class 2606 OID 21320)
-- Name: addresses addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_pkey PRIMARY KEY (id);


--
-- TOC entry 3941 (class 2606 OID 21322)
-- Name: batch_bundles batch_bundles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.batch_bundles
    ADD CONSTRAINT batch_bundles_pkey PRIMARY KEY (id);


--
-- TOC entry 3943 (class 2606 OID 21324)
-- Name: charges_template charges_template_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.charges_template
    ADD CONSTRAINT charges_template_pkey PRIMARY KEY (id);


--
-- TOC entry 3946 (class 2606 OID 21326)
-- Name: companies companies_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT companies_pkey PRIMARY KEY (id);


--
-- TOC entry 3950 (class 2606 OID 21328)
-- Name: company_defaults company_defaults_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.company_defaults
    ADD CONSTRAINT company_defaults_pkey PRIMARY KEY (id);


--
-- TOC entry 3953 (class 2606 OID 21330)
-- Name: company_entities company_entities_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.company_entities
    ADD CONSTRAINT company_entities_pkey PRIMARY KEY (company_id, entity_id);


--
-- TOC entry 3957 (class 2606 OID 21332)
-- Name: contacts contact_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.contacts
    ADD CONSTRAINT contact_pkey PRIMARY KEY (id);


--
-- TOC entry 3961 (class 2606 OID 21334)
-- Name: cost_centers cost_centers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cost_centers
    ADD CONSTRAINT cost_centers_pkey PRIMARY KEY (id);


--
-- TOC entry 4138 (class 2606 OID 21336)
-- Name: r_courts courts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.r_courts
    ADD CONSTRAINT courts_pkey PRIMARY KEY (id);


--
-- TOC entry 3963 (class 2606 OID 21338)
-- Name: currencies currencies_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.currencies
    ADD CONSTRAINT currencies_pkey PRIMARY KEY ("code ");


--
-- TOC entry 3965 (class 2606 OID 21340)
-- Name: currency_exchanges currency_exchanges_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.currency_exchanges
    ADD CONSTRAINT currency_exchanges_pkey PRIMARY KEY (id);


--
-- TOC entry 3968 (class 2606 OID 21342)
-- Name: customers customers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_pkey PRIMARY KEY (id);


--
-- TOC entry 3972 (class 2606 OID 21344)
-- Name: delivery_line_items delivery_line_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.delivery_line_items
    ADD CONSTRAINT delivery_line_items_pkey PRIMARY KEY (item_line_id);


--
-- TOC entry 3974 (class 2606 OID 21346)
-- Name: entities entities_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.entities
    ADD CONSTRAINT entities_pkey PRIMARY KEY (id);


--
-- TOC entry 3976 (class 2606 OID 21348)
-- Name: entity_types entity_types_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.entity_types
    ADD CONSTRAINT entity_types_pkey PRIMARY KEY (entity_type);


--
-- TOC entry 3980 (class 2606 OID 21350)
-- Name: groups groups_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.groups
    ADD CONSTRAINT groups_pkey PRIMARY KEY (id);


--
-- TOC entry 3983 (class 2606 OID 21352)
-- Name: invoiced_item_lines invoiced_item_lines_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.invoiced_item_lines
    ADD CONSTRAINT invoiced_item_lines_pkey PRIMARY KEY (invoice_id, item_line);


--
-- TOC entry 3986 (class 2606 OID 21354)
-- Name: invoices invoices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT invoices_pkey PRIMARY KEY (id);


--
-- TOC entry 3989 (class 2606 OID 21356)
-- Name: item_attribute_values item_attribute_values_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_attribute_values
    ADD CONSTRAINT item_attribute_values_pkey PRIMARY KEY (id);


--
-- TOC entry 3993 (class 2606 OID 21358)
-- Name: item_attributes item_attributes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_attributes
    ADD CONSTRAINT item_attributes_pkey PRIMARY KEY (id);


--
-- TOC entry 3995 (class 2606 OID 21360)
-- Name: item_inventory_settings item_inventory_settings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_inventory_settings
    ADD CONSTRAINT item_inventory_settings_pkey PRIMARY KEY (item_id);


--
-- TOC entry 4004 (class 2606 OID 21362)
-- Name: item_line_stock_entries item_line_stock_entries_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_line_stock_entries
    ADD CONSTRAINT item_line_stock_entries_pkey PRIMARY KEY (id);


--
-- TOC entry 4006 (class 2606 OID 21364)
-- Name: item_line_types item_line_types_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_line_types
    ADD CONSTRAINT item_line_types_pkey PRIMARY KEY (type);


--
-- TOC entry 4001 (class 2606 OID 21366)
-- Name: item_line_receipts item_line_warehouse_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_line_receipts
    ADD CONSTRAINT item_line_warehouse_pkey PRIMARY KEY (item_line);


--
-- TOC entry 4009 (class 2606 OID 21368)
-- Name: item_lines item_lines_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_lines
    ADD CONSTRAINT item_lines_pkey PRIMARY KEY (id);


--
-- TOC entry 4012 (class 2606 OID 21370)
-- Name: price_lists item_price_lists_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.price_lists
    ADD CONSTRAINT item_price_lists_pkey PRIMARY KEY (id);


--
-- TOC entry 4014 (class 2606 OID 21372)
-- Name: item_price_plugins item_price_plugins_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_price_plugins
    ADD CONSTRAINT item_price_plugins_pkey PRIMARY KEY (plugin, base_id);


--
-- TOC entry 4017 (class 2606 OID 21374)
-- Name: item_prices item_prices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_prices
    ADD CONSTRAINT item_prices_pkey PRIMARY KEY (id);


--
-- TOC entry 4019 (class 2606 OID 21376)
-- Name: item_variants item_variants_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_variants
    ADD CONSTRAINT item_variants_pkey PRIMARY KEY (item_id, variant_id);


--
-- TOC entry 4024 (class 2606 OID 21378)
-- Name: items items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_pkey PRIMARY KEY (id);


--
-- TOC entry 4027 (class 2606 OID 21380)
-- Name: journal_entries journal_entries_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.journal_entries
    ADD CONSTRAINT journal_entries_pkey PRIMARY KEY (id);


--
-- TOC entry 4030 (class 2606 OID 21382)
-- Name: journal_entry_lines journal_entry_lines_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.journal_entry_lines
    ADD CONSTRAINT journal_entry_lines_pkey PRIMARY KEY (id);


--
-- TOC entry 4033 (class 2606 OID 21384)
-- Name: key_values key_values_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.key_values
    ADD CONSTRAINT key_values_pkey PRIMARY KEY (id);


--
-- TOC entry 4037 (class 2606 OID 21386)
-- Name: ledger_accounts ledger_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ledger_accounts
    ADD CONSTRAINT ledger_accounts_pkey PRIMARY KEY (ledger_id);


--
-- TOC entry 4039 (class 2606 OID 21388)
-- Name: ledger_statements ledger_statements_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ledger_statements
    ADD CONSTRAINT ledger_statements_pkey PRIMARY KEY (date, ledger_no);


--
-- TOC entry 4041 (class 2606 OID 21390)
-- Name: ledger_types ledger_type_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ledger_types
    ADD CONSTRAINT ledger_type_pkey PRIMARY KEY (ledger_type);


--
-- TOC entry 4048 (class 2606 OID 21392)
-- Name: ledgers ledgers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ledgers
    ADD CONSTRAINT ledgers_pkey PRIMARY KEY (id);


--
-- TOC entry 4050 (class 2606 OID 21394)
-- Name: module_sections module_sections_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.module_sections
    ADD CONSTRAINT module_sections_pkey PRIMARY KEY (id);


--
-- TOC entry 4053 (class 2606 OID 21396)
-- Name: modules modules_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.modules
    ADD CONSTRAINT modules_pkey PRIMARY KEY (id);


--
-- TOC entry 4058 (class 2606 OID 21398)
-- Name: orders order_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT order_pkey PRIMARY KEY (id);


--
-- TOC entry 4061 (class 2606 OID 21400)
-- Name: parties parties_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.parties
    ADD CONSTRAINT parties_pkey PRIMARY KEY (id);


--
-- TOC entry 4065 (class 2606 OID 21402)
-- Name: party_addresses party_addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.party_addresses
    ADD CONSTRAINT party_addresses_pkey PRIMARY KEY (party_id, address_id);


--
-- TOC entry 4069 (class 2606 OID 21404)
-- Name: party_payments party_payments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.party_payments
    ADD CONSTRAINT party_payments_pkey PRIMARY KEY (party_id, payment_id);


--
-- TOC entry 4073 (class 2606 OID 21406)
-- Name: party_references party_references_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.party_references
    ADD CONSTRAINT party_references_pkey PRIMARY KEY (party_id, reference_id);


--
-- TOC entry 4075 (class 2606 OID 21408)
-- Name: party_types party_types_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.party_types
    ADD CONSTRAINT party_types_pkey PRIMARY KEY (code);


--
-- TOC entry 4079 (class 2606 OID 21410)
-- Name: payment_references payment_references_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_references
    ADD CONSTRAINT payment_references_pkey PRIMARY KEY (id);


--
-- TOC entry 4085 (class 2606 OID 21412)
-- Name: payment_terms_lines payment_terms_lines_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_terms_lines
    ADD CONSTRAINT payment_terms_lines_pkey PRIMARY KEY (id);


--
-- TOC entry 4082 (class 2606 OID 21414)
-- Name: payment_terms payment_terms_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_terms
    ADD CONSTRAINT payment_terms_pkey PRIMARY KEY (id);


--
-- TOC entry 4088 (class 2606 OID 21416)
-- Name: payment_terms_templates payment_terms_templates_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_terms_templates
    ADD CONSTRAINT payment_terms_templates_pkey PRIMARY KEY (id);


--
-- TOC entry 4093 (class 2606 OID 21418)
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);


--
-- TOC entry 4097 (class 2606 OID 21420)
-- Name: piano_form piano_form_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.piano_form
    ADD CONSTRAINT piano_form_pkey PRIMARY KEY (id);


--
-- TOC entry 4099 (class 2606 OID 21422)
-- Name: plugins plugins_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugins
    ADD CONSTRAINT plugins_pkey PRIMARY KEY (id);


--
-- TOC entry 4101 (class 2606 OID 21424)
-- Name: pricing_line_items priced_item_lines_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pricing_line_items
    ADD CONSTRAINT priced_item_lines_pkey PRIMARY KEY (id);


--
-- TOC entry 4103 (class 2606 OID 21426)
-- Name: pricing_charges pricer_defaults_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pricing_charges
    ADD CONSTRAINT pricer_defaults_pkey PRIMARY KEY (id);


--
-- TOC entry 4105 (class 2606 OID 21428)
-- Name: pricings pricings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pricings
    ADD CONSTRAINT pricings_pkey PRIMARY KEY (id);


--
-- TOC entry 4109 (class 2606 OID 21430)
-- Name: profiles profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.profiles
    ADD CONSTRAINT profiles_pkey PRIMARY KEY (id);


--
-- TOC entry 4111 (class 2606 OID 21432)
-- Name: progress_invoices progress_invoices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.progress_invoices
    ADD CONSTRAINT progress_invoices_pkey PRIMARY KEY (invoice_id);


--
-- TOC entry 4113 (class 2606 OID 21434)
-- Name: progress_orders progress_orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.progress_orders
    ADD CONSTRAINT progress_orders_pkey PRIMARY KEY (order_id);


--
-- TOC entry 4116 (class 2606 OID 21436)
-- Name: projects project_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.projects
    ADD CONSTRAINT project_pkey PRIMARY KEY (id);


--
-- TOC entry 4121 (class 2606 OID 21438)
-- Name: purchase_records purchase_records_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.purchase_records
    ADD CONSTRAINT purchase_records_pkey PRIMARY KEY (id);


--
-- TOC entry 4124 (class 2606 OID 21440)
-- Name: quotations quotations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.quotations
    ADD CONSTRAINT quotations_pkey PRIMARY KEY (id);


--
-- TOC entry 4126 (class 2606 OID 21442)
-- Name: r_booking_events r_booking_events_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.r_booking_events
    ADD CONSTRAINT r_booking_events_pkey PRIMARY KEY (booking_id, event_id);


--
-- TOC entry 4128 (class 2606 OID 21444)
-- Name: r_booking_prices r_booking_prices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.r_booking_prices
    ADD CONSTRAINT r_booking_prices_pkey PRIMARY KEY (id);


--
-- TOC entry 4130 (class 2606 OID 21446)
-- Name: r_booking_slots r_booking_slots_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.r_booking_slots
    ADD CONSTRAINT r_booking_slots_pkey PRIMARY KEY (id);


--
-- TOC entry 4133 (class 2606 OID 21448)
-- Name: r_bookings r_bookings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.r_bookings
    ADD CONSTRAINT r_bookings_pkey PRIMARY KEY (id);


--
-- TOC entry 4136 (class 2606 OID 21450)
-- Name: r_court_rates r_court_rates_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.r_court_rates
    ADD CONSTRAINT r_court_rates_pkey PRIMARY KEY (id);


--
-- TOC entry 4143 (class 2606 OID 21452)
-- Name: r_events r_events_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.r_events
    ADD CONSTRAINT r_events_pkey PRIMARY KEY (id);


--
-- TOC entry 4148 (class 2606 OID 21454)
-- Name: receipts receipts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT receipts_pkey PRIMARY KEY (id);


--
-- TOC entry 4150 (class 2606 OID 21456)
-- Name: role_actions role_definitions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_actions
    ADD CONSTRAINT role_definitions_pkey PRIMARY KEY (role_id, action_id);


--
-- TOC entry 4152 (class 2606 OID 21458)
-- Name: role_templates role_templates_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_templates
    ADD CONSTRAINT role_templates_pkey PRIMARY KEY (id);


--
-- TOC entry 4156 (class 2606 OID 21460)
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- TOC entry 4161 (class 2606 OID 21462)
-- Name: sales_records sales_records_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales_records
    ADD CONSTRAINT sales_records_pkey PRIMARY KEY (id);


--
-- TOC entry 4163 (class 2606 OID 21464)
-- Name: serial_no_transactions serial_no_transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.serial_no_transactions
    ADD CONSTRAINT serial_no_transactions_pkey PRIMARY KEY (id);


--
-- TOC entry 4167 (class 2606 OID 21466)
-- Name: serial_nos serial_nos_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.serial_nos
    ADD CONSTRAINT serial_nos_pkey PRIMARY KEY (id);


--
-- TOC entry 4169 (class 2606 OID 21468)
-- Name: states states_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.states
    ADD CONSTRAINT states_pkey PRIMARY KEY (state);


--
-- TOC entry 4173 (class 2606 OID 21470)
-- Name: stock_entries stock_entries_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_entries
    ADD CONSTRAINT stock_entries_pkey PRIMARY KEY (id);


--
-- TOC entry 4176 (class 2606 OID 21472)
-- Name: stock_levels stock_levels_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_levels
    ADD CONSTRAINT stock_levels_pkey PRIMARY KEY (ware_house_id, item_id);


--
-- TOC entry 4179 (class 2606 OID 21474)
-- Name: stock_movements stock_movements_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_movements
    ADD CONSTRAINT stock_movements_pkey PRIMARY KEY (id);


--
-- TOC entry 4171 (class 2606 OID 21476)
-- Name: stock_defaults stock_setting_defaults_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_defaults
    ADD CONSTRAINT stock_setting_defaults_pkey PRIMARY KEY (id);


--
-- TOC entry 4182 (class 2606 OID 21478)
-- Name: stock_settings stock_settings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_settings
    ADD CONSTRAINT stock_settings_pkey PRIMARY KEY (id);


--
-- TOC entry 4186 (class 2606 OID 21480)
-- Name: stock_transactions stock_transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_transactions
    ADD CONSTRAINT stock_transactions_pkey PRIMARY KEY (id);


--
-- TOC entry 4190 (class 2606 OID 21482)
-- Name: supplier_orders supplier_orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.supplier_orders
    ADD CONSTRAINT supplier_orders_pkey PRIMARY KEY (supplier_id, order_id);


--
-- TOC entry 4196 (class 2606 OID 21484)
-- Name: suppliers suppliers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_pkey PRIMARY KEY (id);


--
-- TOC entry 4199 (class 2606 OID 21486)
-- Name: tax_and_charge_lines tax_and_charge_lines_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tax_and_charge_lines
    ADD CONSTRAINT tax_and_charge_lines_pkey PRIMARY KEY (id);


--
-- TOC entry 4202 (class 2606 OID 21488)
-- Name: taxes taxes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.taxes
    ADD CONSTRAINT taxes_pkey PRIMARY KEY (id);


--
-- TOC entry 4205 (class 2606 OID 21490)
-- Name: terms_and_conditions term_and_conditions_templates_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.terms_and_conditions
    ADD CONSTRAINT term_and_conditions_templates_pkey PRIMARY KEY (id);


--
-- TOC entry 4207 (class 2606 OID 21492)
-- Name: transaction_accounts transaction_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction_accounts
    ADD CONSTRAINT transaction_accounts_pkey PRIMARY KEY (ledger_no, datetime);


--
-- TOC entry 4218 (class 2606 OID 21494)
-- Name: transaction_type_ext transaction_type_ext_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction_type_ext
    ADD CONSTRAINT transaction_type_ext_pkey PRIMARY KEY (type_code_ext);


--
-- TOC entry 4220 (class 2606 OID 21496)
-- Name: transaction_type_ledgers transaction_type_ledgers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction_type_ledgers
    ADD CONSTRAINT transaction_type_ledgers_pkey PRIMARY KEY (type_code_ledger);


--
-- TOC entry 4216 (class 2606 OID 21498)
-- Name: transaction_type_de transaction_type_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction_type_de
    ADD CONSTRAINT transaction_type_pkey PRIMARY KEY (type_code);


--
-- TOC entry 4214 (class 2606 OID 21500)
-- Name: transaction_ledgers transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction_ledgers
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- TOC entry 4233 (class 2606 OID 21502)
-- Name: users uni_users_identifier; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_identifier UNIQUE (identifier);


--
-- TOC entry 4223 (class 2606 OID 21504)
-- Name: unit_of_measure_translations unit_of_measure_translations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.unit_of_measure_translations
    ADD CONSTRAINT unit_of_measure_translations_pkey PRIMARY KEY (id);


--
-- TOC entry 4226 (class 2606 OID 21506)
-- Name: unit_of_measures unit_of_measures_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.unit_of_measures
    ADD CONSTRAINT unit_of_measures_pkey PRIMARY KEY (id);


--
-- TOC entry 4229 (class 2606 OID 21508)
-- Name: user_relations user_relations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_relations
    ADD CONSTRAINT user_relations_pkey PRIMARY KEY (profile_id, user_id, company_id, role_id);


--
-- TOC entry 4235 (class 2606 OID 21510)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 4238 (class 2606 OID 21512)
-- Name: ware_houses ware_houses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ware_houses
    ADD CONSTRAINT ware_houses_pkey PRIMARY KEY (id);


--
-- TOC entry 4066 (class 1259 OID 21513)
-- Name: fki_;; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "fki_;" ON public.party_payments USING btree (payment_id);


--
-- TOC entry 4007 (class 1259 OID 21514)
-- Name: fki_F; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "fki_F" ON public.item_lines USING btree (id);


--
-- TOC entry 4208 (class 1259 OID 21515)
-- Name: fki_T; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "fki_T" ON public.transaction_ledgers USING btree (ledger);


--
-- TOC entry 4180 (class 1259 OID 21516)
-- Name: fki_com; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_com ON public.stock_settings USING btree (company_id);


--
-- TOC entry 3954 (class 1259 OID 21517)
-- Name: fki_company_en; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_company_en ON public.company_entities USING btree (company_id);


--
-- TOC entry 3969 (class 1259 OID 21518)
-- Name: fki_customers_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_customers_company ON public.customers USING btree (company_id);


--
-- TOC entry 3970 (class 1259 OID 21519)
-- Name: fki_customers_group; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_customers_group ON public.customers USING btree (group_id);


--
-- TOC entry 4203 (class 1259 OID 21520)
-- Name: fki_d; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_d ON public.terms_and_conditions USING btree (company_id);


--
-- TOC entry 3923 (class 1259 OID 21521)
-- Name: fki_fk_acc_setts_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_acc_setts_company ON public.account_settings USING btree (company_id);


--
-- TOC entry 3934 (class 1259 OID 21522)
-- Name: fki_fk_activities_party; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_activities_party ON public.activities USING btree (party_id);


--
-- TOC entry 3935 (class 1259 OID 21523)
-- Name: fki_fk_activities_profile; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_activities_profile ON public.activities USING btree (profile_id);


--
-- TOC entry 3938 (class 1259 OID 21524)
-- Name: fki_fk_addresses_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_addresses_company ON public.addresses USING btree (company_id);


--
-- TOC entry 4131 (class 1259 OID 21525)
-- Name: fki_fk_bookings_court; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_bookings_court ON public.r_bookings USING btree (court_id);


--
-- TOC entry 3951 (class 1259 OID 21526)
-- Name: fki_fk_cd_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_cd_company ON public.company_defaults USING btree (company_id);


--
-- TOC entry 3955 (class 1259 OID 21527)
-- Name: fki_fk_ce_entity; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_ce_entity ON public.company_entities USING btree (entity_id);


--
-- TOC entry 3944 (class 1259 OID 21528)
-- Name: fki_fk_charges_template_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_charges_template_company ON public.charges_template USING btree (company_id);


--
-- TOC entry 3947 (class 1259 OID 21529)
-- Name: fki_fk_companies_party; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_companies_party ON public.companies USING btree (id);


--
-- TOC entry 4140 (class 1259 OID 21530)
-- Name: fki_fk_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_company ON public.r_events USING btree (company_id);


--
-- TOC entry 3958 (class 1259 OID 21531)
-- Name: fki_fk_contacts_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_contacts_company ON public.contacts USING btree (company_id);


--
-- TOC entry 3959 (class 1259 OID 21532)
-- Name: fki_fk_contacts_party; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_contacts_party ON public.contacts USING btree (id);


--
-- TOC entry 4134 (class 1259 OID 21533)
-- Name: fki_fk_court_rates_court; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_court_rates_court ON public.r_court_rates USING btree (court_id);


--
-- TOC entry 4139 (class 1259 OID 21534)
-- Name: fki_fk_courts_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_courts_company ON public.r_courts USING btree (company_id);


--
-- TOC entry 3966 (class 1259 OID 21535)
-- Name: fki_fk_currency_exchanges_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_currency_exchanges_company ON public.currency_exchanges USING btree (company_id);


--
-- TOC entry 3981 (class 1259 OID 21536)
-- Name: fki_fk_inv_item_line; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_inv_item_line ON public.invoiced_item_lines USING btree (item_line);


--
-- TOC entry 3984 (class 1259 OID 21537)
-- Name: fki_fk_invoices_currency; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_invoices_currency ON public.invoices USING btree (currency);


--
-- TOC entry 3990 (class 1259 OID 21538)
-- Name: fki_fk_item_attributes_party; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_item_attributes_party ON public.item_attributes USING btree (id);


--
-- TOC entry 3996 (class 1259 OID 21539)
-- Name: fki_fk_item_lin; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_item_lin ON public.item_line_receipts USING btree (item_line);


--
-- TOC entry 4002 (class 1259 OID 21540)
-- Name: fki_fk_item_line_stock_entries_item_line; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_item_line_stock_entries_item_line ON public.item_line_stock_entries USING btree (item_line);


--
-- TOC entry 3997 (class 1259 OID 21541)
-- Name: fki_fk_item_line_warehouse_accepted; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_item_line_warehouse_accepted ON public.item_line_receipts USING btree (accepted_warehouse);


--
-- TOC entry 3998 (class 1259 OID 21542)
-- Name: fki_fk_item_line_warehouse_item_line; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_item_line_warehouse_item_line ON public.item_line_receipts USING btree (item_line);


--
-- TOC entry 3999 (class 1259 OID 21543)
-- Name: fki_fk_item_line_warehouse_rejected; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_item_line_warehouse_rejected ON public.item_line_receipts USING btree (rejected_warehouse);


--
-- TOC entry 4020 (class 1259 OID 21544)
-- Name: fki_fk_items_group; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_items_group ON public.items USING btree (group_id);


--
-- TOC entry 4021 (class 1259 OID 21545)
-- Name: fki_fk_items_items_party; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_items_items_party ON public.items USING btree (id);


--
-- TOC entry 4028 (class 1259 OID 21546)
-- Name: fki_fk_jel_journal_entry; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_jel_journal_entry ON public.journal_entry_lines USING btree (journal_entry_id);


--
-- TOC entry 4025 (class 1259 OID 21547)
-- Name: fki_fk_journal_entries_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_journal_entries_company ON public.journal_entries USING btree (id);


--
-- TOC entry 4042 (class 1259 OID 21548)
-- Name: fki_fk_ledgers_account_type; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_ledgers_account_type ON public.ledgers USING btree (account_type);


--
-- TOC entry 4034 (class 1259 OID 21549)
-- Name: fki_fk_ledgers_accounts_ledger; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_ledgers_accounts_ledger ON public.ledger_accounts USING btree (ledger_id);


--
-- TOC entry 4043 (class 1259 OID 21550)
-- Name: fki_fk_ledgers_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_ledgers_company ON public.ledgers USING btree (company_id);


--
-- TOC entry 4044 (class 1259 OID 21551)
-- Name: fki_fk_ledgers_parent; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_ledgers_parent ON public.ledgers USING btree (ledger_parent);


--
-- TOC entry 4051 (class 1259 OID 21552)
-- Name: fki_fk_modules_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_modules_company ON public.modules USING btree (company_id);


--
-- TOC entry 4054 (class 1259 OID 21553)
-- Name: fki_fk_orders_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_orders_company ON public.orders USING btree (company_id);


--
-- TOC entry 4055 (class 1259 OID 21554)
-- Name: fki_fk_orders_party; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_orders_party ON public.orders USING btree (id);


--
-- TOC entry 4056 (class 1259 OID 21555)
-- Name: fki_fk_orders_party2; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_orders_party2 ON public.orders USING btree (party_id);


--
-- TOC entry 4062 (class 1259 OID 21556)
-- Name: fki_fk_party_addresses_address; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_party_addresses_address ON public.party_addresses USING btree (address_id);


--
-- TOC entry 4080 (class 1259 OID 21557)
-- Name: fki_fk_payment_tems_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_payment_tems_company ON public.payment_terms USING btree (company_id);


--
-- TOC entry 4086 (class 1259 OID 21558)
-- Name: fki_fk_payment_terms_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_payment_terms_company ON public.payment_terms_templates USING btree (company_id);


--
-- TOC entry 4094 (class 1259 OID 21559)
-- Name: fki_fk_piano_form_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_piano_form_company ON public.piano_form USING btree (company_id);


--
-- TOC entry 4095 (class 1259 OID 21560)
-- Name: fki_fk_piano_form_party_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_piano_form_party_id ON public.piano_form USING btree (id);


--
-- TOC entry 4076 (class 1259 OID 21561)
-- Name: fki_fk_pr_party; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_pr_party ON public.payment_references USING btree (party_id);


--
-- TOC entry 4077 (class 1259 OID 21562)
-- Name: fki_fk_pr_payment; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_pr_payment ON public.payment_references USING btree (payment_id);


--
-- TOC entry 4106 (class 1259 OID 21563)
-- Name: fki_fk_profile_party; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_profile_party ON public.profiles USING btree (id);


--
-- TOC entry 4114 (class 1259 OID 21564)
-- Name: fki_fk_projects_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_projects_company ON public.projects USING btree (company_id);


--
-- TOC entry 4083 (class 1259 OID 21565)
-- Name: fki_fk_ptl_doc_party; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_ptl_doc_party ON public.payment_terms_lines USING btree (document_id);


--
-- TOC entry 4117 (class 1259 OID 21566)
-- Name: fki_fk_purchase_records_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_purchase_records_company ON public.purchase_records USING btree (company_id);


--
-- TOC entry 4118 (class 1259 OID 21567)
-- Name: fki_fk_purchase_records_supplier; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_purchase_records_supplier ON public.purchase_records USING btree (supplier_id);


--
-- TOC entry 4122 (class 1259 OID 21568)
-- Name: fki_fk_quotations_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_quotations_company ON public.quotations USING btree (company_id);


--
-- TOC entry 4141 (class 1259 OID 21569)
-- Name: fki_fk_r_events_party_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_r_events_party_id ON public.r_events USING btree (id);


--
-- TOC entry 4144 (class 1259 OID 21570)
-- Name: fki_fk_receipts_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_receipts_company ON public.receipts USING btree (company_id);


--
-- TOC entry 4145 (class 1259 OID 21571)
-- Name: fki_fk_receipts_party2; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_receipts_party2 ON public.receipts USING btree (party_id);


--
-- TOC entry 4146 (class 1259 OID 21572)
-- Name: fki_fk_receipts_status; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_receipts_status ON public.receipts USING btree (status);


--
-- TOC entry 4070 (class 1259 OID 21573)
-- Name: fki_fk_references_party; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_references_party ON public.party_references USING btree (party_id);


--
-- TOC entry 4071 (class 1259 OID 21574)
-- Name: fki_fk_references_reference; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_references_reference ON public.party_references USING btree (reference_id);


--
-- TOC entry 4153 (class 1259 OID 21575)
-- Name: fki_fk_roles_party; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_roles_party ON public.roles USING btree (id);


--
-- TOC entry 4157 (class 1259 OID 21576)
-- Name: fki_fk_sales_records_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_sales_records_company ON public.sales_records USING btree (company_id);


--
-- TOC entry 4158 (class 1259 OID 21577)
-- Name: fki_fk_sales_records_invoice; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_sales_records_invoice ON public.sales_records USING btree (invoice_id);


--
-- TOC entry 4164 (class 1259 OID 21578)
-- Name: fki_fk_sn_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_sn_id ON public.serial_nos USING btree (id);


--
-- TOC entry 4165 (class 1259 OID 21579)
-- Name: fki_fk_sn_serial_no; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX fki_fk_sn_serial_no ON public.serial_nos USING btree (serial_no);


--
-- TOC entry 4183 (class 1259 OID 21580)
-- Name: fki_fk_stock_tx_item; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_stock_tx_item ON public.stock_transactions USING btree (item_id);


--
-- TOC entry 4184 (class 1259 OID 21581)
-- Name: fki_fk_stock_tx_warehouse; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_stock_tx_warehouse ON public.stock_transactions USING btree (warehouse_id);


--
-- TOC entry 3977 (class 1259 OID 21582)
-- Name: fki_fk_supplier_groups_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_supplier_groups_company ON public.groups USING btree (company_id);


--
-- TOC entry 4187 (class 1259 OID 21583)
-- Name: fki_fk_supplier_orders_supplier; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_supplier_orders_supplier ON public.supplier_orders USING btree (supplier_id);


--
-- TOC entry 4191 (class 1259 OID 21584)
-- Name: fki_fk_suppliers_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_suppliers_company ON public.suppliers USING btree (company_id);


--
-- TOC entry 4192 (class 1259 OID 21585)
-- Name: fki_fk_suppliers_group; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_suppliers_group ON public.suppliers USING btree (group_id);


--
-- TOC entry 4193 (class 1259 OID 21586)
-- Name: fki_fk_suppliers_item_group; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_suppliers_item_group ON public.suppliers USING btree (group_id);


--
-- TOC entry 4188 (class 1259 OID 21587)
-- Name: fki_fk_suppliers_orders_order; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_suppliers_orders_order ON public.supplier_orders USING btree (order_id);


--
-- TOC entry 4194 (class 1259 OID 21588)
-- Name: fki_fk_suppliers_party; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_suppliers_party ON public.suppliers USING btree (id);


--
-- TOC entry 4209 (class 1259 OID 21589)
-- Name: fki_fk_tx_ledger_ledger; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_tx_ledger_ledger ON public.transaction_ledgers USING btree (ledger);


--
-- TOC entry 4210 (class 1259 OID 21590)
-- Name: fki_fk_tx_ledger_ledger_agst; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_tx_ledger_ledger_agst ON public.transaction_ledgers USING btree (ledger_against);


--
-- TOC entry 4211 (class 1259 OID 21591)
-- Name: fki_fk_tx_ledger_party; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_tx_ledger_party ON public.transaction_ledgers USING btree (party_id);


--
-- TOC entry 4212 (class 1259 OID 21592)
-- Name: fki_fk_tx_ledger_project; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_tx_ledger_project ON public.transaction_ledgers USING btree (project_id);


--
-- TOC entry 3978 (class 1259 OID 21593)
-- Name: fki_k; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_k ON public.groups USING btree (parent_id);


--
-- TOC entry 4031 (class 1259 OID 21594)
-- Name: fki_key_values_party; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_key_values_party ON public.key_values USING btree (party_id);


--
-- TOC entry 4045 (class 1259 OID 21595)
-- Name: fki_ledger; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_ledger ON public.ledgers USING btree (account_type);


--
-- TOC entry 4035 (class 1259 OID 21596)
-- Name: fki_ledger_accounts_ledger; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_ledger_accounts_ledger ON public.ledger_accounts USING btree (ledger_id);


--
-- TOC entry 4067 (class 1259 OID 21597)
-- Name: fki_party_payments_party; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_party_payments_party ON public.party_payments USING btree (party_id);


--
-- TOC entry 4089 (class 1259 OID 21598)
-- Name: fki_payments_company; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_payments_company ON public.payments USING btree (id);


--
-- TOC entry 4090 (class 1259 OID 21599)
-- Name: fki_payments_party; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_payments_party ON public.payments USING btree (id);


--
-- TOC entry 4091 (class 1259 OID 21600)
-- Name: fki_payments_party2; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_payments_party2 ON public.payments USING btree (party_id);


--
-- TOC entry 4159 (class 1259 OID 21601)
-- Name: fki_pg_sales_records_customer; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_pg_sales_records_customer ON public.sales_records USING btree (customer_id);


--
-- TOC entry 4119 (class 1259 OID 21602)
-- Name: fki_purchase_records_invoice; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_purchase_records_invoice ON public.purchase_records USING btree (invoice_id);


--
-- TOC entry 4197 (class 1259 OID 21603)
-- Name: fki_tax; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_tax ON public.tax_and_charge_lines USING btree (doc_party_id);


--
-- TOC entry 4230 (class 1259 OID 21604)
-- Name: fki_use; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_use ON public.users USING btree (id);


--
-- TOC entry 4046 (class 1259 OID 21605)
-- Name: fki_v; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_v ON public.ledgers USING btree (company_id);


--
-- TOC entry 3939 (class 1259 OID 21606)
-- Name: idx_addresses_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_addresses_deleted_at ON public.addresses USING btree (deleted_at);


--
-- TOC entry 3948 (class 1259 OID 21607)
-- Name: idx_companies_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_companies_deleted_at ON public.companies USING btree (deleted_at);


--
-- TOC entry 3987 (class 1259 OID 21608)
-- Name: idx_item_attribute_values_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_item_attribute_values_deleted_at ON public.item_attribute_values USING btree (deleted_at);


--
-- TOC entry 3991 (class 1259 OID 21609)
-- Name: idx_item_attributes_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_item_attributes_deleted_at ON public.item_attributes USING btree (deleted_at);


--
-- TOC entry 4010 (class 1259 OID 21610)
-- Name: idx_item_price_lists_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_item_price_lists_deleted_at ON public.price_lists USING btree (deleted_at);


--
-- TOC entry 4015 (class 1259 OID 21611)
-- Name: idx_item_prices_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_item_prices_deleted_at ON public.item_prices USING btree (deleted_at);


--
-- TOC entry 4022 (class 1259 OID 21612)
-- Name: idx_items_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_items_deleted_at ON public.items USING btree (deleted_at);


--
-- TOC entry 4059 (class 1259 OID 21613)
-- Name: idx_parties_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_parties_deleted_at ON public.parties USING btree (deleted_at);


--
-- TOC entry 4063 (class 1259 OID 21614)
-- Name: idx_party_addresses_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_party_addresses_deleted_at ON public.party_addresses USING btree (deleted_at);


--
-- TOC entry 4107 (class 1259 OID 21615)
-- Name: idx_profiles_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_profiles_deleted_at ON public.profiles USING btree (deleted_at);


--
-- TOC entry 4154 (class 1259 OID 21616)
-- Name: idx_roles_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_roles_deleted_at ON public.roles USING btree (deleted_at);


--
-- TOC entry 4174 (class 1259 OID 21617)
-- Name: idx_stock_levels_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_stock_levels_deleted_at ON public.stock_levels USING btree (deleted_at);


--
-- TOC entry 4177 (class 1259 OID 21618)
-- Name: idx_stock_movements_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_stock_movements_deleted_at ON public.stock_movements USING btree (deleted_at);


--
-- TOC entry 4200 (class 1259 OID 21619)
-- Name: idx_taxes_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_taxes_deleted_at ON public.taxes USING btree (deleted_at);


--
-- TOC entry 4221 (class 1259 OID 21620)
-- Name: idx_unit_of_measure_translations_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_unit_of_measure_translations_deleted_at ON public.unit_of_measure_translations USING btree (deleted_at);


--
-- TOC entry 4224 (class 1259 OID 21621)
-- Name: idx_unit_of_measures_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_unit_of_measures_deleted_at ON public.unit_of_measures USING btree (deleted_at);


--
-- TOC entry 4227 (class 1259 OID 21622)
-- Name: idx_user_relations_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_relations_deleted_at ON public.user_relations USING btree (deleted_at);


--
-- TOC entry 4231 (class 1259 OID 21623)
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- TOC entry 4236 (class 1259 OID 21624)
-- Name: idx_ware_houses_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ware_houses_deleted_at ON public.ware_houses USING btree (deleted_at);


--
-- TOC entry 4268 (class 2606 OID 21625)
-- Name: customers customer_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customer_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4269 (class 2606 OID 21630)
-- Name: customers customers_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4270 (class 2606 OID 21635)
-- Name: customers customers_group; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_group FOREIGN KEY (group_id) REFERENCES public.groups(id);


--
-- TOC entry 4239 (class 2606 OID 21640)
-- Name: account_settings fk_acc_setts_acc1; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.account_settings
    ADD CONSTRAINT fk_acc_setts_acc1 FOREIGN KEY (bank_account) REFERENCES public.ledgers(id);


--
-- TOC entry 4240 (class 2606 OID 21645)
-- Name: account_settings fk_acc_setts_acc2; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.account_settings
    ADD CONSTRAINT fk_acc_setts_acc2 FOREIGN KEY (cash_accunt) REFERENCES public.ledgers(id);


--
-- TOC entry 4241 (class 2606 OID 21650)
-- Name: account_settings fk_acc_setts_acc3; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.account_settings
    ADD CONSTRAINT fk_acc_setts_acc3 FOREIGN KEY (payable_account) REFERENCES public.ledgers(id);


--
-- TOC entry 4242 (class 2606 OID 21655)
-- Name: account_settings fk_acc_setts_acc4; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.account_settings
    ADD CONSTRAINT fk_acc_setts_acc4 FOREIGN KEY (cost_of_good_sold_account) REFERENCES public.ledgers(id);


--
-- TOC entry 4243 (class 2606 OID 21660)
-- Name: account_settings fk_acc_setts_acc5; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.account_settings
    ADD CONSTRAINT fk_acc_setts_acc5 FOREIGN KEY (receivable_account) REFERENCES public.ledgers(id);


--
-- TOC entry 4244 (class 2606 OID 21665)
-- Name: account_settings fk_acc_setts_acc6; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.account_settings
    ADD CONSTRAINT fk_acc_setts_acc6 FOREIGN KEY (income_account) REFERENCES public.ledgers(id);


--
-- TOC entry 4245 (class 2606 OID 21670)
-- Name: account_settings fk_acc_setts_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.account_settings
    ADD CONSTRAINT fk_acc_setts_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4246 (class 2606 OID 21675)
-- Name: actions fk_actions_entity; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.actions
    ADD CONSTRAINT fk_actions_entity FOREIGN KEY (entity_id) REFERENCES public.entities(id);


--
-- TOC entry 4247 (class 2606 OID 21680)
-- Name: activities fk_activities_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.activities
    ADD CONSTRAINT fk_activities_party FOREIGN KEY (party_id) REFERENCES public.parties(id);


--
-- TOC entry 4248 (class 2606 OID 21685)
-- Name: activities fk_activities_profile; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.activities
    ADD CONSTRAINT fk_activities_profile FOREIGN KEY (profile_id) REFERENCES public.profiles(id);


--
-- TOC entry 4249 (class 2606 OID 21690)
-- Name: addresses fk_address_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT fk_address_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4250 (class 2606 OID 21695)
-- Name: addresses fk_addresses_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT fk_addresses_company FOREIGN KEY (company_id) REFERENCES public.companies(id) NOT VALID;


--
-- TOC entry 4367 (class 2606 OID 21700)
-- Name: r_bookings fk_bookings_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.r_bookings
    ADD CONSTRAINT fk_bookings_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4368 (class 2606 OID 21705)
-- Name: r_bookings fk_bookings_court; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.r_bookings
    ADD CONSTRAINT fk_bookings_court FOREIGN KEY (court_id) REFERENCES public.r_courts(id);


--
-- TOC entry 4369 (class 2606 OID 21710)
-- Name: r_bookings fk_bookings_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.r_bookings
    ADD CONSTRAINT fk_bookings_party FOREIGN KEY (party) REFERENCES public.parties(id);


--
-- TOC entry 4370 (class 2606 OID 21715)
-- Name: r_bookings fk_bookings_party_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.r_bookings
    ADD CONSTRAINT fk_bookings_party_id FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4366 (class 2606 OID 21720)
-- Name: r_booking_prices fk_bookings_prices_booking; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.r_booking_prices
    ADD CONSTRAINT fk_bookings_prices_booking FOREIGN KEY (id) REFERENCES public.r_bookings(id);


--
-- TOC entry 4371 (class 2606 OID 21725)
-- Name: r_bookings fk_bookings_status; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.r_bookings
    ADD CONSTRAINT fk_bookings_status FOREIGN KEY (status) REFERENCES public.states(state);


--
-- TOC entry 4256 (class 2606 OID 21730)
-- Name: company_defaults fk_cd_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.company_defaults
    ADD CONSTRAINT fk_cd_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4257 (class 2606 OID 21735)
-- Name: company_defaults fk_cd_currency; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.company_defaults
    ADD CONSTRAINT fk_cd_currency FOREIGN KEY (currency) REFERENCES public.currencies("code ");


--
-- TOC entry 4258 (class 2606 OID 21740)
-- Name: company_entities fk_ce_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.company_entities
    ADD CONSTRAINT fk_ce_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4259 (class 2606 OID 21745)
-- Name: company_entities fk_ce_entity; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.company_entities
    ADD CONSTRAINT fk_ce_entity FOREIGN KEY (entity_id) REFERENCES public.entities(id);


--
-- TOC entry 4251 (class 2606 OID 21750)
-- Name: charges_template fk_charges_template_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.charges_template
    ADD CONSTRAINT fk_charges_template_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4252 (class 2606 OID 21755)
-- Name: charges_template fk_charges_template_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.charges_template
    ADD CONSTRAINT fk_charges_template_id FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4253 (class 2606 OID 21760)
-- Name: charges_template fk_charges_template_status; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.charges_template
    ADD CONSTRAINT fk_charges_template_status FOREIGN KEY (status) REFERENCES public.states(state);


--
-- TOC entry 4254 (class 2606 OID 21765)
-- Name: companies fk_companies_parent; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT fk_companies_parent FOREIGN KEY (parent_id) REFERENCES public.companies(id);


--
-- TOC entry 4255 (class 2606 OID 21770)
-- Name: companies fk_companies_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT fk_companies_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4260 (class 2606 OID 21775)
-- Name: contacts fk_contacts_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.contacts
    ADD CONSTRAINT fk_contacts_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4261 (class 2606 OID 21780)
-- Name: contacts fk_contacts_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.contacts
    ADD CONSTRAINT fk_contacts_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4262 (class 2606 OID 21785)
-- Name: cost_centers fk_cost_centers_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cost_centers
    ADD CONSTRAINT fk_cost_centers_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4263 (class 2606 OID 21790)
-- Name: cost_centers fk_cost_centers_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cost_centers
    ADD CONSTRAINT fk_cost_centers_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4264 (class 2606 OID 21795)
-- Name: cost_centers fk_cost_centers_status; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cost_centers
    ADD CONSTRAINT fk_cost_centers_status FOREIGN KEY (status) REFERENCES public.states(state);


--
-- TOC entry 4372 (class 2606 OID 21800)
-- Name: r_court_rates fk_court_rates_court; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.r_court_rates
    ADD CONSTRAINT fk_court_rates_court FOREIGN KEY (court_id) REFERENCES public.r_courts(id);


--
-- TOC entry 4373 (class 2606 OID 21805)
-- Name: r_court_rates fk_court_rates_currency; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.r_court_rates
    ADD CONSTRAINT fk_court_rates_currency FOREIGN KEY (currency) REFERENCES public.currencies("code ");


--
-- TOC entry 4374 (class 2606 OID 21810)
-- Name: r_courts fk_courts_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.r_courts
    ADD CONSTRAINT fk_courts_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4375 (class 2606 OID 21815)
-- Name: r_courts fk_courts_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.r_courts
    ADD CONSTRAINT fk_courts_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4376 (class 2606 OID 21820)
-- Name: r_courts fk_courts_status; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.r_courts
    ADD CONSTRAINT fk_courts_status FOREIGN KEY (status) REFERENCES public.states(state);


--
-- TOC entry 4265 (class 2606 OID 21825)
-- Name: currency_exchanges fk_currency_exchanges_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.currency_exchanges
    ADD CONSTRAINT fk_currency_exchanges_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4266 (class 2606 OID 21830)
-- Name: currency_exchanges fk_currency_exchanges_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.currency_exchanges
    ADD CONSTRAINT fk_currency_exchanges_id FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4267 (class 2606 OID 21835)
-- Name: currency_exchanges fk_currency_exchanges_status; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.currency_exchanges
    ADD CONSTRAINT fk_currency_exchanges_status FOREIGN KEY (status) REFERENCES public.states(state);


--
-- TOC entry 4271 (class 2606 OID 21840)
-- Name: delivery_line_items fk_dli_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.delivery_line_items
    ADD CONSTRAINT fk_dli_warehouse FOREIGN KEY (source_warehouse_id) REFERENCES public.ware_houses(id);


--
-- TOC entry 4377 (class 2606 OID 21845)
-- Name: r_events fk_events_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.r_events
    ADD CONSTRAINT fk_events_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4272 (class 2606 OID 21850)
-- Name: groups fk_groups_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.groups
    ADD CONSTRAINT fk_groups_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4273 (class 2606 OID 21855)
-- Name: groups fk_groups_parent; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.groups
    ADD CONSTRAINT fk_groups_parent FOREIGN KEY (parent_id) REFERENCES public.groups(id);


--
-- TOC entry 4274 (class 2606 OID 21860)
-- Name: groups fk_groups_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.groups
    ADD CONSTRAINT fk_groups_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4285 (class 2606 OID 21865)
-- Name: item_line_receipts fk_il_receipts_line; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_line_receipts
    ADD CONSTRAINT fk_il_receipts_line FOREIGN KEY (item_line) REFERENCES public.item_lines(id) ON DELETE CASCADE;


--
-- TOC entry 4288 (class 2606 OID 21870)
-- Name: item_line_stock_entries fk_ilse_item_line; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_line_stock_entries
    ADD CONSTRAINT fk_ilse_item_line FOREIGN KEY (item_line) REFERENCES public.item_lines(id);


--
-- TOC entry 4289 (class 2606 OID 21875)
-- Name: item_line_stock_entries fk_ilse_source_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_line_stock_entries
    ADD CONSTRAINT fk_ilse_source_warehouse FOREIGN KEY (source_warehouse_id) REFERENCES public.ware_houses(id);


--
-- TOC entry 4290 (class 2606 OID 21880)
-- Name: item_line_stock_entries fk_ilse_target_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_line_stock_entries
    ADD CONSTRAINT fk_ilse_target_warehouse FOREIGN KEY (target_warehouse_id) REFERENCES public.ware_houses(id);


--
-- TOC entry 4275 (class 2606 OID 21885)
-- Name: invoiced_item_lines fk_inv_invoice; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.invoiced_item_lines
    ADD CONSTRAINT fk_inv_invoice FOREIGN KEY (invoice_id) REFERENCES public.invoices(id);


--
-- TOC entry 4276 (class 2606 OID 21890)
-- Name: invoiced_item_lines fk_inv_item_line; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.invoiced_item_lines
    ADD CONSTRAINT fk_inv_item_line FOREIGN KEY (item_line) REFERENCES public.item_lines(id) ON DELETE CASCADE;


--
-- TOC entry 4277 (class 2606 OID 21895)
-- Name: invoices fk_invoices_currency; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT fk_invoices_currency FOREIGN KEY (currency) REFERENCES public.currencies("code ") ON UPDATE CASCADE;


--
-- TOC entry 4278 (class 2606 OID 21900)
-- Name: invoices fk_invoices_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT fk_invoices_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4279 (class 2606 OID 21905)
-- Name: invoices fk_invoices_state; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT fk_invoices_state FOREIGN KEY (status) REFERENCES public.states(state);


--
-- TOC entry 4281 (class 2606 OID 21910)
-- Name: item_attribute_values fk_item_attribute_values_item_attribute; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_attribute_values
    ADD CONSTRAINT fk_item_attribute_values_item_attribute FOREIGN KEY (item_attribute_id) REFERENCES public.item_attributes(id);


--
-- TOC entry 4282 (class 2606 OID 21915)
-- Name: item_attributes fk_item_attributes_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_attributes
    ADD CONSTRAINT fk_item_attributes_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4283 (class 2606 OID 21920)
-- Name: item_attributes fk_item_attributes_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_attributes
    ADD CONSTRAINT fk_item_attributes_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4284 (class 2606 OID 21925)
-- Name: item_inventory_settings fk_item_inventory_item; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_inventory_settings
    ADD CONSTRAINT fk_item_inventory_item FOREIGN KEY (item_id) REFERENCES public.items(id);


--
-- TOC entry 4286 (class 2606 OID 21930)
-- Name: item_line_receipts fk_item_line_warehouse_accepted; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_line_receipts
    ADD CONSTRAINT fk_item_line_warehouse_accepted FOREIGN KEY (accepted_warehouse) REFERENCES public.ware_houses(id);


--
-- TOC entry 4287 (class 2606 OID 21935)
-- Name: item_line_receipts fk_item_line_warehouse_rejected; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_line_receipts
    ADD CONSTRAINT fk_item_line_warehouse_rejected FOREIGN KEY (rejected_warehouse) REFERENCES public.ware_houses(id);


--
-- TOC entry 4291 (class 2606 OID 21940)
-- Name: item_lines fk_item_lines_item; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_lines
    ADD CONSTRAINT fk_item_lines_item FOREIGN KEY (item_id) REFERENCES public.items(id);


--
-- TOC entry 4292 (class 2606 OID 21945)
-- Name: item_lines fk_item_lines_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_lines
    ADD CONSTRAINT fk_item_lines_party FOREIGN KEY (party_id) REFERENCES public.parties(id);


--
-- TOC entry 4293 (class 2606 OID 21950)
-- Name: item_lines fk_item_lines_reference; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_lines
    ADD CONSTRAINT fk_item_lines_reference FOREIGN KEY (id) REFERENCES public.item_lines(id);


--
-- TOC entry 4294 (class 2606 OID 21955)
-- Name: item_lines fk_item_lines_type; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_lines
    ADD CONSTRAINT fk_item_lines_type FOREIGN KEY ("type ") REFERENCES public.item_line_types(type);


--
-- TOC entry 4295 (class 2606 OID 21960)
-- Name: price_lists fk_item_price_lists_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.price_lists
    ADD CONSTRAINT fk_item_price_lists_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4299 (class 2606 OID 21965)
-- Name: item_prices fk_item_prices_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_prices
    ADD CONSTRAINT fk_item_prices_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4300 (class 2606 OID 21970)
-- Name: item_prices fk_item_prices_item; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_prices
    ADD CONSTRAINT fk_item_prices_item FOREIGN KEY (item_id) REFERENCES public.items(id);


--
-- TOC entry 4301 (class 2606 OID 21975)
-- Name: item_prices fk_item_prices_item_price_list; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_prices
    ADD CONSTRAINT fk_item_prices_item_price_list FOREIGN KEY (price_list_id) REFERENCES public.price_lists(id);


--
-- TOC entry 4298 (class 2606 OID 21980)
-- Name: item_price_plugins fk_item_prices_item_price_plugin; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_price_plugins
    ADD CONSTRAINT fk_item_prices_item_price_plugin FOREIGN KEY (base_id) REFERENCES public.item_prices(id);


--
-- TOC entry 4302 (class 2606 OID 21985)
-- Name: item_prices fk_item_prices_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_prices
    ADD CONSTRAINT fk_item_prices_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4303 (class 2606 OID 21990)
-- Name: item_variants fk_item_variants_item; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_variants
    ADD CONSTRAINT fk_item_variants_item FOREIGN KEY (item_id) REFERENCES public.items(id);


--
-- TOC entry 4304 (class 2606 OID 21995)
-- Name: item_variants fk_item_variants_item_attribute_value; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_variants
    ADD CONSTRAINT fk_item_variants_item_attribute_value FOREIGN KEY (item_attribute_value_id) REFERENCES public.item_attribute_values(id);


--
-- TOC entry 4305 (class 2606 OID 22000)
-- Name: item_variants fk_item_variants_variant; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item_variants
    ADD CONSTRAINT fk_item_variants_variant FOREIGN KEY (variant_id) REFERENCES public.items(id);


--
-- TOC entry 4306 (class 2606 OID 22005)
-- Name: items fk_items_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT fk_items_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4307 (class 2606 OID 22010)
-- Name: items fk_items_group; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT fk_items_group FOREIGN KEY (group_id) REFERENCES public.groups(id);


--
-- TOC entry 4308 (class 2606 OID 22015)
-- Name: items fk_items_items_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT fk_items_items_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4309 (class 2606 OID 22020)
-- Name: items fk_items_parent; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT fk_items_parent FOREIGN KEY (parent_id) REFERENCES public.items(id);


--
-- TOC entry 4310 (class 2606 OID 22025)
-- Name: items fk_items_unit_of_measure; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT fk_items_unit_of_measure FOREIGN KEY (unit_of_measure_id) REFERENCES public.unit_of_measures(id);


--
-- TOC entry 4313 (class 2606 OID 22030)
-- Name: journal_entry_lines fk_jel_cost_center; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.journal_entry_lines
    ADD CONSTRAINT fk_jel_cost_center FOREIGN KEY (cost_center_id) REFERENCES public.cost_centers(id);


--
-- TOC entry 4314 (class 2606 OID 22035)
-- Name: journal_entry_lines fk_jel_currency; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.journal_entry_lines
    ADD CONSTRAINT fk_jel_currency FOREIGN KEY (currency) REFERENCES public.currencies("code ");


--
-- TOC entry 4315 (class 2606 OID 22040)
-- Name: journal_entry_lines fk_jel_journal_entry; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.journal_entry_lines
    ADD CONSTRAINT fk_jel_journal_entry FOREIGN KEY (journal_entry_id) REFERENCES public.journal_entries(id);


--
-- TOC entry 4316 (class 2606 OID 22045)
-- Name: journal_entry_lines fk_jel_ledger; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.journal_entry_lines
    ADD CONSTRAINT fk_jel_ledger FOREIGN KEY (ledger_id) REFERENCES public.ledgers(id);


--
-- TOC entry 4317 (class 2606 OID 22050)
-- Name: journal_entry_lines fk_jel_project; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.journal_entry_lines
    ADD CONSTRAINT fk_jel_project FOREIGN KEY (project_id) REFERENCES public.projects(id);


--
-- TOC entry 4311 (class 2606 OID 22055)
-- Name: journal_entries fk_journal_entries_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.journal_entries
    ADD CONSTRAINT fk_journal_entries_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4312 (class 2606 OID 22060)
-- Name: journal_entries fk_journal_entries_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.journal_entries
    ADD CONSTRAINT fk_journal_entries_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4319 (class 2606 OID 22065)
-- Name: ledger_accounts fk_l_acc_currency; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ledger_accounts
    ADD CONSTRAINT fk_l_acc_currency FOREIGN KEY (currency) REFERENCES public.currencies("code ");


--
-- TOC entry 4321 (class 2606 OID 22070)
-- Name: ledgers fk_ledgers_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ledgers
    ADD CONSTRAINT fk_ledgers_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4322 (class 2606 OID 22075)
-- Name: ledgers fk_ledgers_parent; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ledgers
    ADD CONSTRAINT fk_ledgers_parent FOREIGN KEY (ledger_parent) REFERENCES public.ledgers(id);


--
-- TOC entry 4323 (class 2606 OID 22080)
-- Name: ledgers fk_ledgers_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ledgers
    ADD CONSTRAINT fk_ledgers_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4324 (class 2606 OID 22085)
-- Name: modules fk_modules_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.modules
    ADD CONSTRAINT fk_modules_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4325 (class 2606 OID 22090)
-- Name: modules fk_modules_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.modules
    ADD CONSTRAINT fk_modules_id FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4326 (class 2606 OID 22095)
-- Name: orders fk_orders_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_orders_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4327 (class 2606 OID 22100)
-- Name: orders fk_orders_currency; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_orders_currency FOREIGN KEY (currency) REFERENCES public.currencies("code ");


--
-- TOC entry 4328 (class 2606 OID 22105)
-- Name: orders fk_orders_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_orders_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4329 (class 2606 OID 22110)
-- Name: orders fk_orders_party2; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_orders_party2 FOREIGN KEY (party_id) REFERENCES public.parties(id);


--
-- TOC entry 4330 (class 2606 OID 22115)
-- Name: orders fk_orders_state; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_orders_state FOREIGN KEY (status) REFERENCES public.states(state) ON UPDATE CASCADE;


--
-- TOC entry 4331 (class 2606 OID 22120)
-- Name: parties fk_parties_party_type; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.parties
    ADD CONSTRAINT fk_parties_party_type FOREIGN KEY (party_type_code) REFERENCES public.party_types(code) ON UPDATE CASCADE;


--
-- TOC entry 4332 (class 2606 OID 22125)
-- Name: party_addresses fk_party_addresses_address; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.party_addresses
    ADD CONSTRAINT fk_party_addresses_address FOREIGN KEY (address_id) REFERENCES public.parties(id);


--
-- TOC entry 4333 (class 2606 OID 22130)
-- Name: party_addresses fk_party_addresses_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.party_addresses
    ADD CONSTRAINT fk_party_addresses_party FOREIGN KEY (party_id) REFERENCES public.parties(id);


--
-- TOC entry 4280 (class 2606 OID 22135)
-- Name: invoices fk_party_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT fk_party_party FOREIGN KEY (party_id) REFERENCES public.parties(id);


--
-- TOC entry 4341 (class 2606 OID 22140)
-- Name: payment_terms fk_payment_tems_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_terms
    ADD CONSTRAINT fk_payment_tems_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4342 (class 2606 OID 22145)
-- Name: payment_terms fk_payment_tems_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_terms
    ADD CONSTRAINT fk_payment_tems_id FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4345 (class 2606 OID 22150)
-- Name: payment_terms_templates fk_payment_tems_templates_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_terms_templates
    ADD CONSTRAINT fk_payment_tems_templates_id FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4346 (class 2606 OID 22155)
-- Name: payment_terms_templates fk_payment_terms_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_terms_templates
    ADD CONSTRAINT fk_payment_terms_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4351 (class 2606 OID 22160)
-- Name: piano_form fk_piano_form_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.piano_form
    ADD CONSTRAINT fk_piano_form_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4352 (class 2606 OID 22165)
-- Name: piano_form fk_piano_form_party_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.piano_form
    ADD CONSTRAINT fk_piano_form_party_id FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4338 (class 2606 OID 22170)
-- Name: payment_references fk_pr_currency; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_references
    ADD CONSTRAINT fk_pr_currency FOREIGN KEY (currency) REFERENCES public.currencies("code ");


--
-- TOC entry 4339 (class 2606 OID 22175)
-- Name: payment_references fk_pr_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_references
    ADD CONSTRAINT fk_pr_party FOREIGN KEY (party_id) REFERENCES public.parties(id);


--
-- TOC entry 4340 (class 2606 OID 22180)
-- Name: payment_references fk_pr_payment; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_references
    ADD CONSTRAINT fk_pr_payment FOREIGN KEY (payment_id) REFERENCES public.payments(id);


--
-- TOC entry 4296 (class 2606 OID 22185)
-- Name: price_lists fk_price_lists_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.price_lists
    ADD CONSTRAINT fk_price_lists_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4297 (class 2606 OID 22190)
-- Name: price_lists fk_price_lists_status; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.price_lists
    ADD CONSTRAINT fk_price_lists_status FOREIGN KEY (status) REFERENCES public.states(state);


--
-- TOC entry 4353 (class 2606 OID 22195)
-- Name: profiles fk_profile_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.profiles
    ADD CONSTRAINT fk_profile_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4429 (class 2606 OID 22200)
-- Name: user_relations fk_profiles_user_relation; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_relations
    ADD CONSTRAINT fk_profiles_user_relation FOREIGN KEY (profile_id) REFERENCES public.profiles(id);


--
-- TOC entry 4354 (class 2606 OID 22205)
-- Name: progress_invoices fk_progress_invoices_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.progress_invoices
    ADD CONSTRAINT fk_progress_invoices_id FOREIGN KEY (invoice_id) REFERENCES public.invoices(id);


--
-- TOC entry 4355 (class 2606 OID 22210)
-- Name: projects fk_projects_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.projects
    ADD CONSTRAINT fk_projects_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4356 (class 2606 OID 22215)
-- Name: projects fk_projects_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.projects
    ADD CONSTRAINT fk_projects_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4357 (class 2606 OID 22220)
-- Name: projects fk_projects_status; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.projects
    ADD CONSTRAINT fk_projects_status FOREIGN KEY (status) REFERENCES public.states(state);


--
-- TOC entry 4343 (class 2606 OID 22225)
-- Name: payment_terms_lines fk_ptl_doc_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_terms_lines
    ADD CONSTRAINT fk_ptl_doc_party FOREIGN KEY (document_id) REFERENCES public.parties(id);


--
-- TOC entry 4344 (class 2606 OID 22230)
-- Name: payment_terms_lines fk_ptl_payment_terms; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_terms_lines
    ADD CONSTRAINT fk_ptl_payment_terms FOREIGN KEY (payment_terms_id) REFERENCES public.payment_terms(id);


--
-- TOC entry 4358 (class 2606 OID 22235)
-- Name: purchase_records fk_purchase_records_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.purchase_records
    ADD CONSTRAINT fk_purchase_records_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4359 (class 2606 OID 22240)
-- Name: purchase_records fk_purchase_records_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.purchase_records
    ADD CONSTRAINT fk_purchase_records_id FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4360 (class 2606 OID 22245)
-- Name: purchase_records fk_purchase_records_supplier; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.purchase_records
    ADD CONSTRAINT fk_purchase_records_supplier FOREIGN KEY (supplier_id) REFERENCES public.suppliers(id);


--
-- TOC entry 4362 (class 2606 OID 22250)
-- Name: quotations fk_quotations_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.quotations
    ADD CONSTRAINT fk_quotations_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4363 (class 2606 OID 22255)
-- Name: quotations fk_quotations_currency; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.quotations
    ADD CONSTRAINT fk_quotations_currency FOREIGN KEY (currency) REFERENCES public.currencies("code ");


--
-- TOC entry 4364 (class 2606 OID 22260)
-- Name: quotations fk_quotations_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.quotations
    ADD CONSTRAINT fk_quotations_id FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4365 (class 2606 OID 22265)
-- Name: quotations fk_quotations_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.quotations
    ADD CONSTRAINT fk_quotations_party FOREIGN KEY (party_id) REFERENCES public.parties(id);


--
-- TOC entry 4378 (class 2606 OID 22270)
-- Name: r_events fk_r_events_party_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.r_events
    ADD CONSTRAINT fk_r_events_party_id FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4379 (class 2606 OID 22275)
-- Name: receipts fk_receipts_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT fk_receipts_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4380 (class 2606 OID 22280)
-- Name: receipts fk_receipts_currency; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT fk_receipts_currency FOREIGN KEY (currency) REFERENCES public.currencies("code ");


--
-- TOC entry 4381 (class 2606 OID 22285)
-- Name: receipts fk_receipts_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT fk_receipts_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4382 (class 2606 OID 22290)
-- Name: receipts fk_receipts_party2; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT fk_receipts_party2 FOREIGN KEY (party_id) REFERENCES public.parties(id);


--
-- TOC entry 4383 (class 2606 OID 22295)
-- Name: receipts fk_receipts_status; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT fk_receipts_status FOREIGN KEY (status) REFERENCES public.states(state);


--
-- TOC entry 4336 (class 2606 OID 22300)
-- Name: party_references fk_references_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.party_references
    ADD CONSTRAINT fk_references_party FOREIGN KEY (party_id) REFERENCES public.parties(id);


--
-- TOC entry 4337 (class 2606 OID 22305)
-- Name: party_references fk_references_reference; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.party_references
    ADD CONSTRAINT fk_references_reference FOREIGN KEY (reference_id) REFERENCES public.parties(id);


--
-- TOC entry 4384 (class 2606 OID 22310)
-- Name: role_actions fk_role_actions_action; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_actions
    ADD CONSTRAINT fk_role_actions_action FOREIGN KEY (action_id) REFERENCES public.actions(id);


--
-- TOC entry 4386 (class 2606 OID 22315)
-- Name: role_templates fk_role_templates_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_templates
    ADD CONSTRAINT fk_role_templates_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4387 (class 2606 OID 22320)
-- Name: roles fk_roles_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT fk_roles_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4388 (class 2606 OID 22325)
-- Name: roles fk_roles_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT fk_roles_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4385 (class 2606 OID 22330)
-- Name: role_actions fk_roles_role_actions; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_actions
    ADD CONSTRAINT fk_roles_role_actions FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- TOC entry 4389 (class 2606 OID 22335)
-- Name: sales_records fk_sales_records_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales_records
    ADD CONSTRAINT fk_sales_records_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4390 (class 2606 OID 22340)
-- Name: sales_records fk_sales_records_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales_records
    ADD CONSTRAINT fk_sales_records_id FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4391 (class 2606 OID 22345)
-- Name: sales_records fk_sales_records_invoice; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales_records
    ADD CONSTRAINT fk_sales_records_invoice FOREIGN KEY (invoice_id) REFERENCES public.invoices(id);


--
-- TOC entry 4393 (class 2606 OID 22350)
-- Name: serial_nos fk_sn_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.serial_nos
    ADD CONSTRAINT fk_sn_id FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4394 (class 2606 OID 22355)
-- Name: stock_defaults fk_stock_defaults_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_defaults
    ADD CONSTRAINT fk_stock_defaults_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4395 (class 2606 OID 22360)
-- Name: stock_defaults fk_stock_defaults_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_defaults
    ADD CONSTRAINT fk_stock_defaults_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4396 (class 2606 OID 22365)
-- Name: stock_entries fk_stock_entries_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_entries
    ADD CONSTRAINT fk_stock_entries_company FOREIGN KEY (currency) REFERENCES public.currencies("code ");


--
-- TOC entry 4397 (class 2606 OID 22370)
-- Name: stock_entries fk_stock_entries_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_entries
    ADD CONSTRAINT fk_stock_entries_id FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4398 (class 2606 OID 22375)
-- Name: stock_entries fk_stock_entries_project; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_entries
    ADD CONSTRAINT fk_stock_entries_project FOREIGN KEY (project_id) REFERENCES public.projects(id);


--
-- TOC entry 4399 (class 2606 OID 22380)
-- Name: stock_levels fk_stock_levels_item; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_levels
    ADD CONSTRAINT fk_stock_levels_item FOREIGN KEY (item_id) REFERENCES public.items(id);


--
-- TOC entry 4400 (class 2606 OID 22385)
-- Name: stock_levels fk_stock_levels_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_levels
    ADD CONSTRAINT fk_stock_levels_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4401 (class 2606 OID 22390)
-- Name: stock_levels fk_stock_levels_ware_house; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_levels
    ADD CONSTRAINT fk_stock_levels_ware_house FOREIGN KEY (ware_house_id) REFERENCES public.ware_houses(id);


--
-- TOC entry 4402 (class 2606 OID 22395)
-- Name: stock_movements fk_stock_movements_item; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_movements
    ADD CONSTRAINT fk_stock_movements_item FOREIGN KEY (item_id) REFERENCES public.items(id);


--
-- TOC entry 4403 (class 2606 OID 22400)
-- Name: stock_movements fk_stock_movements_ware_house; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_movements
    ADD CONSTRAINT fk_stock_movements_ware_house FOREIGN KEY (ware_house_id) REFERENCES public.ware_houses(id);


--
-- TOC entry 4404 (class 2606 OID 22405)
-- Name: stock_settings fk_stock_settings_acc1; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_settings
    ADD CONSTRAINT fk_stock_settings_acc1 FOREIGN KEY (inventory_account) REFERENCES public.ledgers(id);


--
-- TOC entry 4405 (class 2606 OID 22410)
-- Name: stock_settings fk_stock_settings_acc2; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_settings
    ADD CONSTRAINT fk_stock_settings_acc2 FOREIGN KEY (stock_received_but_not_billed) REFERENCES public.ledgers(id);


--
-- TOC entry 4406 (class 2606 OID 22415)
-- Name: stock_settings fk_stock_settings_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_settings
    ADD CONSTRAINT fk_stock_settings_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4407 (class 2606 OID 22420)
-- Name: stock_transactions fk_stock_tx_currency; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_transactions
    ADD CONSTRAINT fk_stock_tx_currency FOREIGN KEY (currency) REFERENCES public.currencies("code ");


--
-- TOC entry 4408 (class 2606 OID 22425)
-- Name: stock_transactions fk_stock_tx_item; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_transactions
    ADD CONSTRAINT fk_stock_tx_item FOREIGN KEY (item_id) REFERENCES public.items(id);


--
-- TOC entry 4409 (class 2606 OID 22430)
-- Name: stock_transactions fk_stock_tx_warehouse; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_transactions
    ADD CONSTRAINT fk_stock_tx_warehouse FOREIGN KEY (warehouse_id) REFERENCES public.ware_houses(id);


--
-- TOC entry 4410 (class 2606 OID 22435)
-- Name: supplier_orders fk_supplier_orders_supplier; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.supplier_orders
    ADD CONSTRAINT fk_supplier_orders_supplier FOREIGN KEY (supplier_id) REFERENCES public.suppliers(id);


--
-- TOC entry 4412 (class 2606 OID 22440)
-- Name: suppliers fk_suppliers_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT fk_suppliers_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4413 (class 2606 OID 22445)
-- Name: suppliers fk_suppliers_group; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT fk_suppliers_group FOREIGN KEY (group_id) REFERENCES public.groups(id);


--
-- TOC entry 4411 (class 2606 OID 22450)
-- Name: supplier_orders fk_suppliers_orders_order; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.supplier_orders
    ADD CONSTRAINT fk_suppliers_orders_order FOREIGN KEY (order_id) REFERENCES public.orders(id);


--
-- TOC entry 4414 (class 2606 OID 22455)
-- Name: suppliers fk_suppliers_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT fk_suppliers_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4415 (class 2606 OID 22460)
-- Name: tax_and_charge_lines fk_tacl_doc_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tax_and_charge_lines
    ADD CONSTRAINT fk_tacl_doc_party FOREIGN KEY (doc_party_id) REFERENCES public.parties(id);


--
-- TOC entry 4416 (class 2606 OID 22465)
-- Name: tax_and_charge_lines fk_tacl_ledger; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tax_and_charge_lines
    ADD CONSTRAINT fk_tacl_ledger FOREIGN KEY (account_head) REFERENCES public.ledgers(id);


--
-- TOC entry 4417 (class 2606 OID 22470)
-- Name: taxes fk_taxes_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.taxes
    ADD CONSTRAINT fk_taxes_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4418 (class 2606 OID 22475)
-- Name: taxes fk_taxes_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.taxes
    ADD CONSTRAINT fk_taxes_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4419 (class 2606 OID 22480)
-- Name: terms_and_conditions fk_terms_and_conditions_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.terms_and_conditions
    ADD CONSTRAINT fk_terms_and_conditions_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4420 (class 2606 OID 22485)
-- Name: terms_and_conditions fk_terms_and_conditions_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.terms_and_conditions
    ADD CONSTRAINT fk_terms_and_conditions_id FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4421 (class 2606 OID 22490)
-- Name: transaction_ledgers fk_tx_ledger_currency; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction_ledgers
    ADD CONSTRAINT fk_tx_ledger_currency FOREIGN KEY (currency) REFERENCES public.currencies("code ");


--
-- TOC entry 4422 (class 2606 OID 22495)
-- Name: transaction_ledgers fk_tx_ledger_ledger; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction_ledgers
    ADD CONSTRAINT fk_tx_ledger_ledger FOREIGN KEY (ledger) REFERENCES public.ledgers(id);


--
-- TOC entry 4423 (class 2606 OID 22500)
-- Name: transaction_ledgers fk_tx_ledger_ledger_agst; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction_ledgers
    ADD CONSTRAINT fk_tx_ledger_ledger_agst FOREIGN KEY (ledger_against) REFERENCES public.ledgers(id);


--
-- TOC entry 4424 (class 2606 OID 22505)
-- Name: transaction_ledgers fk_tx_ledger_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction_ledgers
    ADD CONSTRAINT fk_tx_ledger_party FOREIGN KEY (party_id) REFERENCES public.parties(id);


--
-- TOC entry 4425 (class 2606 OID 22510)
-- Name: transaction_ledgers fk_tx_ledger_project; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction_ledgers
    ADD CONSTRAINT fk_tx_ledger_project FOREIGN KEY (project_id) REFERENCES public.projects(id);


--
-- TOC entry 4426 (class 2606 OID 22515)
-- Name: unit_of_measure_translations fk_unit_of_measure_translations_unit_of_measure; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.unit_of_measure_translations
    ADD CONSTRAINT fk_unit_of_measure_translations_unit_of_measure FOREIGN KEY (base_id) REFERENCES public.unit_of_measures(id);


--
-- TOC entry 4428 (class 2606 OID 22520)
-- Name: unit_of_measures fk_unit_of_measures_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.unit_of_measures
    ADD CONSTRAINT fk_unit_of_measures_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4427 (class 2606 OID 22525)
-- Name: unit_of_measure_translations fk_unit_of_measures_unit_of_measure_translation; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.unit_of_measure_translations
    ADD CONSTRAINT fk_unit_of_measures_unit_of_measure_translation FOREIGN KEY (base_id) REFERENCES public.unit_of_measures(id);


--
-- TOC entry 4430 (class 2606 OID 22530)
-- Name: user_relations fk_user_ralations_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_relations
    ADD CONSTRAINT fk_user_ralations_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4431 (class 2606 OID 22535)
-- Name: user_relations fk_user_ralations_profile; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_relations
    ADD CONSTRAINT fk_user_ralations_profile FOREIGN KEY (profile_id) REFERENCES public.profiles(id);


--
-- TOC entry 4432 (class 2606 OID 22540)
-- Name: user_relations fk_user_ralations_role; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_relations
    ADD CONSTRAINT fk_user_ralations_role FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- TOC entry 4433 (class 2606 OID 22545)
-- Name: user_relations fk_user_relations_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_relations
    ADD CONSTRAINT fk_user_relations_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4434 (class 2606 OID 22550)
-- Name: user_relations fk_user_relations_profile; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_relations
    ADD CONSTRAINT fk_user_relations_profile FOREIGN KEY (profile_id) REFERENCES public.profiles(id);


--
-- TOC entry 4435 (class 2606 OID 22555)
-- Name: user_relations fk_user_relations_role; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_relations
    ADD CONSTRAINT fk_user_relations_role FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- TOC entry 4438 (class 2606 OID 22560)
-- Name: users fk_users_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT fk_users_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4436 (class 2606 OID 22565)
-- Name: user_relations fk_users_user_ralation; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_relations
    ADD CONSTRAINT fk_users_user_ralation FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 4437 (class 2606 OID 22570)
-- Name: user_relations fk_users_user_relation; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_relations
    ADD CONSTRAINT fk_users_user_relation FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 4439 (class 2606 OID 22575)
-- Name: ware_houses fk_ware_houses_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ware_houses
    ADD CONSTRAINT fk_ware_houses_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4440 (class 2606 OID 22580)
-- Name: ware_houses fk_ware_houses_parent; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ware_houses
    ADD CONSTRAINT fk_ware_houses_parent FOREIGN KEY (parent_id) REFERENCES public.ware_houses(id);


--
-- TOC entry 4441 (class 2606 OID 22585)
-- Name: ware_houses fk_ware_houses_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ware_houses
    ADD CONSTRAINT fk_ware_houses_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4318 (class 2606 OID 22590)
-- Name: key_values key_values_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.key_values
    ADD CONSTRAINT key_values_party FOREIGN KEY (party_id) REFERENCES public.parties(id);


--
-- TOC entry 4320 (class 2606 OID 22595)
-- Name: ledger_accounts ledger_accounts_ledger; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ledger_accounts
    ADD CONSTRAINT ledger_accounts_ledger FOREIGN KEY (ledger_id) REFERENCES public.ledgers(id);


--
-- TOC entry 4334 (class 2606 OID 22600)
-- Name: party_payments party_payments_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.party_payments
    ADD CONSTRAINT party_payments_party FOREIGN KEY (party_id) REFERENCES public.parties(id);


--
-- TOC entry 4335 (class 2606 OID 22605)
-- Name: party_payments party_payments_payment; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.party_payments
    ADD CONSTRAINT party_payments_payment FOREIGN KEY (payment_id) REFERENCES public.parties(id);


--
-- TOC entry 4347 (class 2606 OID 22610)
-- Name: payments payments_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- TOC entry 4348 (class 2606 OID 22615)
-- Name: payments payments_party; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_party FOREIGN KEY (id) REFERENCES public.parties(id);


--
-- TOC entry 4349 (class 2606 OID 22620)
-- Name: payments payments_party2; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_party2 FOREIGN KEY (party_id) REFERENCES public.parties(id);


--
-- TOC entry 4350 (class 2606 OID 22625)
-- Name: payments payments_status; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_status FOREIGN KEY (status) REFERENCES public.states(state);


--
-- TOC entry 4392 (class 2606 OID 22630)
-- Name: sales_records pg_sales_records_customer; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales_records
    ADD CONSTRAINT pg_sales_records_customer FOREIGN KEY (customer_id) REFERENCES public.customers(id);


--
-- TOC entry 4361 (class 2606 OID 22635)
-- Name: purchase_records purchase_records_invoice; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.purchase_records
    ADD CONSTRAINT purchase_records_invoice FOREIGN KEY (invoice_id) REFERENCES public.invoices(id);


--
-- TOC entry 4593 (class 0 OID 0)
-- Dependencies: 7
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE USAGE ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO PUBLIC;


-- Completed on 2025-01-21 12:28:40

--
-- PostgreSQL database dump complete
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 17.2 (Debian 17.2-1.pgdg120+1)
-- Dumped by pg_dump version 17.2

-- Started on 2025-01-21 12:28:36

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 3769 (class 0 OID 20697)
-- Dependencies: 244
-- Data for Name: entities; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.entities (id, name, href) FROM stdin;
5	Item-Stock	
6	Item-Attributes	
8	Tax	
26	Piano-Forms	
36	Financial-Statements	
37	Stock-Setting	
12	Proveedor	buying/supplier
14	Cliente	selling/customer
15	Dirección	address
16	Contacto	contact
57	Pricing	pricing
21	Cancha	court
22	Reserva	booking
23	Evento	event
9	Lista de Precio	stock/priceList
17	Factura de Compra	invoice/purchaseInvoice
18	Pago	payment
19	Plan de Cuentas	accounting/account
20	Recibo de Compra	receipt/purchaseReceipt
1	Empresa	companies
2	Articulo	stock/item
3	Precio de Articulo	stock/itemPrice
4	Grupo de Articulo	group/itemGroup
7	Almacen	stock/warehouse
25	Factura de Venta	invoice/saleInvoice
28	Nota de entrega	receipt/deliveryNote
29	Asiento Contable	accounting/journalEntry
30	Centro de Costo	accounting/costCenter
31	Proyecto	project
32	Ingreso de Stock	stock/stockEntry
33	Libro Mayor	accounting/generalLedger
34	Cuentas por Cobrar	accounting/accountReceivable
35	Cuentas por Pagar	accounting/accountPayable
38	Número de Serie	stock/serialNo
39	Lote de Paquete	stock/batchBundle
40	Cotización del Proveedor	quotation/supplierQuotation
41	Cotización	quotation/salesQuotation
42	Plantilla de Cargos	accounting/chargesTemplate
43	Cambio de Divisa	currencyExchange
44	Registro de Compras	invoicing/purchaseRecord
45	Registro de Ventas	invoicing/salesRecord
46	Libro de Inventario	stock/stockLedger
48	Saldo de Inventario	stock/stockBalance
49	Resumen de Número de Serie	stock/serialNoResume
50	Estado de Resultados	accounting/profitAndLoss
51	Flujo de Efectivo	accounting/cashFlow
52	Balance General	accounting/balanceSheet
53	Resumen de Cuentas por Cobrar	accounting/accountReceivableSumary
54	Resumen de Cuentas por Pagar	accounting/accountReceivableSumary
55	Grupo de Proveedores	group/supplierGroup
56	Grupo de Clientes	group/customerGroup
13	Orden de Compra	order/purchaseOrder
24	Orden de Venta	order/saleOrder
10	Rol	manage/roles
11	Usuario	manage/users
27	Panel Reservas	bookingDashboard
47	Modulo	module
58	Terminos y Condiciones	terms-and-conditions
59	Condiciones de Pago	payment-terms
60	Plantilla de Condiciones de Pago	payment-terms-template
\.


--
-- TOC entry 3766 (class 0 OID 20599)
-- Dependencies: 225
-- Data for Name: actions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.actions (id, name, entity_id) FROM stdin;
1	create	1
2	edit	1
3	view	1
4	view	2
5	view	10
6	edit	10
7	create	2
8	create	10
9	view	11
10	create	11
11	edit	11
12	delete	11
13	view	3
14	create	3
15	edit	3
16	delete	3
17	view	4
18	create	4
19	edit	4
20	delete	4
21	view	5
22	create	5
23	edit	5
24	delete	5
25	view	6
26	create	6
27	edit	6
28	delete	6
29	view	7
30	create	7
31	edit	7
32	delete	7
33	view	8
34	create	8
35	edit	8
36	delete	8
37	view	9
38	create	9
39	edit	9
40	delete	9
41	edit	2
42	delete	2
43	view	12
44	create	12
45	edit	12
46	delete	12
47	view	14
48	create	14
49	edit	14
50	delete	14
51	view	15
52	create	15
53	edit	15
54	delete	15
55	view	16
56	create	16
57	edit	16
58	delete	16
59	view	17
60	create	17
61	edit	17
62	delete	17
63	view	18
64	create	18
65	edit	18
66	delete	18
67	view	19
68	create	19
69	edit	19
70	delete	19
71	view	20
72	create	20
73	edit	20
74	delete	20
75	view	21
76	create	21
77	edit	21
78	delete	21
79	view	22
80	create	22
81	edit	22
82	delete	22
83	view	23
84	create	23
85	edit	23
86	delete	23
87	view	24
88	create	24
89	edit	24
90	delete	24
91	view	25
92	create	25
93	edit	25
94	delete	25
95	view	26
96	create	26
97	edit	26
98	delete	26
99	view	27
100	create	27
101	edit	27
102	delete	27
103	view	28
105	delete	1
106	view	13
107	create	13
108	edit	13
109	delete	13
110	create	28
111	edit	28
112	delete	28
113	view	29
114	edit	29
115	create	29
116	delete	29
117	view	30
118	create	30
119	edit	30
120	delete	30
121	view	31
122	create	31
123	edit	31
124	delete	31
125	view	32
126	create	32
127	edit	32
128	delete	32
129	view	34
133	view	35
137	view	33
141	view	36
145	view	38
146	view	39
147	view	40
148	create	40
149	edit	40
150	view	41
151	create	41
152	edit	41
153	view	42
154	edit	42
155	create	42
156	view	43
157	create	43
158	edit	43
159	view	44
160	edit	44
161	create	44
162	delete	44
163	view	45
164	create	45
165	edit	45
166	delete	45
167	view	46
170	view	57
171	create	57
172	edit	57
173	view	56
174	view	55
175	view	54
176	view	53
177	view	52
178	view	51
179	view	50
180	view	49
181	view	48
182	view	47
183	create	47
184	edit	47
185	delete	47
186	view	58
187	create	58
188	edit	58
189	delete	58
190	view	59
191	create	59
192	edit	59
193	delete	59
194	view	60
195	create	60
196	edit	60
197	delete	60
\.


--
-- TOC entry 3768 (class 0 OID 20672)
-- Dependencies: 240
-- Data for Name: currencies; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.currencies ("code ", "name ") FROM stdin;
USD	USD
BOB	BOB
\.


--
-- TOC entry 3771 (class 0 OID 20919)
-- Dependencies: 288
-- Data for Name: party_types; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.party_types (code, name, created_at, entity_id) FROM stdin;
admin	Admin	2024-12-27 12:41:29.877172	\N
address	Address	2024-12-27 12:41:29.877172	\N
booking	Booking	2024-12-27 12:41:29.877172	\N
client	Client	2024-12-27 12:41:29.877172	\N
company	Company	2024-12-27 12:41:29.877172	\N
contact	Contact	2024-12-27 12:41:29.877172	\N
court	Court	2024-12-27 12:41:29.877172	\N
customer	Customer	2024-12-27 12:41:29.877172	\N
customerGroup	Customer Group	2024-12-27 12:41:29.877172	\N
employee	Employee	2024-12-27 12:41:29.877172	\N
eventBooking	Event Booking	2024-12-27 12:41:29.877172	\N
item	Item	2024-12-27 12:41:29.877172	\N
itemAttribute	Item Attribute	2024-12-27 12:41:29.877172	\N
itemGroup	Item Group	2024-12-27 12:41:29.877172	\N
itemPrice	Item Price	2024-12-27 12:41:29.877172	\N
ledger	Ledger	2024-12-27 12:41:29.877172	\N
pianoForms	Piano Forms	2024-12-27 12:41:29.877172	\N
priceList	Price List	2024-12-27 12:41:29.877172	\N
regateChart	regateChart	2024-12-27 12:41:29.877172	\N
stockLevel	Stock Level	2024-12-27 12:41:29.877172	\N
supplierGroup	Supplier Group	2024-12-27 12:41:29.877172	\N
tax	Tax	2024-12-27 12:41:29.877172	\N
warehouse	Warehouse	2024-12-27 12:41:29.877172	\N
roleTemplate	Role Template	2024-12-27 12:41:29.877172	\N
role	Role	2024-12-27 12:41:29.877172	\N
journalEntry	Journal Entry	2024-12-27 12:41:29.877172	\N
costCenter	Cost Center	2024-12-27 12:41:29.877172	\N
project	Project	2024-12-27 12:41:29.877172	\N
stockEntry	Stock Entry	2024-12-27 12:41:29.877172	\N
generalLedger	General Ledger	2024-12-27 12:41:29.877172	\N
accountReceivable	Account Receivable	2024-12-27 12:41:29.877172	\N
accountPayable	Account Payable	2024-12-27 12:41:29.877172	\N
financialStatements	Financial Statements	2024-12-27 12:41:29.877172	\N
setting	Setting	2024-12-27 12:41:29.877172	\N
serialNo	serialNo	2024-12-27 12:41:29.877172	\N
batchBundle	Batch Bundle	2024-12-27 12:41:29.877172	\N
supplierQuotation	Supplier Quotation	2024-12-27 12:41:29.877172	\N
chargesTemplate	Charges Template	2024-12-27 12:41:29.877172	\N
currencyExchange	Currency Exchange	2024-12-27 12:41:29.877172	\N
module	Module	2024-12-27 12:41:29.877172	\N
supplier	Supplier	2024-12-27 12:41:29.877172	12
pricing	Pricing	2024-12-27 12:41:29.877172	57
purchaseOrder	Purchase Order	2024-12-27 12:41:29.877172	13
salesQuotation	Sales Quotation	2024-12-27 12:41:29.877172	41
payment	Payment	2024-12-27 12:41:29.877172	18
purchaseInvoice	Purchase Invoice	2024-12-27 12:41:29.877172	17
purchaseReceipt	Receipt	2024-12-27 12:41:29.877172	20
purchaseRecord	Purchase Record	2024-12-27 12:41:29.877172	44
saleInvoice	Sale Invoice	2024-12-27 12:41:29.877172	25
saleOrder	Sale Order	2024-12-27 12:41:29.877172	24
salesRecord	Sales Record	2024-12-27 12:41:29.877172	45
deliveryNote	Delivery Note	2024-12-27 12:41:29.877172	28
termsAndConditions	Termisnos y Condiciones	2025-01-14 21:55:26.505318	\N
paymentTerms	Condiciones de Pago	2025-01-14 21:55:26.505318	\N
paymentTermsTemplate	Plantilla de condiciones de pago	2025-01-14 21:55:26.505318	\N
\.


--
-- TOC entry 3772 (class 0 OID 21125)
-- Dependencies: 329
-- Data for Name: states; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.states (state) FROM stdin;
APPROVED
BILLED_AWAITING_PAYMENT
CANCELLED
CLOSED
COMPLETED
DRAFT
ON_HOLD
OVERDUE
PAID
PARTIALLY_PAID
PENDING
RECEIVED
RETURNED
SHIPPED
SUBMITTED
TO_BILL
TO_DELIVER
TO_DELIVER_AND_BILL
TO_RECEIVE
TO_RECEIVE_AND_BILL
UNPAID
UNSPECIFIED
ENABLED
DELETED
DISABLED
\.


--
-- TOC entry 3775 (class 0 OID 21256)
-- Dependencies: 354
-- Data for Name: unit_of_measures; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.unit_of_measures (id, created_at, deleted_at, updated_at, code, enabled, company_id) FROM stdin;
1	2024-11-13 15:39:33.441315+00	\N	2024-11-13 15:39:33.441315+00	MTR	t	\N
2	2024-11-13 15:39:33.441315+00	\N	2024-11-13 15:39:33.441315+00	DAY	t	\N
3	2024-11-13 15:39:33.441315+00	\N	2024-11-13 15:39:33.441315+00	WEE	t	\N
4	2024-11-13 15:39:33.441315+00	\N	2024-11-13 15:39:33.441315+00	MON	t	\N
5	2024-11-13 15:39:33.441315+00	\N	2024-11-13 15:39:33.441315+00	UNIT	t	\N
6	2024-10-31 18:15:08.07301+00	\N	2024-10-31 18:15:08.07301+00	HOUR	t	\N
7	2024-10-31 18:15:08.07301+00	\N	2024-10-31 18:15:08.07301+00	MINUTE	t	\N
8	2024-12-17 20:25:58.014144+00	\N	2024-12-17 20:25:58.014144+00	GRAM	t	\N
\.


--
-- TOC entry 3773 (class 0 OID 21248)
-- Dependencies: 352
-- Data for Name: unit_of_measure_translations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.unit_of_measure_translations (language_code, id, created_at, deleted_at, updated_at, name, base_id) FROM stdin;
en	1	2024-11-13 15:39:33.441315+00	\N	2024-11-13 15:39:33.441315+00	Meter	1
es	2	2024-11-13 15:39:33.441315+00	\N	2024-11-13 15:39:33.441315+00	Metros	1
en	3	2024-11-13 15:39:33.441315+00	\N	2024-11-13 15:39:33.441315+00	Day	2
es	4	2024-11-13 15:39:33.441315+00	\N	2024-11-13 15:39:33.441315+00	Dia	2
en	5	2024-11-13 15:39:33.441315+00	\N	2024-11-13 15:39:33.441315+00	Week	3
es	6	2024-11-13 15:39:33.441315+00	\N	2024-11-13 15:39:33.441315+00	Semana	3
en	7	2024-11-13 15:39:33.441315+00	\N	2024-11-13 15:39:33.441315+00	Months	4
es	8	2024-11-13 15:39:33.441315+00	\N	2024-11-13 15:39:33.441315+00	Mes	4
en	9	2024-11-13 15:39:33.441315+00	\N	2024-11-13 15:39:33.441315+00	Units	5
es	10	2024-11-13 15:39:33.441315+00	\N	2024-11-13 15:39:33.441315+00	Unidades	5
en	11	2024-10-31 18:15:08.07301+00	\N	2024-10-31 18:15:08.07301+00	Hours	6
es	12	2024-10-31 18:15:08.07301+00	\N	2024-10-31 18:15:08.07301+00	Horas	6
en	13	2024-10-31 18:15:08.07301+00	\N	2024-10-31 18:15:08.07301+00	Minutes	7
es	14	2024-10-31 18:15:08.07301+00	\N	2024-10-31 18:15:08.07301+00	Minutos	7
en	15	2024-12-17 20:25:58.014144+00	\N	2024-12-17 20:25:58.014144+00	Gram	8
es	16	2024-12-17 20:25:58.014144+00	\N	2024-12-17 20:25:58.014144+00	Gramo	8
\.


--
-- TOC entry 3782 (class 0 OID 0)
-- Dependencies: 226
-- Name: actions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.actions_id_seq', 197, true);


--
-- TOC entry 3783 (class 0 OID 0)
-- Dependencies: 245
-- Name: entities_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.entities_id_seq', 60, true);


--
-- TOC entry 3784 (class 0 OID 0)
-- Dependencies: 353
-- Name: unit_of_measure_translations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.unit_of_measure_translations_id_seq', 16, true);


--
-- TOC entry 3785 (class 0 OID 0)
-- Dependencies: 355
-- Name: unit_of_measures_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.unit_of_measures_id_seq', 1, true);


-- Completed on 2025-01-21 12:28:37

--
-- PostgreSQL database dump complete
--




insert into parties(party_type_code) values('company');