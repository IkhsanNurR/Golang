CREATE TABLE public.users_privilege (
    id serial PRIMARY KEY,
    id_users INT NOT NULL,
    id_privilege INT NOT NULL,
    "createdAt" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" TIMESTAMPTZ NULL,
    status BOOLEAN DEFAULT true NOT NULL,
    CONSTRAINT fk_users FOREIGN KEY (id_users) REFERENCES public.users(id) ON DELETE CASCADE,
    CONSTRAINT fk_privilege FOREIGN KEY (id_privilege) REFERENCES public.privilege(id) ON DELETE CASCADE
);
