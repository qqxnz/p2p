var dgram = require('dgram');
const urlib = require("url");

var userObj = new Object();

//创建 udp server
var udp_server = dgram.createSocket('udp4');
udp_server.bind(9981); // 绑定端口

// 监听端口
udp_server.on('listening', function () {
    console.log('udp server linstening 9981.');
})
//接收消息
udp_server.on('message', function (msg, rinfo) {
    strmsg = msg.toString();
    console.log(`udp server received data: ${strmsg} from ${rinfo.address}:${rinfo.port}`);
    var myobj = urlib.parse(strmsg,true);
    if(myobj.path.substring(0,11) == '/keepAlive?'){
        setTimeout(() => {
            clientRecall(rinfo.address,rinfo.port);
        }, 1000);
    }else if(myobj.path.substring(0,7) == '/login?'){
        var uid = myobj.query.uid;
        userObj[uid] = {
            'ip':rinfo.address,
            'port':rinfo.port,
        };
        console.log(userObj);
        var result = `/success?ip=${rinfo.address}&port=${rinfo.port}`
        udp_server.send(result, 0, result.length, rinfo.port, rinfo.address);

    }else if(myobj.path.substring(0,9) == '/req_uid?'){
        var uid = myobj.query.uid;
        var obj = userObj[uid];
        if(obj == null){
            console.log(`未找到uid:${uid}对象`);
            return
        }
        console.log(`req_uid:${uid} ip:${obj.ip} port:${obj.port}`);
        var result = `/res_uid?uid=${uid}&ip=${obj.ip}&port=${obj.port}`
        udp_server.send(result, 0, result.length, rinfo.port, rinfo.address);
    }
})
//错误处理
udp_server.on('error', function (err) {
    console.log('some error on udp server.')
    udp_server.close();
})


function clientRecall(ip,port){
    udpClient(ip,port,function(result){
        console.log('recall receive:'+result);
        if(result.length < 1){
            console.log(`${ip}:${port}`+'recall fail');
        }else{
            console.log(`${ip}:${port}`+'recall success');
        }
    });
}

function udpClient(ip,port,callback){
    var udp_client = dgram.createSocket('udp4'); 
    var result = '';
    udp_client.on('close',function(){
        console.log('udp client closed.')
    })
    
    //错误处理
    udp_client.on('error', function () {
        console.log('some error on udp client.')
    })
    
    // 接收消息
    udp_client.on('message', function (msg,rinfo) {
        result = msg.toString();
        console.log(`receive message from ${rinfo.address}:${rinfo.port}：${msg}`);
    })
    var SendBuff = '/serverRecall';
    var SendLen = SendBuff.length;
    udp_client.send(SendBuff, 0, SendLen, port,ip); 
    setTimeout(function(){
        callback(result);
    },1000);
}