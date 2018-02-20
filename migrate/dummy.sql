INSERT INTO account(id, publickey, balance) VALUES ('anzu1', 'pubkey_anzu1', 0);
INSERT INTO account(id, publickey, balance) VALUES ('anzu2', 'pubkey_anzu2', 0);
INSERT INTO account(id, publickey, balance) VALUES ('anzu3', 'pubkey_anzu3', 0);
INSERT INTO account(id, publickey, balance) VALUES ('anzu4', 'pubkey_anzu4', 0);

INSERT INTO tx(tx, StatusID, Status) VALUES ('{"id":"anzu1", "publickey":"pubkey_anzu1"}', 'ANZU001', 0);
INSERT INTO tx(tx, StatusID, Status) VALUES ('{"id":"anzu2", "publickey":"pubkey_anzu2"}', 'ANZU002', 0);
INSERT INTO tx(tx, StatusID, Status) VALUES ('{"id":"anzu3", "publickey":"pubkey_anzu3"}', 'ANZU003', 0);
INSERT INTO tx(tx, StatusID, Status) VALUES ('{"id":"anzu4", "publickey":"pubkey_anzu4"}', 'ANZU004', 0);

INSERT INTO block(prevhash, txs, creator_id, timestamp) VALUES ('zero', '[{"id":"anzu1", "publickey":"pubkey_anzu1"},{"id":"anzu2", "publickey":"pubkey_anzu2"},{"id":"anzu3", "publickey":"pubkey_anzu3"},{"id":"anzu4", "publickey":"pubkey_anzu4"}]', 'kami', 0);
