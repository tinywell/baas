package v1

import (
	"baas/internal/model/request"
	"baas/internal/model/response"
	servicev1 "baas/internal/service/v1"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Network ..
var Network = &apiNetwork{}

type apiNetwork struct {
}

// @Summary 网络初始化
// @Description 根据请求参数对网络进行初始化，生成 fabric 网络
// @Produce json
// @Param request body request.NetInit true "网络初始化请求参数"
// @Success 200 {object} response.Response  "code:0 - 网络成功初始化"
// @Failure 500 {object} response.Response "初始化出错"
// @Router /api/v1/network/init [post]
func (an *apiNetwork) Init(c *gin.Context) {
	req := &request.NetInit{}
	if err := c.ShouldBind(req); err != nil {
		Fail(c, response.Fail(1, errors.WithMessage(err, "解析请求参数错误")))
		return
	}
	checkOrderer(&req.GenesisConfig)

	err := servicev1.Network.Init(req)
	if err != nil {
		Fail(c, response.Fail(1, errors.WithMessage(err, "网络初始化错误")))
		return
	}
	OK(c, response.Success("网络初始化成功", nil))
}

// @Summary 网络查询
// @Description 获取网络信息
// @Param network query string true "网络名称"
// @Success 200 {object} response.Response  "返回网络信息"
// @Router /api/v1/network [get]
func (an *apiNetwork) Info(c *gin.Context) {
	OK(c, response.Success("OK", nil))
}

// @Summary 网络创建
// @Description 创建新的网络基础信息
// @Produce json
// @Param request body request.NetCreate true "网络创建"
// @Success 200 {object} response.Response "创建成功"
// @Router /api/v1/network [post]
func (an *apiNetwork) New(c *gin.Context) {
	OK(c, response.Success("OK", nil))
}

func checkOrderer(ocfg *request.OrdererConfig) {
	if len(ocfg.Type) == 0 {
		ocfg.Type = request.DefOrderer.Type
	}
	if ocfg.BatchTimeout == 0 {
		ocfg.BatchTimeout = request.DefOrderer.BatchTimeout
	}
	if ocfg.AbsoluteMaxBytes == 0 {
		ocfg.AbsoluteMaxBytes = request.DefOrderer.AbsoluteMaxBytes
	}
	if ocfg.PreferredMaxBytes == 0 {
		ocfg.PreferredMaxBytes = request.DefOrderer.PreferredMaxBytes
	}
	if ocfg.MaxMessageCount == 0 {
		ocfg.MaxMessageCount = request.DefOrderer.MaxMessageCount
	}
}
