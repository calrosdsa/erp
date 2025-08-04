--
-- PostgreSQL database dump
--

-- Dumped from database version 17.2 (Debian 17.2-1.pgdg120+1)
-- Dumped by pg_dump version 17.1

-- Started on 2025-01-27 08:24:06

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 361 (class 1259 OID 22641)
-- Name: address_and_contacts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.address_and_contacts (
    doc_id bigint NOT NULL,
    shipping_address_id bigint,
    contact_id bigint,
    billing_address_id bigint
);


ALTER TABLE public.address_and_contacts OWNER TO postgres;

--
-- TOC entry 3605 (class 2606 OID 22645)
-- Name: address_and_contacts address_and_contacts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.address_and_contacts
    ADD CONSTRAINT address_and_contacts_pkey PRIMARY KEY (doc_id);


--
-- TOC entry 3606 (class 1259 OID 22651)
-- Name: fki_fk_aac_doc_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_fk_aac_doc_id ON public.address_and_contacts USING btree (doc_id);


--
-- TOC entry 3607 (class 2606 OID 22664)
-- Name: address_and_contacts fk_aac_billing; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.address_and_contacts
    ADD CONSTRAINT fk_aac_billing FOREIGN KEY (billing_address_id) REFERENCES public.addresses(id);


--
-- TOC entry 3608 (class 2606 OID 22646)
-- Name: address_and_contacts fk_aac_doc_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.address_and_contacts
    ADD CONSTRAINT fk_aac_doc_id FOREIGN KEY (doc_id) REFERENCES public.parties(id) ON DELETE CASCADE;


--
-- TOC entry 3609 (class 2606 OID 22658)
-- Name: address_and_contacts fk_acc_contact; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.address_and_contacts
    ADD CONSTRAINT fk_acc_contact FOREIGN KEY (contact_id) REFERENCES public.contacts(id);


--
-- TOC entry 3610 (class 2606 OID 22652)
-- Name: address_and_contacts fk_acc_shipping; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.address_and_contacts
    ADD CONSTRAINT fk_acc_shipping FOREIGN KEY (shipping_address_id) REFERENCES public.addresses(id);


-- Completed on 2025-01-27 08:24:06

--
-- PostgreSQL database dump complete
--

