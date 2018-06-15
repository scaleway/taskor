# How It Works

![worker](assets/worker.png)

* Get task from a runner and send it to a chan (taskToRun)
* Get task from the chan
    * If the task should be run now send it to the chan for processing (taskToProcess)
    * Else wait the ETA before sending to the chan
* At the end of the task (retry or not) push to a chan (taskDone) to inform runner that this task should be ACK.

# What If worker shutdown / break ?
All task not ACK will be reschedule on other worker.