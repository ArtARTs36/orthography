CREATE TABLE words (
    nominative character varying NOT NULL,
    genitive character varying,
    dative character varying,
    accusative character varying,
    instrumental character varying,
    prepositional character varying,

    CONSTRAINT words_pk PRIMARY KEY (nominative)
);
