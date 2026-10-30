-- 016_transaction_labels.sql
CREATE TABLE IF NOT EXISTS transaction_labels (
  transaction_id TEXT NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
  label_id TEXT NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
  PRIMARY KEY (transaction_id, label_id)
);

CREATE INDEX IF NOT EXISTS idx_transaction_labels_label_id ON transaction_labels(label_id);
