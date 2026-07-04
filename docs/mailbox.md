# Mailbox

Search the documents you have received (inbox) or sent (outbox).

Import: `github.com/brifle-de/brifle-sdk/sdk/endpoints/mailbox`

Both searches are paginated and accept an optional filter on `subject`, `state` and `type`.

## SearchMyInbox

```go
func SearchMyInbox(client *client.BrifleClient, ctx context.Context, inboxSearch *InboxSearch) (*InboxSearchResponse, *api.ResponseStatus, error)
```

Searches the inbox of the authenticated user.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

page := float32(1)
search := mailbox.InboxSearch{
	Page: &page,
	Filter: &mailbox.InboxSearchFilter{
		Subject: sdk.String("Invoice"),
	},
}

res, respStatus, err := mailbox.SearchMyInbox(client, ctx, &search)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}

for _, item := range res.Results {
	fmt.Println(*item.Id, *item.Subject)
}
```

## SearchOutbox

```go
func SearchOutbox(client *client.BrifleClient, ctx context.Context, tenant *string, outboxSearch *OutboxSearch) (*OutboxSearchResponse, *api.ResponseStatus, error)
```

Searches the outbox of the account. By default it returns documents sent by the authenticated user;
set `SenderUser` to filter by a specific user or API key.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

tenant := "567e44de-b6b6-4dac-cbce-c5515031f9ea"
page := float32(1)
search := mailbox.OutboxSearch{
	Page: &page,
}

res, respStatus, err := mailbox.SearchOutbox(client, ctx, &tenant, &search)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}

for _, item := range res.Results {
	fmt.Println(*item.Id, *item.Subject)
}
```

### Filtering

`InboxSearchFilter` / `OutboxFilter` fields:

| Field | Type | Description |
|---|---|---|
| `Subject` | `*string` | Match on subject. |
| `State` | `[]*string` | One or more document states. |
| `Type` | `*string` | Document type, e.g. `letter`, `invoice`, `contract`. |
