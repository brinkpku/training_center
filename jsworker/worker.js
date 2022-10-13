let res = "init";
let timerID = null;
let interval = 2000;
self.onmessage = function(e){
    if(e.data == "start"){
        res = "worker runing";
        timerID = setInterval(() => {
            self.postMessage(res);
        }, interval); 
    } else if (e.data == "stop"){
        res = "worker stop";
        self.postMessage(res);
        if (timerID != null){
            clearInterval(timerID);
        }
    }

}
self.postMessage("worker load");