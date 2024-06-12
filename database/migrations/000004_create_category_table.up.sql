CREATE TABLE public.category (
    id serial PRIMARY KEY,
    category_name TEXT NOT NULL,
    category_detail TEXT NOT NULL,
    status BOOLEAN DEFAULT true NOT NULL,
    "createdAt" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" TIMESTAMPTZ NULL
);
