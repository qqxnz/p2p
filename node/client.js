var dgram = require('dgram');
const urlib = require("url");

var udp_client = dgram.createSocket('udp4'); 

udp_client.on('close',function(){
    console.log('udp client closed.')
})

//错误处理
udp_client.on('error', function () {
    console.log('some error on udp client.')
})

// 接收消息
udp_client.on('message', function (msg,rinfo) {
    var string = msg.toString();
    var urlParas = urlib.parse(string,true);
    if(string == '/serverRecall'){//服务器返call
        var SendBuff = '/recallsuccess';
        var SendLen = SendBuff.length;
        udp_client.send(SendBuff, 0, SendLen, rinfo.port, rinfo.address); 
    }if(string.substring(0,9) == '/message?'){//设备之间使用
        var SendBuff = `/message?data=收到[${string}]`;
        var SendLen = SendBuff.length;
        console.log(SendBuff);
        udp_client.send(SendBuff, 0, SendLen, rinfo.port, rinfo.address); 
    }
    console.log(`receive message from ${rinfo.address}:${rinfo.port}：${msg}`);
})

var SendBuff = '/keepAlive?ts=111';
var SendLen = SendBuff.length;
udp_client.send(SendBuff, 0, SendLen, 9981, '47.106.214.138'); 

//定时向服务器发送消息
setInterval(function(){
    var SendBuff = '/keepAlive?ts=111';
    var SendLen = SendBuff.length;
    udp_client.send(SendBuff, 0, SendLen, 9981, '47.106.214.138'); 
},10000);



//定时向服务器发送消息
setInterval(function(){
    var SendBuff = '/login?uid=100';
    var SendLen = SendBuff.length;
    udp_client.send(SendBuff, 0, SendLen, 9981, '47.106.214.138'); 
},5000);


//定时向服务器发送消息
setInterval(function(){
    var SendBuff = '/req_uid?uid=100';
    var SendLen = SendBuff.length;
    udp_client.send(SendBuff, 0, SendLen, 9981, '47.106.214.138'); 
},10000);