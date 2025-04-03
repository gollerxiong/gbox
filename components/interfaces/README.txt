

接口定义方法返回参数的设定：
1.假如是set开头的方法，一般都具有连贯性，方法返回结构体的实例
  比如：
  type A interface(
    SetName(string) *A
  )

  如果对象B需要实现这个接口，为了能够使得动作是链式的，它可以直接在SetName()后再直接调用其它方法
  type B struct{
     name string
     age int
  }

  func (b *B) SetName(name string) *B {
    b.name = name
    return b
  }

  func (b *B) SetAge(age int) *B {
    b.age = age
    return b
  }

  func (b *B) Desc() string {
    return "My name is " + b.name + " and I'm 20"
  }

  func NewB() *B {
    return &B{}
  }

  调用的时候就可以直接链式调用
  myDescribe := NewB().SetName("bob").SetAge(20).Desc()
  fmt.Printf(myDescribe)

2. 假如是get开头的方法，如果返回结果不跟特定的对象绑定，那么就随心所欲，想返回什么就返回什么！但是需要跟特定对象绑定的时候，返回参数定义为interface{},
这样如果要跟特定对象绑定的时候可以使用门面的方法将返回结果跟特定的对象绑定！
    比如：
  参考baseObject.go文件里面的LoadById()方法，该方法返回一个interface{}对象，主要是方便具体业务继承了baseObject.go方法时，可以复写loadById()方法
来跟具体的对象绑定，方便后续参数的读取和理解。
  举个例子，比如有用户表，我们可以声明用户结构体如下所示
  type UserInfo struct {
  	components.BaseObject
  }

  func (u *UserInfo) GetModel() *models.UserInfoModel {
  	return u.BaseObject.GetModel().(*models.UserInfoModel)
  }

  func (u *UserInfo) LoadById(id int64) *UserInfo {
  	u.BaseObject.LoadById(id)
  	return u
  }

  func (u *UserInfo) LoadByUsername(username string) *UserInfo {
  	model := u.GetModel()
  	connect := model.GetConnect()
  	err := connect.Where("user_name = ?", username).First(u.GetModel()).Error

  	if err == gorm.ErrRecordNotFound {
  		u.SetNew(true)
  	} else {
  		u.SetNew(false)
  		u.SetModel(model)
  		u.SetOldModel(model)
  		u.SetAttributes(components.StructToMap(model))
  		u.SetOldAttributes(components.StructToMap(model))
  	}

  	return u
  }

  /*
  * 实例化对象并注入依赖
   */
  func NewUserInfo() *UserInfo {
  	ins := &UserInfo{}
  	ins.SetField("*")
  	ins.SetModel(models.NewUserInfoModel())
  	ins.SetFormatter(NewUserInfoFormatter())
  	ins.SetHooks(NewUserInfoHooks())
  	return ins
  }


  // 在业务调用的时候，可以如下调用
  userInfoLib := userInfoDao.NewUserInfo().LoadById(int64(10027))
  userInfoLib.GetModel().Icon = "test"
  res := userInfoLib.Save()

  fmt.Println(res)

  这样子我们就可以直接是跟用户表的模型打交道，编辑器也能够识别具体的参数提示。