# Cloud Run Log Streaming

Cloud Run Log Streaming is a small tool used to streaming logs directly from your cloud run services to your terminal.

Log streaming works by utilizing underlying Log Tailing Entries Api 

**Cloud Run Log Streaming is not an officially supported Google product.**

## Usage

Note: you must install and authenticated to the [Google Cloud
SDK](https://cloud.google.com/sdk) (gcloud) for the proxy to pull your
authentication token. You local user must also have Cloud Run Invoker
permissions on the target service.

1.  Run the command:

    ```sh
    cloud-run-log-streaming -projectId=my-projectId -filter="resource.type=cloud_run_revision severity>=DEAFULT"
    ```



## Options

###Filters

To read log entries of a specfic cloud run service, run:

```sh
    cloud-run-log-streaming -projectId=my-projectId -filter="resource.type=cloud_run_revision resource.labels.service_name=my-service-name resource.labels.location=us-west1 serverity>=DEAFULT"
```

To read log entries of a specfic cloud run revision, run:

```sh
    cloud-run-log-streaming -projectId=my-projectId -filter="resource.type=cloud_run_revision resource.labels.revision_name=my-revision-name resource.labels.location=us-west1 serverity>=DEAFULT"
```

To read log entries with severity ERROR or higher, run:

```sh
cloud-run-log-streaming -projectid=my-projectid -filter="resource.type=cloud_run_revision serverity>=ERROR"
```

To read log entries written in a specific time window, run:

```sh
cloud-run-log-streaming -projectid=my-projectid -filter='resource.type=cloud_run_revision timestamp<="2015-05-31T23:59:59Z" AND timestamp>="2015-05-31T00:00:00Z"'
```

Detailed information about filters can be found [here](https://cloud.google.com/logging/docs/view/advanced_filters).




###Process-up Time (optional) 


```sh
cloud-run-log-streaming -projectid=my-projectid -filter="resource.type=cloud_run_revision" -process-up-time=2h"
```

or 
```sh
cloud-run-log-streaming -projectid=my-projectid -filter="resource.type=cloud_run_revision" -process-up-time=1m30s"
```
note: Not setting a Process-up time will keep the process running indefinitely.
