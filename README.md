# The Nola Connect

This repository contains the necessary components to support our client the new orleans connection.

Projects Underway include:

QR Generator

Order Communications

Infrastructure is mostly defined using terraform to interact with google cloud apis. Including networking, file storage buckets, database instances, project, etc.

Applications are deployed under google app engine and must be invoked via the `gcloud` command tooling.

# Getting Started - Frontend

```bash
cd web
npm install
npm run start
```

# Getting Started - Backend

The backend APIs are designed to create and retrieve either a specific or all file mapping information.

To get started in the backend we must run a node and a golang http server. The node server currently fields pdf generation tasks, the golang server fields the rest of the enumerated application logic.

To get connected to an empty database

```bash
source [.local].env.sh
cd pkg
go build ./cmd/cli
./cli db
```

The running database will print out helpful information (connection details, scripts to migrate and seed data), first change directories into the /database directory. Next, execute the migrate and seed commands. It may be helpful to update .local.env.sh with the database port as it isn't expected to say the same between process executions. Otherwise you can pass in the appropriate environment variable as a prefix to running your dependent servers.

```bash
source [.local].env.sh
cd node/order2html
PGPORT=[database port] npm start
```

```bash
source [.local].env.sh
cd pkg/cmd/cli
go build .
DB_PORT=[database port] ./cli server
```

For the server to properly interact with the database we need to sync the database with existing information from our 3rd party providers; twilio & google workspaces. This can be accomplished by running.

```bash
source [.local].env.sh
cd pkg/cmd/cli
go build .
DB_PORT=[database port] ./cli sync
```

For more commands use `./cli help`

Google Forms Prefilled Link:

https://docs.google.com/forms/d/e/1FAIpQLSeKe8iSippG-8wLxdPaGrL2Bpbqw6O8lofNN6gti98MX-YzOw/viewform?usp=pp_url&entry.1000057=New+customer&entry.2062109541=Darius+Calliet&entry.2126282437=3142859562&entry.2027335112=calliet.d@gmail.com&entry.1000026=Email&entry.1000027=2023-01-20&entry.1201531356=23:00&entry.109361733=111+Ex+Ample+Drive+Westwego+LA+12345&entry.2055232012=40+-+60+people&entry.894279446=Yes&entry.889233375=3+Gal+%7C+$20.97&entry.2001086645=Half&entry.1569890877=1+Gal&entry.626484375=50+servings&entry.737218610=Full&entry.1273257173=Full&entry.1000020=Cash