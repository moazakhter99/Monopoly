
create table if not exists game_board (
    block_id varchar(255) primary key not null,
    block_type varchar(255),
    block_name varchar(255),
    colour varchar(255),
    position int,
    price int

);

create table if not exists game(
    game_id varchar(255),
    match_id varchar(255),
    event varchar(255)

);

