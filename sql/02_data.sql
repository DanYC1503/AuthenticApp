
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
-- Name: users trg_create_auth_method; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_create_auth_method AFTER INSERT ON public.users FOR EACH ROW EXECUTE FUNCTION public.create_auth_method_for_user();


--
-- Name: audit_logs audit_logs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: auth_methods auth_methods_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.auth_methods
    ADD CONSTRAINT auth_methods_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: recovery_tokens recovery_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.recovery_tokens
    ADD CONSTRAINT recovery_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;

INSERT INTO users (id_number_encrypted,full_name,email,phone_number,date_of_birth,address,username,password_hash,salt,user_type) VALUES ('\x4b387a32432f544243526d6a2f625a487538534e4c4762324e795271424955666e466248626e45657a6e466a6a6d2b4a5277633d','Administrative Coordinator','admin@gmail.com','0970010020','2003-01-01','Cuenca, Ecuador','admin123','\x61346163326234396537343139333633666632393338323138363239373961383633383563346331396537336638366262363031323464613962663832343531','\xfcb34d8cee32b88da8726079d888c571','admin')
