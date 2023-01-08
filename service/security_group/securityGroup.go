package securityGroup

import (
	alb "go-sdk/service/security_group/alb"
	docdb "go-sdk/service/security_group/docdb"
	elasticache "go-sdk/service/security_group/elasticache"
	msk "go-sdk/service/security_group/msk"
	rds "go-sdk/service/security_group/rds"

	"github.com/gofiber/fiber/v2"
)

func GetSecurityGroup(c *fiber.Ctx) error {
	resourceType := c.Query("type", "")
	resourceName := c.Query("name", "")
	switch resourceType {
	case "alb":
		err, alb := alb.GetSecurityGroup(resourceName)
		if err != nil {
			return c.JSON(fiber.Map{"status": "404", "message": err.Error()})
		}
		return c.JSON(alb)
	case "docdb":
		err, docdb := docdb.GetSecurityGroup(resourceName)
		if err != nil {
			return c.JSON(fiber.Map{"status": "404", "message": err.Error()})
		}
		return c.JSON(docdb)
	case "msk":
		err, msk := msk.GetSecurityGroup(resourceName)
		if err != nil {
			return c.JSON(fiber.Map{"status": "404", "message": err.Error()})
		}
		return c.JSON(msk)
	case "rds":
		err, rds := rds.GetSecurityGroup(resourceName)
		if err != nil {
			return c.JSON(fiber.Map{"status": "404", "message": err.Error()})
		}
		return c.JSON(rds)
	case "elasticache":
		err, elasticache := elasticache.GetSecurityGroup(resourceName)
		if err != nil {
			return c.JSON(fiber.Map{"status": "404", "message": err.Error()})
		}
		return c.JSON(elasticache)
	}

	return c.JSON(fiber.Map{"status": "404", "message": "옳바른 리소스 타입을 입력해주세요(alb, elasticache, msk, rds, docdb, rds)"})
}
