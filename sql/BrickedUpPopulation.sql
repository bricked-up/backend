--SQLite Database Population Script for BrickedUp Project Management System

-- Insert Organizations
INSERT INTO organization (name) VALUES 
('TechCorp'),
('Green Innovations'),
('NextGen Solutions');

-- Insert Users
INSERT INTO user (name) VALUES 
('Alice Johnson', 'password123'),
('Bob Smith', 'password456'),
('Charlie Brown', 'password789'),
('Diana Prince', 'wonderwoman'),
('Ethan Hunt', 'missionimpossible');

-- Insert Organization Roles
INSERT INTO organization_role (organization_id, name, can_read, can_write, can_execute) VALUES
(1, 'Admin', 1, 1, 1),
(1, 'Member', 1, 0, 0),
(2, 'Admin', 1, 1, 1),
(2, 'Member', 1, 0, 0),
(3, 'Admin', 1, 1, 1),
(3, 'Member', 1, 0, 0);

-- Insert Organization Members
INSERT INTO organization_member (user_id, organization_id) VALUES
(1, 1), (2, 1), (3, 2), (4, 2), (5, 3);

-- Insert Projects
INSERT INTO project (organization_id, name, budget, charter, is_archived) VALUES
(1, 'AI Research', 50000, 'Develop AI models for automation.', 0),
(2, 'Eco Farming', 30000, 'Smart farming solutions.', 0),
(3, 'Cyber Security', 70000, 'Enhancing system security.', 0);

-- Insert Project Roles
INSERT INTO project_role (project_id, name, can_read, can_write, can_execute) VALUES
(1, 'Lead', 1, 1, 1),
(1, 'Contributor', 1, 1, 0),
(2, 'Manager', 1, 1, 1),
(2, 'Worker', 1, 0, 0),
(3, 'Security Analyst', 1, 1, 1);

-- Insert Project Members
INSERT INTO project_member (user_id, project_id) VALUES
(1, 1), (2, 1), (3, 2), (4, 2), (5, 3);

-- Insert Tags
INSERT INTO tag (project_id, name, color) VALUES
(1, 'AI', 'blue'),
(2, 'Sustainability', 'green'),
(3, 'Security', 'red');

-- Insert Priorities
INSERT INTO priority (project_id, name, priority_level) VALUES
(1, 'High', 1),
(1, 'Medium', 2),
(2, 'Low', 3),
(3, 'Critical', 1);

-- Insert Issues
INSERT INTO issue (title, description, tag_id, priority_id, created_at, completed_at, cost) VALUES
('Train AI Model', 'Develop a machine learning model.', 1, 1, '2025-02-26', NULL, 10000),
('Optimize Farm System', 'Improve crop monitoring.', 2, 2, '2025-02-25', NULL, 5000),
('Fix Security Breach', 'Patch critical vulnerability.', 3, 4, '2025-02-20', '2025-02-22', 20000),
('Data Cleanup', 'Remove redundant data.', 1, 3, '2025-02-21', NULL, 2000),
('Enhance Firewall', 'Upgrade system defenses.', 3, 1, '2025-02-24', NULL, 15000);

-- Insert Project-Issue Associations
INSERT INTO project_issue (project_id, issue_id) VALUES
(1, 1), (2, 2), (3, 3), (1, 4), (3, 5);

-- Insert User-Issue Assignments
INSERT INTO user_issue (user_id, issue_id) VALUES
(1, 1), (2, 2), (3, 3), (4, 4), (5, 5);

-- Insert Reminders
INSERT INTO reminder (issue_id, user_id) VALUES
(1, 1), (2, 2), (3, 3), (4, 4), (5, 5);
