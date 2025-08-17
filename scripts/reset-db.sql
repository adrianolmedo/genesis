-- DROP DATABASE basename;

SET session_replication_role = replica;

TRUNCATE TABLE
    invoice_item,
    invoice_header,
    product,
    customer,
    "user"
RESTART IDENTITY CASCADE;

SET session_replication_role = DEFAULT;
