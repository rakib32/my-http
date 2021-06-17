## my-http:
this is a small tool which makes http requests and prints the address of the request along with the MD5 hash of the response written in go.

## building

```
git clone https://github.com/rakib32/my-http.git
cd my-http
go build -o myhttp
```

## usage

```
./myhttp [flags] urls
```

with the flags being
```
    -parallel=10: the max numbers of parallel processes
    
```
for example
```
./myhttp http://www.yahoo.com http://google.com
```

**Note: Please `control+c` to stop the tool(Graceful shutdown).**

## Testing
Use following command to run the test
```
go test ./...
```

## Project Workflow

Here, i have used a buffered queue to process the requests in parallel.
* **Dispatcher**: It will initialize the JobQueue and Worker Pool based on the settings.
* **Worker**: Dispatcher will assign the task to worker and worker will process the task and call handler function.

```  
    Job Queue(Buffered)
    =================================
    |<-Job1->| |<-Job2->| |<-Job3->|
    =================================
                  
              ----------
            ├ Dispatcher ├ ---> Pulls job from Job Queue and assign Job to Worker
             ------------
                
    WorkPool Queue(Channel over channel)
    =========================================
    |<-Worker1->| |<-worker2->| |<-worker3->|
    =========================================     
```
