-- Populate ORGANIZATION table
INSERT INTO ORGANIZATION (name) VALUES
('TechCorp Solutions'),
('Creative Designs Inc'),
('Data Innovations LLC');

-- Populate VERIFY_USER table
INSERT INTO VERIFY_USER (code, expires) VALUES
(123456, '2025-04-11'),
(234567, '2025-04-12'),
(345678, '2025-04-13');

-- Populate USER table
INSERT INTO USER (verifyid, email, password, name, avatar) VALUES
(1, 'john.doe@example.com', 'hashed_password_1', 'John Doe', 'avatar1.png'),
(2, 'jane.smith@example.com', 'hashed_password_2', 'Jane Smith', 'avatar2.png'),
(NULL, 'mike.johnson@example.com', 'hashed_password_3', 'Mike Johnson', 'avatar3.png'),
(NULL, 'sarah.williams@example.com', 'hashed_password_4', 'Sarah Williams', 'avatar4.png'),
(NULL, 'alex.brown@example.com', 'hashed_password_5', 'Alex Brown', 'avatar5.png');

-- Populate SESSION table
INSERT INTO SESSION (userid, timestamp) VALUES
(1, '2025-03-10 09:30:00'),
(2, '2025-03-10 10:15:00'),
(3, '2025-03-10 14:22:00'),
(1, '2025-03-11 08:45:00'),
(4, '2025-03-11 11:10:00');

-- Populate ORG_ROLE table
INSERT INTO ORG_ROLE (orgid, name, can_read, can_write, can_exec) VALUES
(1, 'Admin', 1, 1, 1),
(1, 'Developer', 1, 1, 0),
(1, 'Viewer', 1, 0, 0),
(2, 'Admin', 1, 1, 1),
(2, 'Designer', 1, 1, 0),
(3, 'Admin', 1, 1, 1),
(3, 'Analyst', 1, 1, 0);

-- Populate ORG_MEMBER table
INSERT INTO ORG_MEMBER (userid, orgid) VALUES
(1, 1),
(2, 1),
(3, 1),
(2, 2),
(4, 2),
(5, 3),
(3, 3);

-- Populate PROJECT table
INSERT INTO PROJECT (orgid, name, budget, charter, archived) VALUES
(1, 'Web Platform Redesign', 50000, 'Modernize our web platform with updated UX/UI', 0),
(1, 'Mobile App Development', 75000, 'Create native mobile applications for iOS and Android', 0),
(1, 'Legacy System Migration', 120000, 'Migrate legacy systems to cloud infrastructure', 0),
(2, 'Brand Identity Refresh', 35000, 'Update company brand identity and style guides', 0),
(2, 'Marketing Campaign Q2', 25000, 'Design assets for Q2 marketing campaign', 1),
(3, 'Data Warehouse Implementation', 90000, 'Implement enterprise data warehouse solution', 0);

-- Populate PROJECT_ROLE table
INSERT INTO PROJECT_ROLE (projectid, name, can_read, can_write, can_exec) VALUES
(1, 'Project Manager', 1, 1, 1),
(1, 'Developer', 1, 1, 0),
(1, 'QA Tester', 1, 1, 0),
(1, 'Stakeholder', 1, 0, 0),
(2, 'Project Manager', 1, 1, 1),
(2, 'Mobile Developer', 1, 1, 0),
(3, 'Migration Specialist', 1, 1, 1),
(4, 'Design Lead', 1, 1, 1),
(5, 'Campaign Manager', 1, 1, 1),
(6, 'Data Engineer', 1, 1, 1);

-- Populate TAG table
INSERT INTO TAG (projectid, name, color) VALUES
(1, 'Frontend', '#4287f5'),
(1, 'Backend', '#42f59e'),
(1, 'Database', '#f54242'),
(2, 'iOS', '#f5a442'),
(2, 'Android', '#42f5e6'),
(3, 'Migration', '#d142f5'),
(3, 'DevOps', '#9ef542'),
(4, 'Design', '#f542d4'),
(6, 'Data', '#426ff5');

