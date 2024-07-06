# Pantsbase 👖

    "It's fine, I guess" 
      ~ No one

Think of pocketbase, you know? That awesome piece of software? Ok, so pantsbase is like that, but worse.

What I wanted to do by making this project was test how far I can take the standart library and my own go knowledge. I can say it's been a great learning exercise.

As of v1.x.x, pantsbase can **CREATE, VIEW and DELETE tables** and **a full CRUD for table rows** both from the GUI and from the API. It still cannot do relations and does not have a way to protect the API through any type of auth (but this is a priority for the next update).

# API reference

**Any receiving or sending of data is done with json**  
**All table coulmns are nullable**

## Table endpoints

- `GET /api/tables` returns the all existing tables and their schemas
  - **URL parameters:**
  - `?names_only` returns just the names of the tables, without the schemas

- `POST /api/new_table/{table_name}` creates a new table. The schema should be the body of the request where the keys of the json are the column names and the values are the data types (data types com only be "TEXT", "INTEGER" or "REAL").
  - **Example body:** `{"col1": "REAL", "col 2": "INTEGER", "col_3": "TEXT"}`

- `DELETE /api/delete_table/{table_name}` deletes the specified table

## Rows endpoints

- `GET /api/tables/{table_name}` queries for data in a table
  - **URL parameters:**
  - `?limit={int}` limits the amount of data that is queried
  - `?order_by={column_name}` chooses a column to use for ordering the results
  - `?order_direction={DESC|ASC}` only has an effect if order_by if present and valid, the only accepted values are *DESC* and *ASC* (the default if *DESC*)

- `GET /api/tables/{table_name}/{id}` queries for one row that matches the *id*

- `POST /api/tables/{table_name}` creates on or more new rows in the *table*
  - **Example body 1:** `{"col1": 1.2, "col 2": 10, "col_3": "some text"}`
  - **Example body 2:** `[{"col1": null, "col 2": -22, "col_3": "something"}, {"col1": 9328432.4324, "col 2": null, "col_3": null}]`

- `PUT /api/tables/{table_name}/{id}` or `PATCH /api/tables/{table_name}/{id}` updates a row that matches the id, the body should contain only the columns you wnat to udpate
  - **Example body:** `{"col_3": "new updated text"}`

- `DELETE /api/tables/{table_name}/{id}` deletes the row that matches te *id* from the *table*
