CREATE TABLE IF NOT EXISTS inventory_changes (
	id integer PRIMARY KEY AUTOINCREMENT,
	quantity integer NOT NULL,
	direction integer NOT NULL,
	item_code integer NOT NULL,
  created_at integer NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
	code integer PRIMARY KEY,
	name integer NOT NULL
);
