# scrapbox-viz

scrapbox-viz (sbv) is a CLI for fetching and visualizing data from [Scrapbox](https://scrapbox.io) project visualize the.

- Fetch page data (JSON format)
- Aggregate user activities (pages created, views of created page, etc.)
- Generate graph data (as Graphviz dot file)

## Initialize working directory
Data fetched via Scrapbox APIs will be stored in an existing working directory.

If you create `.sbv.yaml` at `${HOME}` and write path in entry `workdir`, working directory will be set to that path.

```yaml
workdir: path/to/workdir
```

If the directory does not exist, it will be created.

Of cource, you can specify the directory every time you execute sub commands with global -d(--workdir) flag.

```
$ sbf fetch -p <project name> -d <path/to/workdir>
```

## Fetch page data of the project
Fetch page data of the Scrapbox project via [Scrapbox APIs](https://scrapbox.io/help-jp/API).

- Page list data will be saved as JSON file at `<WorkDir>/<project name>.json`.
- Each Page data will be saved as JSON file in `<WorkDir>/<project name>`.
  - The file name consists of the page ID.

```
$ sbf fetch -p <project name>
```

To fetch from a private project, you needs to set the cookie to environment variables.

```bash
$ export SB_COOKIE_ID=connect.sid
$ export SB_COOKIE_VALUE=your-fancy-cookie
```

## Aggregate user activites in the project
Parse page data and aggregate activities of the project per user.

- Pages created
- Pages contributed
- Views of created page
- Links of created page

```
$ sbf aggregate -p <project name>
```

CSV will be created at `<WorkDir>/<project name>.csv`.

## Generate graph of the pages and users
Parse page data and generate graph of pages and users.

```
$ sbf graph -p <project name>
```

If you want to include user node to the graph, specify -i(--include) flag.

```
$ sbf graph -p <project name> -i=true
```

If you want to annonymize user name of user node, specify -a(--anonymize) flag.

```
$ sbf graph -p <project name> -i=true -a=true
```

You can reduce number of nodes in the graph by specifying page views as threshold value.

```
$ sbf graph -p <project name> -t 100
```

Graphviz dot file will be created at `<WorkDir>/<project name>.dot`.