-- Populate PRIORITY table
INSERT INTO PRIORITY (projectid, name, priority) VALUES
(1, 'Critical', 1),
(1, 'High', 2),
(1, 'Medium', 3),
(1, 'Low', 4),
(2, 'Critical', 1),
(2, 'High', 2),
(2, 'Medium', 3),
(3, 'Urgent', 1),
(3, 'Important', 2),
(3, 'Normal', 3);

-- Populate ISSUE table
INSERT INTO ISSUE (title, desc, tagid, priorityid, created, completed, cost) VALUES
('Implement user authentication', 'Create secure authentication system with JWT', 2, 1, '2025-02-10', NULL, 8000),
('Design responsive UI', 'Create responsive UI mockups for all screen sizes', 1, 3, '2025-02-12', '2025-03-01', 5000),
('Set up database schema', 'Create initial database schema for user management', 3, 2, '2025-02-15', NULL, 4000),
('iOS app navigation', 'Implement navigation system for iOS app', 4, 5, '2025-02-20', NULL, 6000),
('Android performance optimizations', 'Optimize app performance for low-end Android devices', 5, 6, '2025-02-22', NULL, 7000),
('Server migration plan', 'Create detailed migration plan for server infrastructure', 6, 8, '2025-02-25', '2025-03-05', 10000),
('Configure CI/CD pipeline', 'Set up automated CI/CD pipeline for deployment', 7, 9, '2025-03-01', NULL, 8000),
('Brand color palette', 'Finalize brand color palette for refresh', 8, NULL, '2025-02-15', '2025-02-28', 3000),
('Data warehouse architecture', 'Design data warehouse architecture', 9, NULL, '2025-03-01', NULL, 12000);

-- Populate DEPENDENCY table
INSERT INTO DEPENDENCY (issueid, dependency) VALUES
(3, 1),
(4, 2),
(5, 2),
(7, 6);

-- Populate REMINDER table
INSERT INTO REMINDER (id, issueid, userid) VALUES
(1, 1, 1),
(2, 3, 2),
(3, 6, 3),
(4, 8, 4),
(5, 9, 5);

-- Populate ORG_MEMBER_ROLE table
INSERT INTO ORG_MEMBER_ROLE (memberid, roleid) VALUES
(1, 1),
(2, 2),
(3, 3),
(4, 4),
(5, 6),
(6, 7),
(7, 7);

-- Populate ORG_PROJECTS table
INSERT INTO ORG_PROJECTS (orgid, projectid) VALUES
(1, 1),
(1, 2),
(1, 3),
(2, 4),
(2, 5),
(3, 6);

-- Populate PROJECT_MEMBER table
INSERT INTO PROJECT_MEMBER (userid, projectid) VALUES
(1, 1),
(2, 1),
(3, 1),
(1, 2),
(2, 2),
(3, 3),
(4, 4),
(4, 5),
(5, 6),
(3, 6);

-- Populate PROJECT_MEMBER_ROLE table
INSERT INTO PROJECT_MEMBER_ROLE (memberid, roleid) VALUES
(1, 1),
(2, 2),
(3, 3),
(4, 5),
(5, 6),
(6, 7),
(7, 8),
(8, 9),
(9, 10),
(10, 10);

-- Populate PROJECT_ISSUES table
INSERT INTO PROJECT_ISSUES (projectid, issueid) VALUES
(1, 1),
(1, 2),
(1, 3),
(2, 4),
(2, 5),
(3, 6),
(3, 7),
(4, 8),
(6, 9);

-- Populate USER_ISSUES table
INSERT INTO USER_ISSUES (userid, issueid) VALUES
(1, 1),
(2, 2),
(2, 3),
(3, 4),
(3, 5),
(1, 6),
(2, 7),
(4, 8),
(5, 9);

-- Populate FORGOT_PASSWORD table
INSERT INTO FORGOT_PASSWORD (userid, code, expirationdate) VALUES
(1, 123789, '2025-03-12'),
(3, 456123, '2025-03-13');