# 阿里云短信推送

本项目通过调用aliyun短信接口，实现短信推送

参考文档：

* [接口地址](https://help.aliyun.com/document_detail/101511.html)
* [SendSms接口定义](https://help.aliyun.com/document_detail/101414.html)
* [Go SDK](https://github.com/aliyun/alibaba-cloud-sdk-go/tree/master/services/dysmsapi)

## 接口定义

* 方法：POST
* 地址：/sms
* 参数：采用以下 json 格式（所有字段均为必填）

        {
          "PhoneNumber": "11位手机号码，必须以1开头",
          "User": "用户名，支持中文",
          "DeviceID": "设备代码，长度1~20，只支持大小写字母，数字和字符",
          "State": "设备当前状态，支持中文"
        }
* 例如：

        {
          "PhoneNumber": "13500000000",
          "User": "admin",
          "DeviceID": "K1103-002",
          "State": "已离线3小时，请尽快处理"
        }

会向 13500000000 发送一条内容为“admin，您的设备K1103-002当前已离线3小时，请尽快处理”的短信

## 下载方法

### 1. 使用编译好的二进制程序（推荐）
从[release](https://github.com/newkedison/aliyun-sms-push/releases)页面，下载最新的压缩包，
解压后根据所在的系统，选择对应的程序运行，目前支持
* Linux 64bit
* Windows 32bit
* Windows 64bit

### 2. 下载源码编译（请确保网络畅通）

1. clone 此版本库

    $ git clone git@github.com:newkedison/aliyun-sms-push.git
    
2. 编译

    $ cd aliyun-sms-push
    
    $ go build
    
## 配置说明

有三个配置文件

* config.yaml 配置程序参数，可参考[config.yaml.template](https://github.com/newkedison/aliyun-sms-push/blob/master/config.yaml.template)文件
* [phone.txt](https://github.com/newkedison/aliyun-sms-push/blob/master/phone.txt)
目标手机号码白名单，只有列在白名单中的手机号码，才允许通过此接口发送短信
* [ip.txt](https://github.com/newkedison/aliyun-sms-push/blob/master/ip.txt)
IP地址白名单，只允许白名单中的IP访问此接口，支持[CIDR](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing)，默认允许localhost和三个私有网段访问

### 关于数据库

1. 本程序中使用了 [MongoDB 数据库](https://www.mongodb.com/)
2. 短信的发送记录会自动保存到数据库的 sms_record 集合中
3. 数据库的名称目前固定为sms_push，暂不支持更改

### 关于白名单

1. 除了两个配置文件，还可以通过数据库配置白名单

    * 手机号码白名单保存于 sms_phone_list 集合中
    * IP白名单保存于 sms_ip_list 集合中
    
2. 程序每隔一段时间，会重新读取一遍两个配置文件和两个数据库集合，因此对白名单的更新无需重启程序，只需等待一段时间后，即可生效

### 关于安全性

考虑到本项目主要用于内部，因此没有加上复杂的权限验证，而是通过白名单的方式，避免被滥用。

如果要使用本程序，一般建议只把内网IP加入到白名单里面，这样可以最大限度的保证安全。

另外就是注意 ** 绝对不要通过任何方式，将config.yaml文件公布到公有网络 ** ，
以避免泄漏阿里云的密钥和数据库密钥
