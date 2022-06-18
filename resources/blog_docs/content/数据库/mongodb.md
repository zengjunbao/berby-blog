

https://docs.mongoing.com/mongodb-crud-operations/query-documents

https://www.mongodb.org.cn/tutorial/8.html



## 用法

MongoDB 中默认的数据库为 test，如果你没有创建新的数据库，集合将存放在 test 数据库中。

创建完数据库如果没有数据，则不会显示该数据库。



> 1.  切换数据库，不存在则创建 
>
>    `use DATABASE_NAME`
>
> 2. 查看数据库(当前数据库 db)
>
>    `show dbs`
>
> 3. 删除数据库(默认删除当前数据库)
>
>    `db.dropDatabase()`
>
> 4. 删除集合
>
>    `db.collection.drop()  `
>
> 5. 插入数据 (WriteResult 指定id 存在则更新数据)
>
>    `db.COLLECTION_NAME.insert(document)   WriteResult({ "nInserted" : 1 })   `  
>
> 6. 查看集合数据
>
>    `db.COLLECTION_NAME.find()`
>
> 7. 更新数据
>
>    update
>
>    ```
>    db.collection.update(    
>    	<query>, 
>    	<update>, 
>    	{       
>    		upsert: <boolean>,   
>    		multi: <boolean>,  
>    		writeConcern: <document>
>    	}
>    )
>    
>    query : update的查询条件，类似sql update查询内where后面的。
>    update : update的对象和一些更新的操作符（如$,$inc...）等，也可以理解为sql update查询内set后面的
>    upsert : 可选，这个参数的意思是，如果不存在update的记录，是否插入objNew,true为插入，默认是false，不插入。
>    multi : 可选，mongodb 默认是false,只更新找到的第一条记录，如果这个参数为true,就把按条件查出来多条记录全部更新。
>    writeConcern :可选，抛出异常的级别。
>    ```
>
>    save
>
>    ```
>    db.collection.save(    
>    	<document>,     
>    	{      
>    		writeConcern: <document> 
>    	}  
>    )  
>    document : 文档数据。
>    writeConcern :可选，抛出异常的级别。
>    ```
>
> 8. 删除文档
>
>    ```
>    db.collection.remove(     
>    	<query>,     
>    	{       
>    		justOne: <boolean>,
>    		writeConcern: <document> 
>    	} 
>    )
>    query :（可选）删除的文档的条件。
>    justOne : （可选）如果设为 true 或 1，则只删除一个文档。
>    writeConcern :（可选）抛出异常的级别。
>    ```
>
>    



#### 查询条件



> 查询字段 
>
> db.inventory.find( { status: "A" }, { status: 0, instock: 0 } ) //除status (0忽略 1选择)和 instock 字段之外的所有字段
>
> 
>
> 查询数组
>
> db.col.find( { tags: { $all: ["red", "blank"] } } ) // 查询tags数组中包含这两个元素的
>
> 通过数组下标查询
>
> db.col.find( { "tags.1": { $gt: 25 } } )
>
> 通过数组长度
>
> db.col.find( { "tags": { $size: 3 } } )  // 查询数组长度为3的 
>
> 
>
> 等值匹配
>
> db.inventory.find( { item: null } )  // 返回的是item字段值为null的文档或者不包含item
>
> 类型匹配
>
> db.inventory.find( { item : { $type: 10 } } )  //作为查询条件的时候，仅返回item字段值为null的文档
>
> 存在检查
>
> db.inventory.find( { item : { $exists: false } } )
>
> 
>
> 嵌套查询 
>
> - db.inventory.find( { "instock": { qty: 5, warehouse: "A" } } )
> - db.inventory.find( { "instock.warehouse": "A" ,"instock.qty": 5} )
>
> 



## 查询运算符

https://docs.mongoing.com/can-kao/yun-suan-fu/aggregation-pipeline-operators

> - and : db.col.find({key1:value1, key2:value2}).pretty()  
>
> - or : db.col.find(     {        $or: [  	     {key1: value1}, {key2:value2}        ]     }  ).pretty()  
>
> eg: db.col.find({"likes": {$gt:50}, $or: [{"by": "Mongodb中文网"},{"title": "MongoDB 教程"}]}).pretty()



#### 比较

| 运算符 | 备注                     |
| ------ | ------------------------ |
| $gt    | (>) 大于                 |
| $lt    | (<) 小于                 |
| $gte   | (>=) 大于等于            |
| $lte   | (<= ) 小于等于           |
| $ne    | (!= ) 不等于             |
| $in    | 匹配数组中指定的任何值   |
| $nin   | 不匹配数组中指定的任何值 |



#### 逻辑

| 运算符 | 备注                 |
| ------ | -------------------- |
| $and   | 与                   |
| $or    | 或                   |
| $not   | 非                   |
| $nor   | 或非（多个都不匹配） |



#### 元素

