-- Creating tables for primitive entities

PRAGMA foreign_keys = ON;

-- Table to store organizations
CREATE TABLE ORGANIZATION (
    id INTEGER PRIMARY KEY,    -- Unique ID for each organization
    name TEXT UNIQUE NOT NULL   -- Name of the organization (unique)
);

-- Table to store roles for organizations
CREATE TABLE ORG_ROLE (
    id INTEGER PRIMARY KEY,    -- Unique ID for each role
    orgid INTEGER NOT NULL,    -- Foreign key referencing the organization
    name TEXT NOT NULL,        -- Role name
    can_read BOOLEAN NOT NULL, -- Permission to read
    can_write BOOLEAN NOT NULL, -- Permission to write
    can_exec BOOLEAN NOT NULL, -- Permission to execute
    FOREIGN KEY (orgid) REFERENCES ORGANIZATION(id) -- Enforcing foreign key constraint
);

-- Table to store users
CREATE TABLE USER (
    id INTEGER PRIMARY KEY,    -- Unique ID for each user
    verifyid INTEGER,          -- Foreign key referencing the VERIFY_USER table (verification)
    email TEXT UNIQUE NOT NULL, -- User's email (unique)
    password TEXT NOT NULL,    -- User's password
    name TEXT NOT NULL,        -- User's name
    avatar TEXT UNIQUE NOT NULL, -- User's avatar (unique)
    FOREIGN KEY (verifyid) REFERENCES VERIFY_USER(id) -- Enforcing foreign key constraint
);

-- Table to store session details for users
CREATE TABLE SESSION (
    id INTEGER PRIMARY KEY,    -- Unique ID for each session
    userid INTEGER NOT NULL,   -- Foreign key referencing the USER table
    expires DATE NOT NULL,     -- Expiry date of the session
    FOREIGN KEY (userid) REFERENCES USER(id) -- Enforcing foreign key constraint
);

-- Table to store projects
CREATE TABLE PROJECT (
    id INTEGER PRIMARY KEY,    -- Unique ID for each project
    orgid INTEGER NOT NULL,    -- Foreign key referencing the ORGANIZATION table
    name TEXT NOT NULL,        -- Project name
    budget INTEGER NOT NULL,   -- Project budget
    charter TEXT NOT NULL,     -- Project charter or description
    archived BOOLEAN NOT NULL, -- Flag to indicate if the project is archived
    FOREIGN KEY (orgid) REFERENCES ORGANIZATION(id) -- Enforcing foreign key constraint
);

-- Table to store roles for projects
CREATE TABLE PROJECT_ROLE (
    id INTEGER PRIMARY KEY,    -- Unique ID for each project role
    projectid INTEGER NOT NULL, -- Foreign key referencing the PROJECT table
    name TEXT NOT NULL,         -- Role name
    can_read BOOLEAN NOT NULL,  -- Permission to read
    can_write BOOLEAN NOT NULL, -- Permission to write
    can_exec BOOLEAN NOT NULL,  -- Permission to execute
    FOREIGN KEY (projectid) REFERENCES PROJECT(id) -- Enforcing foreign key constraint
);

-- Table to store issues for projects
CREATE TABLE ISSUE (
    id INTEGER PRIMARY KEY,    -- Unique ID for each issue
    title TEXT NOT NULL,       -- Title of the issue
    desc TEXT NOT NULL,        -- Description of the issue
    tagid INTEGER NOT NULL,    -- Foreign key referencing the TAG table
    priorityid INTEGER NOT NULL, -- Foreign key referencing the PRIORITY table
    created DATE NOT NULL,     -- Date when the issue was created
    completed DATE,            -- Date when the issue was completed
    cost INTEGER NOT NULL,     -- Cost associated with the issue
    FOREIGN KEY (tagid) REFERENCES TAG(id),          -- Enforcing foreign key constraint
    FOREIGN KEY (priorityid) REFERENCES PRIORITY(id) -- Enforcing foreign key constraint
);

-- Table to store tags for issues in projects
CREATE TABLE TAG (
    id INTEGER PRIMARY KEY,    -- Unique ID for each tag
    projectid INTEGER NOT NULL, -- Foreign key referencing the PROJECT table
    name TEXT NOT NULL,         -- Tag name
    color INTEGER NOT NULL,     -- Color for the tag stored as a hex value
    FOREIGN KEY (projectid) REFERENCES PROJECT(id) -- Enforcing foreign key constraint
);

-- Table to store priorities for issues in projects
CREATE TABLE PRIORITY (
    id INTEGER PRIMARY KEY,    -- Unique ID for each priority
    projectid INTEGER NOT NULL, -- Foreign key referencing the PROJECT table
    name TEXT NOT NULL,         -- Name of the priority (e.g., High, Low)
    priority INTEGER NOT NULL,  -- Priority value (1, 2, 3, etc.)
    FOREIGN KEY (projectid) REFERENCES PROJECT(id) -- Enforcing foreign key constraint
);

-- Table to store reminders for issues
CREATE TABLE REMINDER (
    id INTEGER PRIMARY KEY,    -- Unique ID for each reminder
    issueid INTEGER NOT NULL,  -- Foreign key referencing the ISSUE table
    userid INTEGER NOT NULL,   -- Foreign key referencing the USER table
    FOREIGN KEY (issueid) REFERENCES ISSUE(id), -- Enforcing foreign key constraint
    FOREIGN KEY (userid) REFERENCES USER(id)    -- Enforcing foreign key constraint
);

