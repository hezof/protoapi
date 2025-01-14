# protodoc

## go-swagger问题

- 执行"swagger minxin"后“paths”或“tags”失去了原序! 让李同学相当困惑！
- 执行“swagger serve”后采用最新的"redoc.standalone.js". 存在若干无解的BUG. 例如：

```
* OpenAPI 3.1: Missing description when $ref used [#1727](https://github.com/Redocly/redoc/issues/1727) ([fe6909e](https://github.com/Redocly/redoc/commit/fe6909ed80dd6053b48c30f63a2460614bf957a9))
* OpenAPI 3.1: Missing description when $ref used [#1727](https://github.com/Redocly/redoc/issues/1727) ([35f7787](https://github.com/Redocly/redoc/commit/35f77878de7d1dd250040771f17757a5a6ce85f9))
```

