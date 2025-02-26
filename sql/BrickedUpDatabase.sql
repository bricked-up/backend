-- database: :memory:
-- for the database, memory is local on my computer, needs to be changed on project deployment
PRAGMA foreign_keys = ON;

-- Organization table
CREATE TABLE organization (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Organization roles table
CREATE TABLE organization_role (
    id INTEGER PRIMARY KEY,
    organization_id INTEGER,
    name TEXT NOT NULL,
    can_read BOOLEAN DEFAULT 0,
    can_write BOOLEAN DEFAULT 0,
    can_execute BOOLEAN DEFAULT 0,
    FOREIGN KEY (organization_id) REFERENCES organization(id)
);

-- User table
CREATE TABLE user (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

-- Organization members table
CREATE TABLE organization_member (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    organization_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (organization_id) REFERENCES organization(id)
);

-- Project table
CREATE TABLE project (
    id INTEGER PRIMARY KEY,
    organization_id INTEGER,
    name TEXT NOT NULL,
    budget INTEGER CHECK (budget >= 0),
    charter TEXT,
    is_archived BOOLEAN DEFAULT 0,
    FOREIGN KEY (organization_id) REFERENCES organization(id)
);

-- Project roles table
CREATE TABLE project_role (
    id INTEGER PRIMARY KEY,
    project_id INTEGER,
    name TEXT NOT NULL,
    can_read BOOLEAN DEFAULT 0,
    can_write BOOLEAN DEFAULT 0,
    can_execute BOOLEAN DEFAULT 0,
    FOREIGN KEY (project_id) REFERENCES project(id)
);

-- Project members table
CREATE TABLE project_member (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    project_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (project_id) REFERENCES project(id)
);

-- Issues table
CREATE TABLE issue (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    tag_id INTEGER,
    priority_id INTEGER,
    created_at DATE NOT NULL,
    completed_at DATE,
    cost INTEGER CHECK (cost >= 0),
    FOREIGN KEY (tag_id) REFERENCES tag(id),
    FOREIGN KEY (priority_id) REFERENCES priority(id)
);

-- Tags table
CREATE TABLE tag (
    id INTEGER PRIMARY KEY,
    project_id INTEGER,
    name TEXT NOT NULL,
    color TEXT NOT NULL,
    FOREIGN KEY (project_id) REFERENCES project(id)
);

-- Priorities table
CREATE TABLE priority (
    id INTEGER PRIMARY KEY,
    project_id INTEGER,
    name TEXT NOT NULL,
    priority_level INTEGER CHECK (priority_level >= 1),
    FOREIGN KEY (project_id) REFERENCES project(id)
);

-- Organization projects table (many-to-many)
CREATE TABLE organization_project (
    id INTEGER PRIMARY KEY,
    organization_id INTEGER,
    project_id INTEGER,
    FOREIGN KEY (organization_id) REFERENCES organization(id),
    FOREIGN KEY (project_id) REFERENCES project(id)
);

-- Project issues table (many-to-many)
CREATE TABLE project_issue (
    id INTEGER PRIMARY KEY,
    project_id INTEGER,
    issue_id INTEGER,
    FOREIGN KEY (project_id) REFERENCES project(id),
    FOREIGN KEY (issue_id) REFERENCES issue(id)
);

-- User issues table (many-to-many)
CREATE TABLE user_issue (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    issue_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (issue_id) REFERENCES issue(id)
);

-- Reminders table
CREATE TABLE reminder (
    id INTEGER PRIMARY KEY,
    issue_id INTEGER,
    user_id INTEGER,
    FOREIGN KEY (issue_id) REFERENCES issue(id),
    FOREIGN KEY (user_id) REFERENCES user(id)
);


