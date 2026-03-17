# hibiki
The next generation of a free quality music app made with love (to compete with the capitalists)

# Stacks used

## Sveltekit 

### Docs

https://svelte.dev/docs/kit/introduction

## Drizzle Orm 

### Docs 

https://orm.drizzle.team/docs/get-started-postgresql

### CLI 

```bash
npx drizzle-kit push
```
```bash
npx drizzle-kit generate
```
```bash
npx drizzle-kit migrate
```

## PostgreSql

### Docs v17.5 

https://www.postgresql.org/docs/17/index.html


## Better Auth

### Docs

https://www.better-auth.com/docs/integrations/svelte-kit

### CLI

```bash
npx @better-auth/cli@latest secret
```
## MusicBrainz api 

The API root URL is https://musicbrainz.org/ws/2/.

We have 13 resources on our API which represent core entities in our database:

`area, artist, event, genre, instrument, label, place, recording, release, release-group, series, work, url`
We also provide an API interface for the following non-core resources:

`rating, tag, collection`
And we allow you to perform lookups based on other unique identifiers with these resources:

`discid, isrc, iswc`
On each entity resource, you can perform three different GET requests:

lookup:   `/<ENTITY_TYPE>/<MBID>?inc=<INC>`
browse:   `/<RESULT_ENTITY_TYPE>?<BROWSING_ENTITY_TYPE>=<MBID>&limit=<LIMIT>&offset=<OFFSET>&inc=<INC>`
search:   `/<ENTITY_TYPE>?query=<QUERY>&limit=<LIMIT>&offset=<OFFSET>`
... except that browse and search are not implemented for genre entities at this time.