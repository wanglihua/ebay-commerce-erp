package pass

// 有一种做法是登录时把用户的权限全部缓存起来，这样下次检查权限的时候，就不用查询数据库了
// 这里先用不缓存的方法来做

func UserHasPass(pass_id int) bool {
	//获得当前登录用户的Id

	return true
}

func CheckUserPass(pass_id int) {
	panic("你没有当前权限！")
}
