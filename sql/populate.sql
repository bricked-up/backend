-- Enable foreign key constraints
PRAGMA foreign_keys = ON;

-- 1. Insert into ORGANIZATION
INSERT INTO ORGANIZATION (name) 
VALUES ('Tech Corp'), ('Dev Hub');

-- 2. Insert into VERIFY_USER
INSERT INTO VERIFY_USER (code, expires) 
VALUES (123456, '2025-12-31');

-- 3. Insert into USER (verifyid references VERIFY_USER)
INSERT INTO USER (verifyid, email, password, name, avatar) 
VALUES (1, 'user1@example.com', 'password123', 'Alice', 'avatar1.png'), 
       (NULL, 'user2@example.com', 'password456', 'Bob', 'avatar2.png');

-- 4. Insert into SESSION (userid references USER)
INSERT INTO SESSION (userid, expires) 
VALUES (1, '2025-06-30');

-- 5. Insert into ORG_ROLE (orgid references ORGANIZATION)
INSERT INTO ORG_ROLE (orgid, name, can_read, can_write, can_exec) 
VALUES (1, 'Admin', 1, 1, 1), 
       (1, 'Editor', 1, 1, 0), 
       (2, 'Viewer', 1, 0, 0);

-- 6. Insert into ORG_MEMBER (userid, orgid must exist)
INSERT INTO ORG_MEMBER (userid, orgid) 
VALUES (1, 1), (2, 2);

-- 7. Insert into ORG_MEMBER_ROLE (memberid references ORG_MEMBER, roleid references ORG_ROLE)
INSERT INTO ORG_MEMBER_ROLE (memberid, roleid) 
VALUES (1, 1), (2, 3);

-- 8. Insert into PROJECT (orgid must exist)
INSERT INTO PROJECT (orgid, name, budget, charter, archived) 
VALUES (1, 'AI Research', 50000, 'AI development project', 0), 
       (2, 'Web Platform', 30000, 'Web service development', 1);

-- 9. Insert into PROJECT_ROLE (projectid must exist)
INSERT INTO PROJECT_ROLE (projectid, name, can_read, can_write, can_exec) 
VALUES (1, 'Lead', 1, 1, 1), 
       (2, 'Contributor', 1, 1, 0);

-- 10. Insert into PROJECT_MEMBER (userid, projectid must exist)
INSERT INTO PROJECT_MEMBER (userid, projectid) 
VALUES (1, 1), (2, 2);

-- 11. Insert into PROJECT_MEMBER_ROLE (memberid, roleid must exist)
INSERT INTO PROJECT_MEMBER_ROLE (memberid, roleid) 
VALUES (1, 1), (2, 2);

-- 12. Insert into TAG (projectid must exist)
INSERT INTO TAG (projectid, name, color) 
VALUES (1, 'Bug', '#FF0000'), 
       (2, 'Feature', '#00FF00');

-- 13. Insert into PRIORITY (projectid must exist)
INSERT INTO PRIORITY (projectid, name, priority) 
VALUES (1, 'High', 1), 
       (2, 'Medium', 2);

-- 14. Insert into ISSUE (tagid, priorityid must exist)
INSERT INTO ISSUE (title, desc, tagid, priorityid, created, completed, cost) 
VALUES ('Fix Login Bug', 'Users unable to log in', 1, 1, '2025-02-01', NULL, 200), 
       ('Add Dark Mode', 'New UI feature', 2, 2, '2025-02-10', NULL, 500);

-- 15. Insert into PROJECT_ISSUES (projectid, issueid must exist)
INSERT INTO PROJECT_ISSUES (projectid, issueid) 
VALUES (1, 1), (2, 2);

-- 16. Insert into USER_ISSUES (userid, issueid must exist)
INSERT INTO USER_ISSUES (userid, issueid) 
VALUES (1, 1), (2, 2);

-- 17. Insert into REMINDER (userid, issueid must exist)
INSERT INTO REMINDER (issueid, userid) 
VALUES (1, 1), (2, 2);

-- 18. Insert into ORG_PROJECTS (orgid, projectid must exist)
INSERT INTO ORG_PROJECTS (orgid, projectid) 
VALUES (1, 1), (2, 2);
