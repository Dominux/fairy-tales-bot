CREATE TYPE fairy_tale_stage AS ENUM ('inited', 'named', 'created');

CREATE TABLE fairy_tales (
    id                UUID                      PRIMARY KEY,
    name              VARCHAR(255),
    init_msg_id       INTEGER          NOT NULL,
    audio_msg_id      INTEGER,
    stage             fairy_tale_stage NOT NULL
);
