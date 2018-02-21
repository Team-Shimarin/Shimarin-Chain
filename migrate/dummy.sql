INSERT INTO account(id, publickey, balance) VALUES ('shima1', 'pubkey_shima1', 1000);
INSERT INTO account(id, publickey, balance) VALUES ('shima2', 'pubkey_shima2', 1000);
INSERT INTO account(id, publickey, balance) VALUES ('shima3', 'pubkey_shima3', 1000);
INSERT INTO account(id, publickey, balance) VALUES ('shima4', 'pubkey_shima4', 1000);

INSERT INTO tx(tx, StatusID, Status) VALUES ('{"id":"shima1", "publickey":"pubkey_shima1"}', 'shima001', 0);
INSERT INTO tx(tx, StatusID, Status) VALUES ('{"id":"shima2", "publickey":"pubkey_shima2"}', 'shima002', 0);
INSERT INTO tx(tx, StatusID, Status) VALUES ('{"id":"shima3", "publickey":"pubkey_shima3"}', 'shima003', 0);
INSERT INTO tx(tx, StatusID, Status) VALUES ('{"id":"shima4", "publickey":"pubkey_shima4"}', 'shima004', 0);

INSERT INTO block(prevhash, txs, creator_id, timestamp, hash) VALUES ('zero', '[{"id":"shima1", "publickey":"pubkey_shima1"},{"id":"shima2", "publickey":"pubkey_shima2"},{"id":"shima3", "publickey":"pubkey_shima3"},{"id":"shima4", "publickey":"pubkey_shima4"}]', 'kami', 0, 'GENESIS BLOCK HASH');

INSERT INTO health(account_id, hp) VALUES ('shima1', 25);
INSERT INTO health(account_id, hp) VALUES ('shima2', 50);
INSERT INTO health(account_id, hp) VALUES ('shima3', 10);
INSERT INTO health(account_id, hp) VALUES ('shima4', 90);
