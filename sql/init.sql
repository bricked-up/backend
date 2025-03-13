PRAGMA foreign_keys = ON;


CREATE TABLE ORGANIZATION (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE VERIFY_USER (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code INTEGER UNIQUE NOT NULL,
    expires DATE NOT NULL
);

-- Tables that depend only on tables already created
CREATE TABLE USER (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    verifyid INTEGER,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    name TEXT NOT NULL,
    avatar TEXT,
    FOREIGN KEY (verifyid) REFERENCES VERIFY_USER(id) ON DELETE SET NULL
);

CREATE TABLE SESSION (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userid INTEGER NOT NULL,
    timestamp DATE NOT NULL,
    FOREIGN KEY (userid) REFERENCES USER(id) ON DELETE CASCADE
);

CREATE TABLE ORG_ROLE (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    orgid INTEGER NOT NULL,
    name TEXT NOT NULL,
    can_read BOOLEAN NOT NULL,
    can_write BOOLEAN NOT NULL,
    can_exec BOOLEAN NOT NULL,
    FOREIGN KEY (orgid) REFERENCES ORGANIZATION(id) ON DELETE CASCADE
);

CREATE TABLE ORG_MEMBER (
    id INTEGER PRIMARY KEY,
    userid INTEGER NOT NULL,
    orgid INTEGER NOT NULL,
    FOREIGN KEY (userid) REFERENCES USER(id) ON DELETE CASCADE,
    FOREIGN KEY (orgid) REFERENCES ORGANIZATION(id) ON DELETE CASCADE
);

CREATE TABLE PROJECT (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    orgid INTEGER NOT NULL,
    name TEXT NOT NULL,
    budget INTEGER NOT NULL,
    charter TEXT NOT NULL,
    archived BOOLEAN NOT NULL,
    FOREIGN KEY (orgid) REFERENCES ORGANIZATION(id) ON DELETE CASCADE
);

-- Tables dependent on PROJECT
CREATE TABLE PROJECT_ROLE (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    projectid INTEGER NOT NULL,
    name TEXT NOT NULL,
    can_read BOOLEAN NOT NULL,
    can_write BOOLEAN NOT NULL,
    can_exec BOOLEAN NOT NULL,
    FOREIGN KEY (projectid) REFERENCES PROJECT(id) ON DELETE CASCADE
);

CREATE TABLE TAG (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    projectid INTEGER NOT NULL,
    name TEXT NOT NULL,
    color TEXT NOT NULL,
    FOREIGN KEY (projectid) REFERENCES PROJECT(id) ON DELETE CASCADE
);

CREATE TABLE PRIORITY (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    projectid INTEGER NOT NULL,
    name TEXT NOT NULL,
    priority INTEGER NOT NULL,
    FOREIGN KEY (projectid) REFERENCES PROJECT(id) ON DELETE CASCADE
);

CREATE TABLE ISSUE (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    desc TEXT NOT NULL,
    tagid INTEGER,
    priorityid INTEGER,
    created DATE NOT NULL,
    completed DATE,
    cost INTEGER NOT NULL,
    FOREIGN KEY (tagid) REFERENCES TAG(id) ON DELETE SET NULL,
    FOREIGN KEY (priorityid) REFERENCES PRIORITY(id) ON DELETE SET NULL
);


CREATE TABLE DEPENDENCY (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    issueid INTEGER NOT NULL,
    dependency INTEGER NOT NULL,
    FOREIGN KEY (issueid) REFERENCES ISSUE(id) ON DELETE CASCADE,
    FOREIGN KEY (dependency) REFERENCES ISSUE(id) ON DELETE CASCADE
);

CREATE TABLE REMINDER (
    id INTEGER PRIMARY KEY,
    issueid INTEGER NOT NULL,
    userid INTEGER NOT NULL,
    FOREIGN KEY (issueid) REFERENCES ISSUE(id) ON DELETE CASCADE,
    FOREIGN KEY (userid) REFERENCES USER(id) ON DELETE CASCADE
);

CREATE TABLE ORG_MEMBER_ROLE (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberid INTEGER NOT NULL,
    roleid INTEGER NOT NULL,
    FOREIGN KEY (memberid) REFERENCES ORG_MEMBER(id) ON DELETE CASCADE,
    FOREIGN KEY (roleid) REFERENCES ORG_ROLE(id) ON DELETE CASCADE
);

CREATE TABLE ORG_PROJECTS (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    orgid INTEGER NOT NULL,
    projectid INTEGER NOT NULL,
    FOREIGN KEY (orgid) REFERENCES ORGANIZATION(id) ON DELETE CASCADE,
    FOREIGN KEY (projectid) REFERENCES PROJECT(id) ON DELETE CASCADE
);

CREATE TABLE PROJECT_MEMBER (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userid INTEGER NOT NULL,
    projectid INTEGER NOT NULL,
    FOREIGN KEY (userid) REFERENCES USER(id) ON DELETE CASCADE,
    FOREIGN KEY (projectid) REFERENCES PROJECT(id) ON DELETE CASCADE
);

CREATE TABLE PROJECT_MEMBER_ROLE (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberid INTEGER NOT NULL,
    roleid INTEGER NOT NULL,
    FOREIGN KEY (memberid) REFERENCES PROJECT_MEMBER(id) ON DELETE CASCADE,
    FOREIGN KEY (roleid) REFERENCES PROJECT_ROLE(id) ON DELETE CASCADE
);

CREATE TABLE PROJECT_ISSUES (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    projectid INTEGER NOT NULL,
    issueid INTEGER NOT NULL,
    FOREIGN KEY (projectid) REFERENCES PROJECT(id) ON DELETE CASCADE,
    FOREIGN KEY (issueid) REFERENCES ISSUE(id) ON DELETE CASCADE
);

CREATE TABLE USER_ISSUES (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userid INTEGER NOT NULL,
    issueid INTEGER NOT NULL,
    FOREIGN KEY (userid) REFERENCES USER(id) ON DELETE CASCADE,
    FOREIGN KEY (issueid) REFERENCES ISSUE(id) ON DELETE CASCADE,
    UNIQUE (userid, issueid)
);

CREATE TABLE FORGOT_PASSWORD (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userid INTEGER NOT NULL,
    code INTEGER NOT NULL,
    expirationdate DATE NOT NULL,
    FOREIGN KEY (userid) REFERENCES USER(id) ON DELETE CASCADE
);


-- DO NOT REMOVE! I NEED IT FOR TESTING  
----------------------------------------------------------------------

-- Drop user related junction tables
-- DROP TABLE IF EXISTS FORGOT_PASSWORD;
-- DROP TABLE IF EXISTS USER_ISSUES;

-- -- Drop project related junction tables
-- DROP TABLE IF EXISTS PROJECT_ISSUES;
-- DROP TABLE IF EXISTS PROJECT_MEMBER_ROLE;
-- DROP TABLE IF EXISTS PROJECT_MEMBER;
-- DROP TABLE IF EXISTS ORG_PROJECTS;
-- DROP TABLE IF EXISTS ORG_MEMBER_ROLE;

-- -- Drop issue related tables
-- DROP TABLE IF EXISTS REMINDER;
-- DROP TABLE IF EXISTS DEPENDENCY;

-- -- Drop the main ISSUE table
-- DROP TABLE IF EXISTS ISSUE;

-- -- Drop project related tables
-- DROP TABLE IF EXISTS PRIORITY;
-- DROP TABLE IF EXISTS TAG;
-- DROP TABLE IF EXISTS PROJECT_ROLE;

-- -- Drop organization and user related tables
-- DROP TABLE IF EXISTS PROJECT;
-- DROP TABLE IF EXISTS ORG_MEMBER;
-- DROP TABLE IF EXISTS ORG_ROLE;
-- DROP TABLE IF EXISTS SESSION;
-- DROP TABLE IF EXISTS USER;

-- -- Finally drop the base tables with no dependencies
-- DROP TABLE IF EXISTS VERIFY_USER;
-- DROP TABLE IF EXISTS ORGANIZATION;