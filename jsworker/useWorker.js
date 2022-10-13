// run in browser
let running = false;

worker = new Worker("./worker.js");

worker.onmessage = function(e){
    msg = Date()+"\t"+e.data
    // console.log(msg);
    document.getElementById("output").innerText = msg;
};


function run(){
    running = !running;
    if(running){
        this.innerText = "stop";
        worker.postMessage("start");
    }else{
        this.innerText = "start";
        worker.postMessage("stop");
    }
}