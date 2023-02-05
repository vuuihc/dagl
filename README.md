# daglc
## 定义
dagl 一个简洁的领域特定语言（DSL)，用于定义一个有向无环图（DAG)。可以用来描述一个工作流。
编译后的结果可以适用于工作流引擎。

// english
dagl is an easy-to-use domain-specific language (DSL) for defining a directed acyclic graph (DAG). It can be used to describe a workflow.

## 语法
### 基本类型
1. 字符串
```dagl
"hello world"
`hello world`
'hello world'
```
### 数据结构
1. 数组
```dagl
[1,2,3]
```
2. 参数对
```dagl
a=b
```
### 控制流
1. if else
```dagl
if (a==b) {
  a;
}else{
  b;
}
```
### 函数
1. 内置函数 
```dagl
builtin("http", input, endpoint=`http://192002625-146479.Production/suggestion/`,
        method=`post`, max_retry_times="3", default_value=`{"actions":[]}`, timeout="800ms");
```
2. 自定义函数
```dagl
inline func getCacheKey(input) {
  builtin("jq",input,filter=`.suggestion_type+"##"+(.filter_retrievers//[]|join("#"))+"##"+(.context//[]|join("#"))+"##"+.query`);
}
```

### 调用
1. 调用自定义函数
```dagl
@call(getCacheKey, [input]);
```
### 注释
```dagl
// this is a comment
```

### 常量
常量只能在函数外部定义。
1. 常量定义
```dagl
@a=`hello world`;
```
2. 常量引用
```dagl
call(abc,[],a=@a);
```

### 变量
变量只能在函数内部定义。

## 完整示例
```dagl
// a function to get cache key
inline func getCacheKey(input) {
  builtin("jq",input,filter=`.suggestion_type+"##"+(.filter_retrievers//[]|join("#"))+"##"+(.context//[]|join("#"))+"##"+.query`);
}

// a function to set cache
inline func setCache(key, result) {
	cacheReq=builtin("jq",[key,result],filter=`{"key": .[0], "payload": .[1], "ttl": 259200000}`);
	builtin("set_cache", cacheReq, prefix=`ime_rec_bert_ner_v1`);
}

// a function to lookup cache
inline func lookupCache(key){
  builtin("lookup_cache", key, prefix=`ime_rec_bert_ner_v1`);
}

func main(input) {
    input = builtin("jq", input, filter=`.payload | fromjson`);
    key=@call(getCacheKey, [input]);
    cacheRes=@call(lookupCache,[key]);
    result=builtin("http", input, endpoint=`http://192002625-146479.Production/suggestion/`,
        method=`post`, max_retry_times="3", default_value=`{"actions":[]}`, timeout="800ms");
    @call(setCache, [key, result]);
    cacheMiss=builtin("jq", cacheRes, filter=`.found | not`);
    if(cacheMiss){
      result;
    }else{
      builtin("jq", cacheRes, filter=`.payload`);
    }
}

// {"payload": "{\"request_id\":\"1674\",\"request_type\":7,\"context\":[],\"context_interval\":[],\"query\":\"红楼梦小姐姐\",\"uid\":\"1674\",\"api_level\":0}"}
```

# english document
## definition
dagl is an easy-to-use domain-specific language (DSL) for defining a directed acyclic graph (DAG). It can be used to describe a workflow.

## syntax
### basic type
1. string
```dagl
"hello world"
`hello world`
'hello world'
```
### data structure
1. array
```dagl
[1,2,3]
```
2. key-value pair
```dagl
a=b
```
### control flow
1. if else
```dagl
if (a==b) {
  a;
}else{
    b;
}
```
### function
1. builtin function 
```dagl
builtin("http", input, endpoint=`http://192002625-146479.Production/suggestion/`,
        method=`post`, max_retry_times="3", default_value=`{"actions":[]}`, timeout="800ms");
```
2. inline function
```dagl
inline func getCacheKey(input) {
  builtin("jq",input,filter=`.suggestion_type+"##"+(.filter_retrievers//[]|join("#"))+"##"+(.context//[]|join("#"))+"##"+.query`);
}
```

### call
1. call inline function
```dagl
@call(getCacheKey, [input]);
```
### comment
```dagl
// this is a comment
```

### constant
constant can only be defined outside of function.
1. constant definition
```dagl
@a=`hello world`;
```
2. constant reference
```dagl
@b=@a;
```

### variable
variable can only be defined inside of function.

## full example
```dagl
// a function to get cache key
inline func getCacheKey(input) {
  builtin("jq",input,filter=`.suggestion_type+"##"+(.filter_retrievers//[]|join("#"))+"##"+(.context//[]|join("#"))+"##"+.query`);
}

// a function to set cache
inline func setCache(key, result) {
    cacheReq=builtin("jq",[key,result],filter=`{"key": .[0], "payload": .[1], "ttl": 259200000}`);
    builtin("set_cache", cacheReq, prefix=`ime_rec_bert_ner_v1`);
}

// a function to lookup cache
inline func lookupCache(key){
  builtin("lookup_cache", key, prefix=`ime_rec_bert_ner_v1`);
}

func main(input) {
    input = builtin("jq", input, filter=`.payload | fromjson`);
    key=@call(getCacheKey, [input]);
    cacheRes=@call(lookupCache,[key]);
    result=builtin("http", input, endpoint=`http://192002625-146479.Production/suggestion/`,
        method=`post`, max_retry_times="3", default_value=`{"actions":[]}`, timeout="800ms");
    @call(setCache, [key, result]);
    cacheMiss=builtin("jq", cacheRes, filter=`.found | not`);
    if(cacheMiss){
      result;
    }else{
      builtin("jq", cacheRes, filter=`.payload`);
    }
}

// {"payload": "{\"request_id\":\"1674\",\"request_type\":7,\"context\":[],\"context_interval\":[],\"query\":\"红楼梦小姐姐\",\"uid\":\"1674\",\"api_level\":0}"}
```
