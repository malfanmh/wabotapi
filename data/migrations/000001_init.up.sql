create table clients
(
    id                     bigint auto_increment
        primary key,
    name                   varchar(100)                        null,
    finpay_merchant_id     varchar(100)                        null,
    finpay_secret          varchar(255)                        null,
    finpay_callback_url    varchar(255)                        null,
    wa_phone               varchar(30)                         null,
    wa_phone_id            varchar(100)                        null,
    wa_business_account_id varchar(100)                        null,
    created_at             timestamp default CURRENT_TIMESTAMP not null,
    updated_at             timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
);

create table customers
(
    id              bigint auto_increment
        primary key,
    client_id       bigint                                  not null,
    wa_id           varchar(30)                             not null,
    email           varchar(100)  default ''                not null,
    full_name       varchar(100)  default ''                not null,
    address         varchar(1000) default ''                not null,
    birth_date      date                                    null,
    identity_number varchar(100)  default ''                not null,
    identity_type   varchar(20)                             null comment 'oneof[ktp,kitas,sim]',
    gender          varchar(1)    default ''                not null,
    access          int           default 0                 not null comment 'oneof [0: public, 1:registered, 2:activated]',
    created_at      timestamp     default CURRENT_TIMESTAMP not null,
    updated_at      timestamp     default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
);

create table wabot_db.message_actions
(
    id          bigint auto_increment
        primary key,
    message_id  bigint                                 not null,
    slug        varchar(100) default ''                not null,
    title       varchar(20)  default ''                not null,
    description varchar(255) default ''                not null,
    display     tinyint      default 0                 not null,
    access      int          default 0                 not null comment 'oneof [0:public, 1:registered, 2:activated]',
    seq         varchar(20)  default '0'               not null,
    created_at  timestamp    default CURRENT_TIMESTAMP not null,
    updated_at  timestamp    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
);

create table message_flows
(
    id             bigint auto_increment
        primary key,
    client_id      bigint                                 null,
    keyword        varchar(100)                           not null,
    message_id     bigint                                 null,
    access         int          default 0                 not null,
    display        tinyint      default 0                 not null,
    seq            varchar(100) default ''                not null,
    created_at     timestamp    default CURRENT_TIMESTAMP not null,
    updated_at     timestamp    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    validate_input tinyint(1)   default 0                 not null,
    checkout       tinyint(1)   default 0                 not null
);

create table messages
(
    id            bigint auto_increment
        primary key,
    client_id     bigint                                null,
    slug          varchar(100)                          null,
    type          varchar(100)                          null comment 'oneof [text, button, interactive]',
    access        int         default 0                 not null comment 'oneof [0:public, 1:registered, 2:activated]',
    button        varchar(50) default ''                not null,
    header_text   text                                  not null,
    preview_url   tinyint     default 0                 not null,
    body_text     text                                  not null,
    footer_text   text                                  not null,
    with_metadata tinyint(1)                            null,
    created_at    timestamp   default CURRENT_TIMESTAMP not null,
    updated_at    timestamp   default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
);

create table payments
(
    id               bigint auto_increment
        primary key,
    created_at       timestamp      default CURRENT_TIMESTAMP not null,
    updated_at       timestamp      default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    expired_at       timestamp                                null,
    client_id      bigint                                   null,
    customer_id      bigint                                   null,
    ref_id           varchar(50)    default ''                not null,
    payment_provider varchar(50)    default 'finpay'          not null,
    payment_ref_id   varchar(50)                              null,
    payment_type     varchar(20)                              null comment 'oneof[va, bank_transfer, qris, indomaret]',
    payment_code     varchar(100)                             null,
    payment_item     varchar(100)                             null,
    status           varchar(20)                              not null comment 'oneoff [PENDING,PAID,EXPIRED]',
    amount           decimal(30, 2) default 0.00              not null,
    fee              decimal(30, 2) default 0.00              not null
);

create table products
(
    id          bigint auto_increment
        primary key,
    created_at  timestamp      default CURRENT_TIMESTAMP not null,
    updated_at  timestamp      default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    client_id   bigint                                   not null,
    customer_id bigint         default 0                 not null,
    name        varchar(200)   default ''                not null,
    description text                                     null,
    slug        varchar(200)   default ''                not null,
    price       decimal(30, 2) default 0.00              not null,
    stock       decimal(30, 2) default 0.00              not null,
    image       varchar(250)   default ''                not null
);

create table session
(
    id         bigint auto_increment
        primary key,
    client_id  bigint                                 not null,
    wa_id      varchar(30)                            not null,
    access     int          default 0                 not null,
    seq        varchar(10)  default '0'               not null,
    slug       varchar(100) default ''                not null,
    input      varchar(100) default ''                not null,
    created_at timestamp    default CURRENT_TIMESTAMP not null,
    updated_at timestamp    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    expired_at timestamp                              null
);