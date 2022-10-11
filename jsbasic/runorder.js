console.log("start");
function f1(){
    console.log("f1")
}

async function f2(){
    console.log("f2")
}

f3 = ()=>{
    console.log("f3")
}

setTimeout(f1, 0) // async
f2()
setTimeout(f3, 1)