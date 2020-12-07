学习笔记
1. 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

答：不应该。sql.ErrNoRows是database/sql包中一个预定义error，使用Row.Scan返回这个错误，表示QueryRow无法查询到数据库记录。

对于dao层这个错误相当于无法找到数据，如果我们在dao层直接Wrap这个错误返回，那上层代码就可能使用如下代码进行判断

```
if err==sql.ErrNoRows{
    ...
}
```
或者
```
if errors.Is(err,sql.ErrNoRows){
    ...
}
```
这样，上层逻辑处理就会和database/sql包耦合，不利于分层和解耦，日后dao层想要替换其他实现方式，势必会影响到上层逻辑代码。

所以，我认为直接Wrap这个错误抛给上层，不是个好的做法。更好的做法，应该在dao层将这个错误转换成一个与实现无关的错误，比如：dao.ErrRecordNotFound，然后Wrap抛给上层。