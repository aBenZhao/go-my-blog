package service

import (
	"errors"
	"go-my-blog/internal/DTO"
	"go-my-blog/internal/model"
	"go-my-blog/internal/repo"
	"go-my-blog/internal/response"
	"go-my-blog/pkg/jwt"
	"go-my-blog/pkg/logger"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserSevice struct {
	userRepo *repo.UserRepository
}

func NewUserService(userRepo *repo.UserRepository) *UserSevice {
	return &UserSevice{userRepo: userRepo}
}

func (us *UserSevice) GetUserRepo() *repo.UserRepository {
	return us.userRepo
}

// UserRegister 是用户注册的方法
// 接收一个用户注册的数据传输对象(UserRegisterDTO)作为参数
// 返回注册后的用户数据传输对象和可能的错误
// UserRegister 是用户注册的服务方法
// 参数:
//   - userDTO: 用户注册的数据传输对象，包含用户注册所需的信息
//
// 返回值:
//   - *DTO.UserRegisterDTO: 注册成功后的用户数据传输对象
//   - error: 错误信息，如果注册过程中出现错误则返回
func (us *UserSevice) UserRegister(userDTO *DTO.UserRegisterDTO) (*DTO.UserRegisterDTO, error) {
	// 声明一个用户模型变量，用于存储数据库操作的用户数据
	var user model.User
	// 使用copier将DTO数据复制到用户模型中，为数据库操作做准备
	if err := copier.Copy(&user, &userDTO); err != nil {
		// 如果复制失败，记录错误日志并返回错误
		logger.Error("UserSevice.UserRegister copier.Copy is error!", zap.Error(err))
		return nil, err
	}

	// 复制后检查密码是否为空，确保用户设置了密码
	if user.Password == "" {
		// 如果密码为空，记录警告日志并返回错误
		logger.Warn("UserSevice.UserRegister password is empty", zap.Error(errors.New("密码不能为空")))
		return nil, errors.New("密码不能为空")
	}

	// 使用bcrypt对密码进行哈希处理，确保密码安全存储
	afterPassword, err := us.bcryptPassword(user.Password)
	if err != nil {
		// 如果密码加密失败，记录错误日志并返回错误
		logger.Error("UserSevice.UserRegister bcryptPassword is error!", zap.Error(err))
		return nil, err
	}

	// 将哈希后的密码设置回用户模型，准备存储到数据库
	user.Password = afterPassword

	// 调用用户仓库层的注册方法，执行实际的数据库操作
	userResult, err := us.userRepo.UserRegister(&user)
	// 如果注册失败，记录错误日志并返回错误
	if err != nil {
		logger.Error("UserSevice.UserRegister userRepo.UserRegister is error!", zap.Error(err))
		return nil, err
	}

	// 声明一个用户注册结果的数据传输对象，用于返回给调用方
	var userResultDTO DTO.UserRegisterDTO
	// 将用户模型数据复制到结果DTO中，确保返回符合API规范的数据结构
	if copErr := copier.Copy(&userResultDTO, &userResult); copErr != nil {
		// 如果复制失败，记录错误日志并返回错误
		logger.Error("UserSevice.UserRegister copier.Copy is error!", zap.Error(copErr))
		return nil, copErr
	}
	// 返回注册成功的用户数据
	return &userResultDTO, nil
}

// bcryptPassword 使用bcrypt算法对密码进行哈希处理
// 参数:
//   - password: 原始密码字符串
//
// 返回值:
//   - string: 哈希处理后的密码字符串
//   - error: 处理过程中可能出现的错误
func (us *UserSevice) bcryptPassword(password string) (string, error) {
	// 使用bcrypt算法生成密码的哈希值
	// bcrypt.GenerateFromPassword函数将密码转换为哈希值
	// 参数说明:
	//   - []byte(password): 将密码转换为字节数组
	//   - bcrypt.DefaultCost: 使用bcrypt的默认成本值
	afterPassword, bcrErr := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost)
	// 如果生成哈希失败，记录错误日志并返回错误
	// 使用zap记录错误日志，包含错误信息
	if bcrErr != nil {
		logger.Error("UserSevice bcryptPassword is error", zap.Error(bcrErr))
		return "", bcrErr
	}

	// 返回哈希处理后的密码和错误信息
	// 正常情况下错误为nil
	return string(afterPassword), bcrErr
}

func (us *UserSevice) UserLogin(d *DTO.LoginDTO) (*response.LoginResponse, error) {
	user, err := us.userRepo.FindByUserName(d.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误！")
	}

	if user == nil {
		return nil, errors.New("用户名或密码错误！")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(d.Password),
	)

	if err != nil {
		return nil, errors.New("用户名或密码错误！")
	}

	token, expiresAt, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("生成令牌失败！")
	}

	return &response.LoginResponse{AccessToken: token, Username: user.Username, ExpiresAt: expiresAt}, nil
}
