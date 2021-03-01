# cSubsidiary
 利用天眼查查询企业子公司

## 介绍

可以通过两种方式查询自己想要的企业子公司

1. `-n` 参数：利用给出的关键字先进行模糊查询，然后选出第一个匹配的结果，对该公司进行查询。（方便但是不准确，所以不推荐）
2. `-i` 参数：利用给出的公司id对该公司进行查询。（准确，结果唯一，但需要自己先去查找一级公司，比较推荐）
3. `-f` 参数：对文件里的所有关键字和id进行查询。因为我比较推荐用id查询，而且为了方便多次递归查询，读文件时会先去该行尝试匹配是否存在公司id，假如不存在就把该行作为关键字进行查询。因此递归查询可以直接把a次的结果文件当作b次的输入文件。
4. 因为天眼查风控比较严格，所以使用时会出现几种情况。一、因为某个ip一段时间内查询次数过多，所以查询时会自动跳到登陆界面，这种情况需要使用一个手机号进行登陆，然后增加cookie去继续查询。二、海外ip或者云服务器ip访问天眼查会显示海外用户，所以最好使用正常的出口ip进行查询。三、假如短时间呢使用很多很多次查询的话，天眼查会有人机判断的验证码，需要手动打开天眼查网站进行一下人机验证。（情况较少）

## 用法

```
admin@admin cSubsidiary % go run cSubsidiary.go -h
Usage of cSubsidiary:
  -c string
    	天眼查的Cookie
  -f string
    	包含公司ID号码的文件
  -i string
    	公司ID号码
  -n string
    	公司名称
  -no-color
    	No Color
  -o string
    	结果输出的文件(可选)
  -p int
    	显示投资比例为x以上的公司 (default 100)
  -silent
    	Silent mode
  -timeout int
    	连接超时时间 (default 15)
  -verbose
    	详细模式
  -version
    	显示软件版本号
  -w int
    	显示投资资金为x万以上的公司 (default 100)
```

查询子公司

```
admin@admin cSubsidiary % go run cSubsidiary.go -n 字节跳动
[INFO] 正在查询 https://www.tianyancha.com/company/2352987806 的子公司
https://www.tianyancha.com/company/3478099686	小荷健康科技（北京）有限公司
https://www.tianyancha.com/company/3417563763	北京潜龙在渊科技有限公司
https://www.tianyancha.com/company/3414099437	北京游逸科技有限公司
https://www.tianyancha.com/company/3294247465	北京星云创迹科技有限公司
https://www.tianyancha.com/company/3285950613	北京量子跃动科技有限公司
https://www.tianyancha.com/company/3273636827	天津基石科技有限公司
https://www.tianyancha.com/company/3271306567	天津千江科技有限公司
https://www.tianyancha.com/company/3263758754	北京字节新异科技有限公司
https://www.tianyancha.com/company/3255652581	大力创新科技（北京）有限公司
https://www.tianyancha.com/company/3205765844	北京光锥之外科技有限公司
https://www.tianyancha.com/company/3168973708	北京跳动空间科技有限公司
https://www.tianyancha.com/company/2791028018	今日头条有限公司
https://www.tianyancha.com/company/2351052730	江苏今日头条信息科技有限公司
https://www.tianyancha.com/company/515712303	上海图虫网络科技有限公司
https://www.tianyancha.com/company/25174642	北京字节跳动科技有限公司
```

## 其他

软件难免有一些问题，假如大家发现，欢迎大家提意见或者建议。

还有一个工具 `cDomain` 我一般两个一起使用，我后面写个文章，详细写一下