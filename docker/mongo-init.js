db.createUser({
    user: "todo_account",
    pwd: "1234",
    roles: [
        {role: "readWrite", db: "todo_db"}
    ]
});

db = db.getSiblingDB('todo_db');
db.createCollection('todos');