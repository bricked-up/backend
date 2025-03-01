--SQLite Database Population Script for BrickedUp Project Management System

-- turn on foreign key support
PRAGMA foreign_keys = ON;

-- Inserting data into ORGANIZATION table
INSERT INTO ORGANIZATION (name) VALUES 
('Tech Innovators'),
('Health Solutions'),
('Creative Studios');

-- Inserting data into VERIFY_USER table
INSERT INTO VERIFY_USER (code, expires) VALUES
(123456, '2025-12-31'),
(654321, '2025-12-31'),
(111222, '2025-12-31');

-- Inserting data into USER table
INSERT INTO USER (verifyid, email, password, name, avatar) VALUES 
(1, 'alice@techinnovators.com', 'password123', 'Alice', 'alice_avatar.png'),
(2, 'bob@healthsolutions.com', 'password123', 'Bob', 'bob_avatar.png'),
(3, 'carol@creativestudios.com', 'password123', 'Carol', 'carol_avatar.png');

-- Inserting data into SESSION table
INSERT INTO SESSION (userid, expires) VALUES
(1, '2025-12-31'),
(2, '2025-12-31'),
(3, '2025-12-31');

-- Inserting data into PROJECT table
INSERT INTO PROJECT (orgid, name, budget, charter, archived) VALUES 
(1, 'AI Development', 50000, 'Develop AI tools for businesses.', 0),
(2, 'Medical Research', 100000, 'Research new medical treatments.', 0),
(3, 'Website Redesign', 20000, 'Redesign company website for better UX.', 1);

-- Inserting data into PROJECT_ROLE table
INSERT INTO PROJECT_ROLE (projectid, name, can_read, can_write, can_exec) VALUES
(1, 'Project Manager', 1, 1, 0),
(1, 'AI Developer', 1, 1, 1),
(2, 'Research Scientist', 1, 1, 0),
(3, 'Web Designer', 1, 1, 1);

-- Inserting data into TAG table
INSERT INTO TAG (projectid, name, color) VALUES
(1, 'AI', 0xFF5733),
(2, 'Medical', 0x33FF57),
(3, 'Web Design', 0x3357FF);

-- Inserting data into PRIORITY table
INSERT INTO PRIORITY (projectid, name, priority) VALUES
(1, 'High', 1),
(2, 'Medium', 2),
(3, 'Low', 3);

-- Inserting data into ISSUE table
INSERT INTO ISSUE (title, desc, tagid, priorityid, created, completed, cost) VALUES 
('AI Model Optimization', 'Optimize the AI model for better performance.', 1, 1, '2025-02-01', NULL, 1000),
('Clinical Trial Testing', 'Conduct clinical trials for new drug.', 2, 2, '2025-02-05', NULL, 5000),
('Homepage Redesign', 'Redesign the homepage for better user engagement.', 3, 3, '2025-02-10', NULL, 2000);

-- Inserting data into REMINDER table
INSERT INTO REMINDER (issueid, userid) VALUES
(1, 1),
(2, 2),
(3, 3);

-- Inserting data into ORG_MEMBER table
INSERT INTO ORG_MEMBER (userid, orgid) VALUES 
(1, 1),
(2, 2),
(3, 3);

-- Inserting data into ORG_MEMBER_ROLE table
INSERT INTO ORG_MEMBER_ROLE (memberid, roleid) VALUES
(1, 1),  -- Alice is a member of organization 'Tech Innovators' with a role
(2, 2),  -- Bob is a member of 'Health Solutions' with a role
(3, 3);  -- Carol is a member of 'Creative Studios' with a role

-- Inserting data into ORG_PROJECTS table
INSERT INTO ORG_PROJECTS (orgid, projectid) VALUES
(1, 1),  -- 'Tech Innovators' has the 'AI Development' project
(2, 2),  -- 'Health Solutions' has the 'Medical Research' project
(3, 3);  -- 'Creative Studios' has the 'Website Redesign' project

-- Inserting data into PROJECT_MEMBER table
INSERT INTO PROJECT_MEMBER (userid, projectid) VALUES 
(1, 1),
(2, 2),
(3, 3);

-- Inserting data into PROJECT_MEMBER_ROLE table
INSERT INTO PROJECT_MEMBER_ROLE (memberid, roleid) VALUES 
(1, 1),  -- Alice is 'Project Manager' of 'AI Development'
(2, 2),  -- Bob is 'Research Scientist' of 'Medical Research'
(3, 3);  -- Carol is 'Web Designer' of 'Website Redesign'

-- Inserting data into PROJECT_ISSUES table
INSERT INTO PROJECT_ISSUES (projectid, issueid) VALUES 
(1, 1),  -- 'AI Development' has 'AI Model Optimization' issue
(2, 2),  -- 'Medical Research' has 'Clinical Trial Testing' issue
(3, 3);  -- 'Website Redesign' has 'Homepage Redesign' issue

-- Inserting data into USER_ISSUES table
INSERT INTO USER_ISSUES (userid, issueid) VALUES
(1, 1),  -- Alice is responsible for 'AI Model Optimization'
(2, 2),  -- Bob is responsible for 'Clinical Trial Testing'
(3, 3);  -- Carol is responsible for 'Homepage Redesign'

-- Inserting data into USER_REMINDER table
INSERT INTO USER_REMINDER (userid, reminderid) VALUES
(1, 1),  -- Alice gets a reminder for 'AI Model Optimization'
(2, 2),  -- Bob gets a reminder for 'Clinical Trial Testing'
(3, 3);  -- Carol gets a reminder for 'Homepage Redesign'
