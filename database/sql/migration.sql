CREATE TABLE events (
  id INTEGER PRIMARY KEY AUTOINCREMENT, --sqllite bug
  description TEXT NOT NULL,
  state VARCHAR(20) NOT NULL DEFAULT 'in progress' CHECK (state IN ('in progress', 'done'))
);
