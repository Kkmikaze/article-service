CREATE TABLE posts
(
    created_date TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_date TIMESTAMP WITH TIME ZONE DEFAULT now(),
    id           UUID                     DEFAULT uuid_generate_v4() PRIMARY KEY NOT NULL,
    title        VARCHAR(200)                                                    NOT NULL,
    content      TEXT                                                            NOT NULL,
    category     VARCHAR(100)                                                    NOT NULL,
    status       INT                      DEFAULT 0                              NOT NULL
);

CREATE INDEX idx_posts_title ON posts (title);
CREATE INDEX idx_posts_content ON posts (content);
CREATE INDEX idx_posts_category ON posts (category);
CREATE INDEX idx_posts_status ON posts (status);
