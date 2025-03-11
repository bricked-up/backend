package backend

// Commented out the tests as the db is not updated yet 

// import (
//     "database/sql"
//     "testing"

//     _ "github.com/mattn/go-sqlite3"
// )

// // TestUpdateUserName verifies that the ChangeDisplayName function correctly updates a user's name.
// func TestUpdateUserName(t *testing.T) {
//     // Open an in-memory database for testing.
//     dbase, err := sql.Open("sqlite3", ":memory:")
//     if err != nil {
//         t.Fatalf("failed to open in-memory db: %v", err)
//     }
//     defer dbase.Close()

//     // Create the USER table and insert a sample user row.
//     _, err = dbase.Exec(`
//         CREATE TABLE USER (
//             id INTEGER PRIMARY KEY,
//             verifyid INTEGER,
//             email TEXT UNIQUE NOT NULL,
//             password TEXT NOT NULL,
//             name TEXT NOT NULL,
//             avatar TEXT UNIQUE NOT NULL
//         );
//         INSERT INTO USER (id, email, password, name, avatar)
//         VALUES (1, 'alice@example.com', 'password123', 'Alice', 'alice.png');
//     `)
//     if err != nil {
//         t.Fatalf("failed to create schema or insert test data: %v", err)
//     }

//     // Call the function to update the user's name to "Ivan".
//     err = ChangeDisplayName(1, "Ivan")
//     if err != nil {
//         t.Errorf("UpdateUserName returned error: %v", err)
//     }

//     // Query the database to verify that the name was updated.
//     var updatedName string
//     err = dbase.QueryRow("SELECT name FROM USER WHERE id = 1").Scan(&updatedName)
//     if err != nil {
//         t.Errorf("failed to select updated name: %v", err)
//     }
//     if updatedName != "Ivan" {
//         t.Errorf("expected name 'Ivan', got '%s'", updatedName)
//     }
// }