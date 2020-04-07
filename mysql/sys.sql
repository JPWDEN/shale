USE sys;

CREATE TABLE Todos (
    id INT NOT NULL AUTO_INCREMENT,
    title VARCHAR(32) NOT NULL,
    body VARCHAR(255),
    category VARCHAR(255),
    item_priority INT default 0,
    publish_date DATE NOT NULL,
    active BOOLEAN NOT NULL,
    PRIMARY KEY (id)
);
