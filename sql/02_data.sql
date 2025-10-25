DROP EXTENSION IF EXISTS "uuid-ossp";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: create_auth_method_for_user(); Type: FUNCTION; Schema: public; Owner: -
--
CREATE TABLE users (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    id_number_encrypted bytea,
    full_name text NOT NULL,
    email text NOT NULL,
    phone_number text,
    date_of_birth date,
    address text,
    create_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    username character varying(100) NOT NULL,
    password_hash bytea,
    salt bytea,
    is_verified boolean DEFAULT false,
    last_login timestamp without time zone,
    account_status character varying(20) DEFAULT 'active'::character varying,
    oauth_provider character varying(50),
    oauth_id text,
    user_type character varying(20) DEFAULT 'client'::character varying NOT NULL,
    CONSTRAINT users_account_status_check CHECK (((account_status)::text = ANY ((ARRAY['active'::character varying, 'pending'::character varying, 'disabled'::character varying])::text[])))
);

CREATE TABLE audit_logs (
    id bigint NOT NULL,
    user_id uuid,
    action character varying(50) NOT NULL,
    ip_address inet,
    user_agent text,
    "timestamp" timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    metadata jsonb
);

CREATE TABLE recovery_tokens (
    id integer NOT NULL,
    user_id uuid,
    token_hash bytea NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    used boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

--
-- Name: log_user_action(text, text, text, text, text); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.log_user_action(p_username text, p_action text, p_ip_address text, p_user_agent text, p_metadata text) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE 
    v_user_id UUID;
BEGIN
    -- Find the user's ID
    SELECT id INTO v_user_id
    FROM users
    WHERE username = p_username;

    -- Insert into audit_logs with IP cast to inet
    INSERT INTO audit_logs(user_id, action, ip_address, user_agent, metadata)
    VALUES (v_user_id, p_action, p_ip_address::inet, p_user_agent, p_metadata::jsonb);
    
EXCEPTION
    WHEN invalid_text_representation THEN
        -- Handle invalid IP format by using NULL
        INSERT INTO audit_logs(user_id, action, ip_address, user_agent, metadata)
        VALUES (v_user_id, p_action, NULL, p_user_agent, p_metadata::jsonb);
END;
$$;


--
-- Name: log_user_registration(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.log_user_registration() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF NEW.oauth_id IS NOT NULL THEN
        -- Log OAuth registration
        INSERT INTO audit_logs(user_id, action, ip_address, user_agent, metadata)
        VALUES (NEW.id, 'oauth_register', NEW.ip_address, NEW.user_agent, NEW.metadata);
    ELSE
        -- Log password registration
        INSERT INTO audit_logs(user_id, action, ip_address, user_agent, metadata)
        VALUES (NEW.id, 'password_register', NEW.ip_address, NEW.user_agent, NEW.metadata);
    END IF;

    RETURN NEW;
END;
$$;


--
-- Name: retrieve_user_audits(text); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.retrieve_user_audits(p_email text) RETURNS TABLE(user_email text, action text, ip_address text, user_agent text, metadata jsonb, created_at timestamp without time zone)
    LANGUAGE plpgsql
    AS $$
DECLARE
    v_user_id UUID;
BEGIN
    -- Get the user's ID
    SELECT id INTO v_user_id
    FROM users
    WHERE email = p_email;

    -- Return the audit logs for that user, replace NULL ip_address with '0.0.0.0'
    RETURN QUERY
    SELECT 
        u.email AS user_email,            -- fully qualify to avoid ambiguity
        a.action::TEXT,                   -- cast varchar → text
        COALESCE(a.ip_address::TEXT, '0.0.0.0'), -- replace NULL with default
        a.user_agent,
        a.metadata,
        a.timestamp AS created_at
    FROM audit_logs a
    JOIN users u ON u.id = a.user_id
    WHERE a.user_id = v_user_id
    ORDER BY a.timestamp DESC;
END;
$$;


--
-- Name: retrieve_users(text); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.retrieve_users(p_email text) RETURNS TABLE(username text, full_name text, email text, phone_number text, date_of_birth text, address text, create_date text, account_status text, oauth_provider text, oauth_id text, type_name text)
    LANGUAGE plpgsql
    AS $$
DECLARE
    v_user_type TEXT;
    v_username TEXT;
BEGIN
    -- Check the type and username of the user making the request
    SELECT u.user_type, u.username INTO v_user_type, v_username
    FROM users u
    WHERE u.email = retrieve_users.p_email;

    -- If the requesting user is a client, return no rows (deny access)
    IF v_user_type = 'client' THEN
        RETURN;
    END IF;

    -- Otherwise, return ALL users EXCEPT the requesting user
    RETURN QUERY
    SELECT
        u.username::TEXT,
        u.full_name,
        u.email,
        u.phone_number,
        u.date_of_birth::TEXT,
        u.address,
        u.create_date::TEXT,
        u.account_status::TEXT,
        COALESCE(u.oauth_provider::TEXT, ''),
        COALESCE(u.oauth_id, ''),
        u.user_type::TEXT AS type_name
    FROM users u
    WHERE u.email != retrieve_users.p_email;  -- Exclude the requesting user
END;
$$;


--
-- Name: update_auth_method_last_login(text); Type: FUNCTION; Schema: public; Owner: -
--


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: audit_logs; Type: TABLE; Schema: public; Owner: -
--



--
-- Name: audit_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.audit_logs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;



CREATE SEQUENCE public.recovery_tokens_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;



ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_pkey PRIMARY KEY (id);


--
-- Name: recovery_tokens recovery_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.recovery_tokens
    ADD CONSTRAINT recovery_tokens_pkey PRIMARY KEY (id);


--
-- Name: users unique_oauth; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT unique_oauth UNIQUE (oauth_provider, oauth_id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_phone_number_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_phone_number_key UNIQUE (phone_number);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);

--
-- Name: audit_logs audit_logs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: recovery_tokens recovery_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.recovery_tokens
    ADD CONSTRAINT recovery_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;

INSERT INTO users (id, id_number_encrypted,full_name,email,phone_number,date_of_birth,address,username,password_hash,salt,user_type,oauth_id) VALUES (gen_random_uuid(),'\x4b387a32432f544243526d6a2f625a487538534e4c4762324e795271424955666e466248626e45657a6e466a6a6d2b4a5277633d','Administrative Coordinator','admin@gmail.com','0970010020','2003-01-01','Cuenca, Ecuador','admin123','\x61346163326234396537343139333633666632393338323138363239373961383633383563346331396537336638366262363031323464613962663832343531','\xfcb34d8cee32b88da8726079d888c571','admin', NULL)



--
-- PostgreSQL database dump complete
--

