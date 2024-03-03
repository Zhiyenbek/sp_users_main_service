-- Creating tables 
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS users ( 
    id SERIAL PRIMARY KEY, 
    public_id UUID UNIQUE DEFAULT uuid_generate_v4() NOT NULL, 
    first_name TEXT NOT NULL, 
    last_name TEXT
); 
 
CREATE TABLE IF NOT EXISTS candidates ( 
    id SERIAL PRIMARY KEY, 
    public_id UUID UNIQUE NOT NULL, 
    current_position TEXT, 
    education TEXT,
    resume TEXT, 
    bio TEXT 
); 
 
CREATE TABLE IF NOT EXISTS recruiters ( 
    id SERIAL PRIMARY KEY, 
    public_id UUID UNIQUE NOT NULL, 
    company_public_id UUID NOT NULL
); 
 
CREATE TABLE IF NOT EXISTS companies ( 
    id SERIAL PRIMARY KEY, 
    public_id UUID UNIQUE DEFAULT uuid_generate_v4() NOT NULL,
    name TEXT,
    description TEXT
); 
 
CREATE TABLE IF NOT EXISTS positions ( 
    id SERIAL PRIMARY KEY, 
    public_id UUID UNIQUE DEFAULT uuid_generate_v4() NOT NULL,
    name TEXT, 
    recruiters_public_id UUID UNIQUE 
); 
 
CREATE TABLE IF NOT EXISTS skills ( 
    id SERIAL PRIMARY KEY, 
    public_id UUID UNIQUE DEFAULT uuid_generate_v4() NOT NULL,
    name TEXT 
); 
 
CREATE TABLE IF NOT EXISTS areas ( 
    id SERIAL PRIMARY KEY, 
    public_id UUID UNIQUE DEFAULT uuid_generate_v4() NOT NULL,
    position_id INT, 
    name TEXT 
); 
 
CREATE TABLE IF NOT EXISTS interviews ( 
    id SERIAL PRIMARY KEY, 
    public_id UUID UNIQUE DEFAULT uuid_generate_v4() NOT NULL,
    results JSONB 
); 
 
CREATE TABLE IF NOT EXISTS videos ( 
    id SERIAL PRIMARY KEY, 
    public_id UUID UNIQUE DEFAULT uuid_generate_v4() NOT NULL,
    interviews_public_id UUID UNIQUE, 
    path TEXT 
); 
 
CREATE TABLE IF NOT EXISTS auth ( 
    id SERIAL PRIMARY KEY, 
    user_id INT UNIQUE, 
    login TEXT UNIQUE, 
    password TEXT 
); 
 
CREATE TABLE IF NOT EXISTS positions_skills ( 
    position_id INT, 
    skill_id INT, 
    PRIMARY KEY (position_id, skill_id) 
); 
 
CREATE TABLE IF NOT EXISTS candidate_skills ( 
    candidate_id INT, 
    skill_id INT, 
    PRIMARY KEY (candidate_id, skill_id) 
); 
 
CREATE TABLE IF NOT EXISTS user_interviews ( 
    candidate_id INT, 
    position_id INT, 
    interview_id INT, 
    PRIMARY KEY (candidate_id, position_id, interview_id) 
); 
 
