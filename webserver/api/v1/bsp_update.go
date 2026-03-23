package v1

import (
	"os"
	"strconv"

	"github.com/edgehook/ithings/common/dbm/model"
	responce "github.com/edgehook/ithings/webserver/types"
	"github.com/gin-gonic/gin"
)

func GetBspByPage(c *gin.Context) {
	var (
		err   error
		bsps  []*model.Bsp
		count int64
	)
	keywords := c.Query("keywords")
	currentPage := c.Query("currentPage")
	limit := c.Query("limit")
	if currentPage == "" || limit == "" {
		responce.FailWithMessage("Parameter error", c)
		return
	}

	pageInt, err := strconv.Atoi(currentPage)
	if err != nil {
		responce.FailWithMessage("Parameter error", c)
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		responce.FailWithMessage("Parameter error", c)
		return
	}
	if keywords == "" {
		if bsps, err = model.GetBspByPage(pageInt, limitInt); err != nil {
			responce.FailWithMessage("Get alert history error", c)
			return
		}
		if count, err = model.GetBspCount(); err != nil {
			responce.FailWithMessage("Get count error", c)
			return
		}
	} else {
		if bsps, err = model.GetBspByPageAndKeywords(pageInt, limitInt, keywords); err != nil {
			responce.FailWithMessage("Get alert history error", c)
			return
		}
		if count, err = model.GetBspCountByKeywords(keywords); err != nil {
			responce.FailWithMessage("Get count error", c)
			return
		}
	}

	if err != nil {
		responce.FailWithMessage("Get alert log error", c)
		return
	}

	responce.OkWithData(map[string]interface{}{
		"list":  bsps,
		"total": count,
	}, c)
}

func DeleteBsp(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		responce.FailWithMessage("Parameter error", c)
		return
	}
	bsp, err := model.GetBspById(id)
	if err != nil {
		responce.FailWithMessage("Get bsp error", c)
		return
	}
	if err := model.DeleteBsp(id); err != nil {
		responce.FailWithMessage("Delete db error", c)
		return
	}
	os.RemoveAll(bsp.Path)
	responce.Ok(c)
}