| 运算符  | 备注                           |
| ------- | ------------------------------ |
| $exists | 匹配具有指定字段的文档         |
| $type   | 如果字段是指定类型，则选择文档 |



#### 评估

| 运算符      | 备注                                         |
| ----------- | -------------------------------------------- |
| $expr       | 允许在查询语言中使用聚合表达式               |
| $jsonSchema | 根据给定的JSON Schema验证文档                |
| $mod        | 对字段的值执行模运算并选择具有指定结果的文档 |
| $regex      | 选择值与指定的正则表达式匹配的文档           |
| $text       | 执行文本搜索                                 |
| $where      | 匹配满足JavaScript表达式的文档               |



#### 地理空间

| 运算符  | 备注                           |
| ------- | ------------------------------ |
| $exists | 匹配具有指定字段的文档         |
| $type   | 如果字段是指定类型，则选择文档 |



#### 地理空间

| 运算符         | 备注                                                         |
| -------------- | ------------------------------------------------------------ |
| $geoIntersects | 选择与GeoJSON几何形状相交的几何形状。2dsphere索引支持`$geoIntersects`。 |
| $geoWithin     | 选择边界GeoJSON几何内的几何。2dsphere和2D指标支持 `$geoWithin`。 |
| $near          | 返回点附近的地理空间对象。需要地理空间索引。2dsphere和2D指标支持 `$near`。 |
| $nearSphere    | 返回球体上某个点附近的地理空间对象。需要地理空间索引。2dsphere和2D指标支持 `$nearSphere`。 |



#### 数组

| 运算符     | 备注                                                         |
| ---------- | ------------------------------------------------------------ |
| $all       | 匹配包含查询中指定的所有元素的数组                           |
| $elemMatch | 如果array字段中的元素符合所有指定`$elemMatch`条件，则选择文档。 |
| $size      | 如果数组字段为指定大小，则选择文档                           |



#### 按位

| 运算符        | 备注                                                         |
| ------------- | ------------------------------------------------------------ |
| $bitsAllClear | 匹配数字或二进制值，其中一组位的所有值均为`0`                |
| $bitsAllSet   | 匹配数字或二进制值，其中一组位的所有值均为`1`                |
| $bitsAnyClear | 匹配数值或二进制值，在这些数值或二进制值中，一组位的位置中任何位的值为`0` |
| $bitsAnySet   | 匹配数值或二进制值，在这些数值或二进制值中，一组位的位置中任何位的值为`1`。 |



#### 注释

| 运算符   | 备注           |
| -------- | -------------- |
| $comment | 向查询添加注释 |



#### 映射运算符

| 运算符     | 备注                                            |
| ---------- | ----------------------------------------------- |
| $          | 数组中匹配查询条件的第一个元素                  |
| $elemMatch | 符合指定$elemMatch条件的数组中的第一个元素      |
| $meta      | 项目在$text操作期间分配的文档分数               |
| $slice     | 限制从数组中投影的元素数量。支持`limit`和`skip` |



## 更新运算符

#### 字段

| 运算符       | 备注                                                         |
| ------------ | ------------------------------------------------------------ |
| $currentDate | 将字段的值设置为当前日期，即日期或时间戳                     |
| $inc         | 将字段的值增加指定的数量                                     |
| $min         | 仅当指定值小于现有字段值时才更新该字段                       |
| $max         | 仅当指定值大于现有字段值时才更新该字段                       |
| $mul         | 将字段的值乘以指定的数量                                     |
| $rename      | 重命名字段                                                   |
| $set         | 设置文档中字段的值                                           |
| $setOnInsert | 如果更新导致插入文档，则设置字段的值。对修改现有文档的更新操作没有影响 |
| $unset       | 从文档中删除指定的字段                                       |



#### 数组

| 运算符         | 备注                                                         |
| -------------- | ------------------------------------------------------------ |
| $              | 充当占位符，以更新与查询条件匹配的第一个元素                 |
| $[]            | 充当占位符，以更新匹配查询条件的文档的数组中的所有元素       |
| $[<identical>] | 充当占位符，以更新`arrayFilters`与查询条件匹配的文档中所有与条件匹配的元素 |
| $addToSet      | 仅当元素不存在于集合中时才将它们添加到数组中                 |
| $pop           | 删除数组的第一项或最后一项                                   |
| $pull          | 删除与指定查询匹配的所有数组元素                             |
| $push          | 将项目添加到数组                                             |
| $pullAll       | 从数组中删除所有匹配的值                                     |
|                |                                                              |
| $each          | 修改`$push`和`$addToSet`运算符以附加多个项以进行数组更新     |
| $position      | 修改`$push`运算符以指定要添加元素的数组中的位置              |
| $slice         | 修改`$push`运算符以限制更新数组的大小                        |
| $sort          | 修改`$push`运算符以对存储在数组中的文档重新排序              |



#### 按位

| 运算符 | 备注                                   |
| ------ | -------------------------------------- |
| $bit   | 执行按位`AND`，`OR`和`XOR`整数值的更新 |