-- Relationship tables

-- Table to store verification codes for users
CREATE TABLE VERIFY_USER (
    id INTEGER PRIMARY KEY,    -- Unique ID for each verification record
    code INTEGER UNIQUE NOT NULL, -- Unique code for verification
    expires DATE NOT NULL      -- Expiry date for the verification code
);

-- Table to store members of organizations
CREATE TABLE ORG_MEMBER (
    id INTEGER PRIMARY KEY,    -- Unique ID for each organization member
    userid INTEGER NOT NULL,   -- Foreign key referencing the USER table
    orgid INTEGER NOT NULL,    -- Foreign key referencing the ORGANIZATION table
    FOREIGN KEY (userid) REFERENCES USER(id), -- Enforcing foreign key constraint
    FOREIGN KEY (orgid) REFERENCES ORGANIZATION(id) -- Enforcing foreign key constraint
);

-- Table to store the relationship between organization members and roles
CREATE TABLE ORG_MEMBER_ROLE (
    id INTEGER PRIMARY KEY,    -- Unique ID for each organization member's role
    memberid INTEGER NOT NULL, -- Foreign key referencing the ORG_MEMBER table
    roleid INTEGER NOT NULL,   -- Foreign key referencing the ORG_ROLE table
    FOREIGN KEY (memberid) REFERENCES ORG_MEMBER(id), -- Enforcing foreign key constraint
    FOREIGN KEY (roleid) REFERENCES ORG_ROLE(id)     -- Enforcing foreign key constraint
);

-- Table to store the relationship between organizations and their projects
CREATE TABLE ORG_PROJECTS (
    id INTEGER PRIMARY KEY,    -- Unique ID for each organization-project relationship
    orgid INTEGER NOT NULL,    -- Foreign key referencing the ORGANIZATION table
    projectid INTEGER NOT NULL, -- Foreign key referencing the PROJECT table
    FOREIGN KEY (orgid) REFERENCES ORGANIZATION(id), -- Enforcing foreign key constraint
    FOREIGN KEY (projectid) REFERENCES PROJECT(id) -- Enforcing foreign key constraint
);

-- Table to store members of projects
CREATE TABLE PROJECT_MEMBER (
    id INTEGER PRIMARY KEY,    -- Unique ID for each project member
    userid INTEGER NOT NULL,   -- Foreign key referencing the USER table
    projectid INTEGER NOT NULL, -- Foreign key referencing the PROJECT table
    FOREIGN KEY (userid) REFERENCES USER(id), -- Enforcing foreign key constraint
    FOREIGN KEY (projectid) REFERENCES PROJECT(id) -- Enforcing foreign key constraint
);

-- Table to store the roles of project members
CREATE TABLE PROJECT_MEMBER_ROLE (
    id INTEGER PRIMARY KEY,    -- Unique ID for each project member's role
    memberid INTEGER NOT NULL, -- Foreign key referencing the PROJECT_MEMBER table
    roleid INTEGER NOT NULL,   -- Foreign key referencing the PROJECT_ROLE table
    FOREIGN KEY (memberid) REFERENCES PROJECT_MEMBER(id), -- Enforcing foreign key constraint
    FOREIGN KEY (roleid) REFERENCES PROJECT_ROLE(id)      -- Enforcing foreign key constraint
);

-- Table to store the relationship between projects and their issues
CREATE TABLE PROJECT_ISSUES (
    id INTEGER PRIMARY KEY,    -- Unique ID for each project-issue relationship
    projectid INTEGER NOT NULL, -- Foreign key referencing the PROJECT table
    issueid INTEGER NOT NULL,   -- Foreign key referencing the ISSUE table
    FOREIGN KEY (projectid) REFERENCES PROJECT(id), -- Enforcing foreign key constraint
    FOREIGN KEY (issueid) REFERENCES ISSUE(id)     -- Enforcing foreign key constraint
);

-- Table to store the relationship between users and issues they are assigned to
CREATE TABLE USER_ISSUES (
    id INTEGER PRIMARY KEY,    -- Unique ID for each user-issue relationship
    userid INTEGER NOT NULL,   -- Foreign key referencing the USER table
    issueid INTEGER NOT NULL,  -- Foreign key referencing the ISSUE table
    FOREIGN KEY (userid) REFERENCES USER(id), -- Enforcing foreign key constraint
    FOREIGN KEY (issueid) REFERENCES ISSUE(id) -- Enforcing foreign key constraint
);

-- Table to store the relationship between users and their reminders
CREATE TABLE USER_REMINDER (
    id INTEGER PRIMARY KEY,    -- Unique ID for each user-reminder relationship
    userid INTEGER NOT NULL,   -- Foreign key referencing the USER table
    reminderid INTEGER NOT NULL, -- Foreign key referencing the REMINDER table
    FOREIGN KEY (userid) REFERENCES USER(id), -- Enforcing foreign key constraint
    FOREIGN KEY (reminderid) REFERENCES REMINDER(id) -- Enforcing foreign key constraint
);
