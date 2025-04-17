-- DROP SCHEMA
-- DROP SCHEMA IF EXISTS checkout CASCADE;
-- ENV SCHEMA
CREATE SCHEMA IF NOT EXISTS checkout;

-- GATEWAY TABLE
CREATE TABLE
    IF NOT EXISTS checkout.gateways (
        "id" character varying(255) NOT NULL PRIMARY KEY,
        "key" character varying(255) NOT NULL,
        "name" character varying(255) NOT NULL,
        "acronym" character varying(255) NOT NULL,
        "icon" text,
        "type" character varying(255) NOT NULL,
        "can_process" boolean NOT NULL DEFAULT 'FALSE',
        "can_settle" boolean NOT NULL DEFAULT 'FALSE',
        "created_at" timestamp without time zone NOT NULL,
        "updated_at" timestamp without time zone NOT NULL DEFAULT 'NOW()'
    );

-- INSERT GATEWAYS
INSERT INTO
    checkout.gateways (
        "id",
        "key",
        "name",
        "acronym",
        "icon",
        "type",
        "can_process",
        "can_settle",
        "created_at"
    )
VALUES
    (
        'D1D47D97-101A-4DCC-965A-6E7C236BC2EE',
        'TELEBIRR',
        'Telebirr',
        '',
        'https://play-lh.googleusercontent.com/Mtnybz6w7FMdzdQUbc7PWN3_0iLw3t9lUkwjmAa_usFCZ60zS0Xs8o00BW31JDCkAiQk',
        'WALLET',
        'true',
        'false',
        'NOW()'
    ),
    (
        'FD5D56EF-061D-432A-A93C-C57BA3061914',
        'CBE_BIRR',
        'CBE Birr',
        '',
        'https://play-lh.googleusercontent.com/rcSKabjkP2GfX1_I_VXBfhQIPdn_HPXj5kbkDoL4cu5lpvcqPsGmCqfqxaRrSI9h5_A',
        'WALLET',
        'true',
        'false',
        'NOW()'
    ),
    (
        '0EC7A9C8-B40D-4DC8-9082-5A092697958A',
        'CYBERSOURCE',
        'VISA MASTERCARD',
        '',
        'https://getsby.com/wp-content/uploads/2023/01/Visa-Mastercard-1-1024x378.png',
        'CARD',
        'true',
        'false',
        'NOW()'
    ),
    (
        'B078C9AE-045B-4D3E-981F-17D46A1E8F75',
        'AWINETAA',
        'Awash International Bank SC',
        'Awash',
        'https://upload.wikimedia.org/wikipedia/commons/3/33/Awash_International_Bank.png',
        'BANK',
        'false',
        'false',
        'NOW()'
    ),
    (
        '1B2D3252-C25A-4952-BB24-30919FF23C94',
        'BUNAETAA',
        'Bunna International Bank SC',
        'Bunna',
        'https://z-p3-scontent.fadd1-1.fna.fbcdn.net/v/t39.30808-6/272417075_4734443969964854_7325903198322919211_n.jpg?_nc_cat=103&ccb=1-7&_nc_sid=6ee11a&_nc_ohc=Bh4lLYSzbAgQ7kNvgGM6iHo&_nc_oc=AdlnWv8iU5J5E-N4gM0VeNpQsGoDAN1L7yLc8JiLvXdqcxti2w4FQNciLOIH9FysaNg&_nc_zt=23&_nc_ht=z-p3-scontent.fadd1-1.fna&_nc_gid=-wTAXo6UP9h1fEn6JMygEg&oh=00_AYHzaOpQuXSqV4XfHO2bmD95rregJSgUbwhHjJTD1kH-Og&oe=67E1A7AB',
        'BANK',
        'false',
        'false',
        'NOW()'
    ),
    (
        '05E70335-E33A-4470-8C56-315DCB15B5A2',
        'M_PESA',
        '',
        'M-Pesa',
        'https://upload.wikimedia.org/wikipedia/commons/0/0b/M-PESA.png',
        'WALLET',
        'true',
        'false',
        'NOW()'
    )
     ON CONFLICT (id) DO NOTHING;


-- LITE TRANSACTION TABLE
CREATE TABLE
    IF NOT EXISTS checkout.transactions (
        "id" character varying(255) NOT NULL PRIMARY KEY,
        "for" character varying(255) NOT NULL,
        "to" character varying(255) NOT NULL,
        "ttl" integer,
        "pricing" text,
        "status" text,
        "gateway" character varying(255),
        "type" character varying(255),
        "details" text,
        "created_at" timestamp without time zone NOT NULL,
        "updated_at" timestamp without time zone NOT NULL DEFAULT 'NOW()'
    )