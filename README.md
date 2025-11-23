# AgenticAI Orchestrator

This repository own the code for the AgenticAI Orchestrator. The main responsibility of the orchestrator is to respond to the requests that come from connectors (see SlackConnector, etc...) using AI agents that the orchestrator manages.

## How to run the orchestrator

You will need to have `docker` and `docker-compose` installed in your machine. Then, you can run the following command:

```bash
docker-compose up
```

⚠ One of the containers can fail the first time you run it due to concurrent access to the database to create migrations. If this happens, just run the command again.

This command will mount a database, a pgadmin (to visualize the database), a hashicorp vault, a local blockchain and the APIs.

- You can access the core API at `http://localhost:5000/`
- You can access the blockchain at `http://localhost:8545/`
- You can access the vault at `http://localhost:8200/` (use `wBGKZCJjfIGqIKB3HMmvcQVt` as passwords)
- You can access the pgadmin at `http://localhost:5050/` (use `root` as passwords)

#### Command to create a agent postgresql database

```sql
-- Connect to the default database
-- Create the database
CREATE DATABASE "<AGENT_ID>"; -- Run this command separately to avoid transaction error

-- Create the role with login privilege
CREATE ROLE "<AGENT_ID>" WITH LOGIN PASSWORD '<ROLE_PASSWORD>';

-- Grant all privileges on the database to the role
GRANT ALL PRIVILEGES ON DATABASE "<AGENT_ID>" TO "<AGENT_ID>";

-- Then connect to the created database and run the following command
GRANT ALL PRIVILEGES ON SCHEMA public TO "<AGENT_ID>";
```
