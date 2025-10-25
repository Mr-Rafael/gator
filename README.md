## Prerequisites

### PostgreSQL
Postgres v15 or later.

**macOS** with brew

```brew install postgresql@15```

**Linux / WSL (Debian)**

```sudo apt update```

```sudo apt install postgresql postgresql-contrib```

Check installation with:

```psql --version```

(Linux only) Update postgres password:

```sudo passwd postgres```

### Go

Follow the instructions at

https://go.dev/doc/install

## Installation

Install using ```go install```:

```
go install github.com/username/git@github.com:Mr-Rafael/gator.git
```

## Configuration

Create a ```.gatorconfig.json``` file on your ```~/``` directory. 

Configure the database URL:

```
{
  "db_url": "postgres://<user>:<password>@localhost:5432/gator?sslmode=disable"
}
```

## Use

Build the program from the ```<root>/cmd/gator``` folder:

```go build```

Then, execute the program with any of the valid commands:

```./gator <command> <parameters>```

## Valid Commands

### Register

```register  [username]```

Creates a new user with the specified username.

### Login

```login [username]```

Sets the current user to the specified one (by username)

### Users

```users```

Displays a list of all existing users.

### Login

```login [username]```

Sets the current user to the specified one (by username)

### Add Feed

```addfeed [name] [url]```

Creates a new feed with the specified URL and Name, and sets the current logged user to follow that feed.

### Feeds

```feeds```

Displays a list of all registered feeds.

### Aggregate

```agg [time between scrapes]```

Starts continuous loop that scrapes (updates) all feeds periodically, with the frequency specified. The program will continue running until stopped.

### Follow

```follow [url]```

Sets the current User to follow the feed specifid by URL. The Feed should have been added first with the Add Feed command.

### Following

```following```

Displays all the Feeds that the current User is following.

### Unfollow

```unfollow [url]```

Makes the current user unfollow the Feed specified by URL.

### Browse

```browse [limit (default 2)]```

Displays the most recent posts from the feeds the current user is following. Displays the N most recent posts if specified, or 2 by default.

### Reset

```reset```

Resets the data in Gator (users, feeds, posts and follows)