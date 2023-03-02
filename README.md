# QR Generator

This QR generator borrows techniques from URL shortening systems and will allow for files within the original QR to be updated. The api service will support the generation and tracking of QR codes with their current and historical attachments.

Expectations around beta release:

User can store a file and receive a QR code that will return selected file.
User can view all QR codes currently available.
User can update the file that a QR will return.

# Getting Started - Frontend

The Frontend interface is designed to support manipulating your store of QR codes and the items they map to. QR codes can be given a label, will automatically generate an appropriate QR code and will allow storage to an s3 bucket any items you would wish the QR code to return.

To get started in the frontend:

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
cd pkg
source ../export_env_local.sh
go build ./cmd/cli; ./cli db
```

The running database will print out helpful information (credentials, scripts to migrate and seed data), first change directories into the /database directory. Next, execute the commands sequentially.

```bash
cd node/order2html
source ./my-env.sh
PGPORT=<check database initialization logs> npm start
```

```bash
cd pkg/cmd/cli
source ../export_env_local.sh
go build .; DB_PORT=<check database initialization logs> ./cli server
```

For the server to properly interact with the database we need to sync the database with existing information from our 3rd party providers; twilio & google workspaces. This can be accomplished by running.

```bash
cd pkg/cmd/cli
source ../export_env_local.sh
go build .; DB_PORT=<check database initialization logs> ./cli sync
```

For more commands use `help`

Google Forms Prefilled Link:

https://docs.google.com/forms/d/e/1FAIpQLSeKe8iSippG-8wLxdPaGrL2Bpbqw6O8lofNN6gti98MX-YzOw/viewform?usp=pp_url&entry.1000057=New+customer&entry.2062109541=Darius+Calliet&entry.2126282437=3142859562&entry.2027335112=calliet.d@gmail.com&entry.1000026=Email&entry.1000027=2023-01-20&entry.1201531356=23:00&entry.109361733=111+Ex+Ample+Drive+Westwego+LA+12345&entry.2055232012=40+-+60+people&entry.894279446=Yes&entry.889233375=3+Gal+%7C+$20.97&entry.2001086645=Half&entry.1569890877=1+Gal&entry.626484375=50+servings&entry.737218610=Full&entry.1273257173=Full&entry.1000020=Cash