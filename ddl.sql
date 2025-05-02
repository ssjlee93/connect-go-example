CREATE KEYSPACE IF NOT EXISTS greetings WITH REPLICATION = {
  'class': 'SimpleStrategy',
  'replication_factor': 1
};

CREATE TABLE IF NOT EXISTS greetings.greeting (
  id UUID PRIMARY KEY,
  name TEXT,
  message TEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
);

CREATE INDEX IF NOT EXISTS name, message ON greetings.greeting (name, message);