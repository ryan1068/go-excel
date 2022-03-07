<p>本项目是一个中间件服务，本身不实现创建资源和获取资源业务逻辑，只提供RESTful接口，通过传参的方式实现导入/导出功能。</p>

Example：

导入：
```go
POST https://{{host}}/excel/import {
    "taskType": "import_account",
    "filePath": "host/import/clue/2021/10/28141280193823672622088628622460301cc63e61f24c425616e379c08ce8af.xlsx",
    "headerKey": "nickname,account,password,department_name,is_manager,role_name,tel",
    "apiPath": "{{host}}/v1/accounts"
}
```

导出：
```go
POST https://{{host}}/excel/export {
    "taskType": "export_account",
    "headerTitle": "门店编号,门店ID,门店名称,经营品牌,销售城市,经营地址,省份,城市,行政区,门店状态,总经理",
    "headerKey": "number,id,name,car_brands_name,sale_city_name,address,province_name,city_name,region_name,status_name,admins_name",
    "apiPath": "{{host}}/v1/accounts?group_id=1"
}
```
