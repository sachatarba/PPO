--
-- PostgreSQL database dump
--

-- Dumped from database version 16.1 (Debian 16.1-1.pgdg120+1)
-- Dumped by pg_dump version 16.1 (Debian 16.1-1.pgdg120+1)

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
-- Name: client_memberships; Type: TABLE; Schema: public; Owner: gym
--

CREATE TABLE public.client_memberships (
    id uuid NOT NULL,
    start_date date,
    end_date date,
    membership_type_id uuid,
    client_id uuid
);


ALTER TABLE public.client_memberships OWNER TO gym;

--
-- Name: clients; Type: TABLE; Schema: public; Owner: gym
--

CREATE TABLE public.clients (
    id uuid NOT NULL,
    login text,
    password text,
    fullname text,
    email text,
    phone text,
    birthdate date,
    CONSTRAINT check_valid_email CHECK ((email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}$'::text)),
    CONSTRAINT check_valid_international_phone CHECK ((phone ~ '^\+[0-9]+-[0-9]+-[0-9]+-[0-9]+-[0-9]+'::text))
);


ALTER TABLE public.clients OWNER TO gym;

--
-- Name: equipment; Type: TABLE; Schema: public; Owner: gym
--

CREATE TABLE public.equipment (
    id uuid NOT NULL,
    name text,
    description text,
    gym_id uuid
);


ALTER TABLE public.equipment OWNER TO gym;

--
-- Name: gym_trainers; Type: TABLE; Schema: public; Owner: gym
--

CREATE TABLE public.gym_trainers (
    trainer_id uuid NOT NULL,
    gym_id uuid NOT NULL
);


ALTER TABLE public.gym_trainers OWNER TO gym;

--
-- Name: gyms; Type: TABLE; Schema: public; Owner: gym
--

CREATE TABLE public.gyms (
    id uuid NOT NULL,
    name text,
    phone text,
    city text,
    addres text,
    is_chain boolean,
    CONSTRAINT check_valid_international_phone CHECK ((phone ~ '^\+[0-9]+-[0-9]+-[0-9]+-[0-9]+-[0-9]+'::text))
);


ALTER TABLE public.gyms OWNER TO gym;

--
-- Name: membership_types; Type: TABLE; Schema: public; Owner: gym
--

CREATE TABLE public.membership_types (
    id uuid NOT NULL,
    type text,
    description text,
    price real,
    days_duration bigint,
    gym_id uuid
);


ALTER TABLE public.membership_types OWNER TO gym;

--
-- Name: schedules; Type: TABLE; Schema: public; Owner: gym
--

CREATE TABLE public.schedules (
    id uuid NOT NULL,
    day_of_the_week date,
    start_time timestamp with time zone,
    end_time timestamp with time zone,
    client_id uuid,
    training_id uuid
);


ALTER TABLE public.schedules OWNER TO gym;

--
-- Name: trainers; Type: TABLE; Schema: public; Owner: gym
--

CREATE TABLE public.trainers (
    id uuid NOT NULL,
    fullname text,
    email text,
    phone text,
    qualification text,
    unit_price real,
    CONSTRAINT check_valid_email CHECK ((email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}$'::text)),
    CONSTRAINT check_valid_international_phone CHECK ((phone ~ '^\+[0-9]+-[0-9]+-[0-9]+-[0-9]+-[0-9]+'::text))
);


ALTER TABLE public.trainers OWNER TO gym;

--
-- Name: trainings; Type: TABLE; Schema: public; Owner: gym
--

CREATE TABLE public.trainings (
    id uuid NOT NULL,
    title text,
    description text,
    training_type text,
    trainer_id uuid
);


ALTER TABLE public.trainings OWNER TO gym;

--
-- Name: client_memberships client_memberships_pkey; Type: CONSTRAINT; Schema: public; Owner: gym
--

ALTER TABLE ONLY public.client_memberships
    ADD CONSTRAINT client_memberships_pkey PRIMARY KEY (id);


--
-- Name: clients clients_pkey; Type: CONSTRAINT; Schema: public; Owner: gym
--

ALTER TABLE ONLY public.clients
    ADD CONSTRAINT clients_pkey PRIMARY KEY (id);


--
-- Name: equipment equipment_pkey; Type: CONSTRAINT; Schema: public; Owner: gym
--

ALTER TABLE ONLY public.equipment
    ADD CONSTRAINT equipment_pkey PRIMARY KEY (id);


--
-- Name: gym_trainers gym_trainers_pkey; Type: CONSTRAINT; Schema: public; Owner: gym
--

ALTER TABLE ONLY public.gym_trainers
    ADD CONSTRAINT gym_trainers_pkey PRIMARY KEY (trainer_id, gym_id);


--
-- Name: gyms gyms_pkey; Type: CONSTRAINT; Schema: public; Owner: gym
--

ALTER TABLE ONLY public.gyms
    ADD CONSTRAINT gyms_pkey PRIMARY KEY (id);


--
-- Name: membership_types membership_types_pkey; Type: CONSTRAINT; Schema: public; Owner: gym
--

ALTER TABLE ONLY public.membership_types
    ADD CONSTRAINT membership_types_pkey PRIMARY KEY (id);


--
-- Name: schedules schedules_pkey; Type: CONSTRAINT; Schema: public; Owner: gym
--

ALTER TABLE ONLY public.schedules
    ADD CONSTRAINT schedules_pkey PRIMARY KEY (id);


--
-- Name: trainers trainers_pkey; Type: CONSTRAINT; Schema: public; Owner: gym
--

ALTER TABLE ONLY public.trainers
    ADD CONSTRAINT trainers_pkey PRIMARY KEY (id);


--
-- Name: trainings trainings_pkey; Type: CONSTRAINT; Schema: public; Owner: gym
--

ALTER TABLE ONLY public.trainings
    ADD CONSTRAINT trainings_pkey PRIMARY KEY (id);


--
-- Name: clients uni_clients_login; Type: CONSTRAINT; Schema: public; Owner: gym
--

ALTER TABLE ONLY public.clients
    ADD CONSTRAINT uni_clients_login UNIQUE (login);


--
-- Name: client_memberships fk_client_memberships_membership_type; Type: FK CONSTRAINT; Schema: public; Owner: gym
--

ALTER TABLE ONLY public.client_memberships
    ADD CONSTRAINT fk_client_memberships_membership_type FOREIGN KEY (membership_type_id) REFERENCES public.membership_types(id);


--
-- Name: client_memberships fk_clients_client_memberships; Type: FK CONSTRAINT; Schema: public; Owner: gym
--

ALTER TABLE ONLY public.client_memberships
    ADD CONSTRAINT fk_clients_client_memberships FOREIGN KEY (client_id) REFERENCES public.clients(id);


--
-- Name: schedules fk_clients_schedules; Type: FK CONSTRAINT; Schema: public; Owner: gym
--

ALTER TABLE ONLY public.schedules
    ADD CONSTRAINT fk_clients_schedules FOREIGN KEY (client_id) REFERENCES public.clients(id);


--
-- Name: gym_trainers fk_gym_trainers_gym; Type: FK CONSTRAINT; Schema: public; Owner: gym
--

ALTER TABLE ONLY public.gym_trainers
    ADD CONSTRAINT fk_gym_trainers_gym FOREIGN KEY (gym_id) REFERENCES public.gyms(id);


--
-- Name: gym_trainers fk_gym_trainers_trainer; Type: FK CONSTRAINT; Schema: public; Owner: gym
--

ALTER TABLE ONLY public.gym_trainers
    ADD CONSTRAINT fk_gym_trainers_trainer FOREIGN KEY (trainer_id) REFERENCES public.trainers(id);


--
-- Name: equipment fk_gyms_equipments; Type: FK CONSTRAINT; Schema: public; Owner: gym
--

ALTER TABLE ONLY public.equipment
    ADD CONSTRAINT fk_gyms_equipments FOREIGN KEY (gym_id) REFERENCES public.gyms(id);


--
-- Name: membership_types fk_gyms_membership_types; Type: FK CONSTRAINT; Schema: public; Owner: gym
--

ALTER TABLE ONLY public.membership_types
    ADD CONSTRAINT fk_gyms_membership_types FOREIGN KEY (gym_id) REFERENCES public.gyms(id);


--
-- Name: schedules fk_schedules_training; Type: FK CONSTRAINT; Schema: public; Owner: gym
--

ALTER TABLE ONLY public.schedules
    ADD CONSTRAINT fk_schedules_training FOREIGN KEY (training_id) REFERENCES public.trainings(id);


--
-- Name: trainings fk_trainers_trainings; Type: FK CONSTRAINT; Schema: public; Owner: gym
--

ALTER TABLE ONLY public.trainings
    ADD CONSTRAINT fk_trainers_trainings FOREIGN KEY (trainer_id) REFERENCES public.trainers(id);


--
-- PostgreSQL database dump complete
--

