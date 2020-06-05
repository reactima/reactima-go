# pgbroadcaster
Package that listen to PostgresSQL JSON notifications and exposes a websocket handler to broadcast the notification to client.

## Usage

### Postgres

Create a trigger that nofifies a certain channel, here called events, in the following form:

```SQL
CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS $$

    DECLARE 
        data json;
        notification json;
    
    BEGIN
    
        -- Convert the old or new row to JSON, based on the kind of action.
        -- Action = DELETE?             -> OLD row
        -- Action = INSERT or UPDATE?   -> NEW row
        IF (TG_OP = 'DELETE') THEN
            data = row_to_json(OLD);
        ELSE
            data = row_to_json(NEW);
        END IF;
        
        -- Contruct the notification as a JSON string.
        notification = json_build_object(
                          'table',TG_TABLE_NAME,
                          'action', TG_OP,
                          'data', data);
        
                        
        -- Execute pg_notify(channel, notification)
        PERFORM pg_notify('events',notification::text);
        
        -- Result is ignored since this is an AFTER trigger
        RETURN NULL; 
    END;
    
$$ LANGUAGE plpgsql;
```

Add the trigger to your table:

```SQL
CREATE TRIGGER exampletable_notify_event
AFTER INSERT OR UPDATE OR DELETE ON exampletable
    FOR EACH ROW EXECUTE PROCEDURE notify_event();
```

### Server

In your server app, just start the listener and expose the provided handler in your web app.

``` go
package main

import (
	"log"
	"net/http"
	"fmt"

	"github.com/coussej/pgbroadcast"
)

func main() {
	// Create a new broadcaster
	pb, err := pgbroadcaster.NewPgBroadcaster("dbname=exampledb user=webapp password=webapp")

	// listen to the events channel
	err = pb.Listen("events")
	if err != nil {
		fmt.Println(err)
	}

	// server the websocket	
	http.HandleFunc("/ws", pb.ServeWs)
	err = http.ListenAndServe(":6060", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
```

### Client

In the client app, open a websocket and send the table(s) you want to subscribe to as a string

```javascript
    $(function() {
    var conn;
    if (window["WebSocket"]) {
        conn = new WebSocket("ws://localhost/ws");
        conn.onopen = function(evt) {
          // SUBSCRIBE TO TABLE
          conn.send("exampletable");
        }  
        conn.onclose = function(evt) {
            // HANDLE CONNECTION CLOSE
        }  
        conn.onmessage = function(evt) {
            var data = JSON.parse(evt.data);
            // HANDLE DATA
        }
    } 
    });
```
