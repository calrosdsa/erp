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

