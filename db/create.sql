--
-- PostgreSQL database dump
--

-- Dumped from database version 15.3
-- Dumped by pg_dump version 15.3

-- Started on 2025-01-19 13:52:46

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
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
-- TOC entry 217 (class 1259 OID 18081)
-- Name: cred; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cred (
    id bigint NOT NULL,
    key character varying,
    value character varying,
    res bigint
);


ALTER TABLE public.cred OWNER TO postgres;

--
-- TOC entry 216 (class 1259 OID 18080)
-- Name: cred_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.cred_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.cred_id_seq OWNER TO postgres;

--
-- TOC entry 3342 (class 0 OID 0)
-- Dependencies: 216
-- Name: cred_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.cred_id_seq OWNED BY public.cred.id;


--
-- TOC entry 219 (class 1259 OID 18090)
-- Name: results; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.results (
    id bigint NOT NULL,
    csp boolean,
    xframe boolean,
    url character varying,
    scan bigint
);


ALTER TABLE public.results OWNER TO postgres;

--
-- TOC entry 218 (class 1259 OID 18089)
-- Name: results_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.results_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.results_id_seq OWNER TO postgres;

--
-- TOC entry 3343 (class 0 OID 0)
-- Dependencies: 218
-- Name: results_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.results_id_seq OWNED BY public.results.id;


--
-- TOC entry 215 (class 1259 OID 18074)
-- Name: scans; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.scans (
    id bigint NOT NULL,
    start timestamp without time zone,
    "end" timestamp without time zone
);


ALTER TABLE public.scans OWNER TO postgres;

--
-- TOC entry 214 (class 1259 OID 18073)
-- Name: scans_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.scans_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.scans_id_seq OWNER TO postgres;

--
-- TOC entry 3344 (class 0 OID 0)
-- Dependencies: 214
-- Name: scans_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.scans_id_seq OWNED BY public.scans.id;


--
-- TOC entry 3184 (class 2604 OID 18084)
-- Name: cred id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cred ALTER COLUMN id SET DEFAULT nextval('public.cred_id_seq'::regclass);


--
-- TOC entry 3185 (class 2604 OID 18093)
-- Name: results id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.results ALTER COLUMN id SET DEFAULT nextval('public.results_id_seq'::regclass);


--
-- TOC entry 3183 (class 2604 OID 18077)
-- Name: scans id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.scans ALTER COLUMN id SET DEFAULT nextval('public.scans_id_seq'::regclass);


--
-- TOC entry 3189 (class 2606 OID 18088)
-- Name: cred cred_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cred
    ADD CONSTRAINT cred_pkey PRIMARY KEY (id);


--
-- TOC entry 3191 (class 2606 OID 18097)
-- Name: results results_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.results
    ADD CONSTRAINT results_pkey PRIMARY KEY (id);


--
-- TOC entry 3187 (class 2606 OID 18079)
-- Name: scans scans_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.scans
    ADD CONSTRAINT scans_pkey PRIMARY KEY (id);


--
-- TOC entry 3192 (class 1259 OID 18108)
-- Name: results_url_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX results_url_idx ON public.results USING btree (url varchar_ops);


--
-- TOC entry 3193 (class 2606 OID 18103)
-- Name: cred cred_res_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cred
    ADD CONSTRAINT cred_res_fkey FOREIGN KEY (res) REFERENCES public.results(id);


--
-- TOC entry 3194 (class 2606 OID 18098)
-- Name: results results_scan_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.results
    ADD CONSTRAINT results_scan_fkey FOREIGN KEY (scan) REFERENCES public.scans(id);


-- Completed on 2025-01-19 13:52:46

--
-- PostgreSQL database dump complete
--