-- Creating references 
ALTER TABLE recruiters ADD CONSTRAINT fk_users_candidates FOREIGN KEY (public_id) REFERENCES users(public_id); 
ALTER TABLE candidates ADD CONSTRAINT fk_users_recruiters FOREIGN KEY (public_id) REFERENCES users(public_id); 
ALTER TABLE user_interviews ADD CONSTRAINT fk_user_interviews_candidates FOREIGN KEY (candidate_id) REFERENCES candidates(id); 
ALTER TABLE user_interviews ADD CONSTRAINT fk_user_interviews_positions FOREIGN KEY (position_id) REFERENCES positions(id); 
ALTER TABLE user_interviews ADD CONSTRAINT fk_user_interviews_interviews FOREIGN KEY (interview_id) REFERENCES interviews(id); 
ALTER TABLE auth ADD CONSTRAINT fk_auth_users FOREIGN KEY (user_id) REFERENCES users(id); 
ALTER TABLE positions_skills ADD CONSTRAINT fk_positions_skills_skills FOREIGN KEY (skill_id) REFERENCES skills(id); 
ALTER TABLE positions_skills ADD CONSTRAINT fk_positions_skills_positions FOREIGN KEY (position_id) REFERENCES positions(id); 
ALTER TABLE candidate_skills ADD CONSTRAINT fk_candidate_skills_skills FOREIGN KEY (skill_id) REFERENCES skills(id); 
ALTER TABLE candidate_skills ADD CONSTRAINT fk_candidate_skills_candidates FOREIGN KEY (candidate_id) REFERENCES candidates(id); 
ALTER TABLE positions ADD CONSTRAINT fk_recruiters_positions FOREIGN KEY (recruiters_public_id) REFERENCES recruiters(public_id); 
ALTER TABLE videos ADD CONSTRAINT fk_interviews_videos FOREIGN KEY (interviews_public_id) REFERENCES interviews(public_id);




INSERT INTO users (first_name, last_name)
VALUES
    ('John', 'Doe'),
    ('Jane', 'Smith'),
    ('Michael', 'Johnson'),
    ('Emily', 'Williams'),
    ('David', 'Brown'),
    ('Olivia', 'Jones'),
    ('Daniel', 'Miller'),
    ('Sophia', 'Taylor'),
    ('Matthew', 'Anderson'),
    ('Ava', 'Thomas');


INSERT INTO candidates (public_id, current_position, resume, bio)
SELECT public_id, 'Software Engineer', 'John Doe Resume', 'John Doe Bio'
FROM users
WHERE id <= 5;

INSERT INTO companies (name, description)
VALUES
    ('Company A', 'A technology company that specializes in software development.'),
    ('Company B', 'A global retail company with a focus on e-commerce.'),
    ('Company C', 'A financial services company providing investment and banking solutions.');

INSERT INTO recruiters (public_id, company_public_id)
SELECT public_id, (SELECT public_id FROM companies WHERE name = 'Company A')
FROM users
WHERE id > 28;




INSERT INTO positions (public_id, name, recruiters_public_id)
SELECT public_id, 'Software Engineer', (SELECT public_id FROM recruiters WHERE id = 1)
FROM candidates;


INSERT INTO skills (name)
VALUES
    ('Java'),
    ('Python'),
    ('JavaScript'),
    ('SQL'),
    ('HTML'),
    ('CSS'),
    ('React'),
    ('Node.js'),
    ('AWS'),
    ('Agile Methodology');

INSERT INTO areas (position_id, name)
SELECT id, 'Area ' || id
FROM positions;


INSERT INTO interviews (public_id, results)
SELECT public_id, '{"result": "Pass"}'
FROM candidates;


INSERT INTO videos (public_id, interviews_public_id, path)
SELECT public_id, (SELECT public_id FROM interviews WHERE id = 1), '/path/to/video'
FROM candidates;


INSERT INTO auth (user_id, login, password)
SELECT id, 'user' || id, 'password' || id
FROM users;


INSERT INTO positions_skills (position_id, skill_id)
SELECT position_id, skill_id
FROM (
    SELECT p.id AS position_id, s.id AS skill_id, ROW_NUMBER() OVER () AS rn
    FROM positions p
    CROSS JOIN skills s
) AS sub
WHERE sub.rn <= 5;


INSERT INTO candidate_skills (candidate_id, skill_id)
SELECT candidate_id, skill_id
FROM (
    SELECT c.id AS candidate_id, s.id AS skill_id, ROW_NUMBER() OVER () AS rn
    FROM candidates c
    CROSS JOIN skills s
) AS sub
WHERE sub.rn <= 5;

insert into candidate_skills values (1,2);
insert into candidate_skills values (2,2);
insert into candidate_skills values (1,3);
insert into candidate_skills values (3,4);

INSERT INTO user_interviews (candidate_id, position_id, interview_id)
SELECT c.id AS candidate_id, p.id AS position_id, i.id AS interview_id
FROM candidates c
CROSS JOIN positions p
CROSS JOIN interviews i
WHERE c.id <= 5;