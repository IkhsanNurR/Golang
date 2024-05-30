CREATE TABLE public."users" (
	id serial4 NOT NULL,
	full_name text NOT NULL,
	email text NOT NULL,
	"password" text NOT NULL,
	"createdAt" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"updatedAt" timestamptz NULL,
	status boolean DEFAULT true NOT NULL,
	profile_picture TEXT NULL,
	CONSTRAINT "users_pkey" PRIMARY KEY (id)
);
CREATE UNIQUE INDEX "users_email_key" ON public."users" USING btree (email);

INSERT INTO public.users
(full_name, email, "password", "createdAt", "updatedAt", status)
VALUES('Sabin Adams', 'sabin@adams.com', '$2b$10$QJXakrvxPIdhQGcg8RsI4.0gtqOGjjWB4rJ1HlkcpqFMK2YljhZhy', '2024-04-15 03:19:28.180', '2024-04-15 03:20:47.048', true);
INSERT INTO public.users
(full_name, email, "password", "createdAt", "updatedAt", status)
VALUES('Alex Ruheni', 'alex@ruheni.com', '$2b$10$rra3fk1.4e.x1LFttlXFp.4qPEXN.KQzHuulZnN5eYWgcVmvd5SJe', '2024-04-15 03:19:28.197', '2024-04-15 03:20:47.053',true);
INSERT INTO public.users
(full_name, email, "password", "createdAt", "updatedAt", status)
VALUES('Sabin Adams2', 'test2@gmail.com', '$2a$14$5LWwhogxLczy27FnkukkjuTNNE4u5z1U9FoEKqoeT4VpiGQ3Ye5fC', '2024-05-28 17:16:57.465', '2024-05-28 17:16:57.465',true);
INSERT INTO public.users
(full_name, email, "password", "createdAt", "updatedAt", status)
VALUES('Sabin Adams3', 'test3@gmail.com', '$2a$14$FofRjdtBcX6bTkqGW5eqNOeiScvdvxnhEC39o8FdEfXM.fjTFpxU.', '2024-05-28 17:17:31.248', '2024-05-28 17:17:31.248',true);