<p>本项目是一个导入/导出中间件服务，只需传入创建资源和获取资源的api地址，就能轻松实现导入和导出功能，而无需修改该服务本身。本项目提供了RESTful风格的导入/导出接口，通过传入指定参数就能简单实现10w+数据导入/导出功能。</p>

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
