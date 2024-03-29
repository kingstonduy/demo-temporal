CREATE TABLE public.money_transfer (
	from_account_id varchar NOT NULL,
	to_account_id varchar NOT NULL,
	amount int NOT NULL,
	state varchar NOT NULL,
	transaction_id varchar NOT NULL,
	CONSTRAINT money_transfer_pk PRIMARY KEY (transaction_id)
);

CREATE TABLE public.limit_manage (
	account_id varchar NOT NULL,
	amount int NOT NULL,
	CONSTRAINT limit_manage_pk PRIMARY KEY (account_id)
);

CREATE TABLE public.napas (
	account_id varchar NOT NULL,
	account_name varchar NOT NULL,
	amount int NOT NULL,
	CONSTRAINT napas_pk PRIMARY KEY (account_id)
);

CREATE TABLE public.t24 (
	account_id varchar NOT NULL,
	amount int NULL,
	CONSTRAINT t24_pk PRIMARY KEY (account_id)
);

CREATE TABLE public."transaction" (
	transaction_id varchar NOT NULL,
	from_account_id varchar NOT NULL,
	to_account_id varchar NOT NULL,
	amount int NOT NULL,
	state varchar NOT NULL,
	CONSTRAINT transaction_pk PRIMARY KEY (transaction_id)
);
