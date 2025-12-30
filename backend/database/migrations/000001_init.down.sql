-- Drop tables in reverse order of creation (respecting foreign key constraints)
DROP TABLE IF EXISTS health_checks;
DROP TABLE IF EXISTS deployments;
DROP TABLE IF EXISTS releases;
DROP TABLE IF EXISTS environments;

-- Drop custom enum types
DROP TYPE IF EXISTS health_status;
DROP TYPE IF EXISTS health_check_type;
DROP TYPE IF EXISTS trigger_type;
DROP TYPE IF EXISTS deployment_status;

-- Drop UUID extension (optional, can be kept for other tables)
-- DROP EXTENSION IF EXISTS "uuid-ossp";
