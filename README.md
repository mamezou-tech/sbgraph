# sbgraph

![Go](https://github.com/mamezou-tech/sbgraph/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/mamezou-tech/sbgraph)](https://goreportcard.com/report/github.com/mamezou-tech/sbgraph)

sbgraph is a CLI for fetching and visualizing data from [Scrapbox](https://scrapbox.io) project.

- Fetch page data (JSON format)
- Aggregate user activities (pages created, views of created page, etc.)
- Generate graph data (as Graphviz dot file)

![20200306204303](https://user-images.githubusercontent.com/2092183/79331841-ca874880-7f56-11ea-9127-c1f249742028.png)


## Installing

```
go get -u github.com/mamezou-tech/sbgraph
```
or
```
curl -LO https://github.com/mamezou-tech/sbgraph/releases/download/<version>/sbgraph-<platform>-amd64.tar.gz
tar xvf sbgraph-<platform>-amd64.tar.gz
sudo mv sbgraph /usr/local/bin
```

## Sub commands

### Initialize config & working directory

```
sbgraph init
```

Config file (.sbgraph.yaml) will be created in users home directory.

Data fetched via Scrapbox APIs will be stored in an existing working directory.

By default working directory will be set `$(pwd)/_work`

In config file working directory will be set to that path.

```yaml
workdir: path/to/workdir
```

If the working directory does not exist, it will be created.

Of cource, you can specify the directory every time you execute sub commands with global -d(--workdir) flag.

```
sbgraph fetch -d <path/to/workdir>
```

### Set target project

```
sbgraph project -p <project name>
```

```yaml
currentproject: project-name
```

### Print configuration status

```
sbgraph status
```

Config file path and current settings will be printed.

### Fetch page data of the project
Fetch page data of the Scrapbox project via [Scrapbox APIs](https://scrapbox.io/help-jp/API).

- Page list data will be saved as JSON file at `<WorkDir>/<project name>.json`.
- Each Page data will be saved as JSON file in `<WorkDir>/<project name>`.
  - The file name consists of the page ID.

```
sbgraph fetch
```

To fetch from a private project, you needs to set the cookie to environment variables.

```
export SB_COOKIE_ID=connect.sid
export SB_COOKIE_VALUE=your-fancy-cookie
```

### Aggregate user activites in the project
Parse page data and aggregate activities of the project per user.

- Pages created
- Pages contributed
- Views of created page
- Links of created page

```
sbgraph aggregate
```

CSV will be created at `<WorkDir>/<project name>.csv`.

### Generate graph of the pages and users
Parse page data and generate graph of pages and users.

```
sbgraph graph
```

If you want to include user node to the graph, specify -i(--include) flag.

```
sbgraph graph -i=true
```

If you want to annonymize user name of user node, specify -a(--anonymize) flag.

```
sbgraph graph -i=true -a=true
```

You can reduce number of nodes in the graph by specifying page views as threshold value.

```
sbgraph graph -t 100
```

Graphviz dot file will be created at `<WorkDir>/<project name>.dot`.


To generate graph data as JSON, specify -j(--json) flag.

```
sbgraph graph -j=true
```

To generate graph Image as SVG format, specify -m(--image) flag.

```
sbgraph graph -m=true
```

Graphviz needs to be installed. You can not specify layout engine.

### Visualize with Graphviz

You can generate SVG with graph sub command, but to specify more options such as layout engine, you can use dot command directly.

e.g.

```
dot <project name>.dot -Tsvg -Kfdp -Goverlap=prism -<project name>.svg
```
