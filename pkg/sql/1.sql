CREATE TABLE spots (
    name varchar(255) NOT NULL,
    created_at bigint NOT NULL,
    PRIMARY KEY (name)
);

CREATE TABLE photos (
    id varchar(36) NOT NULL,
    poi_id varchar(36),
    src text NOT NULL,
    thunbnail text NOT NULL,
    spot varchar(255) NOT NULL,
    created_at bigint NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE pois (
    id varchar(36) NOT NULL,
    photo_id varchar(36),
    latitude decimal NOT NULL,
    longitude decimal NOT NULL,
    comment text,
    created_at bigint NOT NULL,
    updated_at bigint NOT NULL,
    PRIMARY KEY (id)
);
