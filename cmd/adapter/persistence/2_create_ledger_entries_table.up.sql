SET statement_timeout = 60000;
set lock_timeout = 60000;

CREATE TABLE ledger_entries
(
  id VARCHAR(256),
  item VARCHAR(256),
  receiver VARCHAR(256),
  amount JSONB,
  due_date TIMESTAMP,
  paid_date TIMESTAMP,
  notes TEXT,
  PRIMARY KEY (id)
);
