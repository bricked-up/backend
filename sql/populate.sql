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
INSERT INTO USER (verifyid, email, password, name, avatar, verified) VALUES
(1, 'john.doe@example.com', '$2a$10$N7GjZkymUNWzmWpsLvTMpeUAEf8fZ4AzWP1KOIWY4IaxBeXfOpYmm', 'John Doe', 'avatar1.png', 1),
(2, 'jane.smith@example.com', '$2a$10$GZTCATobqiHPNknznX8CbuiFW/4Lr91rAfW/DaFmbviXuevBLDoGu', 'Jane Smith', 'avatar2.png', 1),
(NULL, 'mike.johnson@example.com', '$2a$10$wlsTt32wUFeEl89imMLxPeiGAnnVSbw1eVFaVC/jyviLB4nZxt4.K', 'Mike Johnson', 'avatar3.png', 1),
(NULL, 'sarah.williams@example.com', '$2a$10$6lnVZ6Po41a8WTn9qJyfMeiNhiZkjx/A4cYIRR7dxICgl7pk8vCra', 'Sarah Williams', 'avatar4.png', 0),
(NULL, 'alex.brown@example.com', '$2a$10$vsQ0I0vp7bINyqe77WaqcOlB2vUXgZaC4JhNr1.6sb36N8xekHuqO', 'Alex Brown', 'avatar5.png', 0);

-- Populate SESSION table
INSERT INTO SESSION (userid, expires) VALUES
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

-- Populate ISSUE table (with tagid from TAG table)
INSERT INTO ISSUE (title, desc, tagid, created, cost, priority) VALUES
('Setup Development Environment', 'Install and configure all necessary tools', 1, '2023-01-01 10:00:00', 500, 1),
('Design Database Schema', 'Create ERD and implement tables', 3, '2023-01-02 09:00:00', 1000, 2),
('Implement User Authentication', 'Add login and registration system', 2, '2023-01-03 14:00:00', 1500, 1),
('Create API Documentation', 'Document all endpoints and parameters', 3, '2023-01-04 11:00:00', 800, 3),
('Bug Fix: Login Page', 'Fix validation errors on login form', 2, '2023-01-05 16:00:00', 300, 2);

-- Populate DEPENDENCY table (updated to match actual issue IDs)
INSERT INTO DEPENDENCY (issueid, dependency) VALUES
(3, 1),
(4, 2),
(5, 2);

-- Populate REMINDER table (updated to match actual issue and user IDs)
INSERT INTO REMINDER (issueid, userid) VALUES
(1, 1),
(3, 2);

-- Rest of the script remains the same as in the original populate script...

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
(1, 4),
(1, 5);

-- Populate USER_ISSUES table
INSERT INTO USER_ISSUES (userid, issueid) VALUES
(1, 1),
(2, 2),
(2, 3),
(3, 4),
(3, 5);

-- Populate FORGOT_PASSWORD table
INSERT INTO FORGOT_PASSWORD (userid, code, expirationdate) VALUES
(1, 123456, '2025-03-12'),
(2, 234567, '2025-03-13'),
(4, 345678, '2025-03-14'),
(5, 456789, '2025-03-15');