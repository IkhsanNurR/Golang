CREATE TABLE public.privilege (
	id serial4 NOT NULL,
	privilege_name varchar NOT NULL,
	"createdAt" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"updatedAt" timestamptz NULL,
	status boolean DEFAULT true NOT NULL,
	CONSTRAINT "privilege_pkey" PRIMARY KEY (id)
);