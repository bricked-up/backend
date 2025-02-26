-- database: :memory:

-- 1. Get All Organizations and Their Roles
SELECT 
    o.name AS organization_name, 
    orl.name AS role_name, 
    orl.can_read, 
    orl.can_write, 
    orl.can_execute
FROM 
    organization o
JOIN 
    organization_role orl ON o.id = orl.organization_id;

-- 2. Get All Users and Their Assigned Organizations
SELECT 
    u.name AS user_name, 
    o.name AS organization_name
FROM 
    user u
JOIN 
    organization_member om ON u.id = om.user_id
JOIN 
    organization o ON om.organization_id = o.id;

-- 3. Get All Projects and Their Associated Tags
SELECT 
    p.name AS project_name, 
    t.name AS tag_name
FROM 
    project p
JOIN 
    tag t ON p.id = t.project_id;

-- 4. Get All Projects and Their Members
SELECT 
    p.name AS project_name, 
    u.name AS user_name
FROM 
    project p
JOIN 
    project_member pm ON p.id = pm.project_id
JOIN 
    user u ON pm.user_id = u.id;

-- 5. Get All Issues with Their Associated Priorities
SELECT 
    i.title AS issue_title, 
    p.name AS priority_name, 
    p.priority_level
FROM 
    issue i
JOIN 
    priority p ON i.priority_id = p.id;

-- 6. Get All Issues Assigned to a Specific User (e.g., Alice)
SELECT 
    i.title AS issue_title, 
    i.description AS issue_description
FROM 
    user u
JOIN 
    user_issue ui ON u.id = ui.user_id
JOIN 
    issue i ON ui.issue_id = i.id
WHERE 
    u.name = 'Alice Johnson';

-- 7. Get All Issues in a Specific Project (e.g., AI Research)
SELECT 
    i.title AS issue_title, 
    i.description AS issue_description
FROM 
    project p
JOIN 
    project_issue pi ON p.id = pi.project_id
JOIN 
    issue i ON pi.issue_id = i.id
WHERE 
    p.name = 'AI Research';

-- 8. Get All Reminders and Their Associated Issues
SELECT 
    r.id AS reminder_id, 
    i.title AS issue_title, 
    u.name AS user_name
FROM 
    reminder r
JOIN 
    issue i ON r.issue_id = i.id
JOIN 
    user u ON r.user_id = u.id;

-- 9. Get All Active Projects (Not Archived)
SELECT 
    p.name AS project_name, 
    p.budget, 
    p.charter
FROM 
    project p
WHERE 
    p.is_archived = 0;

-- 10. Get All Issues with Their Associated Costs
SELECT 
    i.title AS issue_title, 
    i.cost
FROM 
    issue i;

-- 11. Get Users with Roles in Specific Project (e.g., AI Research)
SELECT 
    u.name AS user_name, 
    pr.name AS role_name
FROM 
    user u
JOIN 
    project_member pm ON u.id = pm.user_id
JOIN 
    project p ON pm.project_id = p.id
JOIN 
    project_role pr ON p.id = pr.project_id
WHERE 
    p.name = 'AI Research';

-- 12. Get Total Budget of Projects in Each Organization
SELECT 
    o.name AS organization_name, 
    SUM(p.budget) AS total_budget
FROM 
    organization o
JOIN 
    project p ON o.id = p.organization_id
GROUP BY 
    o.id;
