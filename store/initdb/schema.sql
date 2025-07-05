CREATE TABLE Things (
  id    INTEGER PRIMARY KEY,
  What TEXT NOT NULL,
  -- Store dates as string parts.
  -- This is simple, sortable, formattable at will after fetching,
  -- no need to care about timezones or date format if we were
  -- storing an integer based timestamp or a string like yyyy-mm-dd
  Year INTEGER NOT NULL,
  -- [1, 12]
  Month INTEGER NOT NULL,
  -- [1, 31]
  Day INTEGER NOT NULL
);

CREATE INDEX OftenDoneThingsWhat ON Things(What);